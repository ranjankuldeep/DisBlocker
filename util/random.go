package util

import (
	"crypto/rand"
	"io"
	randM "math/rand"
	"time"

	"github.com/ranjankuldeep/DisBlocker/proto"
)

func RandomHash() []byte {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return buf
}

func RandomBlock() *proto.Block {
	header := &proto.Header{
		Version:   1,
		Height:    int32(randM.Intn(1000)),
		PrevHash:  RandomHash(),
		RootHash:  RandomHash(),
		TimeStamp: time.Now().UnixNano(),
	}
	return &proto.Block{
		Header: header,
	}
}
