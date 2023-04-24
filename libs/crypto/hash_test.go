package crypto_test

import (
	"testing"

	vgcrypto "code.zetaprotocol.io/vega/libs/crypto"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	t.Run("Hashing data succeeds", testHashingDataSucceeds)
}

func testHashingDataSucceeds(t *testing.T) {
	data := []byte("Hello, World!")
	hashedData := vgcrypto.Hash(data)
	assert.Equal(t, []byte{0x1a, 0xf1, 0x7a, 0x66, 0x4e, 0x3f, 0xa8, 0xe4, 0x19, 0xb8, 0xba, 0x5, 0xc2, 0xa1, 0x73, 0x16, 0x9d, 0xf7, 0x61, 0x62, 0xa5, 0xa2, 0x86, 0xe0, 0xc4, 0x5, 0xb4, 0x60, 0xd4, 0x78, 0xf7, 0xef}, hashedData)
}
