package p2p

type GRPCTransportOpts struct {
	ListenAddr    string // Holds the address where a peer is listening.
	HandshakeFunc HandshakeFunc
	OnPeer        func(Peer) error
}

type GRPCTransport struct {
}
