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
	"testing"
	"time"

	"zuluprotocol/zeta/zeta/datanode/entities"
	"zuluprotocol/zeta/zeta/datanode/sqlstore"
	"zuluprotocol/zeta/zeta/datanode/sqlstore/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addTestParty(t *testing.T, ctx context.Context, ps *sqlstore.Parties, block entities.Block) entities.Party {
	t.Helper()
	party := entities.Party{
		ID:       entities.PartyID(helpers.GenerateID()),
		ZetaTime: &block.ZetaTime,
	}

	err := ps.Add(ctx, party)
	require.NoError(t, err)
	return party
}

func TestParty(t *testing.T) {
	ctx, rollback := tempTransaction(t)
	defer rollback()
	ps := sqlstore.NewParties(connectionSource)
	ps.Initialise(ctx)
	bs := sqlstore.NewBlocks(connectionSource)
	block := addTestBlock(t, ctx, bs)

	// Make sure we're starting with an empty set of parties (except network party)
	parties, err := ps.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, parties, 1)
	assert.Equal(t, "network", parties[0].ID.String())

	// Make a new party
	party := addTestParty(t, ctx, ps, block)

	// Add it again, we shouldn't get a primary key violation (we just ignore)
	err = ps.Add(ctx, party)
	require.NoError(t, err)

	// Query and check we've got back a party the same as the one we put in
	fetchedParty, err := ps.GetByID(ctx, party.ID.String())
	require.NoError(t, err)
	assert.Equal(t, party, fetchedParty)

	// Get all assets and make sure ours is in there (along with built in network party)
	parties, err = ps.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, parties, 2)

	// Check we get the right error if we ask for a non-existent party
	_, err = ps.GetByID(ctx, "beef")
	assert.ErrorIs(t, err, entities.ErrNotFound)
}

func setupPartyTest(t *testing.T) (*sqlstore.Blocks, *sqlstore.Parties) {
	t.Helper()
	bs := sqlstore.NewBlocks(connectionSource)
	pt := sqlstore.NewParties(connectionSource)

	return bs, pt
}

func populateTestParties(ctx context.Context, t *testing.T, bs *sqlstore.Blocks, ps *sqlstore.Parties, blockTimes map[string]time.Time) {
	t.Helper()
	parties := []entities.Party{
		{
			ID: entities.PartyID("02a16077"),
		},
		{
			ID: entities.PartyID("44eea1bc"),
		},
		{
			ID: entities.PartyID("65be62cd"),
		},
		{
			ID: entities.PartyID("7a797e0e"),
		},
		{
			ID: entities.PartyID("7bb2356e"),
		},
		{
			ID: entities.PartyID("b7c84b8e"),
		},
		{
			ID: entities.PartyID("c612300d"),
		},
		{
			ID: entities.PartyID("c8744329"),
		},
		{
			ID: entities.PartyID("da8d1803"),
		},
		{
			ID: entities.PartyID("fb1528a5"),
		},
	}

	for _, party := range parties {
		block := addTestBlock(t, ctx, bs)
		party.ZetaTime = &block.ZetaTime
		blockTimes[party.ID.String()] = block.ZetaTime
		err := ps.Add(ctx, party)
		require.NoError(t, err)
		time.Sleep(time.Microsecond * 100)
	}
}

func TestPartyPagination(t *testing.T) {
	t.Run("CursorPagination should return the party if Party ID is provided", testPartyPaginationReturnsTheSpecifiedParty)
	t.Run("CursorPagination should return all parties if no party ID and no cursor is provided", testPartyPaginationReturnAllParties)
	t.Run("CursorPagination should return the first page when first limit is provided with no after cursor", testPartyPaginationReturnsFirstPage)
	t.Run("CursorPagination should return the last page when last limit is provided with no before cursor", testPartyPaginationReturnsLastPage)
	t.Run("CursorPagination should return the page specified by the first limit and after cursor", testPartyPaginationReturnsPageTraversingForward)
	t.Run("CursorPagination should return the page specified by the last limit and before cursor", testPartyPaginationReturnsPageTraversingBackward)

	t.Run("CursorPagination should return the party if Party ID is provided - newest first", testPartyPaginationReturnsTheSpecifiedPartyNewestFirst)
	t.Run("CursorPagination should return all parties if no party ID and no cursor is provided - newest first", testPartyPaginationReturnAllPartiesNewestFirst)
	t.Run("CursorPagination should return the first page when first limit is provided with no after cursor - newest first", testPartyPaginationReturnsFirstPageNewestFirst)
	t.Run("CursorPagination should return the last page when last limit is provided with no before cursor - newest first", testPartyPaginationReturnsLastPageNewestFirst)
	t.Run("CursorPagination should return the page specified by the first limit and after cursor - newest first", testPartyPaginationReturnsPageTraversingForwardNewestFirst)
	t.Run("CursorPagination should return the page specified by the last limit and before cursor - newest first", testPartyPaginationReturnsPageTraversingBackwardNewestFirst)
}

func testPartyPaginationReturnsTheSpecifiedParty(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := pt.GetAllPaged(ctx, "c612300d", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "c612300d", got[0].ID.String())

	zetaTime := blockTimes["c612300d"]
	party := entities.Party{
		ID:       "c612300d",
		ZetaTime: &zetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(party).Encode()
	wantEndCursor := entities.NewCursor(party).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnAllParties(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, false)
	require.NoError(t, err)
	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 10)
	assert.Equal(t, "02a16077", got[0].ID.String())
	assert.Equal(t, "fb1528a5", got[9].ID.String())

	startZetaTime := blockTimes["02a16077"]
	startParty := entities.Party{
		ID:       "02a16077",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["fb1528a5"]
	endParty := entities.Party{
		ID:       "fb1528a5",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsFirstPage(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	first := int32(3)
	pagination, err := entities.NewCursorPagination(&first, nil, nil, nil, false)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "02a16077", got[0].ID.String())
	assert.Equal(t, "65be62cd", got[2].ID.String())

	startZetaTime := blockTimes["02a16077"]
	startParty := entities.Party{
		ID:       "02a16077",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["65be62cd"]
	endParty := entities.Party{
		ID:       "65be62cd",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: false,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsLastPage(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	last := int32(3)
	pagination, err := entities.NewCursorPagination(nil, nil, &last, nil, false)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "c8744329", got[0].ID.String())
	assert.Equal(t, "fb1528a5", got[2].ID.String())

	startZetaTime := blockTimes["c8744329"]
	startParty := entities.Party{
		ID:       "c8744329",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["fb1528a5"]
	endParty := entities.Party{
		ID:       "fb1528a5",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: true,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsPageTraversingForward(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	first := int32(3)
	afterZetaTime := blockTimes["65be62cd"]
	afterParty := entities.Party{
		ID:       "65be62cd",
		ZetaTime: &afterZetaTime,
	}.String()
	after := entities.NewCursor(afterParty).Encode()
	pagination, err := entities.NewCursorPagination(&first, &after, nil, nil, false)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "7a797e0e", got[0].ID.String())
	assert.Equal(t, "b7c84b8e", got[2].ID.String())

	startZetaTime := blockTimes["7a797e0e"]
	startParty := entities.Party{
		ID:       "7a797e0e",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["b7c84b8e"]
	endParty := entities.Party{
		ID:       "b7c84b8e",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsPageTraversingBackward(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	last := int32(3)
	beforeZetaTime := blockTimes["c8744329"]
	beforeParty := entities.Party{
		ID:       "c8744329",
		ZetaTime: &beforeZetaTime,
	}.String()
	before := entities.NewCursor(beforeParty).Encode()
	pagination, err := entities.NewCursorPagination(nil, nil, &last, &before, false)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "7bb2356e", got[0].ID.String())
	assert.Equal(t, "c612300d", got[2].ID.String())

	startZetaTime := blockTimes["7bb2356e"]
	startParty := entities.Party{
		ID:       "7bb2356e",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["c612300d"]
	endParty := entities.Party{
		ID:       "c612300d",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsTheSpecifiedPartyNewestFirst(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, true)
	require.NoError(t, err)
	got, pageInfo, err := pt.GetAllPaged(ctx, "c612300d", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "c612300d", got[0].ID.String())

	startZetaTime := blockTimes["c612300d"]
	startParty := entities.Party{
		ID:       "c612300d",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["c612300d"]
	endParty := entities.Party{
		ID:       "c612300d",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnAllPartiesNewestFirst(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	pagination, err := entities.NewCursorPagination(nil, nil, nil, nil, true)
	require.NoError(t, err)
	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 10)
	assert.Equal(t, "fb1528a5", got[0].ID.String())
	assert.Equal(t, "02a16077", got[9].ID.String())

	startZetaTime := blockTimes["fb1528a5"]
	startParty := entities.Party{
		ID:       "fb1528a5",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["02a16077"]
	endParty := entities.Party{
		ID:       "02a16077",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsFirstPageNewestFirst(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	first := int32(3)
	pagination, err := entities.NewCursorPagination(&first, nil, nil, nil, true)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "fb1528a5", got[0].ID.String())
	assert.Equal(t, "c8744329", got[2].ID.String())

	startZetaTime := blockTimes["fb1528a5"]
	startParty := entities.Party{
		ID:       "fb1528a5",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["c8744329"]
	endParty := entities.Party{
		ID:       "c8744329",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: false,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsLastPageNewestFirst(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	last := int32(3)
	pagination, err := entities.NewCursorPagination(nil, nil, &last, nil, true)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "65be62cd", got[0].ID.String())
	assert.Equal(t, "02a16077", got[2].ID.String())

	startZetaTime := blockTimes["65be62cd"]
	startParty := entities.Party{
		ID:       "65be62cd",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["02a16077"]
	endParty := entities.Party{
		ID:       "02a16077",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: true,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsPageTraversingForwardNewestFirst(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	first := int32(3)
	afterZetaTime := blockTimes["c8744329"]
	afterParty := entities.Party{
		ID:       "c8744329",
		ZetaTime: &afterZetaTime,
	}.String()
	after := entities.NewCursor(afterParty).Encode()
	pagination, err := entities.NewCursorPagination(&first, &after, nil, nil, true)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "c612300d", got[0].ID.String())
	assert.Equal(t, "7bb2356e", got[2].ID.String())

	startZetaTime := blockTimes["c612300d"]
	startParty := entities.Party{
		ID:       "c612300d",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["7bb2356e"]
	endParty := entities.Party{
		ID:       "7bb2356e",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}

func testPartyPaginationReturnsPageTraversingBackwardNewestFirst(t *testing.T) {
	bs, pt := setupPartyTest(t)
	ctx, rollback := tempTransaction(t)
	defer rollback()

	blockTimes := make(map[string]time.Time)
	populateTestParties(ctx, t, bs, pt, blockTimes)
	last := int32(3)
	beforeZetaTime := blockTimes["65be62cd"]
	beforeParty := entities.Party{
		ID:       "65be62cd",
		ZetaTime: &beforeZetaTime,
	}.String()
	before := entities.NewCursor(beforeParty).Encode()
	pagination, err := entities.NewCursorPagination(nil, nil, &last, &before, true)
	require.NoError(t, err)

	got, pageInfo, err := pt.GetAllPaged(ctx, "", pagination)
	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "b7c84b8e", got[0].ID.String())
	assert.Equal(t, "7a797e0e", got[2].ID.String())

	startZetaTime := blockTimes["b7c84b8e"]
	startParty := entities.Party{
		ID:       "b7c84b8e",
		ZetaTime: &startZetaTime,
	}.String()
	endZetaTime := blockTimes["7a797e0e"]
	endParty := entities.Party{
		ID:       "7a797e0e",
		ZetaTime: &endZetaTime,
	}.String()
	wantStartCursor := entities.NewCursor(startParty).Encode()
	wantEndCursor := entities.NewCursor(endParty).Encode()
	assert.Equal(t, entities.PageInfo{
		HasNextPage:     true,
		HasPreviousPage: true,
		StartCursor:     wantStartCursor,
		EndCursor:       wantEndCursor,
	}, pageInfo)
}
