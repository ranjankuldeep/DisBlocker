package node

import (
	"context"

	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.None, error) {
	perr, ok := peer.FromContext(ctx)
	if !ok {
		logs.Logger.Errorf("No peer exist")
		return nil, nil
	}
	logs.Logger.Infof("Peer: %+v", perr)
	return nil, nil
}
