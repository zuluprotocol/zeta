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

package sqlstore_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"zuluprotocol/zeta/datanode/entities"
	"zuluprotocol/zeta/datanode/sqlstore"
	"zuluprotocol/zeta/logging"
	"zuluprotocol/zeta/protos/zeta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiquidityProvision(t *testing.T) {
	t.Run("Upsert should insert a liquidity provision record if the id doesn't exist in the current block", testInsertNewInCurrentBlock)
	t.Run("Upsert should update a liquidity provision record if the id already exists in the current block", testUpdateExistingInCurrentBlock)
	t.Run("Get should return all LP for a given party if no market is provided", testGetLPByPartyOnly)
	t.Run("Get should return all LP for a given party and market if both are provided", testGetLPByPartyAndMarket)
	t.Run("Get should error if no party and market are provided", testGetLPNoPartyAndMarketErrors)
	t.Run("Get should return all LP for a given market if no party id is provided", testGetLPNoPartyWithMarket)
	t.Run("Get should return LP with the corresponding reference", testGetLPByReferenceAndParty)
}

func TestLiquidityProvisionPagination(t *testing.T) {
	t.Run("should return all liquidity provisions if no pagination is specified", testLiquidityProvisionPaginationNoPagination)
	t.Run("should return the first page of results if first is provided", testLiquidityProvisionPaginationFirst)
	t.Run("should return the last page of results if last is provided", testLiquidityProvisionPaginationLast)
	t.Run("should return the specified page of results if first and after are provided", testLiquidityProvisionPaginationFirstAfter)
	t.Run("should return the specified page of results if last and before are provided", testLiquidityProvisionPaginationLastBefore)
}

func setupLPTests(t *testing.T) (*sqlstore.Blocks, *sqlstore.LiquidityProvision, sqlstore.Connection) {
	t.Helper()

	bs := sqlstore.NewBlocks(connectionSource)
	lp := sqlstore.NewLiquidityProvision(connectionSource, logging.NewTestLogger())

	return bs, lp, connectionSource.Connection
}

func testInsertNewInCurrentBlock(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, lp, conn := setupLPTests(t)

	var rowCount int
	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 0, rowCount)

	block := addTestBlock(t, ctx, bs)
	lpProto := getTestLiquidityProvision()

	data, err := entities.LiquidityProvisionFromProto(lpProto[0], generateTxHash(), block.ZetaTime)
	require.NoError(t, err)
	assert.NoError(t, lp.Upsert(ctx, data))
	err = lp.Flush(ctx)
	require.NoError(t, err)

	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 1, rowCount)
}

func testUpdateExistingInCurrentBlock(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, lp, conn := setupLPTests(t)

	var rowCount int
	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 0, rowCount)

	block := addTestBlock(t, ctx, bs)
	lpProto := getTestLiquidityProvision()

	data, err := entities.LiquidityProvisionFromProto(lpProto[0], generateTxHash(), block.ZetaTime)
	require.NoError(t, err)
	assert.NoError(t, lp.Upsert(ctx, data))

	data.Reference = "Updated"
	assert.NoError(t, lp.Upsert(ctx, data))
	err = lp.Flush(ctx)
	require.NoError(t, err)

	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 1, rowCount)
}

func testGetLPByReferenceAndParty(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, lp, conn := setupLPTests(t)

	var rowCount int
	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 0, rowCount)

	lpProto := getTestLiquidityProvision()

	for _, lpp := range lpProto {
		block := addTestBlock(t, ctx, bs)

		data, err := entities.LiquidityProvisionFromProto(lpp, generateTxHash(), block.ZetaTime)
		require.NoError(t, err)
		assert.NoError(t, lp.Upsert(ctx, data))
		err = lp.Flush(ctx)
		require.NoError(t, err)

		data.CreatedAt = data.CreatedAt.Truncate(time.Microsecond)
		data.UpdatedAt = data.UpdatedAt.Truncate(time.Microsecond)

		time.Sleep(100 * time.Millisecond)
	}

	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 3, rowCount)

	partyID := entities.PartyID("deadbaad")
	marketID := entities.MarketID("")
	got, _, err := lp.Get(ctx, partyID, marketID, "TEST1", entities.OffsetPagination{})
	require.NoError(t, err)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, got[0].Reference, "TEST1")
}

func testGetLPByPartyOnly(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, lp, conn := setupLPTests(t)

	var rowCount int
	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 0, rowCount)

	lpProto := getTestLiquidityProvision()

	want := make([]entities.LiquidityProvision, 0)

	for _, lpp := range lpProto {
		block := addTestBlock(t, ctx, bs)

		data, err := entities.LiquidityProvisionFromProto(lpp, generateTxHash(), block.ZetaTime)
		require.NoError(t, err)
		assert.NoError(t, lp.Upsert(ctx, data))
		err = lp.Flush(ctx)
		require.NoError(t, err)

		data.CreatedAt = data.CreatedAt.Truncate(time.Microsecond)
		data.UpdatedAt = data.UpdatedAt.Truncate(time.Microsecond)

		want = append(want, data)

		time.Sleep(100 * time.Millisecond)
	}

	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 3, rowCount)

	partyID := entities.PartyID("deadbaad")
	marketID := entities.MarketID("")
	got, _, err := lp.Get(ctx, partyID, marketID, "", entities.OffsetPagination{})
	require.NoError(t, err)
	assert.Equal(t, len(want), len(got))
	assert.ElementsMatch(t, want, got)
}

func testGetLPByPartyAndMarket(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, lp, conn := setupLPTests(t)

	var rowCount int
	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 0, rowCount)

	lpProto := getTestLiquidityProvision()

	wantMarketID := "dabbad00"

	want := make([]entities.LiquidityProvision, 0)

	for _, lpp := range lpProto {
		block := addTestBlock(t, ctx, bs)

		data, err := entities.LiquidityProvisionFromProto(lpp, generateTxHash(), block.ZetaTime)
		require.NoError(t, err)
		assert.NoError(t, lp.Upsert(ctx, data))
		err = lp.Flush(ctx)
		require.NoError(t, err)

		data.CreatedAt = data.CreatedAt.Truncate(time.Microsecond)
		data.UpdatedAt = data.UpdatedAt.Truncate(time.Microsecond)

		if data.MarketID.String() == wantMarketID {
			want = append(want, data)
		}

		time.Sleep(100 * time.Millisecond)
	}

	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 3, rowCount)

	partyID := entities.PartyID("DEADBAAD")
	marketID := entities.MarketID(wantMarketID)
	got, _, err := lp.Get(ctx, partyID, marketID, "", entities.OffsetPagination{})
	require.NoError(t, err)
	assert.Equal(t, len(want), len(got))
	assert.ElementsMatch(t, want, got)
}

func testGetLPNoPartyAndMarketErrors(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	_, lp, _ := setupLPTests(t)
	partyID := entities.PartyID("")
	marketID := entities.MarketID("")
	_, _, err := lp.Get(ctx, partyID, marketID, "", entities.OffsetPagination{})
	assert.Error(t, err)
}

func testGetLPNoPartyWithMarket(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, lp, conn := setupLPTests(t)

	var rowCount int
	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 0, rowCount)

	lpProto := getTestLiquidityProvision()
	wantMarketID := "dabbad00"
	want := make([]entities.LiquidityProvision, 0)

	for _, lpp := range lpProto {
		block := addTestBlock(t, ctx, bs)

		data, err := entities.LiquidityProvisionFromProto(lpp, generateTxHash(), block.ZetaTime)
		require.NoError(t, err)
		assert.NoError(t, lp.Upsert(ctx, data))
		err = lp.Flush(ctx)
		require.NoError(t, err)

		data.CreatedAt = data.CreatedAt.Truncate(time.Microsecond)
		data.UpdatedAt = data.UpdatedAt.Truncate(time.Microsecond)

		if data.MarketID.String() == wantMarketID {
			want = append(want, data)
		}

		time.Sleep(100 * time.Millisecond)
	}

	assert.NoError(t, conn.QueryRow(ctx, "select count(*) from liquidity_provisions").Scan(&rowCount))
	assert.Equal(t, 3, rowCount)
	partyID := entities.PartyID("")
	marketID := entities.MarketID(wantMarketID)
	got, _, err := lp.Get(ctx, partyID, marketID, "", entities.OffsetPagination{})
	require.NoError(t, err)
	assert.Equal(t, len(want), len(got))
	assert.ElementsMatch(t, want, got)
}

func getTestLiquidityProvision() []*zeta.LiquidityProvision {
	return []*zeta.LiquidityProvision{
		{
			Id:               "deadbeef",
			PartyId:          "deadbaad",
			CreatedAt:        time.Now().UnixNano(),
			UpdatedAt:        time.Now().UnixNano(),
			MarketId:         "cafed00d",
			CommitmentAmount: "100000",
			Fee:              "0.3",
			Sells:            nil,
			Buys:             nil,
			Version:          0,
			Status:           zeta.LiquidityProvision_STATUS_ACTIVE,
			Reference:        "TEST1",
		},
		{
			Id:               "0d15ea5e",
			PartyId:          "deadbaad",
			CreatedAt:        time.Now().UnixNano(),
			UpdatedAt:        time.Now().UnixNano(),
			MarketId:         "dabbad00",
			CommitmentAmount: "100000",
			Fee:              "0.3",
			Sells:            nil,
			Buys:             nil,
			Version:          0,
			Status:           zeta.LiquidityProvision_STATUS_ACTIVE,
			Reference:        "TEST",
		},
		{
			Id:               "deadc0de",
			PartyId:          "deadbaad",
			CreatedAt:        time.Now().UnixNano(),
			UpdatedAt:        time.Now().UnixNano(),
			MarketId:         "deadd00d",
			CommitmentAmount: "100000",
			Fee:              "0.3",
			Sells:            nil,
			Buys:             nil,
			Version:          0,
			Status:           zeta.LiquidityProvision_STATUS_ACTIVE,
			Reference:        "TEST",
		},
	}
}

func testLiquidityProvisionPaginationNoPagination(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, lpStore, _ := setupLPTests(t)
	testLps := addLiquidityProvisions(ctx, t, bs, lpStore)

	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := lpStore.Get(ctx, entities.PartyID("deadbaad"), entities.MarketID(""), "", pagination)

	require.NoError(t, err)
	assert.Equal(t, testLps, got)
	assert.False(t, pageInfo.HasPreviousPage)
	assert.False(t, pageInfo.HasNextPage)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[0].ZetaTime,
		ID:       testLps[0].ID,
	}.String()).Encode(), pageInfo.StartCursor)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[9].ZetaTime,
		ID:       testLps[9].ID,
	}.String()).Encode(), pageInfo.EndCursor)
}

func testLiquidityProvisionPaginationFirst(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()

	bs, lpStore, _ := setupLPTests(t)
	testLps := addLiquidityProvisions(ctx, t, bs, lpStore)

	first := int32(3)
	pagination, err := entities.NewCursorPagination(&first, nil, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := lpStore.Get(ctx, entities.PartyID("deadbaad"), entities.MarketID(""), "", pagination)

	require.NoError(t, err)
	want := testLps[:3]
	assert.Equal(t, want, got)
	assert.False(t, pageInfo.HasPreviousPage)
	assert.True(t, pageInfo.HasNextPage)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[0].ZetaTime,
		ID:       testLps[0].ID,
	}.String()).Encode(), pageInfo.StartCursor)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[2].ZetaTime,
		ID:       testLps[2].ID,
	}.String()).Encode(), pageInfo.EndCursor)
}

func testLiquidityProvisionPaginationLast(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, lpStore, _ := setupLPTests(t)
	testLps := addLiquidityProvisions(ctx, t, bs, lpStore)

	last := int32(3)
	pagination, err := entities.NewCursorPagination(nil, nil, &last, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := lpStore.Get(ctx, entities.PartyID("deadbaad"), entities.MarketID(""), "", pagination)

	require.NoError(t, err)
	want := testLps[7:]
	assert.Equal(t, want, got)
	assert.True(t, pageInfo.HasPreviousPage)
	assert.False(t, pageInfo.HasNextPage)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[7].ZetaTime,
		ID:       testLps[7].ID,
	}.String()).Encode(), pageInfo.StartCursor)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[9].ZetaTime,
		ID:       testLps[9].ID,
	}.String()).Encode(), pageInfo.EndCursor)
}

func testLiquidityProvisionPaginationFirstAfter(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, lpStore, _ := setupLPTests(t)
	testLps := addLiquidityProvisions(ctx, t, bs, lpStore)

	first := int32(3)
	after := testLps[2].Cursor().Encode()
	pagination, err := entities.NewCursorPagination(&first, &after, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := lpStore.Get(ctx, entities.PartyID("deadbaad"), entities.MarketID(""), "", pagination)

	require.NoError(t, err)
	want := testLps[3:6]
	assert.Equal(t, want, got)
	assert.True(t, pageInfo.HasPreviousPage)
	assert.True(t, pageInfo.HasNextPage)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[3].ZetaTime,
		ID:       testLps[3].ID,
	}.String()).Encode(), pageInfo.StartCursor)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[5].ZetaTime,
		ID:       testLps[5].ID,
	}.String()).Encode(), pageInfo.EndCursor)
}

func testLiquidityProvisionPaginationLastBefore(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	bs, lsStore, _ := setupLPTests(t)
	testLps := addLiquidityProvisions(ctx, t, bs, lsStore)

	last := int32(3)
	before := entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[7].ZetaTime,
		ID:       testLps[7].ID,
	}.String()).Encode()
	pagination, err := entities.NewCursorPagination(nil, nil, &last, &before, false)
	require.NoError(t, err)
	got, pageInfo, err := lsStore.Get(ctx, entities.PartyID("deadbaad"), entities.MarketID(""), "", pagination)

	require.NoError(t, err)
	want := testLps[4:7]
	assert.Equal(t, want, got)
	assert.True(t, pageInfo.HasPreviousPage)
	assert.True(t, pageInfo.HasNextPage)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[4].ZetaTime,
		ID:       testLps[4].ID,
	}.String()).Encode(), pageInfo.StartCursor)
	assert.Equal(t, entities.NewCursor(entities.LiquidityProvisionCursor{
		ZetaTime: testLps[6].ZetaTime,
		ID:       testLps[6].ID,
	}.String()).Encode(), pageInfo.EndCursor)
}

func addLiquidityProvisions(ctx context.Context, t *testing.T, bs *sqlstore.Blocks, lpstore *sqlstore.LiquidityProvision) []entities.LiquidityProvision {
	t.Helper()
	zetaTime := time.Now().Truncate(time.Microsecond)
	amount := int64(1000)
	lps := make([]entities.LiquidityProvision, 0, 10)
	for i := 0; i < 10; i++ {
		addTestBlockForTime(t, ctx, bs, zetaTime)

		lp := &zeta.LiquidityProvision{
			Id:               fmt.Sprintf("deadbeef%02d", i+1),
			PartyId:          "deadbaad",
			CreatedAt:        zetaTime.UnixNano(),
			UpdatedAt:        zetaTime.UnixNano(),
			MarketId:         "cafed00d",
			CommitmentAmount: "100000",
			Fee:              "0.3",
			Sells:            nil,
			Buys:             nil,
			Version:          0,
			Status:           zeta.LiquidityProvision_STATUS_ACTIVE,
			Reference:        "TEST1",
		}

		withdrawal, err := entities.LiquidityProvisionFromProto(lp, generateTxHash(), zetaTime)
		require.NoError(t, err, "Converting withdrawal proto to database entity")
		err = lpstore.Upsert(ctx, withdrawal)
		require.NoError(t, err)
		require.NoError(t, err)
		lpstore.Flush(ctx)
		lps = append(lps, withdrawal)
		require.NoError(t, err)

		zetaTime = zetaTime.Add(time.Second)
		amount += 100
	}

	return lps
}
