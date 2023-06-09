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

package gql

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc"

	"zuluprotocol/zeta/datanode/gateway"
	"zuluprotocol/zeta/datanode/zetatime"
	"zuluprotocol/zeta/libs/num"
	"zuluprotocol/zeta/libs/ptr"
	"zuluprotocol/zeta/logging"
	v2 "zuluprotocol/zeta/protos/data-node/api/v2"
	"zuluprotocol/zeta/protos/zeta"
	types "zuluprotocol/zeta/protos/zeta"
	zetaprotoapi "code.zetaprotocol.io/zeta/protos/zeta/api/v1"
	commandspb "zuluprotocol/zeta/protos/zeta/commands/v1"
	data "zuluprotocol/zeta/protos/zeta/data/v1"
	eventspb "zuluprotocol/zeta/protos/zeta/events/v1"
)

var (
	// ErrMissingIDOrReference is returned when neither id nor reference has been supplied in the query.
	ErrMissingIDOrReference = errors.New("missing id or reference")
	// ErrMissingNodeID is returned when no node id has been supplied in the query.
	ErrMissingNodeID = errors.New("missing node id")
	// ErrInvalidVotesSubscription is returned if neither proposal ID nor party ID is specified.
	ErrInvalidVotesSubscription = errors.New("invalid subscription, either proposal or party ID required")
	// ErrInvalidProposal is returned when invalid governance data is received by proposal resolver.
	ErrInvalidProposal = errors.New("invalid proposal")
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/mocks.go -package mocks zuluprotocol/zeta/datanode/gateway/graphql CoreProxyServiceClient,TradingDataServiceClientV2

// CoreProxyServiceClient ...
type CoreProxyServiceClient interface {
	zetaprotoapi.CoreServiceClient
}

type TradingDataServiceClientV2 interface {
	v2.TradingDataServiceClient
}

// ZetaResolverRoot is the root resolver for all graphql types.
type ZetaResolverRoot struct {
	gateway.Config

	log                 *logging.Logger
	tradingProxyClient  CoreProxyServiceClient
	tradingDataClientV2 TradingDataServiceClientV2
	r                   allResolver
}

// NewResolverRoot instantiate a graphql root resolver.
func NewResolverRoot(
	log *logging.Logger,
	config gateway.Config,
	tradingClient CoreProxyServiceClient,
	tradingDataClientV2 TradingDataServiceClientV2,
) *ZetaResolverRoot {
	return &ZetaResolverRoot{
		log:                 log,
		Config:              config,
		tradingProxyClient:  tradingClient,
		tradingDataClientV2: tradingDataClientV2,
		r:                   allResolver{log, tradingDataClientV2},
	}
}

// Query returns the query resolver.
func (r *ZetaResolverRoot) Query() QueryResolver {
	return (*myQueryResolver)(r)
}

// Candle returns the candles resolver.
func (r *ZetaResolverRoot) Candle() CandleResolver {
	return (*myCandleResolver)(r)
}

func (r *ZetaResolverRoot) DataSourceSpecConfiguration() DataSourceSpecConfigurationResolver {
	return (*myDataSourceSpecConfigurationResolver)(r)
}

// MarginLevels returns the market levels resolver.
func (r *ZetaResolverRoot) MarginLevels() MarginLevelsResolver {
	return (*myMarginLevelsResolver)(r)
}

// MarginLevelsUpdate returns the market levels resolver.
func (r *ZetaResolverRoot) MarginLevelsUpdate() MarginLevelsUpdateResolver {
	return (*myMarginLevelsUpdateResolver)(r)
}

// PriceLevel returns the price levels resolver.
func (r *ZetaResolverRoot) PriceLevel() PriceLevelResolver {
	return (*myPriceLevelResolver)(r)
}

// Market returns the markets resolver.
func (r *ZetaResolverRoot) Market() MarketResolver {
	return (*myMarketResolver)(r)
}

// Order returns the order resolver.
func (r *ZetaResolverRoot) Order() OrderResolver {
	return (*myOrderResolver)(r)
}

// OrderUpdate returns the order resolver.
func (r *ZetaResolverRoot) OrderUpdate() OrderUpdateResolver {
	return (*myOrderUpdateResolver)(r)
}

// Trade returns the trades resolver.
func (r *ZetaResolverRoot) Trade() TradeResolver {
	return (*myTradeResolver)(r)
}

// Position returns the positions resolver.
func (r *ZetaResolverRoot) Position() PositionResolver {
	return (*myPositionResolver)(r)
}

// PositionUpdate returns the positionUpdate resolver.
func (r *ZetaResolverRoot) PositionUpdate() PositionUpdateResolver {
	return (*positionUpdateResolver)(r)
}

// Party returns the parties resolver.
func (r *ZetaResolverRoot) Party() PartyResolver {
	return (*myPartyResolver)(r)
}

// Subscription returns the subscriptions resolver.
func (r *ZetaResolverRoot) Subscription() SubscriptionResolver {
	return (*mySubscriptionResolver)(r)
}

// Account returns the accounts resolver.
func (r *ZetaResolverRoot) AccountEvent() AccountEventResolver {
	return (*myAccountEventResolver)(r)
}

// Account returns the accounts resolver.
func (r *ZetaResolverRoot) AccountBalance() AccountBalanceResolver {
	return (*myAccountResolver)(r)
}

// Account returns the accounts resolver.
func (r *ZetaResolverRoot) AccountDetails() AccountDetailsResolver {
	return (*myAccountDetailsResolver)(r)
}

// Proposal returns the proposal resolver.
func (r *ZetaResolverRoot) Proposal() ProposalResolver {
	return (*proposalResolver)(r)
}

// NodeSignature ...
func (r *ZetaResolverRoot) NodeSignature() NodeSignatureResolver {
	return (*myNodeSignatureResolver)(r)
}

// Asset ...
func (r *ZetaResolverRoot) Asset() AssetResolver {
	return (*myAssetResolver)(r)
}

// Deposit ...
func (r *ZetaResolverRoot) Deposit() DepositResolver {
	return (*myDepositResolver)(r)
}

// Withdrawal ...
func (r *ZetaResolverRoot) Withdrawal() WithdrawalResolver {
	return (*myWithdrawalResolver)(r)
}

func (r *ZetaResolverRoot) PropertyKey() PropertyKeyResolver {
	return (*myPropertyKeyResolver)(r)
}

func (r *ZetaResolverRoot) LiquidityOrderReference() LiquidityOrderReferenceResolver {
	return (*myLiquidityOrderReferenceResolver)(r)
}

func (r *ZetaResolverRoot) LiquidityProvision() LiquidityProvisionResolver {
	return (*myLiquidityProvisionResolver)(r)
}

func (r *ZetaResolverRoot) Future() FutureResolver {
	return (*myFutureResolver)(r)
}

func (r *ZetaResolverRoot) FutureProduct() FutureProductResolver {
	return (*myFutureProductResolver)(r)
}

func (r *ZetaResolverRoot) Instrument() InstrumentResolver {
	return (*myInstrumentResolver)(r)
}

func (r *ZetaResolverRoot) InstrumentConfiguration() InstrumentConfigurationResolver {
	return (*myInstrumentConfigurationResolver)(r)
}

func (r *ZetaResolverRoot) TradableInstrument() TradableInstrumentResolver {
	return (*myTradableInstrumentResolver)(r)
}

func (r *ZetaResolverRoot) NewAsset() NewAssetResolver {
	return (*newAssetResolver)(r)
}

func (r *ZetaResolverRoot) UpdateAsset() UpdateAssetResolver {
	return (*updateAssetResolver)(r)
}

func (r *ZetaResolverRoot) UpdateFutureProduct() UpdateFutureProductResolver {
	return (*updateFutureProductResolver)(r)
}

func (r *ZetaResolverRoot) NewMarket() NewMarketResolver {
	return (*newMarketResolver)(r)
}

func (r *ZetaResolverRoot) ProposalTerms() ProposalTermsResolver {
	return (*proposalTermsResolver)(r)
}

func (r *ZetaResolverRoot) UpdateMarket() UpdateMarketResolver {
	return (*updateMarketResolver)(r)
}

func (r *ZetaResolverRoot) UpdateNetworkParameter() UpdateNetworkParameterResolver {
	return (*updateNetworkParameterResolver)(r)
}

func (r *ZetaResolverRoot) NewFreeform() NewFreeformResolver {
	return (*newFreeformResolver)(r)
}

func (r *ZetaResolverRoot) OracleSpec() OracleSpecResolver {
	return (*oracleSpecResolver)(r)
}

func (r *ZetaResolverRoot) OracleData() OracleDataResolver {
	return (*oracleDataResolver)(r)
}

func (r *ZetaResolverRoot) AuctionEvent() AuctionEventResolver {
	return (*auctionEventResolver)(r)
}

func (r *ZetaResolverRoot) Vote() VoteResolver {
	return (*voteResolver)(r)
}

func (r *ZetaResolverRoot) NodeData() NodeDataResolver {
	return (*nodeDataResolver)(r)
}

func (r *ZetaResolverRoot) Node() NodeResolver {
	return (*nodeResolver)(r)
}

func (r *ZetaResolverRoot) RankingScore() RankingScoreResolver {
	return (*rankingScoreResolver)(r)
}

func (r *ZetaResolverRoot) KeyRotation() KeyRotationResolver {
	return (*keyRotationResolver)(r)
}

func (r *ZetaResolverRoot) EthereumKeyRotation() EthereumKeyRotationResolver {
	return (*ethereumKeyRotationResolver)(r)
}

func (r *ZetaResolverRoot) Delegation() DelegationResolver {
	return (*delegationResolver)(r)
}

func (r *ZetaResolverRoot) DateRange() DateRangeResolver {
	return (*dateRangeResolver)(r)
}

func (r *ZetaResolverRoot) Epoch() EpochResolver {
	return (*epochResolver)(r)
}

func (r *ZetaResolverRoot) EpochTimestamps() EpochTimestampsResolver {
	return (*epochTimestampsResolver)(r)
}

func (r *ZetaResolverRoot) Reward() RewardResolver {
	return (*rewardResolver)(r)
}

func (r *ZetaResolverRoot) RewardSummary() RewardSummaryResolver {
	return (*rewardSummaryResolver)(r)
}

func (r *ZetaResolverRoot) StakeLinking() StakeLinkingResolver {
	return (*stakeLinkingResolver)(r)
}

func (r *ZetaResolverRoot) PartyStake() PartyStakeResolver {
	return (*partyStakeResolver)(r)
}

func (r *ZetaResolverRoot) Statistics() StatisticsResolver {
	return (*statisticsResolver)(r)
}

func (r *ZetaResolverRoot) Transfer() TransferResolver {
	return (*transferResolver)(r)
}

func (r *ZetaResolverRoot) RecurringTransfer() RecurringTransferResolver {
	return (*recurringTransferResolver)(r)
}

func (r *ZetaResolverRoot) UpdateMarketConfiguration() UpdateMarketConfigurationResolver {
	return (*updateMarketConfigurationResolver)(r)
}

func (r *ZetaResolverRoot) AccountUpdate() AccountUpdateResolver {
	return (*accountUpdateResolver)(r)
}

func (r *ZetaResolverRoot) TradeUpdate() TradeUpdateResolver {
	return (*tradeUpdateResolver)(r)
}

func (r *ZetaResolverRoot) LiquidityProvisionUpdate() LiquidityProvisionUpdateResolver {
	return (*liquidityProvisionUpdateResolver)(r)
}

func (r *ZetaResolverRoot) TransactionResult() TransactionResultResolver {
	return (*transactionResultResolver)(r)
}

func (r *ZetaResolverRoot) ProtocolUpgradeProposal() ProtocolUpgradeProposalResolver {
	return (*protocolUpgradeProposalResolver)(r)
}

func (r *ZetaResolverRoot) CoreSnapshotData() CoreSnapshotDataResolver {
	return (*coreDataSnapshotResolver)(r)
}

func (r *ZetaResolverRoot) EpochRewardSummary() EpochRewardSummaryResolver {
	return (*epochRewardSummaryResolver)(r)
}

func (r *ZetaResolverRoot) OrderFilter() OrderFilterResolver {
	return (*orderFilterResolver)(r)
}

// RewardSummaryFilter returns RewardSummaryFilterResolver implementation.
func (r *ZetaResolverRoot) RewardSummaryFilter() RewardSummaryFilterResolver {
	return (*rewardSummaryFilterResolver)(r)
}

type protocolUpgradeProposalResolver ZetaResolverRoot

func (r *protocolUpgradeProposalResolver) UpgradeBlockHeight(ctx context.Context, obj *eventspb.ProtocolUpgradeEvent) (string, error) {
	return fmt.Sprintf("%d", obj.UpgradeBlockHeight), nil
}

type coreDataSnapshotResolver ZetaResolverRoot

func (r *coreDataSnapshotResolver) BlockHeight(ctx context.Context, obj *eventspb.CoreSnapshotData) (string, error) {
	return fmt.Sprintf("%d", obj.BlockHeight), nil
}

func (r *coreDataSnapshotResolver) ZetaCoreVersion(ctx context.Context, obj *eventspb.CoreSnapshotData) (string, error) {
	return obj.CoreVersion, nil
}

type epochRewardSummaryResolver ZetaResolverRoot

func (r *epochRewardSummaryResolver) RewardType(ctx context.Context, obj *zeta.EpochRewardSummary) (zeta.AccountType, error) {
	accountType, ok := zeta.AccountType_value[obj.RewardType]
	if !ok {
		return zeta.AccountType_ACCOUNT_TYPE_UNSPECIFIED, fmt.Errorf("Unknown account type %v", obj.RewardType)
	}

	return zeta.AccountType(accountType), nil
}

func (r *epochRewardSummaryResolver) Epoch(ctx context.Context, obj *zeta.EpochRewardSummary) (int, error) {
	return int(obj.Epoch), nil
}

type transactionResultResolver ZetaResolverRoot

func (r *transactionResultResolver) Error(ctx context.Context, tr *eventspb.TransactionResult) (*string, error) {
	if tr == nil || tr.Status {
		return nil, nil
	}

	return &tr.GetFailure().Error, nil
}

type accountUpdateResolver ZetaResolverRoot

func (r *accountUpdateResolver) AssetID(ctx context.Context, obj *v2.AccountBalance) (string, error) {
	return obj.Asset, nil
}

func (r *accountUpdateResolver) PartyID(ctx context.Context, obj *v2.AccountBalance) (string, error) {
	return obj.Owner, nil
}

// AggregatedLedgerEntriesResolver resolver.
type aggregatedLedgerEntriesResolver ZetaResolverRoot

func (r *ZetaResolverRoot) AggregatedLedgerEntry() AggregatedLedgerEntryResolver {
	return (*aggregatedLedgerEntriesResolver)(r)
}

func (r *aggregatedLedgerEntriesResolver) ZetaTime(ctx context.Context, obj *v2.AggregatedLedgerEntry) (int64, error) {
	return obj.Timestamp, nil
}

// LiquidityOrderReference resolver.

type myLiquidityOrderReferenceResolver ZetaResolverRoot

func (r *myLiquidityOrderReferenceResolver) Order(ctx context.Context, obj *types.LiquidityOrderReference) (*types.Order, error) {
	if len(obj.OrderId) <= 0 {
		return nil, nil
	}
	return r.r.getOrderByID(ctx, obj.OrderId, nil)
}

// deposit resolver

type myDepositResolver ZetaResolverRoot

func (r *myDepositResolver) Asset(ctx context.Context, obj *types.Deposit) (*types.Asset, error) {
	return r.r.getAssetByID(ctx, obj.Asset)
}

func (r *myDepositResolver) Party(ctx context.Context, obj *types.Deposit) (*types.Party, error) {
	if len(obj.PartyId) <= 0 {
		return nil, errors.New("missing party ID")
	}
	return &types.Party{Id: obj.PartyId}, nil
}

func (r *myDepositResolver) CreatedTimestamp(ctx context.Context, obj *types.Deposit) (string, error) {
	if obj.CreatedTimestamp == 0 {
		return "", errors.New("invalid timestamp")
	}
	return zetatime.Format(zetatime.UnixNano(obj.CreatedTimestamp)), nil
}

func (r *myDepositResolver) CreditedTimestamp(ctx context.Context, obj *types.Deposit) (*string, error) {
	if obj.CreditedTimestamp == 0 {
		return nil, nil
	}
	t := zetatime.Format(zetatime.UnixNano(obj.CreditedTimestamp))
	return &t, nil
}

// BEGIN: Query Resolver

type myQueryResolver ZetaResolverRoot

func (r *myQueryResolver) Positions(ctx context.Context, filter *v2.PositionsFilter, pagination *v2.Pagination) (*v2.PositionConnection, error) {
	resp, err := r.tradingDataClientV2.ListAllPositions(ctx, &v2.ListAllPositionsRequest{
		Filter:     filter,
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}

	return resp.Positions, nil
}

func (r *myQueryResolver) TransfersConnection(ctx context.Context, partyID *string, direction *TransferDirection,
	pagination *v2.Pagination,
) (*v2.TransferConnection, error) {
	return r.r.transfersConnection(ctx, partyID, direction, pagination)
}

func (r *myQueryResolver) LastBlockHeight(ctx context.Context) (string, error) {
	resp, err := r.tradingProxyClient.LastBlockHeight(ctx, &zetaprotoapi.LastBlockHeightRequest{})
	if err != nil {
		return "0", err
	}

	return strconv.FormatUint(resp.Height, 10), nil
}

func (r *myQueryResolver) OracleSpecsConnection(ctx context.Context, pagination *v2.Pagination) (*v2.OracleSpecsConnection, error) {
	req := v2.ListOracleSpecsRequest{
		Pagination: pagination,
	}
	res, err := r.tradingDataClientV2.ListOracleSpecs(ctx, &req)
	if err != nil {
		return nil, err
	}

	return res.OracleSpecs, nil
}

func (r *myQueryResolver) OracleSpec(ctx context.Context, id string) (*types.OracleSpec, error) {
	res, err := r.tradingDataClientV2.GetOracleSpec(
		ctx, &v2.GetOracleSpecRequest{OracleSpecId: id},
	)
	if err != nil {
		return nil, err
	}

	return res.OracleSpec, nil
}

func (r *myQueryResolver) OracleDataBySpecConnection(ctx context.Context, oracleSpecID string,
	pagination *v2.Pagination,
) (*v2.OracleDataConnection, error) {
	var specID *string
	if oracleSpecID != "" {
		specID = &oracleSpecID
	}
	req := v2.ListOracleDataRequest{
		OracleSpecId: specID,
		Pagination:   pagination,
	}

	resp, err := r.tradingDataClientV2.ListOracleData(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.OracleData, nil
}

func (r *myQueryResolver) OracleDataConnection(ctx context.Context, pagination *v2.Pagination) (*v2.OracleDataConnection, error) {
	req := v2.ListOracleDataRequest{
		Pagination: pagination,
	}

	resp, err := r.tradingDataClientV2.ListOracleData(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.OracleData, nil
}

func (r *myQueryResolver) NetworkParametersConnection(ctx context.Context, pagination *v2.Pagination) (*v2.NetworkParameterConnection, error) {
	res, err := r.tradingDataClientV2.ListNetworkParameters(ctx, &v2.ListNetworkParametersRequest{
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}
	return res.NetworkParameters, nil
}

func (r *myQueryResolver) NetworkParameter(ctx context.Context, key string) (*types.NetworkParameter, error) {
	res, err := r.tradingDataClientV2.GetNetworkParameter(
		ctx, &v2.GetNetworkParameterRequest{Key: key},
	)
	if err != nil {
		return nil, err
	}

	return res.NetworkParameter, nil
}

func (r *myQueryResolver) Erc20WithdrawalApproval(ctx context.Context, wid string) (*Erc20WithdrawalApproval, error) {
	res, err := r.tradingDataClientV2.GetERC20WithdrawalApproval(
		ctx, &v2.GetERC20WithdrawalApprovalRequest{WithdrawalId: wid},
	)
	if err != nil {
		return nil, err
	}

	return &Erc20WithdrawalApproval{
		AssetSource:   res.AssetSource,
		Amount:        res.Amount,
		Nonce:         res.Nonce,
		Signatures:    res.Signatures,
		TargetAddress: res.TargetAddress,
		Creation:      fmt.Sprintf("%d", res.Creation),
	}, nil
}

func (r *myQueryResolver) Erc20ListAssetBundle(ctx context.Context, assetID string) (*Erc20ListAssetBundle, error) {
	res, err := r.tradingDataClientV2.GetERC20ListAssetBundle(
		ctx, &v2.GetERC20ListAssetBundleRequest{AssetId: assetID})
	if err != nil {
		return nil, err
	}

	return &Erc20ListAssetBundle{
		AssetSource: res.AssetSource,
		ZetaAssetID: res.ZetaAssetId,
		Nonce:       res.Nonce,
		Signatures:  res.Signatures,
	}, nil
}

func (r *myQueryResolver) Erc20SetAssetLimitsBundle(ctx context.Context, proposalID string) (*ERC20SetAssetLimitsBundle, error) {
	res, err := r.tradingDataClientV2.GetERC20SetAssetLimitsBundle(
		ctx, &v2.GetERC20SetAssetLimitsBundleRequest{ProposalId: proposalID})
	if err != nil {
		return nil, err
	}

	return &ERC20SetAssetLimitsBundle{
		AssetSource:   res.AssetSource,
		ZetaAssetID:   res.ZetaAssetId,
		Nonce:         res.Nonce,
		LifetimeLimit: res.LifetimeLimit,
		Threshold:     res.Threshold,
		Signatures:    res.Signatures,
	}, nil
}

func (r *myQueryResolver) Erc20MultiSigSignerAddedBundles(ctx context.Context, nodeID string, submitter, epochSeq *string, pagination *v2.Pagination) (*ERC20MultiSigSignerAddedConnection, error) {
	res, err := r.tradingDataClientV2.ListERC20MultiSigSignerAddedBundles(
		ctx, &v2.ListERC20MultiSigSignerAddedBundlesRequest{
			NodeId:     nodeID,
			Submitter:  ptr.UnBox(submitter),
			EpochSeq:   ptr.UnBox(epochSeq),
			Pagination: pagination,
		})
	if err != nil {
		return nil, err
	}

	edges := make([]*ERC20MultiSigSignerAddedBundleEdge, 0, len(res.Bundles.Edges))

	for _, edge := range res.Bundles.Edges {
		edges = append(edges, &ERC20MultiSigSignerAddedBundleEdge{
			Node: &ERC20MultiSigSignerAddedBundle{
				NewSigner:  edge.Node.NewSigner,
				Submitter:  edge.Node.Submitter,
				Nonce:      edge.Node.Nonce,
				Timestamp:  fmt.Sprint(edge.Node.Timestamp),
				Signatures: edge.Node.Signatures,
				EpochSeq:   edge.Node.EpochSeq,
			},
			Cursor: edge.Cursor,
		})
	}

	return &ERC20MultiSigSignerAddedConnection{
		Edges:    edges,
		PageInfo: res.Bundles.PageInfo,
	}, nil
}

func (r *myQueryResolver) Erc20MultiSigSignerRemovedBundles(ctx context.Context, nodeID string, submitter, epochSeq *string, pagination *v2.Pagination) (*ERC20MultiSigSignerRemovedConnection, error) {
	res, err := r.tradingDataClientV2.ListERC20MultiSigSignerRemovedBundles(
		ctx, &v2.ListERC20MultiSigSignerRemovedBundlesRequest{
			NodeId:     nodeID,
			Submitter:  ptr.UnBox(submitter),
			EpochSeq:   ptr.UnBox(epochSeq),
			Pagination: pagination,
		})
	if err != nil {
		return nil, err
	}

	edges := make([]*ERC20MultiSigSignerRemovedBundleEdge, 0, len(res.Bundles.Edges))

	for _, edge := range res.Bundles.Edges {
		edges = append(edges, &ERC20MultiSigSignerRemovedBundleEdge{
			Node: &ERC20MultiSigSignerRemovedBundle{
				OldSigner:  edge.Node.OldSigner,
				Submitter:  edge.Node.Submitter,
				Nonce:      edge.Node.Nonce,
				Timestamp:  fmt.Sprint(edge.Node.Timestamp),
				Signatures: edge.Node.Signatures,
				EpochSeq:   edge.Node.EpochSeq,
			},
			Cursor: edge.Cursor,
		})
	}

	return &ERC20MultiSigSignerRemovedConnection{
		Edges:    edges,
		PageInfo: res.Bundles.PageInfo,
	}, nil
}

func (r *myQueryResolver) Withdrawal(ctx context.Context, wid string) (*types.Withdrawal, error) {
	res, err := r.tradingDataClientV2.GetWithdrawal(
		ctx, &v2.GetWithdrawalRequest{Id: wid},
	)
	if err != nil {
		return nil, err
	}

	return res.Withdrawal, nil
}

func (r *myQueryResolver) Withdrawals(ctx context.Context, dateRange *v2.DateRange, pagination *v2.Pagination) (*v2.WithdrawalsConnection, error) {
	res, err := r.tradingDataClientV2.ListWithdrawals(
		ctx, &v2.ListWithdrawalsRequest{
			DateRange:  dateRange,
			Pagination: pagination,
		},
	)
	if err != nil {
		return nil, err
	}

	return res.Withdrawals, nil
}

func (r *myQueryResolver) Deposit(ctx context.Context, did string) (*types.Deposit, error) {
	res, err := r.tradingDataClientV2.GetDeposit(
		ctx, &v2.GetDepositRequest{Id: did},
	)
	if err != nil {
		return nil, err
	}

	return res.Deposit, nil
}

func (r *myQueryResolver) Deposits(ctx context.Context, dateRange *v2.DateRange, pagination *v2.Pagination) (*v2.DepositsConnection, error) {
	res, err := r.tradingDataClientV2.ListDeposits(
		ctx, &v2.ListDepositsRequest{DateRange: dateRange, Pagination: pagination},
	)
	if err != nil {
		return nil, err
	}

	return res.Deposits, nil
}

func (r *myQueryResolver) EstimateOrder(
	ctx context.Context,
	market, party string,
	price *string,
	size string,
	side zeta.Side,
	timeInForce zeta.Order_TimeInForce,
	expiration *int64,
	ty zeta.Order_Type,
) (*OrderEstimate, error) {
	order := &types.Order{}

	var err error

	// We need to convert strings to uint64 (JS doesn't yet support uint64)
	if price != nil {
		order.Price = *price
	}
	s, err := safeStringUint64(size)
	if err != nil {
		return nil, err
	}
	order.Size = s
	if len(market) <= 0 {
		return nil, errors.New("market missing or empty")
	}
	order.MarketId = market
	if len(party) <= 0 {
		return nil, errors.New("party missing or empty")
	}

	order.PartyId = party
	order.TimeInForce = timeInForce
	order.Side = side
	order.Type = ty

	// GTT must have an expiration value
	if order.TimeInForce == types.Order_TIME_IN_FORCE_GTT && expiration != nil {
		order.ExpiresAt = zetatime.UnixNano(*expiration).UnixNano()
	}

	req := v2.EstimateFeeRequest{
		MarketId: order.MarketId,
		Price:    order.Price,
		Size:     order.Size,
	}

	// Pass the order over for consensus (service layer will use RPC client internally and handle errors etc)
	resp, err := r.tradingDataClientV2.EstimateFee(ctx, &req)
	if err != nil {
		r.log.Error("Failed to get fee estimates using rpc client in graphQL resolver", logging.Error(err))
		return nil, err
	}

	// calclate the fee total amount
	var mfee, ifee, lfee num.Decimal
	// errors doesn't matter here, they just give us zero values anyway for the decimals
	if len(resp.Fee.MakerFee) > 0 {
		mfee, _ = num.DecimalFromString(resp.Fee.MakerFee)
	}
	if len(resp.Fee.InfrastructureFee) > 0 {
		ifee, _ = num.DecimalFromString(resp.Fee.InfrastructureFee)
	}
	if len(resp.Fee.LiquidityFee) > 0 {
		lfee, _ = num.DecimalFromString(resp.Fee.LiquidityFee)
	}

	fee := TradeFee{
		MakerFee:          resp.Fee.MakerFee,
		InfrastructureFee: resp.Fee.InfrastructureFee,
		LiquidityFee:      resp.Fee.LiquidityFee,
	}

	// now we calculate the margins
	reqm := v2.EstimateMarginRequest{
		MarketId: order.MarketId,
		PartyId:  order.PartyId,
		Price:    order.Price,
		Size:     order.Size,
		Side:     order.Side,
		Type:     order.Type,
	}

	// Pass the order over for consensus (service layer will use RPC client internally and handle errors etc)
	respm, err := r.tradingDataClientV2.EstimateMargin(ctx, &reqm)
	if err != nil {
		r.log.Error("Failed to get margin estimates using rpc client in graphQL resolver", logging.Error(err))
		return nil, err
	}

	return &OrderEstimate{
		Fee:            &fee,
		TotalFeeAmount: decimal.Sum(mfee, ifee, lfee).String(),
		MarginLevels:   respm.MarginLevels,
	}, nil
}

func (r *myQueryResolver) Asset(ctx context.Context, id string) (*types.Asset, error) {
	return r.r.getAssetByID(ctx, id)
}

func (r *myQueryResolver) AssetsConnection(ctx context.Context, id *string, pagination *v2.Pagination) (*v2.AssetsConnection, error) {
	req := &v2.ListAssetsRequest{
		AssetId:    id,
		Pagination: pagination,
	}
	resp, err := r.tradingDataClientV2.ListAssets(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Assets, nil
}

func (r *myQueryResolver) NodeSignaturesConnection(ctx context.Context, resourceID string, pagination *v2.Pagination) (*v2.NodeSignaturesConnection, error) {
	if len(resourceID) <= 0 {
		return nil, ErrMissingIDOrReference
	}

	req := &v2.ListNodeSignaturesRequest{
		Id:         resourceID,
		Pagination: pagination,
	}
	res, err := r.tradingDataClientV2.ListNodeSignatures(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Signatures, nil
}

func (r *myQueryResolver) Market(ctx context.Context, id string) (*types.Market, error) {
	return r.r.getMarketByID(ctx, id)
}

func (r *myQueryResolver) Party(ctx context.Context, name string) (*types.Party, error) {
	return getParty(ctx, r.log, r.tradingDataClientV2, name)
}

func (r *myQueryResolver) OrderByID(ctx context.Context, orderID string, version *int) (*types.Order, error) {
	return r.r.getOrderByID(ctx, orderID, version)
}

func (r *myQueryResolver) OrderVersionsConnection(ctx context.Context, orderID *string, pagination *v2.Pagination) (*v2.OrderConnection, error) {
	if orderID == nil {
		return nil, ErrMissingIDOrReference
	}
	req := &v2.ListOrderVersionsRequest{
		OrderId:    *orderID,
		Pagination: pagination,
	}

	resp, err := r.tradingDataClientV2.ListOrderVersions(ctx, req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}
	return resp.Orders, nil
}

func (r *myQueryResolver) OrderByReference(ctx context.Context, reference string) (*types.Order, error) {
	req := &v2.ListOrdersRequest{
		Reference: &reference,
		Filter:    &v2.OrderFilter{},
	}
	res, err := r.tradingDataClientV2.ListOrders(ctx, req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	if len(res.Orders.Edges) == 0 {
		return nil, fmt.Errorf("order reference not found: %s", reference)
	}

	return res.Orders.Edges[0].Node, nil
}

func (r *myQueryResolver) ProposalsConnection(ctx context.Context, proposalType *v2.ListGovernanceDataRequest_Type, inState *zeta.Proposal_State,
	pagination *v2.Pagination,
) (*v2.GovernanceDataConnection, error) {
	return handleProposalsRequest(ctx, r.tradingDataClientV2, nil, nil, proposalType, inState, pagination)
}

func (r *myQueryResolver) Proposal(ctx context.Context, id *string, reference *string) (*types.GovernanceData, error) {
	if id != nil {
		resp, err := r.tradingDataClientV2.GetGovernanceData(ctx, &v2.GetGovernanceDataRequest{
			ProposalId: id,
		})
		if err != nil {
			return nil, err
		}
		return resp.Data, nil
	} else if reference != nil {
		resp, err := r.tradingDataClientV2.GetGovernanceData(ctx, &v2.GetGovernanceDataRequest{
			Reference: reference,
		})
		if err != nil {
			return nil, err
		}
		return resp.Data, nil
	}

	return nil, ErrMissingIDOrReference
}

func (r *myQueryResolver) ProtocolUpgradeStatus(ctx context.Context) (*ProtocolUpgradeStatus, error) {
	status, err := r.tradingDataClientV2.GetProtocolUpgradeStatus(ctx, &v2.GetProtocolUpgradeStatusRequest{})
	if err != nil {
		return nil, err
	}

	return &ProtocolUpgradeStatus{
		Ready: status.Ready,
	}, nil
}

func (r *myQueryResolver) CoreSnapshots(ctx context.Context, pagination *v2.Pagination) (*v2.CoreSnapshotConnection, error) {
	req := v2.ListCoreSnapshotsRequest{Pagination: pagination}
	resp, err := r.tradingDataClientV2.ListCoreSnapshots(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.CoreSnapshots, nil
}

func (r *myQueryResolver) EpochRewardSummaries(
	ctx context.Context,
	filter *v2.RewardSummaryFilter,
	pagination *v2.Pagination,
) (*v2.EpochRewardSummaryConnection, error) {
	req := v2.ListEpochRewardSummariesRequest{
		Filter:     filter,
		Pagination: pagination,
	}
	resp, err := r.tradingDataClientV2.ListEpochRewardSummaries(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp.Summaries, nil
}

func (r *myQueryResolver) ProtocolUpgradeProposals(
	ctx context.Context,
	inState *eventspb.ProtocolUpgradeProposalStatus,
	approvedBy *string,
	pagination *v2.Pagination,
) (
	*v2.ProtocolUpgradeProposalConnection, error,
) {
	req := v2.ListProtocolUpgradeProposalsRequest{Status: inState, ApprovedBy: approvedBy, Pagination: pagination}
	resp, err := r.tradingDataClientV2.ListProtocolUpgradeProposals(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.ProtocolUpgradeProposals, nil
}

func (r *myQueryResolver) NodeData(ctx context.Context) (*types.NodeData, error) {
	resp, err := r.tradingDataClientV2.GetNetworkData(ctx, &v2.GetNetworkDataRequest{})
	if err != nil {
		return nil, err
	}

	return resp.NodeData, nil
}

func (r *myQueryResolver) NodesConnection(ctx context.Context, pagination *v2.Pagination) (*v2.NodesConnection, error) {
	req := &v2.ListNodesRequest{
		Pagination: pagination,
	}
	resp, err := r.tradingDataClientV2.ListNodes(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Nodes, nil
}

func (r *myQueryResolver) Node(ctx context.Context, id string) (*types.Node, error) {
	resp, err := r.tradingDataClientV2.GetNode(ctx, &v2.GetNodeRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Node, nil
}

func (r *myQueryResolver) KeyRotationsConnection(ctx context.Context, id *string, pagination *v2.Pagination) (*v2.KeyRotationConnection, error) {
	resp, err := r.tradingDataClientV2.ListKeyRotations(ctx, &v2.ListKeyRotationsRequest{NodeId: id, Pagination: pagination})
	if err != nil {
		return nil, err
	}

	return resp.Rotations, nil
}

func (r *myQueryResolver) EthereumKeyRotations(ctx context.Context, nodeID *string) (*v2.EthereumKeyRotationsConnection, error) {
	resp, err := r.tradingDataClientV2.ListEthereumKeyRotations(ctx, &v2.ListEthereumKeyRotationsRequest{NodeId: nodeID})
	if err != nil {
		return nil, err
	}

	return resp.KeyRotations, nil
}

func (r *myQueryResolver) Epoch(ctx context.Context, id *string) (*types.Epoch, error) {
	var epochID *uint64
	if id != nil {
		parsedID, err := strconv.ParseUint(*id, 10, 64)
		if err != nil {
			return nil, err
		}

		epochID = &parsedID
	}

	resp, err := r.tradingDataClientV2.GetEpoch(ctx, &v2.GetEpochRequest{Id: epochID})
	if err != nil {
		return nil, err
	}

	return resp.Epoch, nil
}

func (r *myQueryResolver) Statistics(ctx context.Context) (*zetaprotoapi.Statistics, error) {
	req := &zetaprotoapi.StatisticsRequest{}
	resp, err := r.tradingProxyClient.Statistics(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetStatistics(), nil
}

func (r *myQueryResolver) BalanceChanges(
	ctx context.Context,
	filter *v2.AccountFilter,
	dateRange *v2.DateRange,
	pagination *v2.Pagination,
) (*v2.AggregatedBalanceConnection, error) {
	req := &v2.ListBalanceChangesRequest{
		Filter:     filter,
		DateRange:  dateRange,
		Pagination: pagination,
	}

	resp, err := r.tradingDataClientV2.ListBalanceChanges(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetBalances(), nil
}

func (r *myQueryResolver) LedgerEntries(
	ctx context.Context,
	filter *v2.LedgerEntryFilter,
	dateRange *v2.DateRange,
	pagination *v2.Pagination,
) (*v2.AggregatedLedgerEntriesConnection, error) {
	req := &v2.ListLedgerEntriesRequest{}
	req.Filter = filter

	req.DateRange = dateRange
	req.Pagination = pagination

	resp, err := r.tradingDataClientV2.ListLedgerEntries(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetLedgerEntries(), nil
}

func (r *myQueryResolver) NetworkLimits(ctx context.Context) (*types.NetworkLimits, error) {
	req := &v2.GetNetworkLimitsRequest{}
	resp, err := r.tradingDataClientV2.GetNetworkLimits(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetLimits(), nil
}

func (r *myQueryResolver) MostRecentHistorySegment(ctx context.Context) (*v2.HistorySegment, error) {
	req := &v2.GetMostRecentNetworkHistorySegmentRequest{}

	resp, err := r.tradingDataClientV2.GetMostRecentNetworkHistorySegment(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetSegment(), nil
}

// END: Root Resolver

type myNodeSignatureResolver ZetaResolverRoot

func (r *myNodeSignatureResolver) Signature(ctx context.Context, obj *commandspb.NodeSignature) (*string, error) {
	sig := base64.StdEncoding.EncodeToString(obj.Sig)
	return &sig, nil
}

// BEGIN: Party Resolver

type myPartyResolver ZetaResolverRoot

func (r *myPartyResolver) TransfersConnection(
	ctx context.Context,
	party *types.Party,
	direction *TransferDirection,
	pagination *v2.Pagination,
) (*v2.TransferConnection, error) {
	return r.r.transfersConnection(ctx, &party.Id, direction, pagination)
}

func (r *myPartyResolver) RewardsConnection(ctx context.Context, party *types.Party, assetID *string, pagination *v2.Pagination, fromEpoch *int, toEpoch *int) (*v2.RewardsConnection, error) {
	var from, to *uint64

	if fromEpoch != nil {
		from = new(uint64)
		if *fromEpoch < 0 {
			return nil, errors.New("invalid fromEpoch for reward query - must be positive")
		}
		*from = uint64(*fromEpoch)
	}
	if toEpoch != nil {
		to = new(uint64)
		if *toEpoch < 0 {
			return nil, errors.New("invalid toEpoch for reward query - must be positive")
		}
		*to = uint64(*toEpoch)
	}

	req := v2.ListRewardsRequest{
		PartyId:    party.Id,
		AssetId:    assetID,
		Pagination: pagination,
		FromEpoch:  from,
		ToEpoch:    to,
	}
	resp, err := r.tradingDataClientV2.ListRewards(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rewards information: %w", err)
	}

	return resp.Rewards, nil
}

func (r *myPartyResolver) RewardSummaries(
	ctx context.Context,
	party *types.Party,
	asset *string,
) ([]*types.RewardSummary, error) {
	var assetID string
	if asset != nil {
		assetID = *asset
	}

	req := &v2.ListRewardSummariesRequest{
		PartyId: &party.Id,
		AssetId: &assetID,
	}

	resp, err := r.tradingDataClientV2.ListRewardSummaries(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Summaries, err
}

func (r *myPartyResolver) StakingSummary(ctx context.Context, party *types.Party, pagination *v2.Pagination) (*StakingSummary, error) {
	if party == nil {
		return nil, errors.New("party must not be nil")
	}

	req := &v2.GetStakeRequest{
		PartyId:    party.Id,
		Pagination: pagination,
	}

	resp, err := r.tradingDataClientV2.GetStake(ctx, req)
	if err != nil {
		return nil, err
	}

	return &StakingSummary{
		CurrentStakeAvailable: resp.CurrentStakeAvailable,
		Linkings:              resp.StakeLinkings,
	}, nil
}

func (r *myPartyResolver) LiquidityProvisionsConnection(
	ctx context.Context,
	party *types.Party,
	market, ref *string,
	pagination *v2.Pagination,
) (*v2.LiquidityProvisionsConnection, error) {
	var partyID string
	if party != nil {
		partyID = party.Id
	}
	var mid string
	if market != nil {
		mid = *market
	}

	var refID string
	if ref != nil {
		refID = *ref
	}

	req := v2.ListLiquidityProvisionsRequest{
		PartyId:    &partyID,
		MarketId:   &mid,
		Reference:  &refID,
		Pagination: pagination,
	}

	res, err := r.tradingDataClientV2.ListLiquidityProvisions(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	return res.LiquidityProvisions, nil
}

func (r *myPartyResolver) MarginsConnection(ctx context.Context, party *types.Party, marketID *string,
	pagination *v2.Pagination,
) (*v2.MarginConnection, error) {
	if party == nil {
		return nil, errors.New("party is nil")
	}

	market := ""

	if marketID != nil {
		market = *marketID
	}

	req := v2.ListMarginLevelsRequest{
		PartyId:    party.Id,
		MarketId:   market,
		Pagination: pagination,
	}

	res, err := r.tradingDataClientV2.ListMarginLevels(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	return res.MarginLevels, nil
}

func (r *myPartyResolver) OrdersConnection(ctx context.Context, party *types.Party, dateRange *v2.DateRange,
	pagination *v2.Pagination, filter *v2.OrderFilter, marketID *string,
) (*v2.OrderConnection, error) {
	if party == nil {
		return nil, errors.New("party is required")
	}
	req := v2.ListOrdersRequest{
		PartyId:    &party.Id,
		Pagination: pagination,
		DateRange:  dateRange,
		Filter:     filter,
		MarketId:   marketID,
	}
	res, err := r.tradingDataClientV2.ListOrders(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}
	return res.Orders, nil
}

func (r *myPartyResolver) TradesConnection(ctx context.Context, party *types.Party, market *string, dateRange *v2.DateRange, pagination *v2.Pagination) (*v2.TradeConnection, error) {
	req := v2.ListTradesRequest{
		PartyId:    &party.Id,
		MarketId:   market,
		Pagination: pagination,
		DateRange:  dateRange,
	}

	res, err := r.tradingDataClientV2.ListTrades(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}
	return res.Trades, nil
}

func (r *myPartyResolver) PositionsConnection(ctx context.Context, party *types.Party, market *string, pagination *v2.Pagination) (*v2.PositionConnection, error) {
	partyID := ""
	if party != nil {
		partyID = party.Id
	}

	marketID := ""
	if market != nil {
		marketID = *market
	}

	req := v2.ListPositionsRequest{
		PartyId:    partyID,
		MarketId:   marketID,
		Pagination: pagination,
	}

	res, err := r.tradingDataClientV2.ListPositions(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	return res.Positions, nil
}

func (r *myPartyResolver) AccountsConnection(ctx context.Context, party *types.Party, marketID *string, asset *string, accType *types.AccountType, pagination *v2.Pagination) (*v2.AccountsConnection, error) {
	if party == nil {
		return nil, errors.New("a party must be specified when querying accounts")
	}
	var (
		marketIDs    = []string{}
		mktID        = ""
		asst         = ""
		accountTypes = []types.AccountType{}
		accTy        = types.AccountType_ACCOUNT_TYPE_UNSPECIFIED
		err          error
	)

	if marketID != nil {
		marketIDs = []string{*marketID}
		mktID = *marketID
	}

	if asset != nil {
		asst = *asset
	}
	if accType != nil {
		accTy = *accType
		if accTy != types.AccountType_ACCOUNT_TYPE_GENERAL &&
			accTy != types.AccountType_ACCOUNT_TYPE_MARGIN &&
			accTy != types.AccountType_ACCOUNT_TYPE_BOND {
			return nil, fmt.Errorf("invalid account type for party %v", accType)
		}
		accountTypes = []types.AccountType{accTy}
	}

	filter := v2.AccountFilter{
		AssetId:      asst,
		PartyIds:     []string{party.Id},
		MarketIds:    marketIDs,
		AccountTypes: accountTypes,
	}

	req := v2.ListAccountsRequest{Filter: &filter, Pagination: pagination}
	res, err := r.tradingDataClientV2.ListAccounts(ctx, &req)
	if err != nil {
		r.log.Error("unable to get Party account",
			logging.Error(err),
			logging.String("party-id", party.Id),
			logging.String("market-id", mktID),
			logging.String("asset", asst),
			logging.String("type", accTy.String()))
		return nil, err
	}

	return res.Accounts, nil
}

func (r *myPartyResolver) ProposalsConnection(ctx context.Context, party *types.Party, proposalType *v2.ListGovernanceDataRequest_Type, inState *zeta.Proposal_State,
	pagination *v2.Pagination,
) (*v2.GovernanceDataConnection, error) {
	return handleProposalsRequest(ctx, r.tradingDataClientV2, party, nil, proposalType, inState, pagination)
}

func (r *myPartyResolver) WithdrawalsConnection(ctx context.Context, party *types.Party, dateRange *v2.DateRange, pagination *v2.Pagination) (*v2.WithdrawalsConnection, error) {
	return handleWithdrawalsConnectionRequest(ctx, r.tradingDataClientV2, party, dateRange, pagination)
}

func (r *myPartyResolver) DepositsConnection(ctx context.Context, party *types.Party, dateRange *v2.DateRange, pagination *v2.Pagination) (*v2.DepositsConnection, error) {
	return handleDepositsConnectionRequest(ctx, r.tradingDataClientV2, party, dateRange, pagination)
}

func (r *myPartyResolver) VotesConnection(ctx context.Context, party *types.Party, pagination *v2.Pagination) (*ProposalVoteConnection, error) {
	req := v2.ListVotesRequest{
		PartyId:    &party.Id,
		Pagination: pagination,
	}

	res, err := r.tradingDataClientV2.ListVotes(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	edges := make([]*ProposalVoteEdge, 0, len(res.Votes.Edges))

	for _, edge := range res.Votes.Edges {
		edges = append(edges, &ProposalVoteEdge{
			Cursor: &edge.Cursor,
			Node:   ProposalVoteFromProto(edge.Node),
		})
	}

	connection := &ProposalVoteConnection{
		Edges:    edges,
		PageInfo: res.Votes.PageInfo,
	}

	return connection, nil
}

func (r *myPartyResolver) DelegationsConnection(ctx context.Context, party *types.Party, nodeID *string, pagination *v2.Pagination) (*v2.DelegationsConnection, error) {
	var partyID *string
	if party != nil {
		partyID = &party.Id
	}

	return handleDelegationConnectionRequest(ctx, r.tradingDataClientV2, partyID, nodeID, nil, pagination)
}

// END: Party Resolver

type myMarginLevelsUpdateResolver ZetaResolverRoot

func (r *myMarginLevelsUpdateResolver) InitialLevel(_ context.Context, m *types.MarginLevels) (string, error) {
	return m.InitialMargin, nil
}

func (r *myMarginLevelsUpdateResolver) SearchLevel(_ context.Context, m *types.MarginLevels) (string, error) {
	return m.SearchLevel, nil
}

func (r *myMarginLevelsUpdateResolver) MaintenanceLevel(_ context.Context, m *types.MarginLevels) (string, error) {
	return m.MaintenanceMargin, nil
}

// BEGIN: MarginLevels Resolver

type myMarginLevelsResolver ZetaResolverRoot

func (r *myMarginLevelsResolver) Market(ctx context.Context, m *types.MarginLevels) (*types.Market, error) {
	return r.r.getMarketByID(ctx, m.MarketId)
}

func (r *myMarginLevelsResolver) Party(ctx context.Context, m *types.MarginLevels) (*types.Party, error) {
	if m == nil {
		return nil, errors.New("nil order")
	}
	if len(m.PartyId) == 0 {
		return nil, errors.New("invalid party")
	}
	req := v2.GetPartyRequest{PartyId: m.PartyId}
	res, err := r.tradingDataClientV2.GetParty(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}
	return res.Party, nil
}

func (r *myMarginLevelsResolver) Asset(ctx context.Context, m *types.MarginLevels) (*types.Asset, error) {
	return r.r.getAssetByID(ctx, m.Asset)
}

func (r *myMarginLevelsResolver) CollateralReleaseLevel(_ context.Context, m *types.MarginLevels) (string, error) {
	return m.CollateralReleaseLevel, nil
}

func (r *myMarginLevelsResolver) InitialLevel(_ context.Context, m *types.MarginLevels) (string, error) {
	return m.InitialMargin, nil
}

func (r *myMarginLevelsResolver) SearchLevel(_ context.Context, m *types.MarginLevels) (string, error) {
	return m.SearchLevel, nil
}

func (r *myMarginLevelsResolver) MaintenanceLevel(_ context.Context, m *types.MarginLevels) (string, error) {
	return m.MaintenanceMargin, nil
}

// END: MarginLevels Resolver

type myOrderUpdateResolver ZetaResolverRoot

func (r *myOrderUpdateResolver) Price(ctx context.Context, obj *types.Order) (string, error) {
	return obj.Price, nil
}

func (r *myOrderUpdateResolver) Size(ctx context.Context, obj *types.Order) (string, error) {
	return strconv.FormatUint(obj.Size, 10), nil
}

func (r *myOrderUpdateResolver) Remaining(ctx context.Context, obj *types.Order) (string, error) {
	return strconv.FormatUint(obj.Remaining, 10), nil
}

func (r *myOrderUpdateResolver) CreatedAt(ctx context.Context, obj *types.Order) (int64, error) {
	return obj.CreatedAt, nil
}

func (r *myOrderUpdateResolver) UpdatedAt(ctx context.Context, obj *types.Order) (*int64, error) {
	var updatedAt *int64
	if obj.UpdatedAt > 0 {
		t := obj.UpdatedAt
		updatedAt = &t
	}
	return updatedAt, nil
}

func (r *myOrderUpdateResolver) Version(ctx context.Context, obj *types.Order) (string, error) {
	return strconv.FormatUint(obj.Version, 10), nil
}

func (r *myOrderUpdateResolver) ExpiresAt(ctx context.Context, obj *types.Order) (*string, error) {
	if obj.ExpiresAt <= 0 {
		return nil, nil
	}
	expiresAt := zetatime.Format(zetatime.UnixNano(obj.ExpiresAt))
	return &expiresAt, nil
}

func (r *myOrderUpdateResolver) RejectionReason(_ context.Context, o *types.Order) (*zeta.OrderError, error) {
	return o.Reason, nil
}

// BEGIN: Order Resolver

type myOrderResolver ZetaResolverRoot

func (r *myOrderResolver) RejectionReason(_ context.Context, o *types.Order) (*zeta.OrderError, error) {
	return o.Reason, nil
}

func (r *myOrderResolver) Price(ctx context.Context, obj *types.Order) (string, error) {
	return obj.Price, nil
}

func (r *myOrderResolver) Market(ctx context.Context, obj *types.Order) (*types.Market, error) {
	return r.r.getMarketByID(ctx, obj.MarketId)
}

func (r *myOrderResolver) Size(ctx context.Context, obj *types.Order) (string, error) {
	return strconv.FormatUint(obj.Size, 10), nil
}

func (r *myOrderResolver) Remaining(ctx context.Context, obj *types.Order) (string, error) {
	return strconv.FormatUint(obj.Remaining, 10), nil
}

func (r *myOrderResolver) CreatedAt(ctx context.Context, obj *types.Order) (int64, error) {
	return obj.CreatedAt, nil
}

func (r *myOrderResolver) UpdatedAt(ctx context.Context, obj *types.Order) (*int64, error) {
	var updatedAt *int64
	if obj.UpdatedAt > 0 {
		t := obj.UpdatedAt
		updatedAt = &t
	}
	return updatedAt, nil
}

func (r *myOrderResolver) Version(ctx context.Context, obj *types.Order) (string, error) {
	return strconv.FormatUint(obj.Version, 10), nil
}

func (r *myOrderResolver) ExpiresAt(ctx context.Context, obj *types.Order) (*string, error) {
	if obj.ExpiresAt <= 0 {
		return nil, nil
	}
	expiresAt := zetatime.Format(zetatime.UnixNano(obj.ExpiresAt))
	return &expiresAt, nil
}

func (r *myOrderResolver) TradesConnection(ctx context.Context, ord *types.Order, dateRange *v2.DateRange, pagination *v2.Pagination) (*v2.TradeConnection, error) {
	if ord == nil {
		return nil, errors.New("nil order")
	}
	req := v2.ListTradesRequest{OrderId: &ord.Id, Pagination: pagination, DateRange: dateRange}
	res, err := r.tradingDataClientV2.ListTrades(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}
	return res.Trades, nil
}

func (r *myOrderResolver) Party(ctx context.Context, order *types.Order) (*types.Party, error) {
	if order == nil {
		return nil, errors.New("nil order")
	}
	if len(order.PartyId) == 0 {
		return nil, errors.New("invalid party")
	}
	return &types.Party{Id: order.PartyId}, nil
}

func (r *myOrderResolver) PeggedOrder(ctx context.Context, order *types.Order) (*types.PeggedOrder, error) {
	return order.PeggedOrder, nil
}

func (r *myOrderResolver) LiquidityProvision(ctx context.Context, obj *types.Order) (*types.LiquidityProvision, error) {
	if obj == nil || len(obj.LiquidityProvisionId) <= 0 {
		return nil, nil
	}

	req := v2.ListLiquidityProvisionsRequest{
		PartyId:  &obj.PartyId,
		MarketId: &obj.MarketId,
	}
	res, err := r.tradingDataClientV2.ListLiquidityProvisions(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	if len(res.LiquidityProvisions.Edges) <= 0 {
		return nil, nil
	}

	return res.LiquidityProvisions.Edges[0].Node, nil
}

// END: Order Resolver

// BEGIN: Trade Resolver

type myTradeResolver ZetaResolverRoot

func (r *myTradeResolver) Market(ctx context.Context, obj *types.Trade) (*types.Market, error) {
	return r.r.getMarketByID(ctx, obj.MarketId)
}

func (r *myTradeResolver) Price(ctx context.Context, obj *types.Trade) (string, error) {
	return obj.Price, nil
}

func (r *myTradeResolver) Size(ctx context.Context, obj *types.Trade) (string, error) {
	return strconv.FormatUint(obj.Size, 10), nil
}

func (r *myTradeResolver) CreatedAt(ctx context.Context, obj *types.Trade) (int64, error) {
	return obj.Timestamp, nil
}

func (r *myTradeResolver) Buyer(ctx context.Context, obj *types.Trade) (*types.Party, error) {
	if obj == nil {
		return nil, errors.New("invalid trade")
	}
	if len(obj.Buyer) == 0 {
		return nil, errors.New("invalid buyer")
	}
	req := v2.GetPartyRequest{PartyId: obj.Buyer}
	res, err := r.tradingDataClientV2.GetParty(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}
	return res.Party, nil
}

func (r *myTradeResolver) Seller(ctx context.Context, obj *types.Trade) (*types.Party, error) {
	if obj == nil {
		return nil, errors.New("invalid trade")
	}
	if len(obj.Seller) == 0 {
		return nil, errors.New("invalid seller")
	}
	req := v2.GetPartyRequest{PartyId: obj.Seller}
	res, err := r.tradingDataClientV2.GetParty(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}
	return res.Party, nil
}

func (r *myTradeResolver) BuyerAuctionBatch(ctx context.Context, obj *types.Trade) (*int, error) {
	i := int(obj.BuyerAuctionBatch)
	return &i, nil
}

func (r *myTradeResolver) BuyerFee(ctx context.Context, obj *types.Trade) (*TradeFee, error) {
	fee := TradeFee{
		MakerFee:          "0",
		InfrastructureFee: "0",
		LiquidityFee:      "0",
	}
	if obj.BuyerFee != nil {
		fee.MakerFee = obj.BuyerFee.MakerFee
		fee.InfrastructureFee = obj.BuyerFee.InfrastructureFee
		fee.LiquidityFee = obj.BuyerFee.LiquidityFee
	}
	return &fee, nil
}

func (r *myTradeResolver) SellerAuctionBatch(ctx context.Context, obj *types.Trade) (*int, error) {
	i := int(obj.SellerAuctionBatch)
	return &i, nil
}

func (r *myTradeResolver) SellerFee(ctx context.Context, obj *types.Trade) (*TradeFee, error) {
	fee := TradeFee{
		MakerFee:          "0",
		InfrastructureFee: "0",
		LiquidityFee:      "0",
	}
	if obj.SellerFee != nil {
		fee.MakerFee = obj.SellerFee.MakerFee
		fee.InfrastructureFee = obj.SellerFee.InfrastructureFee
		fee.LiquidityFee = obj.SellerFee.LiquidityFee
	}
	return &fee, nil
}

// END: Trade Resolver

// BEGIN: Candle Resolver

type myCandleResolver ZetaResolverRoot

func (r *myCandleResolver) PeriodStart(_ context.Context, obj *v2.Candle) (int64, error) {
	return obj.Start, nil
}

func (r *myCandleResolver) LastUpdateInPeriod(_ context.Context, obj *v2.Candle) (int64, error) {
	return obj.LastUpdate, nil
}

func (r *myCandleResolver) Volume(_ context.Context, obj *v2.Candle) (string, error) {
	return strconv.FormatUint(obj.Volume, 10), nil
}

// END: Candle Resolver

// BEGIN: DataSourceSpecConfiguration Resolver.
type myDataSourceSpecConfigurationResolver ZetaResolverRoot

func (m *myDataSourceSpecConfigurationResolver) Signers(_ context.Context, obj *types.DataSourceSpecConfiguration) ([]*Signer, error) {
	return resolveSigners(obj.Signers), nil
}

// END: DataSourceSpecConfiguration Resolver

// BEGIN: Price Level Resolver

type myPriceLevelResolver ZetaResolverRoot

func (r *myPriceLevelResolver) Price(ctx context.Context, obj *types.PriceLevel) (string, error) {
	return obj.Price, nil
}

func (r *myPriceLevelResolver) Volume(ctx context.Context, obj *types.PriceLevel) (string, error) {
	return strconv.FormatUint(obj.Volume, 10), nil
}

func (r *myPriceLevelResolver) NumberOfOrders(ctx context.Context, obj *types.PriceLevel) (string, error) {
	return strconv.FormatUint(obj.NumberOfOrders, 10), nil
}

// END: Price Level Resolver

type positionUpdateResolver ZetaResolverRoot

func (r *positionUpdateResolver) OpenVolume(ctx context.Context, obj *types.Position) (string, error) {
	return strconv.FormatInt(obj.OpenVolume, 10), nil
}

func (r *positionUpdateResolver) UpdatedAt(ctx context.Context, obj *types.Position) (*string, error) {
	var updatedAt *string
	if obj.UpdatedAt > 0 {
		t := zetatime.Format(zetatime.UnixNano(obj.UpdatedAt))
		updatedAt = &t
	}
	return updatedAt, nil
}

func (r *positionUpdateResolver) LossSocializationAmount(ctx context.Context, obj *types.Position) (string, error) {
	return obj.LossSocialisationAmount, nil
}

// BEGIN: Position Resolver

type myPositionResolver ZetaResolverRoot

func (r *myPositionResolver) Market(ctx context.Context, obj *types.Position) (*types.Market, error) {
	return r.r.getMarketByID(ctx, obj.MarketId)
}

func (r *myPositionResolver) UpdatedAt(ctx context.Context, obj *types.Position) (*string, error) {
	var updatedAt *string
	if obj.UpdatedAt > 0 {
		t := zetatime.Format(zetatime.UnixNano(obj.UpdatedAt))
		updatedAt = &t
	}
	return updatedAt, nil
}

func (r *myPositionResolver) OpenVolume(ctx context.Context, obj *types.Position) (string, error) {
	return strconv.FormatInt(obj.OpenVolume, 10), nil
}

func (r *myPositionResolver) RealisedPnl(ctx context.Context, obj *types.Position) (string, error) {
	return obj.RealisedPnl, nil
}

func (r *myPositionResolver) UnrealisedPnl(ctx context.Context, obj *types.Position) (string, error) {
	return obj.UnrealisedPnl, nil
}

func (r *myPositionResolver) AverageEntryPrice(ctx context.Context, obj *types.Position) (string, error) {
	return obj.AverageEntryPrice, nil
}

func (r *myPositionResolver) LossSocializationAmount(ctx context.Context, obj *types.Position) (string, error) {
	return obj.LossSocialisationAmount, nil
}

func (r *myPositionResolver) Party(ctx context.Context, obj *types.Position) (*types.Party, error) {
	return getParty(ctx, r.log, r.tradingDataClientV2, obj.PartyId)
}

func (r *myPositionResolver) MarginsConnection(ctx context.Context, pos *types.Position, pagination *v2.Pagination) (*v2.MarginConnection, error) {
	req := v2.ListMarginLevelsRequest{
		PartyId:    pos.PartyId,
		MarketId:   pos.MarketId,
		Pagination: pagination,
	}

	res, err := r.tradingDataClientV2.ListMarginLevels(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	return res.MarginLevels, nil
}

// END: Position Resolver

// BEGIN: Subscription Resolver

type mySubscriptionResolver ZetaResolverRoot

func (r *mySubscriptionResolver) Delegations(ctx context.Context, party, nodeID *string) (<-chan *types.Delegation, error) {
	req := &v2.ObserveDelegationsRequest{
		PartyId: party,
		NodeId:  nodeID,
	}
	stream, err := r.tradingDataClientV2.ObserveDelegations(ctx, req)
	if err != nil {
		return nil, err
	}

	sCtx := stream.Context()
	ch := make(chan *types.Delegation)
	go func() {
		defer func() {
			stream.CloseSend()
			close(ch)
		}()
		for {
			dl, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("delegations: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("delegations levls: stream closed", logging.Error(err))
				break
			}
			select {
			case ch <- dl.Delegation:
				r.log.Debug("delegations: data sent")
			case <-ctx.Done():
				r.log.Error("delegations: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("delegations: stream closed by server")
				break
			}
		}
	}()

	return ch, nil
}

func (r *mySubscriptionResolver) Rewards(ctx context.Context, assetID, party *string) (<-chan *types.Reward, error) {
	req := &v2.ObserveRewardsRequest{
		AssetId: assetID,
		PartyId: party,
	}
	stream, err := r.tradingDataClientV2.ObserveRewards(ctx, req)
	if err != nil {
		return nil, err
	}

	sCtx := stream.Context()
	ch := make(chan *types.Reward)
	go func() {
		defer func() {
			stream.CloseSend()
			close(ch)
		}()
		for {
			rd, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("reward details: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("reward details: stream closed", logging.Error(err))
				break
			}
			select {
			case ch <- rd.Reward:
				r.log.Debug("rewards: data sent")
			case <-ctx.Done():
				r.log.Error("rewards: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("rewards: stream closed by server")
				break
			}
		}
	}()

	return ch, nil
}

func (r *mySubscriptionResolver) Margins(ctx context.Context, partyID string, marketID *string) (<-chan *types.MarginLevels, error) {
	req := &v2.ObserveMarginLevelsRequest{
		MarketId: marketID,
		PartyId:  partyID,
	}
	stream, err := r.tradingDataClientV2.ObserveMarginLevels(ctx, req)
	if err != nil {
		return nil, err
	}

	sCtx := stream.Context()
	ch := make(chan *types.MarginLevels)
	go func() {
		defer func() {
			stream.CloseSend()
			close(ch)
		}()
		for {
			m, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("margin levels: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("margin levls: stream closed", logging.Error(err))
				break
			}
			select {
			case ch <- m.MarginLevels:
				r.log.Debug("margin levels: data sent")
			case <-ctx.Done():
				r.log.Error("margin levels: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("margin levels: stream closed by server")
				break
			}
		}
	}()

	return ch, nil
}

func (r *mySubscriptionResolver) Accounts(ctx context.Context, marketID *string, partyID *string, asset *string, typeArg *types.AccountType) (<-chan []*v2.AccountBalance, error) {
	var (
		mkt, pty, ast string
		ty            types.AccountType
	)

	if marketID == nil && partyID == nil && asset == nil && typeArg == nil {
		// Updates on every balance update, on every account, for everyone and shouldn't be allowed for GraphQL.
		return nil, errors.New("at least one query filter must be applied for this subscription")
	}
	if asset != nil {
		ast = *asset
	}
	if marketID != nil {
		mkt = *marketID
	}
	if partyID != nil {
		pty = *partyID
	}
	if typeArg != nil {
		ty = *typeArg
	}

	req := &v2.ObserveAccountsRequest{
		Asset:    ast,
		MarketId: mkt,
		PartyId:  pty,
		Type:     ty,
	}
	stream, err := r.tradingDataClientV2.ObserveAccounts(ctx, req)
	if err != nil {
		return nil, err
	}

	c := make(chan []*v2.AccountBalance)
	data := []*v2.AccountBalance{}
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(c)
		}()
		for {
			a, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("accounts: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("accounts: stream closed", logging.Error(err))
				break
			}

			// empty slice, but preserve cap to avoid excessive reallocation
			data = data[:0]
			if snapshot := a.GetSnapshot(); snapshot != nil {
				data = append(data, snapshot.Accounts...)
			}

			if updates := a.GetUpdates(); updates != nil {
				data = append(data, updates.Accounts...)
			}
			select {
			case c <- data:
				r.log.Debug("accounts: data sent")
			case <-ctx.Done():
				r.log.Error("accounts: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("accounts: stream closed by server")
				break
			}
		}
	}()

	return c, nil
}

func (r *mySubscriptionResolver) Orders(ctx context.Context, market *string, party *string) (<-chan []*types.Order, error) {
	req := &v2.ObserveOrdersRequest{
		MarketId: market,
		PartyId:  party,
	}
	stream, err := r.tradingDataClientV2.ObserveOrders(ctx, req)
	if err != nil {
		return nil, err
	}

	c := make(chan []*types.Order)
	data := []*types.Order{}
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(c)
		}()
		for {
			o, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("orders: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("orders: stream closed", logging.Error(err))
				break
			}
			data = data[:0]
			if snapshot := o.GetSnapshot(); snapshot != nil {
				data = append(data, snapshot.Orders...)
			}
			if updates := o.GetUpdates(); updates != nil {
				data = append(data, updates.Orders...)
			}
			select {
			case c <- data:
				r.log.Debug("orders: data sent")
			case <-ctx.Done():
				r.log.Error("orders: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("orders: stream closed by server")
				break
			}
		}
	}()

	return c, nil
}

func (r *mySubscriptionResolver) Trades(ctx context.Context, market *string, party *string) (<-chan []*types.Trade, error) {
	req := &v2.ObserveTradesRequest{
		MarketId: market,
		PartyId:  party,
	}
	stream, err := r.tradingDataClientV2.ObserveTrades(ctx, req)
	if err != nil {
		return nil, err
	}

	c := make(chan []*types.Trade)
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(c)
		}()
		for {
			t, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("trades: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("trades: stream closed", logging.Error(err))
				break
			}
			select {
			case c <- t.Trades:
				r.log.Debug("trades: data sent")
			case <-ctx.Done():
				r.log.Error("trades: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("trades: stream closed by server")
				break
			}
		}
	}()

	return c, nil
}

func (r *mySubscriptionResolver) Positions(ctx context.Context, party, market *string) (<-chan []*types.Position, error) {
	req := &v2.ObservePositionsRequest{
		PartyId:  party,
		MarketId: market,
	}
	stream, err := r.tradingDataClientV2.ObservePositions(ctx, req)
	if err != nil {
		return nil, err
	}

	c := make(chan []*types.Position)
	data := []*types.Position{}
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(c)
		}()
		for {
			t, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("positions: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("positions: stream closed", logging.Error(err))
				break
			}
			data = data[:0]
			if snapshot := t.GetSnapshot(); snapshot != nil {
				data = append(data, snapshot.Positions...)
			}

			if updates := t.GetUpdates(); updates != nil {
				data = append(data, updates.Positions...)
			}
			select {
			case c <- data:
				r.log.Debug("positions: data sent")
			case <-ctx.Done():
				r.log.Error("positions: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("positions: stream closed by server")
				break
			}
		}
	}()

	return c, nil
}

func (r *mySubscriptionResolver) Candles(ctx context.Context, market string, interval zeta.Interval) (<-chan *v2.Candle, error) {
	intervalToCandleIDs, err := r.tradingDataClientV2.ListCandleIntervals(ctx, &v2.ListCandleIntervalsRequest{
		MarketId: market,
	})
	if err != nil {
		return nil, err
	}

	candleID := ""
	var candleInterval types.Interval
	for _, ic := range intervalToCandleIDs.IntervalToCandleId {
		candleInterval, err = convertDataNodeIntervalToProto(ic.Interval)
		if err != nil {
			r.log.Errorf("convert interval to candle id failed: %v", err)
			continue
		}
		if candleInterval == interval {
			candleID = ic.CandleId
			break
		}
	}

	if candleID == "" {
		return nil, fmt.Errorf("candle information not found for market: %s, interval: %s", market, interval)
	}

	req := &v2.ObserveCandleDataRequest{
		CandleId: candleID,
	}
	stream, err := r.tradingDataClientV2.ObserveCandleData(ctx, req)
	if err != nil {
		return nil, err
	}

	sCtx := stream.Context()
	c := make(chan *v2.Candle)
	go func() {
		defer func() {
			stream.CloseSend()
			close(c)
		}()
		for {
			cdl, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("candles: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("candles: stream closed", logging.Error(err))
				break
			}

			select {
			case c <- cdl.Candle:
				r.log.Debug("candles: data sent")
			case <-ctx.Done():
				r.log.Error("candles: stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("candles: stream closed by server")
				break
			}
		}
	}()
	return c, nil
}

func isStreamClosed(err error, log *logging.Logger) bool {
	if err == io.EOF {
		log.Error("stream closed by server", logging.Error(err))
		return true
	}
	if err != nil {
		log.Error("stream closed", logging.Error(err))
		return true
	}
	return false
}

func (r *mySubscriptionResolver) subscribeAllProposals(ctx context.Context) (<-chan *types.GovernanceData, error) {
	stream, err := r.tradingDataClientV2.ObserveGovernance(ctx, &v2.ObserveGovernanceRequest{})
	if err != nil {
		return nil, err
	}
	output := make(chan *types.GovernanceData)
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(output)
		}()
		for data, err := stream.Recv(); !isStreamClosed(err, r.log); data, err = stream.Recv() {
			select {
			case output <- data.Data:
				r.log.Debug("governance (all): data sent")
			case <-ctx.Done():
				r.log.Error("governance (all): stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("governance (all): stream closed by server")
				break
			}
		}
	}()
	return output, nil
}

func (r *mySubscriptionResolver) subscribePartyProposals(ctx context.Context, partyID string) (<-chan *types.GovernanceData, error) {
	stream, err := r.tradingDataClientV2.ObserveGovernance(ctx, &v2.ObserveGovernanceRequest{
		PartyId: &partyID,
	})
	if err != nil {
		return nil, err
	}
	sCtx := stream.Context()
	output := make(chan *types.GovernanceData)
	go func() {
		defer func() {
			stream.CloseSend()
			close(output)
		}()
		for data, err := stream.Recv(); !isStreamClosed(err, r.log); data, err = stream.Recv() {
			select {
			case output <- data.Data:
				r.log.Debug("governance (party): data sent")
			case <-ctx.Done():
				r.log.Error("governance (party): stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("governance (party): stream closed by server")
				break
			}
		}
	}()
	return output, nil
}

func (r *mySubscriptionResolver) Proposals(ctx context.Context, partyID *string) (<-chan *types.GovernanceData, error) {
	if partyID != nil && len(*partyID) > 0 {
		return r.subscribePartyProposals(ctx, *partyID)
	}
	return r.subscribeAllProposals(ctx)
}

func (r *mySubscriptionResolver) subscribeProposalVotes(ctx context.Context, proposalID string) (<-chan *ProposalVote, error) {
	output := make(chan *ProposalVote)
	stream, err := r.tradingDataClientV2.ObserveVotes(ctx, &v2.ObserveVotesRequest{
		ProposalId: &proposalID,
	})
	if err != nil {
		return nil, err
	}
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(output)
		}()
		for {
			data, err := stream.Recv()
			if isStreamClosed(err, r.log) {
				break
			}
			select {
			case output <- ProposalVoteFromProto(data.Vote):
				r.log.Debug("votes (proposal): data sent")
			case <-ctx.Done():
				r.log.Error("votes (proposal): stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("votes (proposal): stream closed by server")
				break
			}
		}
	}()
	return output, nil
}

func (r *mySubscriptionResolver) subscribePartyVotes(ctx context.Context, partyID string) (<-chan *ProposalVote, error) {
	output := make(chan *ProposalVote)
	stream, err := r.tradingDataClientV2.ObserveVotes(ctx, &v2.ObserveVotesRequest{
		PartyId: &partyID,
	})
	if err != nil {
		return nil, err
	}
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(output)
		}()
		for {
			data, err := stream.Recv()
			if isStreamClosed(err, r.log) {
				break
			}
			select {
			case output <- ProposalVoteFromProto(data.Vote):
				r.log.Debug("votes (party): data sent")
			case <-ctx.Done():
				r.log.Error("votes (party): stream closed")
				break
			case <-sCtx.Done():
				r.log.Error("votes (party): stream closed by server")
				break
			}
		}
	}()
	return output, nil
}

func (r *mySubscriptionResolver) Votes(ctx context.Context, proposalID *string, partyID *string) (<-chan *ProposalVote, error) {
	if proposalID != nil && len(*proposalID) != 0 {
		return r.subscribeProposalVotes(ctx, *proposalID)
	} else if partyID != nil && len(*partyID) != 0 {
		return r.subscribePartyVotes(ctx, *partyID)
	}
	return nil, ErrInvalidVotesSubscription
}

func (r *mySubscriptionResolver) BusEvents(ctx context.Context, types []BusEventType, marketID, partyID *string, batchSize int) (<-chan []*BusEvent, error) {
	if len(types) > 1 {
		return nil, errors.New("busEvents subscription support streaming 1 event at a time for now")
	}
	if len(types) <= 0 {
		return nil, errors.New("busEvents subscription requires 1 event type")
	}
	t := eventTypeToProto(types...)
	req := v2.ObserveEventBusRequest{
		Type:      t,
		BatchSize: int64(batchSize),
	}
	if req.BatchSize == 0 {
		// req.BatchSize = -1 // sending this with -1 to indicate to underlying gRPC call this is a special case: GQL
		batchSize = 0
	}
	if marketID != nil {
		req.MarketId = *marketID
	}
	if partyID != nil {
		req.PartyId = *partyID
	}
	mb := 10
	// about 10MB message size allowed
	msgSize := grpc.MaxCallRecvMsgSize(mb * 10e6)

	// build the bidirectional stream connection
	stream, err := r.tradingDataClientV2.ObserveEventBus(ctx, msgSize)
	if err != nil {
		return nil, err
	}

	// send our initial message to initialize the connection
	if err := stream.Send(&req); err != nil {
		return nil, err
	}

	// we no longer buffer this channel. Client receives batch, then we request the next batch
	out := make(chan []*BusEvent)

	go func() {
		defer func() {
			stream.CloseSend()
			close(out)
		}()

		if batchSize == 0 {
			r.busEvents(ctx, stream, out)
		} else {
			r.busEventsWithBatch(ctx, int64(batchSize), stream, out)
		}
	}()

	return out, nil
}

func (r *mySubscriptionResolver) busEvents(ctx context.Context, stream v2.TradingDataService_ObserveEventBusClient, out chan []*BusEvent) {
	sCtx := stream.Context()
	for {
		// receive batch
		data, err := stream.Recv()
		if isStreamClosed(err, r.log) {
			return
		}
		if err != nil {
			r.log.Error("Event bus stream error", logging.Error(err))
			return
		}
		select {
		case out <- busEventFromProto(data.Events...):
			r.log.Debug("bus events: data sent")
		case <-ctx.Done():
			r.log.Debug("bus events: stream closed")
			return
		case <-sCtx.Done():
			r.log.Debug("bus events: stream closed by server")
			return
		}
	}
}

func (r *mySubscriptionResolver) busEventsWithBatch(ctx context.Context, batchSize int64, stream v2.TradingDataService_ObserveEventBusClient, out chan []*BusEvent) {
	sCtx := stream.Context()
	poll := &v2.ObserveEventBusRequest{
		BatchSize: batchSize,
	}
	for {
		// receive batch
		data, err := stream.Recv()
		if isStreamClosed(err, r.log) {
			return
		}
		if err != nil {
			r.log.Error("Event bus stream error", logging.Error(err))
			return
		}
		select {
		case out <- busEventFromProto(data.Events...):
			r.log.Debug("bus events: data sent")
		case <-ctx.Done():
			r.log.Debug("bus events: stream closed")
			return
		case <-sCtx.Done():
			r.log.Debug("bus events: stream closed by server")
			return
		}
		// send request for the next batch
		if err := stream.SendMsg(poll); err != nil {
			r.log.Error("Failed to poll next event batch", logging.Error(err))
			return
		}
	}
}

func (r *mySubscriptionResolver) LiquidityProvisions(ctx context.Context, partyID *string, marketID *string) (<-chan []*types.LiquidityProvision, error) {
	req := &v2.ObserveLiquidityProvisionsRequest{
		MarketId: marketID,
		PartyId:  partyID,
	}
	stream, err := r.tradingDataClientV2.ObserveLiquidityProvisions(ctx, req)
	if err != nil {
		return nil, err
	}

	c := make(chan []*types.LiquidityProvision)
	sCtx := stream.Context()
	go func() {
		defer func() {
			stream.CloseSend()
			close(c)
		}()
		for {
			received, err := stream.Recv()
			if err == io.EOF {
				r.log.Error("orders: stream closed by server", logging.Error(err))
				break
			}
			if err != nil {
				r.log.Error("orders: stream closed", logging.Error(err))
				break
			}
			lps := received.LiquidityProvisions
			if len(lps) == 0 {
				continue
			}
			select {
			case c <- lps:
				r.log.Debug("liquidity provisions: data sent")
			case <-sCtx.Done():
				r.log.Debug("liquidity provisions: stream closed by server")
				break
			case <-ctx.Done():
				r.log.Debug("liquidity provisions: stream closed")
				break
			}
		}
	}()

	return c, nil
}

type myAccountDetailsResolver ZetaResolverRoot

func (r *myAccountDetailsResolver) PartyID(ctx context.Context, acc *types.AccountDetails) (*string, error) {
	if acc.Owner != nil {
		return acc.Owner, nil
	}
	return nil, nil
}

// START: Account Resolver

type myAccountResolver ZetaResolverRoot

func (r *myAccountResolver) Balance(ctx context.Context, acc *v2.AccountBalance) (string, error) {
	return acc.Balance, nil
}

func (r *myAccountResolver) Market(ctx context.Context, acc *v2.AccountBalance) (*types.Market, error) {
	if acc.MarketId == "" {
		return nil, nil
	}
	return r.r.getMarketByID(ctx, acc.MarketId)
}

func (r *myAccountResolver) Party(ctx context.Context, acc *v2.AccountBalance) (*types.Party, error) {
	if acc.Owner == "" {
		return nil, nil
	}
	return getParty(ctx, r.log, r.r.clt2, acc.Owner)
}

func (r *myAccountResolver) Asset(ctx context.Context, obj *v2.AccountBalance) (*types.Asset, error) {
	return r.r.getAssetByID(ctx, obj.Asset)
}

// START: Account Resolver

type myAccountEventResolver ZetaResolverRoot

func (r *myAccountEventResolver) Balance(ctx context.Context, acc *zeta.Account) (string, error) {
	return acc.Balance, nil
}

func (r *myAccountEventResolver) Market(ctx context.Context, acc *zeta.Account) (*types.Market, error) {
	if acc.MarketId == "" {
		return nil, nil
	}
	return r.r.getMarketByID(ctx, acc.MarketId)
}

func (r *myAccountEventResolver) Party(ctx context.Context, acc *zeta.Account) (*types.Party, error) {
	if acc.Owner == "" {
		return nil, nil
	}
	return getParty(ctx, r.log, r.r.clt2, acc.Owner)
}

func (r *myAccountEventResolver) Asset(ctx context.Context, obj *zeta.Account) (*types.Asset, error) {
	return r.r.getAssetByID(ctx, obj.Asset)
}

// END: Account Resolver

func getParty(ctx context.Context, _ *logging.Logger, client TradingDataServiceClientV2, id string) (*types.Party, error) {
	if len(id) == 0 {
		return nil, nil
	}
	res, err := client.GetParty(ctx, &v2.GetPartyRequest{PartyId: id})
	if err != nil {
		return nil, err
	}
	return res.Party, nil
}

// Market Data Resolvers.
type myPropertyKeyResolver ZetaResolverRoot

func (r *myPropertyKeyResolver) NumberDecimalPlaces(ctx context.Context, obj *data.PropertyKey) (*int, error) {
	ndp := obj.NumberDecimalPlaces
	if ndp == nil {
		return nil, nil
	}
	indp := new(int)
	*indp = int(*ndp)
	return indp, nil
}

// GetMarketDataHistoryByID returns all the market data information for a given market between the dates specified.
func (r *myQueryResolver) GetMarketDataHistoryByID(ctx context.Context, id string, start, end *int64, skip, first, last *int) ([]*types.MarketData, error) {
	pagination := makeAPIV2Pagination(skip, first, last)

	return r.getMarketDataHistoryByID(ctx, id, start, end, pagination)
}

func makeAPIV2Pagination(skip, first, last *int) *v2.OffsetPagination {
	var (
		offset, limit uint64
		descending    bool
	)
	if skip != nil {
		offset = uint64(*skip)
	}
	if last != nil {
		limit = uint64(*last)
		descending = true
	} else if first != nil {
		limit = uint64(*first)
	}
	return &v2.OffsetPagination{
		Skip:       offset,
		Limit:      limit,
		Descending: descending,
	}
}

func (r *myQueryResolver) getMarketData(ctx context.Context, req *v2.GetMarketDataHistoryByIDRequest) ([]*types.MarketData, error) {
	resp, err := r.tradingDataClientV2.GetMarketDataHistoryByID(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.MarketData == nil {
		return nil, errors.New("no market data not found")
	}

	results := make([]*types.MarketData, 0, len(resp.MarketData.Edges))

	for _, edge := range resp.MarketData.Edges {
		results = append(results, edge.Node)
	}

	return results, nil
}

func (r *myQueryResolver) getMarketDataHistoryByID(ctx context.Context, id string, start, end *int64, pagination *v2.OffsetPagination) ([]*types.MarketData, error) {
	var startTime, endTime *int64

	if start != nil {
		s := time.Unix(*start, 0).UnixNano()
		startTime = &s
	}

	if end != nil {
		e := time.Unix(*end, 0).UnixNano()
		endTime = &e
	}

	req := v2.GetMarketDataHistoryByIDRequest{
		MarketId:         id,
		StartTimestamp:   startTime,
		EndTimestamp:     endTime,
		OffsetPagination: pagination,
	}

	return r.getMarketData(ctx, &req)
}

func (r *myQueryResolver) GetMarketDataHistoryConnectionByID(ctx context.Context, marketID string, start, end *int64, pagination *v2.Pagination) (*v2.MarketDataConnection, error) {
	req := v2.GetMarketDataHistoryByIDRequest{
		MarketId:       marketID,
		StartTimestamp: start,
		EndTimestamp:   end,
		Pagination:     pagination,
	}

	resp, err := r.tradingDataClientV2.GetMarketDataHistoryByID(ctx, &req)
	if err != nil {
		r.log.Error("tradingData client", logging.Error(err))
		return nil, err
	}

	return resp.GetMarketData(), nil
}

func (r *myQueryResolver) MarketsConnection(ctx context.Context, id *string, pagination *v2.Pagination, includeSettled *bool) (*v2.MarketConnection, error) {
	var marketID string

	if id != nil {
		marketID = *id

		resp, err := r.tradingDataClientV2.GetMarket(ctx, &v2.GetMarketRequest{MarketId: marketID})
		if err != nil {
			return nil, err
		}

		connection := &v2.MarketConnection{
			Edges: []*v2.MarketEdge{
				{
					Node:   resp.Market,
					Cursor: "",
				},
			},
			PageInfo: &v2.PageInfo{
				HasNextPage:     false,
				HasPreviousPage: false,
				StartCursor:     "",
				EndCursor:       "",
			},
		}

		return connection, nil
	}

	resp, err := r.tradingDataClientV2.ListMarkets(ctx, &v2.ListMarketsRequest{
		Pagination:     pagination,
		IncludeSettled: includeSettled,
	})
	if err != nil {
		return nil, err
	}

	return resp.Markets, nil
}

func (r *myQueryResolver) PartiesConnection(ctx context.Context, id *string, pagination *v2.Pagination) (*v2.PartyConnection, error) {
	var partyID string
	if id != nil {
		partyID = *id
	}
	resp, err := r.tradingDataClientV2.ListParties(ctx, &v2.ListPartiesRequest{
		PartyId:    partyID,
		Pagination: pagination,
	})
	if err != nil {
		return nil, err
	}

	return resp.Parties, nil
}
