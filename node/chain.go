package node

import (
	"encoding/hex"

	"github.com/ranjankuldeep/DisBlocker/proto"
	"github.com/ranjankuldeep/DisBlocker/store"
)

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}

func (h HeaderList) GetLen() int {
	return len(h.headers)
}

func (h HeaderList) GetHeight() int {
	return h.GetLen() - 1
}

func (h *HeaderList) AddHeader(header *proto.Header) {
	h.headers = append(h.headers, header)
}

type Chain struct {
	headers     *HeaderList
	blockStorer store.BlockStorer
}

func NewChain(bs store.BlockStorer) *Chain {
	return &Chain{
		blockStorer: bs,
		headers:     NewHeaderList(),
	}
}

func (c *Chain) AddBlock(block *proto.Block) error {
	// Need to add the header to the header list
	c.headers.AddHeader(block.Header)
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

func (c *Chain) Height() int {
	return c.headers.GetHeight()
}
