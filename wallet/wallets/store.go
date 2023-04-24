package wallets

import (
	"fmt"

	"zuluprotocol/zeta/zeta/paths"
	wstorev1 "zuluprotocol/zeta/zeta/wallet/wallet/store/v1"
)

// InitialiseStore builds a wallet Store specifically for users wallets.
func InitialiseStore(zetaHome string) (*wstorev1.FileStore, error) {
	p := paths.New(zetaHome)
	return InitialiseStoreFromPaths(p)
}

// InitialiseStoreFromPaths builds a wallet Store specifically for users wallets.
func InitialiseStoreFromPaths(zetaPaths paths.Paths) (*wstorev1.FileStore, error) {
	walletsHome, err := zetaPaths.CreateDataPathFor(paths.WalletsDataHome)
	if err != nil {
		return nil, fmt.Errorf("couldn't get wallets data home path: %w", err)
	}
	return wstorev1.InitialiseStore(walletsHome)
}
