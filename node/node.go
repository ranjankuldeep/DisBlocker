package node

import (
	"context"

	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		version: "blocker-0.1",
	}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	perr, ok := peer.FromContext(ctx)
	if !ok {
		logs.Logger.Errorf("No peer exist")
		return nil, nil
	}
	logs.Logger.Infof("Peer: %+v", perr)
	return nil, nil
}

func (n *Node) HandShake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}
	peer, ok := peer.FromContext(ctx)
	if !ok {
		logs.Logger.Info("Unable to fetch peer from the context")
	}
	logs.Logger.Infof("Recieved version from %s: %v ", peer, v)
	return ourVersion, nil
}
