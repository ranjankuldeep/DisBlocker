package p2p

import "net"

// Transport is anything that handles the communication
// between the nodes in the network. This can be of the
// form (TCP, UDP, websockets, ...)
type Transport interface {
	Listen() (net.Listener, error) // Should be conitnously listening for any incoming connection and accept the connection with handshaking
}
