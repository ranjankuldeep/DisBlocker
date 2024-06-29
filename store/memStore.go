package store

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ranjankuldeep/DisBlocker/proto"
	"github.com/ranjankuldeep/DisBlocker/types.go"
)

type MemoryBlockStore struct {
	lock   sync.RWMutex
	blocks map[string]*proto.Block
}

func NewMemoryBlockStore() *MemoryBlockStore {
	return &MemoryBlockStore{
		lock:   sync.RWMutex{},
		blocks: make(map[string]*proto.Block),
	}
}

func (m *MemoryBlockStore) Put(block *proto.Block) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	hash := types.HashBlock(block)
	hexHash := hex.EncodeToString(hash)

	m.blocks[hexHash] = block
	return nil
}

func (m *MemoryBlockStore) Get(hash string) (*proto.Block, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	block, ok := m.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash [%s] does not exist", hash)
	}
	return block, nil
}
