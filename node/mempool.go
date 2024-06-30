package node

import (
	"encoding/hex"

	"github.com/ranjankuldeep/DisBlocker/proto"
	"github.com/ranjankuldeep/DisBlocker/types.go"
)

type MemPool interface {
	Has(*proto.Transaction) bool
	Add(*proto.Transaction) bool
}

type InMemPool struct {
	txx map[string]*proto.Transaction
}

func NewInMemPool() *InMemPool {
	return &InMemPool{
		txx: make(map[string]*proto.Transaction),
	}
}

func (pool *InMemPool) Has(tx *proto.Transaction) bool {
	hashHex := hex.EncodeToString(types.HashTransaction(tx))
	_, ok := pool.txx[hashHex]
	return ok
}

func (pool *InMemPool) Add(tx *proto.Transaction) bool {
	if pool.Has(tx) {
		return false
	}
	hashHex := hex.EncodeToString(types.HashTransaction(tx))
	pool.txx[hashHex] = tx
	return true
}
