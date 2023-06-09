// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package visor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"zuluprotocol/zeta/core/types"
	"zuluprotocol/zeta/logging"
	"zuluprotocol/zeta/visor/config"
	"zuluprotocol/zeta/visor/utils"

	"golang.org/x/sync/errgroup"
)

const snapshotBlockHeightFlagName = "--snapshot.load-from-block-height"

type BinariesRunner struct {
	mut         sync.RWMutex
	running     map[int]*exec.Cmd
	binsFolder  string
	log         *logging.Logger
	stopTimeout time.Duration
	releaseInfo *types.ReleaseInfo
}

func NewBinariesRunner(log *logging.Logger, binsFolder string, stopTimeout time.Duration, rInfo *types.ReleaseInfo) *BinariesRunner {
	return &BinariesRunner{
		binsFolder:  binsFolder,
		running:     map[int]*exec.Cmd{},
		log:         log,
		stopTimeout: stopTimeout,
		releaseInfo: rInfo,
	}
}

func (r *BinariesRunner) cleanBinaryPath(binPath string) string {
	if !filepath.IsAbs(binPath) {
		return path.Join(r.binsFolder, binPath)
	}

	return binPath
}

func (r *BinariesRunner) runBinary(ctx context.Context, binPath string, args []string) error {
	binPath = r.cleanBinaryPath(binPath)

	if err := utils.EnsureBinary(binPath); err != nil {
		return fmt.Errorf("failed to locate binary %s %v: %w", binPath, args, err)
	}

	if r.releaseInfo != nil {
		if err := ensureBinaryVersion(binPath, r.releaseInfo.ZetaReleaseTag); err != nil {
			return err
		}
	}

	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	r.log.Debug("Starting binary",
		logging.String("binaryPath", binPath),
		logging.Strings("args", args),
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to execute binary %s %v: %w", binPath, args, err)
	}

	// Ensures that if one binary failes all of them are killed
	go func() {
		<-ctx.Done()

		if cmd.Process == nil {
			return
		}

		// Process has already exited - no need to kill it
		if cmd.ProcessState != nil {
			return
		}

		r.log.Debug("Killing binary", logging.String("binaryPath", binPath))

		if err := cmd.Process.Kill(); err != nil {
			r.log.Debug("Failed to kill binary",
				logging.String("binaryPath", binPath),
				logging.Error(err),
			)
		}
	}()

	processID := cmd.Process.Pid

	r.mut.Lock()
	r.running[processID] = cmd
	r.mut.Unlock()

	defer func() {
		r.mut.Lock()
		delete(r.running, processID)
		r.mut.Unlock()
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to execute binary %s %v: %w", binPath, args, err)
	}

	return nil
}

func (r *BinariesRunner) prepareZetaArgs(runConf *config.RunConfig, isRestart bool) (Args, error) {
	args := Args(runConf.Zeta.Binary.Args)

	// if a node restart happens (not due protocol upgrade) and data node is present
	// we need to make sure that they will start on the block that data node has already processed.
	if isRestart && runConf.DataNode != nil {
		latestSegment, err := latestDataNodeHistorySegment(
			r.cleanBinaryPath(runConf.DataNode.Binary.Path),
			runConf.DataNode.Binary.Args,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest history segment from data node: %w", err)
		}

		args.Set(snapshotBlockHeightFlagName, strconv.FormatUint(uint64(latestSegment.LatestSegment.Height), 10))
		return args, nil
	}

	if r.releaseInfo != nil {
		args.Set(snapshotBlockHeightFlagName, strconv.FormatUint(r.releaseInfo.UpgradeBlockHeight, 10))
	}

	return args, nil
}

func (r *BinariesRunner) Run(ctx context.Context, runConf *config.RunConfig, isRestart bool) chan error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		args, err := r.prepareZetaArgs(runConf, isRestart)
		if err != nil {
			return fmt.Errorf("failed to prepare args for Zeta binary: %w", err)
		}

		return r.runBinary(ctx, runConf.Zeta.Binary.Path, args)
	})

	if runConf.DataNode != nil {
		eg.Go(func() error {
			return r.runBinary(ctx, runConf.DataNode.Binary.Path, runConf.DataNode.Binary.Args)
		})
	}

	errChan := make(chan error)

	go func() {
		err := eg.Wait()
		if err != nil {
			errChan <- err
		}
	}()

	return errChan
}

func (r *BinariesRunner) signal(signal syscall.Signal) error {
	r.mut.RLock()
	defer r.mut.RUnlock()

	var err error
	for _, c := range r.running {
		r.log.Info("Signaling process",
			logging.String("binaryName", c.Path),
			logging.String("signal", signal.String()),
			logging.Strings("args", c.Args),
		)

		err = c.Process.Signal(signal)
		if err != nil {
			r.log.Error("Failed to signal running binary",
				logging.String("binaryPath", c.Path),
				logging.Strings("args", c.Args),
				logging.Error(err),
			)
		}
	}

	return err
}

func (r *BinariesRunner) Stop() error {
	if err := r.signal(syscall.SIGTERM); err != nil {
		return err
	}

	r.mut.RLock()
	timeout := time.After(r.stopTimeout)
	r.mut.RUnlock()

	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("failed to gracefully shut down processes: timed out")
		case <-ticker.C:
			r.mut.RLock()
			if len(r.running) == 0 {
				r.mut.RUnlock()
				return nil
			}
			r.mut.RUnlock()
		}
	}
}

func (r *BinariesRunner) Kill() error {
	return r.signal(syscall.SIGKILL)
}
