package types

import (
	"crypto/sha256"

	"github.com/ranjankuldeep/DisBlocker/crypto"
	"github.com/ranjankuldeep/DisBlocker/proto"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}
func HashTransaction(tx *proto.Transaction) []byte {
	buf, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(buf)
	return hash[:]
}
