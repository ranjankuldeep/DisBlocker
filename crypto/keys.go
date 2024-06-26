package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"

	"github.com/ranjankuldeep/DisBlocker/logs"
)

const (
	privKeyLen = 64
	pubKeyLen  = 32
	bufLen     = 32
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func GeneratePrivateKey() *PrivateKey {
	seed := make([]byte, bufLen)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		logs.Logger.Errorf("Error Generating Private Key: %+v", err)
		panic(err)
	}
	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func (pk *PrivateKey) Bytes() []byte {
	return pk.key
}

func (pk *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		Value: ed25519.Sign(pk.key, msg),
	}
}

func (pk *PrivateKey) Public() *PublicKey {
	buf := make([]byte, pubKeyLen)
	copy(buf, pk.key[32:])
	return &PublicKey{
		key: buf,
	}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func (pk *PublicKey) Bytes() []byte {
	return pk.key
}

type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}
