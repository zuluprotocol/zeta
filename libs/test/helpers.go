package test

import (
	"path/filepath"

	vgrand "code.zetaprotocol.io/vega/libs/rand"
)

func RandomPath() string {
	return filepath.Join("/tmp", "zeta_tests", vgrand.RandomStr(10))
}
