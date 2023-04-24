package client

import (
	"context"

	"zuluprotocol/zeta/zeta/core/admin"
	"zuluprotocol/zeta/zeta/core/types"
	"zuluprotocol/zeta/zeta/logging"
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
