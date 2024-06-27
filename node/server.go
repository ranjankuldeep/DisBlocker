package node

import (
	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/p2p"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"google.golang.org/grpc"
)

type NodeServerOpts struct {
	Transport p2p.Transport
}
type NodeServer struct {
	NodeServerOpts
	Node
}

func NewNodeServer(node Node, opts *NodeServerOpts) *NodeServer {
	return &NodeServer{
		NodeServerOpts: *opts,
		Node:           node,
	}
}

func (n *NodeServer) StartServer() error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	listener, err := n.Transport.Listen()
	if err != nil {
		logs.Logger.Errorf("Error Listening to Port %+v", err)
		return err
	}
	proto.RegisterNodeServer(grpcServer, &n.Node)
	if err := grpcServer.Serve(listener); err != nil {
		logs.Logger.Errorf("Error Spinning up grpc server %+v", err)
		return err
	}
	return nil
}
