// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package sqlstore

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	"code.zetaprotocol.io/vega/datanode/entities"
	"code.zetaprotocol.io/vega/datanode/metrics"
	v2 "code.zetaprotocol.io/vega/protos/data-node/api/v2"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

var accountOrdering = TableOrdering{
	ColumnOrdering{Name: "account_id", Sorting: ASC},
}

type Accounts struct {
	*ConnectionSource
	idToAccount map[entities.AccountID]entities.Account
	cacheLock   sync.RWMutex
}

func NewAccounts(connectionSource *ConnectionSource) *Accounts {
	a := &Accounts{
		ConnectionSource: connectionSource,
		idToAccount:      make(map[entities.AccountID]entities.Account),
	}
	return a
}

// Add inserts a row and updates supplied struct with autogenerated ID.
func (as *Accounts) Add(ctx context.Context, a *entities.Account) error {
	defer metrics.StartSQLQuery("Accounts", "Add")()

	err := as.Connection.QueryRow(ctx,
		`INSERT INTO accounts(id, party_id, asset_id, market_id, type, tx_hash, zeta_time)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		deterministicIDFromAccount(a),
		a.PartyID,
		a.AssetID,
		a.MarketID,
		a.Type,
		a.TxHash,
		a.ZetaTime).Scan(&a.ID)
	return err
}

func (as *Accounts) GetByID(ctx context.Context, accountID entities.AccountID) (entities.Account, error) {
	if account, ok := as.getAccountFromCache(accountID); ok {
		return account, nil
	}

	as.cacheLock.Lock()
	defer as.cacheLock.Unlock()

	// It's possible that in-between releasing the read lock and obtaining the write lock that the account has been
	// added to cache, so we need to check here and return the cached account if that's the case.
	if account, ok := as.idToAccount[accountID]; ok {
		return account, nil
	}

	a := entities.Account{}
	defer metrics.StartSQLQuery("Accounts", "GetByID")()

	if err := pgxscan.Get(ctx, as.Connection, &a,
		`SELECT id, party_id, asset_id, market_id, type, tx_hash, zeta_time
		 FROM accounts WHERE id=$1`,
		accountID,
	); err != nil {
		return a, as.wrapE(err)
	}

	as.idToAccount[accountID] = a
	return a, nil
}

func (as *Accounts) GetAll(ctx context.Context) ([]entities.Account, error) {
	accounts := []entities.Account{}
	defer metrics.StartSQLQuery("Accounts", "GetAll")()
	err := pgxscan.Select(ctx, as.Connection, &accounts, `
		SELECT id, party_id, asset_id, market_id, type, tx_hash, zeta_time
		FROM accounts`)
	return accounts, err
}

// Obtain will either fetch or create an account in the database.
// If an account with matching party/asset/market/type does not exist in the database, create one.
// If an account already exists, fetch that one.
// In either case, update the entities.Account object passed with an ID from the database.
func (as *Accounts) Obtain(ctx context.Context, a *entities.Account) error {
	accountID := deterministicIDFromAccount(a)
	if account, ok := as.getAccountFromCache(accountID); ok {
		a.ID = account.ID
		a.ZetaTime = account.VegaTime
		a.TxHash = account.TxHash
		return nil
	}

	as.cacheLock.Lock()
	defer as.cacheLock.Unlock()

	// It's possible that in-between releasing the cache read lock and obtaining the cache write lock that the account has been
	// added to the cache, so we need to check here and return the cached account if that's the case.
	if account, ok := as.idToAccount[accountID]; ok {
		a.ID = account.ID
		a.ZetaTime = account.VegaTime
		a.TxHash = account.TxHash
		return nil
	}

	insertQuery := `INSERT INTO accounts(id, party_id, asset_id, market_id, type, tx_hash, zeta_time)
                           VALUES ($1, $2, $3, $4, $5, $6, $7)
                           ON CONFLICT (party_id, asset_id, market_id, type) DO NOTHING`

	selectQuery := `SELECT id, party_id, asset_id, market_id, type, tx_hash, zeta_time
	                FROM accounts
	                WHERE party_id=$1 AND asset_id=$2 AND market_id=$3 AND type=$4`

	batch := pgx.Batch{}

	batch.Queue(insertQuery, accountID, a.PartyID, a.AssetID, a.MarketID, a.Type, a.TxHash, a.ZetaTime)
	batch.Queue(selectQuery, a.PartyID, a.AssetID, a.MarketID, a.Type)
	defer metrics.StartSQLQuery("Accounts", "Obtain")()
	results := as.Connection.SendBatch(ctx, &batch)
	defer results.Close()

	if _, err := results.Exec(); err != nil {
		return fmt.Errorf("inserting account: %w", err)
	}

	rows, err := results.Query()
	if err != nil {
		return fmt.Errorf("querying accounts: %w", err)
	}

	if err = pgxscan.ScanOne(a, rows); err != nil {
		return fmt.Errorf("scanning account: %w", err)
	}

	as.idToAccount[accountID] = *a
	return nil
}

func (as *Accounts) getAccountFromCache(id entities.AccountID) (entities.Account, bool) {
	as.cacheLock.RLock()
	defer as.cacheLock.RUnlock()

	if account, ok := as.idToAccount[id]; ok {
		return account, true
	}
	return entities.Account{}, false
}

func deterministicIDFromAccount(a *entities.Account) entities.AccountID {
	idAsBytes := sha256.Sum256([]byte(a.AssetID.String() + a.PartyID.String() + a.MarketID.String() + a.Type.String()))
	accountID := hex.EncodeToString(idAsBytes[:])
	return entities.AccountID(accountID)
}

func (as *Accounts) Query(ctx context.Context, filter entities.AccountFilter) ([]entities.Account, error) {
	query, args, err := filterAccountsQuery(filter, true)
	if err != nil {
		return nil, err
	}
	accs := []entities.Account{}

	defer metrics.StartSQLQuery("Accounts", "Query")()
	rows, err := as.Connection.Query(ctx, query, args...)
	if err != nil {
		return accs, fmt.Errorf("querying accounts: %w", err)
	}
	defer rows.Close()

	if err = pgxscan.ScanAll(&accs, rows); err != nil {
		return accs, fmt.Errorf("scanning account: %w", err)
	}

	return accs, nil
}

func (as *Accounts) QueryBalancesV1(ctx context.Context, filter entities.AccountFilter, pagination entities.OffsetPagination) ([]entities.AccountBalance, error) {
	query, args, err := filterAccountBalancesQuery(filter)
	if err != nil {
		return nil, fmt.Errorf("querying account balances: %w", err)
	}

	query, args = orderAndPaginateQuery(query, nil, pagination, args...)

	accountBalances := make([]entities.AccountBalance, 0)

	defer metrics.StartSQLQuery("Accounts", "QueryBalancesV1")()
	rows, err := as.Connection.Query(ctx, query, args...)
	if err != nil {
		return accountBalances, fmt.Errorf("querying account balances: %w", err)
	}
	defer rows.Close()

	if err = pgxscan.ScanAll(&accountBalances, rows); err != nil {
		return accountBalances, fmt.Errorf("parsing account balances: %w", err)
	}

	return accountBalances, nil
}

func (as *Accounts) QueryBalances(ctx context.Context,
	filter entities.AccountFilter,
	pagination entities.CursorPagination,
) ([]entities.AccountBalance, entities.PageInfo, error) {
	query, args, err := filterAccountBalancesQuery(filter)
	if err != nil {
		return nil, entities.PageInfo{}, fmt.Errorf("querying account balances: %w", err)
	}

	query, args, err = PaginateQuery[entities.AccountCursor](query, args, accountOrdering, pagination)
	if err != nil {
		return nil, entities.PageInfo{}, fmt.Errorf("querying account balances: %w", err)
	}

	defer metrics.StartSQLQuery("Accounts", "QueryBalances")()

	accountBalances := make([]entities.AccountBalance, 0)
	rows, err := as.Connection.Query(ctx, query, args...)
	if err != nil {
		return accountBalances, entities.PageInfo{}, fmt.Errorf("querying account balances: %w", err)
	}
	defer rows.Close()

	if err = pgxscan.ScanAll(&accountBalances, rows); err != nil {
		return accountBalances, entities.PageInfo{}, fmt.Errorf("parsing account balances: %w", err)
	}

	pagedAccountBalances, pageInfo := entities.PageEntities[*v2.AccountEdge](accountBalances, pagination)
	return pagedAccountBalances, pageInfo, nil
}
