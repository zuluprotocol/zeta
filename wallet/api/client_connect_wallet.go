package api

import (
	"context"
	"errors"
	"fmt"

	"zuluprotocol/zeta/libs/jsonrpc"
	"zuluprotocol/zeta/wallet/preferences"
	"zuluprotocol/zeta/wallet/wallet"
)

const WalletConnectionSuccessfullyEstablished = "The connection to the wallet has been successfully established."

type ClientConnectWallet struct {
	walletStore WalletStore
	interactor  Interactor
}

// Handle initiates the wallet connection between the API and a third-party
// application.
//
// It triggers a selection of the wallet the client wants to use for this
// connection. The wallet is then loaded in memory. All changes done to that wallet
// will start in-memory, and then, be saved in the wallet file. Any changes done
// to the wallet outside the JSON-RPC session (via the command-line for example)
// will be overridden. For the effects to be taken into account, the wallet has
// to be disconnected first, and then re-connected.
//
// All sessions have to be initialized by using this handler. Otherwise, a call
// to any other handlers will be rejected.
func (h *ClientConnectWallet) Handle(ctx context.Context, hostname string) (wallet.Wallet, *jsonrpc.ErrorDetails) {
	traceID := jsonrpc.TraceIDFromContext(ctx)

	if err := h.interactor.NotifyInteractionSessionBegan(ctx, traceID); err != nil {
		return nil, requestNotPermittedError(err)
	}
	defer h.interactor.NotifyInteractionSessionEnded(ctx, traceID)

	availableWallets, err := h.walletStore.ListWallets(ctx)
	if err != nil {
		h.interactor.NotifyError(ctx, traceID, InternalError, fmt.Errorf("could not list the available wallets: %w", err))
		return nil, internalError(ErrCouldNotConnectToWallet)
	}
	if len(availableWallets) == 0 {
		h.interactor.NotifyError(ctx, traceID, ApplicationError, ErrNoWalletToConnectTo)
		return nil, applicationCancellationError(ErrApplicationCancelledTheRequest)
	}

	var approval preferences.ConnectionApproval
	for {
		rawApproval, err := h.interactor.RequestWalletConnectionReview(ctx, traceID, hostname)
		if err != nil {
			if errDetails := handleRequestFlowError(ctx, traceID, h.interactor, err); errDetails != nil {
				return nil, errDetails
			}
			h.interactor.NotifyError(ctx, traceID, InternalError, fmt.Errorf("reviewing the wallet connection failed: %w", err))
			return nil, internalError(ErrCouldNotConnectToWallet)
		}

		a, err := preferences.ParseConnectionApproval(rawApproval)
		if err != nil {
			h.interactor.NotifyError(ctx, traceID, UserError, err)
			continue
		}
		approval = a
		break
	}

	if isConnectionRejected(approval) {
		return nil, userRejectionError(ErrUserRejectedWalletConnection)
	}

	// Wallet selection process.
	var loadedWallet wallet.Wallet
	for {
		if ctx.Err() != nil {
			return nil, requestInterruptedError(ErrRequestInterrupted)
		}

		selectedWallet, err := h.interactor.RequestWalletSelection(ctx, traceID, hostname, availableWallets)
		if err != nil {
			if errDetails := handleRequestFlowError(ctx, traceID, h.interactor, err); errDetails != nil {
				return nil, errDetails
			}
			h.interactor.NotifyError(ctx, traceID, InternalError, fmt.Errorf("requesting the wallet selection failed: %w", err))
			return nil, internalError(ErrCouldNotConnectToWallet)
		}

		if exist, err := h.walletStore.WalletExists(ctx, selectedWallet.Wallet); err != nil {
			h.interactor.NotifyError(ctx, traceID, InternalError, fmt.Errorf("could not verify the wallet existence: %w", err))
			return nil, internalError(ErrCouldNotConnectToWallet)
		} else if !exist {
			h.interactor.NotifyError(ctx, traceID, UserError, ErrWalletDoesNotExist)
			continue
		}

		if err := h.walletStore.UnlockWallet(ctx, selectedWallet.Wallet, selectedWallet.Passphrase); err != nil {
			if errors.Is(err, wallet.ErrWrongPassphrase) {
				h.interactor.NotifyError(ctx, traceID, UserError, wallet.ErrWrongPassphrase)
				continue
			}
			h.interactor.NotifyError(ctx, traceID, InternalError, fmt.Errorf("could not unlock the wallet: %w", err))
			return nil, internalError(ErrCouldNotConnectToWallet)
		}

		w, err := h.walletStore.GetWallet(ctx, selectedWallet.Wallet)
		if err != nil {
			h.interactor.NotifyError(ctx, traceID, InternalError, fmt.Errorf("could not retrieve the wallet: %w", err))
			return nil, internalError(ErrCouldNotConnectToWallet)
		}

		loadedWallet = w
		break
	}

	h.interactor.NotifySuccessfulRequest(ctx, traceID, WalletConnectionSuccessfullyEstablished)

	return loadedWallet, nil
}

func isConnectionRejected(approval preferences.ConnectionApproval) bool {
	return approval != preferences.ApprovedOnlyThisTime
}

func NewConnectWallet(walletStore WalletStore, interactor Interactor) *ClientConnectWallet {
	return &ClientConnectWallet{
		walletStore: walletStore,
		interactor:  interactor,
	}
}
