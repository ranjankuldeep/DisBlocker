package store

import "github.com/ranjankuldeep/DisBlocker/proto"

type BlockStorer interface {
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}
