package networkhistory

import (
	"context"
	"fmt"
	"os"

	coreConfig "zuluprotocol/zeta/zeta/core/config"
	"zuluprotocol/zeta/zeta/datanode/networkhistory"
	vgjson "zuluprotocol/zeta/zeta/libs/json"

	"zuluprotocol/zeta/zeta/datanode/config"
	"zuluprotocol/zeta/zeta/logging"
	"zuluprotocol/zeta/zeta/paths"
	v2 "zuluprotocol/zeta/zeta/protos/data-node/api/v2"
)

type latestHistorySegmentFromPeers struct {
	config.ZetaHomeFlag
	coreConfig.OutputFlag
	config.Config
}

type segmentInfo struct {
	Peer         string
	SwarmKeySeed string
	Segment      *v2.HistorySegment
}

type latestHistorFromPeersyOutput struct {
	Segments              []segmentInfo
	SuggestedFetchSegment *v2.HistorySegment
}

func (o *latestHistorFromPeersyOutput) printHuman() {
	segmentsInfo := "Most Recent History Segments:\n\n"
	for _, segment := range o.Segments {
		segmentsInfo += fmt.Sprintf("Peer:%-39s,  Swarm Key:%s, Segment{%s}\n\n", segment.Peer, segment.SwarmKeySeed, segment.Segment)
	}
	fmt.Println(segmentsInfo)
	fmt.Printf("Suggested segment to use to fetch network history data {%s}\n\n", o.SuggestedFetchSegment)
}

func (cmd *latestHistorySegmentFromPeers) Execute(_ []string) error {
	cfg := logging.NewDefaultConfig()
	cfg.Custom.Zap.Level = logging.InfoLevel
	cfg.Environment = "custom"
	log := logging.NewLoggerFromConfig(
		cfg,
	)
	defer log.AtExit()

	zetaPaths := paths.New(cmd.ZetaHome)
	err := fixConfig(&cmd.Config, zetaPaths)
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to fix config", err)
		os.Exit(1)
	}

	if !datanodeLive(cmd.Config) {
		return fmt.Errorf("datanode must be running for this command to work")
	}

	client, conn, err := getDatanodeClient(cmd.Config)
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to get datanode client", err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	resp, err := client.GetActiveNetworkHistoryPeerAddresses(context.Background(), &v2.GetActiveNetworkHistoryPeerAddressesRequest{})
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to get active peer addresses", errorFromGrpcError("", err))
		os.Exit(1)
	}

	peerAddresses := resp.IpAddresses

	grpcAPIPorts := []int{cmd.Config.API.Port}
	grpcAPIPorts = append(grpcAPIPorts, cmd.Config.NetworkHistory.Initialise.GrpcAPIPorts...)
	selectedResponse, peerToResponse, err := networkhistory.GetMostRecentHistorySegmentFromPeersAddresses(context.Background(), peerAddresses,
		cmd.Config.NetworkHistory.Store.GetSwarmKeySeed(log, cmd.Config.ChainID), grpcAPIPorts)
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to get most recent history segment from peers", err)
		os.Exit(1)
	}

	output := latestHistorFromPeersyOutput{}
	output.Segments = []segmentInfo{}

	for peer, segment := range peerToResponse {
		output.Segments = append(output.Segments, segmentInfo{
			Peer:         peer,
			SwarmKeySeed: segment.SwarmKeySeed,
			Segment:      segment.Segment,
		})
	}

	if selectedResponse != nil {
		output.SuggestedFetchSegment = selectedResponse.Response.Segment
	}

	if cmd.Output.IsJSON() {
		if err := vgjson.Print(&output); err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to marshal output", err)
			os.Exit(1)
		}
	} else {
		output.printHuman()
	}

	return nil
}
