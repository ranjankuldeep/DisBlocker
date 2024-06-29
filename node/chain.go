package node

import (
	"encoding/hex"

	"github.com/ranjankuldeep/DisBlocker/proto"
	"github.com/ranjankuldeep/DisBlocker/store"
)

type Chain struct {
	blockStorer store.BlockStorer
}

func NewChain(bs store.BlockStorer) *Chain {
	return &Chain{
		blockStorer: bs,
	}
}

func (c *Chain) AddBlock(block *proto.Block) error {
	// Do the validation here before adding block to the chain
	return c.blockStorer.Put(block)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	// Implement the logic here.
	return nil, nil
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStorer.Get(hashHex)
}
