package p2p

import (
	"net"
)

type TCPTransportOpts struct {
	ListenAddr string // Holds the address where a peer is listening.
}

type TCPTransport struct {
	listener net.Listener
	TCPTransportOpts
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) Listen() (net.Listener, error) {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return nil, err
	}
	return t.listener, err
}
