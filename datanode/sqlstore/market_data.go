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
	"errors"
	"fmt"
	"time"

	"zuluprotocol/zeta/zeta/datanode/entities"
	"zuluprotocol/zeta/zeta/datanode/metrics"
	"zuluprotocol/zeta/zeta/logging"
	v2 "zuluprotocol/zeta/zeta/protos/data-node/api/v2"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

var marketdataOrdering = TableOrdering{
	ColumnOrdering{Name: "synthetic_time", Sorting: ASC},
}

type MarketData struct {
	*ConnectionSource
	columns    []string
	marketData []*entities.MarketData
}

var ErrInvalidDateRange = errors.New("invalid date range, end date must be after start date")

func NewMarketData(connectionSource *ConnectionSource) *MarketData {
	return &MarketData{
		ConnectionSource: connectionSource,
		columns: []string{
			"synthetic_time", "tx_hash", "zeta_time", "seq_num",
			"market", "mark_price", "best_bid_price", "best_bid_volume",
			"best_offer_price", "best_offer_volume", "best_static_bid_price", "best_static_bid_volume",
			"best_static_offer_price", "best_static_offer_volume", "mid_price", "static_mid_price",
			"open_interest", "auction_end", "auction_start", "indicative_price", "indicative_volume",
			"market_trading_mode", "auction_trigger", "extension_trigger", "target_stake",
			"supplied_stake", "price_monitoring_bounds", "market_value_proxy", "liquidity_provider_fee_shares", "market_state", "next_mark_to_market",
		},
	}
}

func (md *MarketData) Add(data *entities.MarketData) error {
	md.marketData = append(md.marketData, data)
	return nil
}

func (md *MarketData) Flush(ctx context.Context) ([]*entities.MarketData, error) {
	rows := make([][]interface{}, 0, len(md.marketData))
	for _, data := range md.marketData {
		rows = append(rows, []interface{}{
			data.SyntheticTime, data.TxHash, data.ZetaTime, data.SeqNum,
			data.Market, data.MarkPrice,
			data.BestBidPrice, data.BestBidVolume, data.BestOfferPrice, data.BestOfferVolume,
			data.BestStaticBidPrice, data.BestStaticBidVolume, data.BestStaticOfferPrice, data.BestStaticOfferVolume,
			data.MidPrice, data.StaticMidPrice, data.OpenInterest, data.AuctionEnd,
			data.AuctionStart, data.IndicativePrice, data.IndicativeVolume, data.MarketTradingMode,
			data.AuctionTrigger, data.ExtensionTrigger, data.TargetStake, data.SuppliedStake,
			data.PriceMonitoringBounds, data.MarketValueProxy, data.LiquidityProviderFeeShares, data.MarketState, data.NextMarkToMarket,
		})
	}
	defer metrics.StartSQLQuery("MarketData", "Flush")()
	if rows != nil {
		copyCount, err := md.Connection.CopyFrom(
			ctx,
			pgx.Identifier{"market_data"}, md.columns, pgx.CopyFromRows(rows),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to copy market data into database:%w", err)
		}

		if copyCount != int64(len(rows)) {
			return nil, fmt.Errorf("copied %d market data rows into the database, expected to copy %d", copyCount, len(rows))
		}
	}

	flushed := md.marketData
	md.marketData = nil

	return flushed, nil
}

func (md *MarketData) GetMarketDataByID(ctx context.Context, marketID string) (entities.MarketData, error) {
	defer metrics.StartSQLQuery("MarketData", "GetMarketDataByID")()
	md.log.Debug("Retrieving market data from Postgres", logging.String("market-id", marketID))

	var marketData entities.MarketData
	query := "select * from market_data_snapshot where market = $1"

	return marketData, md.wrapE(pgxscan.Get(ctx, md.Connection, &marketData, query, entities.MarketID(marketID)))
}

func (md *MarketData) GetMarketsData(ctx context.Context) ([]entities.MarketData, error) {
	md.log.Debug("Retrieving markets data from Postgres")

	var marketData []entities.MarketData
	query := "select * from market_data_snapshot"

	defer metrics.StartSQLQuery("MarketData", "GetMarketsData")()
	err := pgxscan.Select(ctx, md.Connection, &marketData, query)

	return marketData, err
}

func (md *MarketData) GetBetweenDatesByID(ctx context.Context, marketID string, start, end time.Time, pagination entities.Pagination) ([]entities.MarketData, entities.PageInfo, error) {
	if end.Before(start) {
		return nil, entities.PageInfo{}, ErrInvalidDateRange
	}

	switch p := pagination.(type) {
	case entities.OffsetPagination:
		return md.getBetweenDatesByIDOffset(ctx, marketID, &start, &end, p)
	case entities.CursorPagination:
		return md.getBetweenDatesByID(ctx, marketID, &start, &end, p)
	}

	return nil, entities.PageInfo{}, errors.New("invalid pagination")
}

func (md *MarketData) GetFromDateByID(ctx context.Context, marketID string, start time.Time, pagination entities.Pagination) ([]entities.MarketData, entities.PageInfo, error) {
	switch p := pagination.(type) {
	case entities.OffsetPagination:
		return md.getBetweenDatesByIDOffset(ctx, marketID, &start, nil, p)
	case entities.CursorPagination:
		return md.getBetweenDatesByID(ctx, marketID, &start, nil, p)
	}
	return nil, entities.PageInfo{}, errors.New("invalid pagination")
}

func (md *MarketData) GetToDateByID(ctx context.Context, marketID string, end time.Time, pagination entities.Pagination) ([]entities.MarketData, entities.PageInfo, error) {
	switch p := pagination.(type) {
	case entities.OffsetPagination:
		return md.getBetweenDatesByIDOffset(ctx, marketID, nil, &end, p)
	case entities.CursorPagination:
		return md.getBetweenDatesByID(ctx, marketID, nil, &end, p)
	}
	return nil, entities.PageInfo{}, errors.New("invalid pagination")
}

func (md *MarketData) getBetweenDatesByIDOffset(ctx context.Context, marketID string, start, end *time.Time, pagination entities.OffsetPagination) ([]entities.MarketData, entities.PageInfo, error) {
	defer metrics.StartSQLQuery("MarketData", "getBetweenDatesByID")()
	market := entities.MarketID(marketID)

	selectStatement := `select * from market_data`

	var results []entities.MarketData
	var err error

	if start != nil && end != nil {
		query, args := orderAndPaginateQuery(
			fmt.Sprintf(`%s where market = $1 and synthetic_time between $2 and $3`, selectStatement),
			[]string{"synthetic_time"}, pagination,
			market, *start, *end)

		err = pgxscan.Select(ctx, md.Connection, &results, query, args...)
	} else if start != nil && end == nil {
		query, args := orderAndPaginateQuery(fmt.Sprintf(`%s where market = $1 and zeta_time >= $2`, selectStatement),
			[]string{"synthetic_time"}, pagination,
			market, *start)

		err = pgxscan.Select(ctx, md.Connection, &results, query, args...)
	} else if start == nil && end != nil {
		query, args := orderAndPaginateQuery(fmt.Sprintf(`%s where market = $1 and zeta_time <= $2`, selectStatement),
			[]string{"synthetic_time"}, pagination,
			market, *end)

		err = pgxscan.Select(ctx, md.Connection, &results, query, args...)
	}

	return results, entities.PageInfo{}, err
}

func (md *MarketData) getBetweenDatesByID(ctx context.Context, marketID string, start, end *time.Time, pagination entities.CursorPagination) ([]entities.MarketData, entities.PageInfo, error) {
	defer metrics.StartSQLQuery("MarketData", "getBetweenDatesByID")()
	market := entities.MarketID(marketID)

	selectStatement := `select * from market_data`
	args := make([]interface{}, 0)

	var (
		query    string
		err      error
		pageInfo entities.PageInfo
	)

	if start != nil && end != nil {
		query = fmt.Sprintf(`%s where market = %s and zeta_time between %s and %s`, selectStatement,
			nextBindVar(&args, market),
			nextBindVar(&args, *start),
			nextBindVar(&args, *end),
		)
	} else if start != nil && end == nil {
		query = fmt.Sprintf(`%s where market = %s and zeta_time >= %s`, selectStatement,
			nextBindVar(&args, market),
			nextBindVar(&args, *start))
	} else if start == nil && end != nil {
		query = fmt.Sprintf(`%s where market = %s and zeta_time <= %s`, selectStatement,
			nextBindVar(&args, market),
			nextBindVar(&args, *end))
	}

	query, args, err = PaginateQuery[entities.MarketDataCursor](query, args, marketdataOrdering, pagination)
	if err != nil {
		return nil, pageInfo, err
	}
	var pagedData []entities.MarketData

	if err = pgxscan.Select(ctx, md.Connection, &pagedData, query, args...); err != nil {
		return pagedData, pageInfo, err
	}

	pagedData, pageInfo = entities.PageEntities[*v2.MarketDataEdge](pagedData, pagination)

	return pagedData, pageInfo, nil
}
