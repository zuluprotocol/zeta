package networkhistory

import (
	"context"
	"errors"
	"fmt"

	"code.zetaprotocol.io/vega/datanode/config"
	"code.zetaprotocol.io/vega/logging"
)

type copyCmd struct {
	config.ZetaHomeFlag
	config.Config
}

func (cmd *copyCmd) Execute(args []string) error {
	cfg := logging.NewDefaultConfig()
	cfg.Custom.Zap.Level = logging.InfoLevel
	cfg.Environment = "custom"
	log := logging.NewLoggerFromConfig(
		cfg,
	)
	defer log.AtExit()

	if len(args) != 2 {
		return errors.New("expected <history-segment-id> <output-file>")
	}

	segmentID := args[0]
	outputPath := args[1]

	client := getDatanodeAdminClient(log, cmd.Config)
	reply, err := client.CopyHistorySegmentToFile(context.Background(), segmentID, outputPath)
	if err != nil {
		return fmt.Errorf("failed to copy history segment to file: %w", err)
	}

	if reply.Err != nil {
		return reply.Err
	}

	log.Info(reply.Reply)

	return nil
}
