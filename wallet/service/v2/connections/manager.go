package connections

import (
	"context"
	"fmt"
	"sort"
	"sync"

	vgcrypto "zuluprotocol/zeta/libs/crypto"
	"zuluprotocol/zeta/wallet/api"
	"zuluprotocol/zeta/wallet/wallet"
)

// Manager holds the opened connections between the third-party
// applications and the wallets.
type Manager struct {
	// tokenToConnection maps the token to the connection. It is the base
	// registry for all opened connections.
	tokenToConnection map[Token]*walletConnection

	// sessionFingerprintToToken maps the session fingerprint (wallet + hostname
	// on a session connection). This is used to determine whether a session
	// connection already exists for a given wallet and hostname.
	// It only holds sessions fingerprints.
	// Long-living connections are not tracked.
	sessionFingerprintToToken map[string]Token

	// timeService is used to resolve the current time to update the last activity
	// time on the token, and figure out their expiration.
	timeService TimeService

	walletStore WalletStore

	tokenStore TokenStore

	mu sync.Mutex
}

type walletConnection struct {
	// connectedWallet is the projection of the wallet through the permissions
	// and authentication system. On a regular wallet, there are no restrictions
	// on what we can call, which doesn't fit the model of having allowed
	// access, so we wrap the "regular wallet" behind the "connected wallet".
	connectedWallet api.ConnectedWallet

	policy connectionPolicy
}

// StartSession initializes a connection between a wallet and a third-party
// application.
// If a connection already exists, it's disconnected and a new token is
// generated.
func (m *Manager) StartSession(hostname string, w wallet.Wallet) (Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cw, err := api.NewConnectedWallet(hostname, w)
	if err != nil {
		return "", fmt.Errorf("could not instantiate the connected wallet for a session connection: %w", err)
	}

	newToken := m.generateToken()

	sessionFingerprint := asSessionFingerprint(hostname, w.Name())
	if previousToken, sessionAlreadyExists := m.sessionFingerprintToToken[sessionFingerprint]; sessionAlreadyExists {
		m.destroySessionToken(previousToken)
	}
	m.sessionFingerprintToToken[sessionFingerprint] = newToken

	m.tokenToConnection[newToken] = &walletConnection{
		connectedWallet: cw,
		policy: &sessionPolicy{
			lastActivityDate: m.timeService.Now(),
		},
	}

	return newToken, nil
}

func (m *Manager) EndSessionConnectionWithToken(token Token) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.destroySessionToken(token)
}

func (m *Manager) EndSessionConnection(hostname, walletName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fingerprint := asSessionFingerprint(hostname, walletName)

	token, exists := m.sessionFingerprintToToken[fingerprint]
	if !exists {
		return
	}

	m.destroySessionToken(token)
}

func (m *Manager) EndAllSessionConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for token := range m.tokenToConnection {
		m.destroySessionToken(token)
	}
}

// ConnectedWallet retrieves the wallet associated to the specified token.
func (m *Manager) ConnectedWallet(hostname string, token Token) (api.ConnectedWallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	connection, exists := m.tokenToConnection[token]
	if !exists {
		return api.ConnectedWallet{}, ErrNoConnectionAssociatedThisAuthenticationToken
	}

	if !connection.policy.IsLongLivingConnection() && connection.connectedWallet.Hostname() != hostname {
		return api.ConnectedWallet{}, ErrHostnamesMismatchForThisToken
	}

	now := m.timeService.Now()

	if connection.policy.HasConnectionExpired(now) {
		return api.ConnectedWallet{}, ErrTokenHasExpired
	}

	if connection.policy.IsClosed() {
		return api.ConnectedWallet{}, ErrConnectionHasBeenClosed
	}

	connection.policy.UpdateActivityDate(now)

	return connection.connectedWallet, nil
}

// ListSessionConnections lists all the session connections as a list of pairs of
// hostname/wallet name.
// The list is sorted, first, by hostname, and, then, by wallet name.
func (m *Manager) ListSessionConnections() []api.Connection {
	m.mu.Lock()
	defer m.mu.Unlock()

	connections := []api.Connection{}
	for _, connection := range m.tokenToConnection {
		if connection.policy.IsLongLivingConnection() {
			continue
		}

		connections = append(connections, api.Connection{
			Hostname: connection.connectedWallet.Hostname(),
			Wallet:   connection.connectedWallet.Name(),
		})
	}

	sort.SliceStable(connections, func(i, j int) bool {
		if connections[i].Hostname == connections[j].Hostname {
			return connections[i].Wallet < connections[j].Wallet
		}

		return connections[i].Hostname < connections[j].Hostname
	})

	return connections
}

// generateToken generates a new token and ensure it is not already in use to
// avoid collisions.
func (m *Manager) generateToken() Token {
	for {
		token := GenerateToken()
		if _, alreadyExistingToken := m.tokenToConnection[token]; !alreadyExistingToken {
			return token
		}
	}
}

func (m *Manager) destroySessionToken(tokenToDestroy Token) {
	connection, exists := m.tokenToConnection[tokenToDestroy]
	if !exists || connection.policy.IsLongLivingConnection() {
		return
	}

	// Remove the session fingerprint associated to the session token.
	for sessionFingerprint, t := range m.sessionFingerprintToToken {
		if t == tokenToDestroy {
			delete(m.sessionFingerprintToToken, sessionFingerprint)
			break
		}
	}

	// Break the link between a token and its associated wallet.
	m.tokenToConnection[tokenToDestroy] = nil
	delete(m.tokenToConnection, tokenToDestroy)
}

func (m *Manager) loadLongLivingConnections() error {
	ctx := context.Background()

	tokenSummaries, err := m.tokenStore.ListTokens()
	if err != nil {
		return err
	}

	for _, tokenSummary := range tokenSummaries {
		tokenDescription, err := m.tokenStore.DescribeToken(tokenSummary.Token)
		if err != nil {
			return fmt.Errorf("could not get information associated to the token %q: %w", tokenDescription.Token.Short(), err)
		}

		if err := m.loadLongLivingConnection(ctx, tokenDescription); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) loadLongLivingConnection(ctx context.Context, tokenDescription TokenDescription) error {
	if err := m.walletStore.UnlockWallet(ctx, tokenDescription.Wallet.Name, tokenDescription.Wallet.Passphrase); err != nil {
		// We don't properly handle wallets renaming, nor wallets passphrase
		// update in the token file automatically. We only support a direct
		// update of the token file.
		return fmt.Errorf("could not unlock the wallet %q associated to the token %q: %w",
			tokenDescription.Wallet.Name,
			tokenDescription.Token.Short(),
			err)
	}

	w, err := m.walletStore.GetWallet(ctx, tokenDescription.Wallet.Name)
	if err != nil {
		// This should not happen because we just unlocked the wallet.
		return fmt.Errorf("could not get the information for the wallet %q associated to the token %q: %w",
			tokenDescription.Wallet.Name,
			tokenDescription.Token.Short(),
			err)
	}

	m.tokenToConnection[tokenDescription.Token] = &walletConnection{
		connectedWallet: api.NewLongLivingConnectedWallet(w),
		policy: &longLivingConnectionPolicy{
			expirationDate: tokenDescription.ExpirationDate,
		},
	}

	return nil
}

// refreshConnections is called when the wallet store notices a change in
// the wallets. This way the connection manager is able to reload the connected
// wallets.
func (m *Manager) refreshConnections(_ context.Context, event wallet.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if event.Type == wallet.WalletRemovedEventType {
		m.closeConnectionsUsingThisWallet(event.Data.(wallet.WalletRemovedEventData).Name)
	} else if event.Type == wallet.UnlockedWalletUpdatedEventType {
		m.updateConnectionsUsingThisWallet(event.Data.(wallet.UnlockedWalletUpdatedEventData).UpdatedWallet)
	} else if event.Type == wallet.WalletHasBeenLockedEventType {
		m.expireConnectionsUsingThisWallet(event.Data.(wallet.WalletHasBeenLockedEventData).Name)
	} else if event.Type == wallet.WalletRenamedEventType {
		data := event.Data.(wallet.WalletRenamedEventData)
		m.updateConnectionsUsingThisRenamedWallet(data.PreviousName, data.NewName)
	}
}

func (m *Manager) closeConnectionsUsingThisWallet(walletName string) {
	for token, connection := range m.tokenToConnection {
		if connection.connectedWallet.Name() != walletName {
			continue
		}

		connection.policy.SetAsForcefullyClose()

		delete(m.tokenToConnection, token)
	}
}

func (m *Manager) updateConnectionsUsingThisWallet(w wallet.Wallet) {
	for token, connection := range m.tokenToConnection {
		if connection.connectedWallet.Name() != w.Name() {
			continue
		}

		var updatedConnectedWallet api.ConnectedWallet
		if connection.policy.IsLongLivingConnection() {
			updatedConnectedWallet = api.NewLongLivingConnectedWallet(w)
		} else {
			updatedConnectedWallet, _ = api.NewConnectedWallet(connection.connectedWallet.Hostname(), w)
		}

		m.tokenToConnection[token].connectedWallet = updatedConnectedWallet
	}
}

func (m *Manager) refreshLongLivingTokens(ctx context.Context, activeTokensDescriptions ...TokenDescription) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// We need to find the new tokens among the active ones, so, we build a
	// registry with all the active tokens. Then, we remove the tracked token
	// when found. We will end up with the new tokens, only.
	activeTokens := map[Token]TokenDescription{}
	for _, tokenDescription := range activeTokensDescriptions {
		activeTokens[tokenDescription.Token] = tokenDescription
	}

	isActiveToken := func(token Token) (TokenDescription, bool) {
		activeTokenDescription, isTracked := activeTokens[token]
		if isTracked {
			delete(activeTokens, token)
		}
		return activeTokenDescription, isTracked
	}

	// First, we address the update of the tokens we already track.
	for token, connection := range m.tokenToConnection {
		if !connection.policy.IsLongLivingConnection() {
			continue
		}

		activeToken, isActive := isActiveToken(token)
		if !isActive {
			// If the token could not be found in the active tokens, this means
			// the token has been deleted from the token store. Thus, we close the
			// connection.
			delete(m.tokenToConnection, token)
			continue
		}

		_ = m.loadLongLivingConnection(ctx, activeToken)
	}

	for _, tokenDescription := range activeTokens {
		_ = m.loadLongLivingConnection(ctx, tokenDescription)
	}
}

func (m *Manager) expireConnectionsUsingThisWallet(walletName string) {
	for _, connection := range m.tokenToConnection {
		if connection.connectedWallet.Name() != walletName {
			continue
		}

		connection.policy.SetAsForcefullyClose()
	}
}

func (m *Manager) updateConnectionsUsingThisRenamedWallet(previousWalletName, newWalletName string) {
	var _updatedWallet wallet.Wallet

	// This acts as a cached getter, to avoid multiple or useless fetch.
	getUpdatedWallet := func() wallet.Wallet {
		if _updatedWallet == nil {
			w, _ := m.walletStore.GetWallet(context.Background(), newWalletName)
			_updatedWallet = w
		}
		return _updatedWallet
	}

	for _, connection := range m.tokenToConnection {
		if connection.connectedWallet.Name() != previousWalletName {
			continue
		} else if connection.policy.IsLongLivingConnection() {
			// We don't update the long-living connections because they come
			// from a configuration file that doesn't support token updates, for
			// now.
			continue
		}

		connection.connectedWallet, _ = api.NewConnectedWallet(connection.connectedWallet.Hostname(), getUpdatedWallet())
	}
}

func NewManager(timeService TimeService, walletStore WalletStore, tokenStore TokenStore) (*Manager, error) {
	m := &Manager{
		sessionFingerprintToToken: map[string]Token{},
		tokenToConnection:         map[Token]*walletConnection{},
		timeService:               timeService,
		walletStore:               walletStore,
		tokenStore:                tokenStore,
	}

	walletStore.OnUpdate(m.refreshConnections)
	tokenStore.OnUpdate(m.refreshLongLivingTokens)

	if err := m.loadLongLivingConnections(); err != nil {
		return nil, fmt.Errorf("could not load the long-living connections: %w", err)
	}

	return m, nil
}

func asSessionFingerprint(hostname string, walletName string) string {
	return vgcrypto.HashStrToHex(hostname + "::" + walletName)
}
