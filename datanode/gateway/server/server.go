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

package server

import (
	"context"

	"golang.org/x/sync/errgroup"

	"zuluprotocol/zeta/datanode/gateway"
	gql "zuluprotocol/zeta/datanode/gateway/graphql"
	"zuluprotocol/zeta/datanode/gateway/rest"
	"zuluprotocol/zeta/logging"
	"zuluprotocol/zeta/paths"
)

type Server struct {
	cfg       *gateway.Config
	log       *logging.Logger
	zetaPaths paths.Paths

	rest *rest.ProxyServer
	gql  *gql.GraphServer
}

const namedLogger = "gateway"

func New(cfg gateway.Config, log *logging.Logger, zetaPaths paths.Paths) *Server {
	log = log.Named(namedLogger)
	log.SetLevel(cfg.Level.Get())

	return &Server{
		log:       log,
		cfg:       &cfg,
		zetaPaths: zetaPaths,
	}
}

func (srv *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	if srv.cfg.GraphQL.Enabled {
		var err error
		srv.gql, err = gql.New(srv.log, *srv.cfg, srv.zetaPaths)
		if err != nil {
			return err
		}
		eg.Go(func() error { return srv.gql.Start() })
	}

	if srv.cfg.REST.Enabled {
		srv.rest = rest.NewProxyServer(srv.log, *srv.cfg, srv.zetaPaths)
		eg.Go(func() error { return srv.rest.Start() })
	}

	if srv.cfg.REST.Enabled || srv.cfg.GraphQL.Enabled {
		eg.Go(func() error {
			<-ctx.Done()
			srv.stop()
			return nil
		})
	}

	return eg.Wait()
}

func (srv *Server) stop() {
	if s := srv.rest; s != nil {
		s.Stop()
	}

	if s := srv.gql; s != nil {
		s.Stop()
	}
}
