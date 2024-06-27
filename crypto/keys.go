package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/ranjankuldeep/DisBlocker/logs"
)

const (
	privKeyLen = 64
	pubKeyLen  = 32
	bufLen     = 32
	sigLen     = 64
	addressLen = 20
	seedLen    = 32
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return NewPrivateKeyFromSeed(b)
}
func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != seedLen {
		logs.Logger.Error("Invalid Seed length")
		panic("Invalid Seed Length")
	}
	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
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

func SignatureFromBytes(b []byte) *Signature {
	if len(b) != sigLen {
		panic("Invalid Siganture")
	}
	return &Signature{
		value: b,
	}
}
func (pk *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(pk.key, msg),
	}
}

func (pk *PrivateKey) Public() *PublicKey {
	buf := make([]byte, pubKeyLen)
	copy(buf, pk.key[32:])
	return &PublicKey{
		key: buf,
	}
}

func NewPublicKeyFromBytes(b []byte) *PublicKey {
	if len(b) != pubKeyLen {
		panic("Invalid Public Key length")
	}
	return &PublicKey{
		key: b,
	}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func (pk *PublicKey) Bytes() []byte {
	return pk.key
}

func (pk *PublicKey) Address() Address {
	return Address{
		value: pk.key[len(pk.key)-addressLen:],
	}
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

type Address struct {
	value []byte
}

func (a Address) Bytes() []byte {
	return a.value
}
func (a Address) String() string {
	return hex.EncodeToString(a.value)
}
