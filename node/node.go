package node

import (
	"context"

	"github.com/ranjankuldeep/DisBlocker/proto"
)

type Node struct {
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.None, error) {
	return nil, nil
}
