package networkhistory

import (
	"context"
	"fmt"
	"os"

	coreConfig "zuluprotocol/zeta/core/config"
	vgjson "zuluprotocol/zeta/libs/json"
	v2 "zuluprotocol/zeta/protos/data-node/api/v2"

	"zuluprotocol/zeta/datanode/config"
	"zuluprotocol/zeta/datanode/networkhistory/store"
	"zuluprotocol/zeta/logging"
	"zuluprotocol/zeta/paths"
)

var errNoHistorySegmentFound = fmt.Errorf("no history segments found")

type latestHistorySegment struct {
	config.ZetaHomeFlag
	coreConfig.OutputFlag
	config.Config
}

type latestHistoryOutput struct {
	LatestSegment *v2.HistorySegment
}

func (o *latestHistoryOutput) printHuman() {
	fmt.Printf("Latest segment to use data {%s}\n\n", o.LatestSegment)
}

func (cmd *latestHistorySegment) Execute(_ []string) error {
	cfg := logging.NewDefaultConfig()
	cfg.Custom.Zap.Level = logging.ErrorLevel
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

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	var latestSegment *v2.HistorySegment
	if datanodeLive(cmd.Config) {
		client, conn, err := getDatanodeClient(cmd.Config)
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to get datanode client", err)
			os.Exit(1)
		}
		defer func() { _ = conn.Close() }()

		response, err := client.ListAllNetworkHistorySegments(ctx, &v2.ListAllNetworkHistorySegmentsRequest{})
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to list all network history segments", errorFromGrpcError("", err))
			os.Exit(1)
		}

		if len(response.Segments) < 1 {
			handleErr(log, cmd.Output.IsJSON(), errNoHistorySegmentFound.Error(), errNoHistorySegmentFound)
			os.Exit(1)
		}

		latestSegment = response.Segments[0]
	} else {
		networkHistoryStore, err := store.New(ctx, log, cmd.Config.ChainID, cmd.Config.NetworkHistory.Store, zetaPaths.StatePathFor(paths.DataNodeNetworkHistoryHome), false)
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to create network history store", err)
			os.Exit(1)
		}

		segments, err := networkHistoryStore.ListAllIndexEntriesOldestFirst()
		if err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to list all network history segments", err)
			os.Exit(1)
		}

		if len(segments) < 1 {
			handleErr(log, cmd.Output.IsJSON(), errNoHistorySegmentFound.Error(), errNoHistorySegmentFound)
			os.Exit(1)
		}

		latestSegmentIndex := segments[len(segments)-1]

		latestSegment = &v2.HistorySegment{
			FromHeight:               latestSegmentIndex.GetFromHeight(),
			ToHeight:                 latestSegmentIndex.GetToHeight(),
			HistorySegmentId:         latestSegmentIndex.GetHistorySegmentId(),
			PreviousHistorySegmentId: latestSegmentIndex.GetPreviousHistorySegmentId(),
		}
	}

	output := latestHistoryOutput{
		LatestSegment: latestSegment,
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
