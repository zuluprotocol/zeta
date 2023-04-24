package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	vgcrypto "zuluprotocol/zeta/zeta/libs/crypto"
)

type HDPublicKey struct {
	Idx          uint32     `json:"index"`
	PublicKey    string     `json:"key"`
	KeyName      string     `json:"name"`
	Algorithm    Algorithm  `json:"algorithm"`
	Tainted      bool       `json:"tainted"`
	MetadataList []Metadata `json:"metadata"`
}

func (k *HDPublicKey) Index() uint32 {
	return k.Idx
}

func (k *HDPublicKey) Key() string {
	return k.PublicKey
}

func (k *HDPublicKey) Name() string {
	return k.KeyName
}

func (k *HDPublicKey) IsTainted() bool {
	return k.Tainted
}

func (k *HDPublicKey) Metadata() []Metadata {
	return k.MetadataList
}

func (k *HDPublicKey) AlgorithmVersion() uint32 {
	return k.Algorithm.Version
}

func (k *HDPublicKey) AlgorithmName() string {
	return k.Algorithm.Name
}

func (k *HDPublicKey) Hash() (string, error) {
	decoded, err := hex.DecodeString(k.PublicKey)
	if err != nil {
		return "", fmt.Errorf("couldn't decode public key: %w", err)
	}

	return hex.EncodeToString(vgcrypto.Hash(decoded)), nil
}

func (k *HDPublicKey) MarshalJSON() ([]byte, error) {
	type alias HDPublicKey
	aliasPublicKey := (*alias)(k)
	return json.Marshal(aliasPublicKey)
}

func (k *HDPublicKey) UnmarshalJSON(data []byte) error {
	type alias HDPublicKey
	aliasPublicKey := (*alias)(k)
	return json.Unmarshal(data, aliasPublicKey)
}
