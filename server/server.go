package server

import (
	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/p2p"
)

type BlockServerOpts struct {
	Transport p2p.Transport
}
type BlockServer struct {
	BlockServerOpts
}

func NewBlockServer(opts *BlockServerOpts) *BlockServer {
	return &BlockServer{
		BlockServerOpts: *opts,
	}
}

func (b *BlockServer) StartBlockServer() error {
	err := b.Transport.ListenAndAccept()
	if err != nil {
		logs.Logger.Errorf("Error Listening to Port %+v", err)
		return err
	}
	return nil
}
