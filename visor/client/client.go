package client

import (
	"context"

	"code.zetaprotocol.io/vega/core/admin"
	"code.zetaprotocol.io/vega/core/types"
	"code.zetaprotocol.io/vega/logging"
)

type AdminClient interface {
	UpgradeStatus(ctx context.Context) (*types.UpgradeStatus, error)
}

type Factory interface {
	GetClient(socketPath, httpPath string) AdminClient
}

type clientFactory struct {
	log *logging.Logger
}

func NewClientFactory(log *logging.Logger) Factory {
	return &clientFactory{
		log: log,
	}
}

func (cf *clientFactory) GetClient(socketPath, httpPath string) AdminClient {
	return admin.NewClient(cf.log, admin.Config{
		Server: admin.ServerConfig{
			SocketPath: socketPath,
			HTTPPath:   httpPath,
		},
	})
}
