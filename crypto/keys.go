package crypto

import "crypto/ed25519"

type PrivateKey struct {
	key ed25519.PrivateKey
}

func (pk *PrivateKey) Bytes() []byte {
	return pk.key
}

func (pk *PrivateKey) Sign(msg []byte) []byte {
	return ed25519.Sign(pk.key, msg)
}
