package networkhistory

import (
	"context"
	"fmt"
	"os"
	"sort"

	"zuluprotocol/zeta/datanode/networkhistory"

	"github.com/jackc/pgx/v4/pgxpool"

	coreConfig "zuluprotocol/zeta/core/config"
	"zuluprotocol/zeta/datanode/sqlstore"
	vgjson "zuluprotocol/zeta/libs/json"
	"zuluprotocol/zeta/paths"

	"zuluprotocol/zeta/logging"
	v2 "zuluprotocol/zeta/protos/data-node/api/v2"

	"zuluprotocol/zeta/datanode/config"
)

type showCmd struct {
	config.ZetaHomeFlag
	config.Config
	coreConfig.OutputFlag

	AllSegments bool `short:"s" long:"segments" description:"show all segments for each contiguous history"`
}

type showOutput struct {
	Segments            []*v2.HistorySegment
	ContiguousHistories []*networkhistory.ContiguousHistory
	DataNodeBlockStart  int64
	DataNodeBlockEnd    int64
}

func (o *showOutput) printHuman(allSegments bool) {
	if len(o.ContiguousHistories) > 0 {
		fmt.Printf("Available contiguous history spans:")
		for _, contiguousHistory := range o.ContiguousHistories {
			fmt.Printf("\n\nContiguous history from block height %d to %d, from segment id: %s to %s\n",
				contiguousHistory.HeightFrom,
				contiguousHistory.HeightTo,
				contiguousHistory.SegmentsOldestFirst[0].GetHistorySegmentId(),
				contiguousHistory.SegmentsOldestFirst[len(contiguousHistory.SegmentsOldestFirst)-1].GetHistorySegmentId(),
			)

			if allSegments {
				for _, segment := range contiguousHistory.SegmentsOldestFirst {
					fmt.Printf("\n%d to %d, id: %s, previous segment id: %s",
						segment.GetFromHeight(),
						segment.GetToHeight(),
						segment.GetHistorySegmentId(),
						segment.GetPreviousHistorySegmentId())
				}
			}
		}
	} else {
		fmt.Printf("\nNo network history is available.  Use the fetch command to fetch network history\n")
	}

	if o.DataNodeBlockEnd > 0 {
		fmt.Printf("\n\nDatanode currently has data from block height %d to %d\n", o.DataNodeBlockStart, o.DataNodeBlockEnd)
	} else {
		fmt.Printf("\n\nDatanode contains no data\n")
	}
}

func (cmd *showCmd) Execute(_ []string) error {
	cfg := logging.NewDefaultConfig()
	cfg.Custom.Zap.Level = logging.WarnLevel
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
		handleErr(log,
			cmd.Output.IsJSON(),
			"datanode must be running for this command to work",
			fmt.Errorf("couldn't connect to datanode on %v:%v", cmd.Config.API.IP, cmd.Config.API.Port))
		os.Exit(1)
	}

	client, conn, err := getDatanodeClient(cmd.Config)
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to get datanode client", err)
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	response, err := client.ListAllNetworkHistorySegments(context.Background(), &v2.ListAllNetworkHistorySegmentsRequest{})
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to list all network history segments", err)
		os.Exit(1)
	}

	output := showOutput{}
	output.Segments = response.Segments

	sort.Slice(output.Segments, func(i int, j int) bool {
		return output.Segments[i].ToHeight < output.Segments[j].ToHeight
	})

	output.ContiguousHistories = networkhistory.GetContiguousHistories(response.Segments)

	pool, err := getCommandConnPool(cmd.Config.SQLStore.ConnectionConfig)
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to get command conn pool", err)
	}
	defer pool.Close()

	span, err := sqlstore.GetDatanodeBlockSpan(context.Background(), pool)
	if err != nil {
		handleErr(log, cmd.Output.IsJSON(), "failed to get datanode block span", err)
		os.Exit(1)
	}

	if span.HasData {
		output.DataNodeBlockStart = span.FromHeight
		output.DataNodeBlockEnd = span.ToHeight
	}

	if cmd.Output.IsJSON() {
		if err := vgjson.Print(&output); err != nil {
			handleErr(log, cmd.Output.IsJSON(), "failed to marshal output", err)
			os.Exit(1)
		}
	} else {
		output.printHuman(cmd.AllSegments)
	}

	return nil
}

func getCommandConnPool(conf sqlstore.ConnectionConfig) (*pgxpool.Pool, error) {
	conf.MaxConnPoolSize = 3

	connPool, err := sqlstore.CreateConnectionPool(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return connPool, nil
}
