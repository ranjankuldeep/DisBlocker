package types

import (
	"crypto/sha256"

	"github.com/ranjankuldeep/DisBlocker/proto"
	pb "google.golang.org/protobuf/proto"
)

// Hashblock returns a SHA256 of the Header.
func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}
