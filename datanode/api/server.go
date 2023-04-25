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

package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"zuluprotocol/zeta/libs/subscribers"

	"zuluprotocol/zeta/datanode/networkhistory"
	"zuluprotocol/zeta/datanode/ratelimit"

	"zuluprotocol/zeta/core/events"
	"zuluprotocol/zeta/datanode/candlesv2"
	"zuluprotocol/zeta/datanode/contextutil"
	"zuluprotocol/zeta/datanode/entities"
	"zuluprotocol/zeta/datanode/service"
	"zuluprotocol/zeta/logging"
	protoapi "zuluprotocol/zeta/protos/data-node/api/v2"
	zetaprotoapi "code.zetaprotocol.io/zeta/protos/zeta/api/v1"
	eventspb "zuluprotocol/zeta/protos/zeta/events/v1"

	"github.com/fullstorydev/grpcui/standalone"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

// EventService ...
//
//go:generate go run github.com/golang/mock/mockgen -destination mocks/event_service_mock.go -package mocks zuluprotocol/zeta/datanode/api EventService
type EventService interface {
	ObserveEvents(ctx context.Context, retries int, eTypes []events.Type, batchSize int, filters ...subscribers.EventFilter) (<-chan []*eventspb.BusEvent, chan<- int)
}

// BlockService ...
//
//go:generate go run github.com/golang/mock/mockgen -destination mocks/block_service_mock.go -package mocks zuluprotocol/zeta/datanode/api BlockService
type BlockService interface {
	GetLastBlock(ctx context.Context) (entities.Block, error)
}

// NetworkHistoryService ...
//
//go:generate go run github.com/golang/mock/mockgen -destination mocks/networkhistory_service_mock.go -package mocks zuluprotocol/zeta/datanode/api NetworkHistoryService
type NetworkHistoryService interface {
	GetHighestBlockHeightHistorySegment() (networkhistory.Segment, error)
	ListAllHistorySegments() ([]networkhistory.Segment, error)
	FetchHistorySegment(ctx context.Context, historySegmentID string) (networkhistory.Segment, error)
	GetActivePeerIPAddresses() []string
	CopyHistorySegmentToFile(ctx context.Context, historySegmentID string, outFile string) error
	GetSwarmKeySeed() string
	GetConnectedPeerAddresses() ([]string, error)
	GetIpfsAddress() (string, error)
	GetSwarmKey() string
	GetBootstrapPeers() []string
}

// GRPCServer represent the grpc api provided by the zeta node.
type GRPCServer struct {
	Config
	log                   *logging.Logger
	srv                   *grpc.Server
	zetaCoreServiceClient CoreServiceClient

	eventService               *subscribers.Service
	coreProxySvc               *coreProxyService
	orderService               *service.Order
	candleService              *candlesv2.Svc
	networkLimitsService       *service.NetworkLimits
	marketDataService          *service.MarketData
	tradeService               *service.Trade
	assetService               *service.Asset
	accountService             *service.Account
	rewardService              *service.Reward
	marketsService             *service.Markets
	delegationService          *service.Delegation
	epochService               *service.Epoch
	depositService             *service.Deposit
	withdrawalService          *service.Withdrawal
	governanceService          *service.Governance
	riskFactorService          *service.RiskFactor
	riskService                *service.Risk
	networkParameterService    *service.NetworkParameter
	blockService               BlockService
	partyService               *service.Party
	checkpointService          *service.Checkpoint
	oracleSpecService          *service.OracleSpec
	oracleDataService          *service.OracleData
	liquidityProvisionService  *service.LiquidityProvision
	positionService            *service.Position
	transferService            *service.Transfer
	stakeLinkingService        *service.StakeLinking
	notaryService              *service.Notary
	multiSigService            *service.MultiSig
	keyRotationService         *service.KeyRotations
	ethereumKeyRotationService *service.EthereumKeyRotation
	nodeService                *service.Node
	marketDepthService         *service.MarketDepth
	ledgerService              *service.Ledger
	protocolUpgradeService     *service.ProtocolUpgrade
	networkHistoryService      NetworkHistoryService
	coreSnapshotService        *service.SnapshotData

	eventObserver *eventObserver

	// used in order to gracefully close streams
	ctx   context.Context
	cfunc context.CancelFunc
}

// NewGRPCServer create a new instance of the GPRC api for the zeta node.
func NewGRPCServer(
	log *logging.Logger,
	config Config,
	coreServiceClient CoreServiceClient,
	eventService *subscribers.Service,
	orderService *service.Order,
	networkLimitsService *service.NetworkLimits,
	marketDataService *service.MarketData,
	tradeService *service.Trade,
	assetService *service.Asset,
	accountService *service.Account,
	rewardService *service.Reward,
	marketsService *service.Markets,
	delegationService *service.Delegation,
	epochService *service.Epoch,
	depositService *service.Deposit,
	withdrawalService *service.Withdrawal,
	governanceService *service.Governance,
	riskFactorService *service.RiskFactor,
	riskService *service.Risk,
	networkParameterService *service.NetworkParameter,
	blockService BlockService,
	checkpointService *service.Checkpoint,
	partyService *service.Party,
	candleService *candlesv2.Svc,
	oracleSpecService *service.OracleSpec,
	oracleDataService *service.OracleData,
	liquidityProvisionService *service.LiquidityProvision,
	positionService *service.Position,
	transferService *service.Transfer,
	stakeLinkingService *service.StakeLinking,
	notaryService *service.Notary,
	multiSigService *service.MultiSig,
	keyRotationService *service.KeyRotations,
	ethereumKeyRotationService *service.EthereumKeyRotation,
	nodeService *service.Node,
	marketDepthService *service.MarketDepth,
	ledgerService *service.Ledger,
	protocolUpgradeService *service.ProtocolUpgrade,
	networkHistoryService NetworkHistoryService,
	coreSnapshotService *service.SnapshotData,
) *GRPCServer {
	// setup logger
	log = log.Named(namedLogger)
	log.SetLevel(config.Level.Get())
	ctx, cfunc := context.WithCancel(context.Background())

	return &GRPCServer{
		log:                        log,
		Config:                     config,
		zetaCoreServiceClient:      coreServiceClient,
		eventService:               eventService,
		orderService:               orderService,
		networkLimitsService:       networkLimitsService,
		tradeService:               tradeService,
		assetService:               assetService,
		accountService:             accountService,
		rewardService:              rewardService,
		marketsService:             marketsService,
		delegationService:          delegationService,
		epochService:               epochService,
		depositService:             depositService,
		withdrawalService:          withdrawalService,
		multiSigService:            multiSigService,
		governanceService:          governanceService,
		riskFactorService:          riskFactorService,
		networkParameterService:    networkParameterService,
		blockService:               blockService,
		checkpointService:          checkpointService,
		partyService:               partyService,
		candleService:              candleService,
		oracleSpecService:          oracleSpecService,
		oracleDataService:          oracleDataService,
		liquidityProvisionService:  liquidityProvisionService,
		positionService:            positionService,
		transferService:            transferService,
		stakeLinkingService:        stakeLinkingService,
		notaryService:              notaryService,
		keyRotationService:         keyRotationService,
		ethereumKeyRotationService: ethereumKeyRotationService,
		nodeService:                nodeService,
		marketDepthService:         marketDepthService,
		riskService:                riskService,
		marketDataService:          marketDataService,
		ledgerService:              ledgerService,
		protocolUpgradeService:     protocolUpgradeService,
		networkHistoryService:      networkHistoryService,
		coreSnapshotService:        coreSnapshotService,

		eventObserver: &eventObserver{
			log:          log,
			eventService: eventService,
			Config:       config,
		},
		ctx:   ctx,
		cfunc: cfunc,
	}
}

// ReloadConf update the internal configuration of the GRPC server.
func (g *GRPCServer) ReloadConf(cfg Config) {
	g.log.Info("reloading configuration")
	if g.log.GetLevel() != cfg.Level.Get() {
		g.log.Info("updating log level",
			logging.String("old", g.log.GetLevel().String()),
			logging.String("new", cfg.Level.String()),
		)
		g.log.SetLevel(cfg.Level.Get())
	}

	// TODO(): not updating the actual server for now, may need to look at this later
	// e.g restart the http server on another port or whatever
	g.Config = cfg
}

func remoteAddrInterceptor(log *logging.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// first check if the request is forwarded from our restproxy
		// get the metadata
		var ip string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			forwardedFor, ok := md["x-forwarded-for"]
			if ok && len(forwardedFor) > 0 {
				log.Debug("grpc request x-forwarded-for",
					logging.String("method", info.FullMethod),
					logging.String("remote-ip-addr", forwardedFor[0]),
				)
				ip = forwardedFor[0]
			}
		}

		// if the request is not forwarded let's get it from the peer infos
		if len(ip) <= 0 {
			p, ok := peer.FromContext(ctx)
			if ok && p != nil {
				log.Debug("grpc peer client request",
					logging.String("method", info.FullMethod),
					logging.String("remote-ip-addr", p.Addr.String()))
				ip = p.Addr.String()
			}
		}

		ctx = contextutil.WithRemoteIPAddr(ctx, ip)

		// Calls the handler
		h, err := handler(ctx, req)

		log.Debug("Invoked RPC call",
			logging.String("method", info.FullMethod),
			logging.Error(err),
		)

		return h, err
	}
}

func headersInterceptor(
	getLastBlock func(context.Context) (entities.Block, error),
	log *logging.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var (
			height    int64
			timestamp int64
		)

		block, bErr := getLastBlock(ctx)
		if bErr != nil {
			log.Debug("failed to get last block", logging.Error(bErr))
		} else {
			height = block.Height
			timestamp = block.ZetaTime.UnixNano()
		}

		for _, h := range []metadata.MD{
			// Deprecated: use 'X-Block-Height' and 'X-Block-Timestamp' instead to determine if data is fresh.
			metadata.Pairs("X-Zeta-Connection", "CONNECTED"),
			metadata.Pairs("X-Block-Height", strconv.FormatInt(height, 10)),
			metadata.Pairs("X-Block-Timestamp", strconv.FormatInt(timestamp, 10)),
			// TODO: remove warning once deprecated header is gone.
			metadata.Pairs("Warning", "199 - \"The header 'X-Zeta-Connection' is deprecated and now defaults to 'CONNECTED'. It will be removed in a future version. See https://github.com/zetaprotocol/zeta/issues/7385#issuecomment-1398719810\""),
		} {
			if errH := grpc.SetHeader(ctx, h); errH != nil {
				log.Error("failed to set header", logging.Error(errH))
			}
		}

		return handler(ctx, req)
	}
}

func (g *GRPCServer) getTCPListener() (net.Listener, error) {
	ip := g.IP
	port := strconv.Itoa(g.Port)

	g.log.Info("Starting gRPC based API", logging.String("addr", ip), logging.String("port", port))

	tpcLis, err := net.Listen("tcp", net.JoinHostPort(ip, port))
	if err != nil {
		return nil, err
	}

	return tpcLis, nil
}

// Start starts the grpc server.
// Uses default TCP listener if no provided.
func (g *GRPCServer) Start(ctx context.Context, lis net.Listener) error {
	if lis == nil {
		tpcLis, err := g.getTCPListener()
		if err != nil {
			return err
		}

		lis = tpcLis
	}

	rateLimit := ratelimit.NewFromConfig(&g.RateLimit, g.log)
	intercept := grpc.ChainUnaryInterceptor(
		remoteAddrInterceptor(g.log),
		headersInterceptor(g.blockService.GetLastBlock, g.log),
		rateLimit.GRPCInterceptor,
	)

	g.srv = grpc.NewServer(intercept)

	coreProxySvc := &coreProxyService{
		conf:              g.Config,
		coreServiceClient: g.zetaCoreServiceClient,
		eventObserver:     g.eventObserver,
	}
	g.coreProxySvc = coreProxySvc
	zetaprotoapi.RegisterCoreServiceServer(g.srv, coreProxySvc)

	tradingDataSvcV2 := &tradingDataServiceV2{
		config:                     g.Config,
		log:                        g.log,
		orderService:               g.orderService,
		networkLimitsService:       g.networkLimitsService,
		marketDataService:          g.marketDataService,
		tradeService:               g.tradeService,
		multiSigService:            g.multiSigService,
		notaryService:              g.notaryService,
		assetService:               g.assetService,
		candleService:              g.candleService,
		marketsService:             g.marketsService,
		partyService:               g.partyService,
		riskService:                g.riskService,
		positionService:            g.positionService,
		accountService:             g.accountService,
		rewardService:              g.rewardService,
		depositService:             g.depositService,
		withdrawalService:          g.withdrawalService,
		oracleSpecService:          g.oracleSpecService,
		oracleDataService:          g.oracleDataService,
		liquidityProvisionService:  g.liquidityProvisionService,
		governanceService:          g.governanceService,
		transfersService:           g.transferService,
		delegationService:          g.delegationService,
		marketService:              g.marketsService,
		marketDepthService:         g.marketDepthService,
		nodeService:                g.nodeService,
		epochService:               g.epochService,
		riskFactorService:          g.riskFactorService,
		networkParameterService:    g.networkParameterService,
		checkpointService:          g.checkpointService,
		stakeLinkingService:        g.stakeLinkingService,
		eventService:               g.eventService,
		ledgerService:              g.ledgerService,
		keyRotationService:         g.keyRotationService,
		ethereumKeyRotationService: g.ethereumKeyRotationService,
		blockService:               g.blockService,
		protocolUpgradeService:     g.protocolUpgradeService,
		networkHistoryService:      g.networkHistoryService,
		coreSnapshotService:        g.coreSnapshotService,
	}

	protoapi.RegisterTradingDataServiceServer(g.srv, tradingDataSvcV2)

	eg, ctx := errgroup.WithContext(ctx)

	if g.Reflection || g.WebUIEnabled {
		reflection.Register(g.srv)
	}

	eg.Go(func() error {
		<-ctx.Done()
		g.stop()
		return ctx.Err()
	})

	eg.Go(func() error {
		return g.srv.Serve(lis)
	})

	if g.WebUIEnabled {
		g.startWebUI(ctx)
	}

	return eg.Wait()
}

func (g *GRPCServer) stop() {
	if g.srv == nil {
		return
	}

	done := make(chan struct{})
	go func() {
		g.log.Info("Gracefully stopping gRPC based API")
		g.srv.GracefulStop()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		g.log.Info("Force stopping gRPC based API")
		g.srv.Stop()
	}
}

func (g *GRPCServer) startWebUI(ctx context.Context) {
	cc, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", g.Port), grpc.WithInsecure())
	if err != nil {
		g.log.Error("failed to create client to local grpc server", logging.Error(err))
		return
	}

	uiHandler, err := standalone.HandlerViaReflection(ctx, cc, "zeta data node")
	if err != nil {
		g.log.Error("failed to create grpc-ui server", logging.Error(err))
		return
	}

	uiListener, err := net.Listen("tcp", net.JoinHostPort(g.IP, strconv.Itoa(g.WebUIPort)))
	if err != nil {
		g.log.Error("failed to open listen socket on port", logging.Int("port", g.WebUIPort), logging.Error(err))
		return
	}

	g.log.Info("Starting gRPC Web UI", logging.String("addr", g.IP), logging.Int("port", g.WebUIPort))
	go http.Serve(uiListener, uiHandler)
}
