package p2p

// This represents a grpc message sent over the wir
type RPC struct {
	From    string
	Payload []byte
}
