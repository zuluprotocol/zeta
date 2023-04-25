package admin

import (
	"os"
	"path"

	"zuluprotocol/zeta/libs/config/encoding"
	"zuluprotocol/zeta/logging"
)

const namedLogger = "admin.server"

// ServerConfig contains the configuration for the admin server.
type ServerConfig struct {
	SocketPath string `long:"socket-path" description:"Listen for connection on UNIX socket path <file-path>"`
	HTTPPath   string `long:"http-path" description:"Http path of the socket HTTP RPC server"`
}

// Config represents the configuration of the admin package.
type Config struct {
	Level  encoding.LogLevel `long:"log-level"`
	Server ServerConfig      `group:"Server" namespace:"server"`
}

// NewDefaultConfig creates an instance of the package specific configuration.
func NewDefaultConfig() Config {
	return Config{
		Level: encoding.LogLevel{Level: logging.InfoLevel},
		Server: ServerConfig{
			SocketPath: path.Join(os.TempDir(), "datanode.sock"),
			HTTPPath:   "/datanode/rpc",
		},
	}
}
