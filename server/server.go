package server

import "github.com/ranjankuldeep/DisBlocker/p2p"

type BlockServerOpts struct {
	Transport p2p.GRPCTransport
}
type BlockServer struct {
	BlockServerOpts
}
