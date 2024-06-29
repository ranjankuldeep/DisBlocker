package node

import (
	"testing"

	"github.com/ranjankuldeep/DisBlocker/store"
	"github.com/ranjankuldeep/DisBlocker/types.go"
	"github.com/ranjankuldeep/DisBlocker/util"
	"github.com/stretchr/testify/assert"
)

func TestChainHeight(t *testing.T) {
	chain := NewChain(store.NewMemoryBlockStore())
	for i := 0; i < 100; i++ {

		b := util.RandomBlock()
		assert.Nil(t, chain.AddBlock(b))
		assert.Equal(t, chain.Height(), i)
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(store.NewMemoryBlockStore())
	block := util.RandomBlock()
	blockHash := types.HashBlock(block)

	assert.Nil(t, chain.AddBlock(block))
	fetchedBlock, err := chain.GetBlockByHash(blockHash)
	assert.Nil(t, err)
	assert.Equal(t, fetchedBlock, block)
}
