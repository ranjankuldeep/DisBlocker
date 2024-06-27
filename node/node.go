package node

import (
	"context"
	"sync"

	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string

	PeerLock sync.RWMutex
	Peers    map[proto.NodeClient]*proto.Version
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		PeerLock: sync.RWMutex{},
		Peers:    make(map[proto.NodeClient]*proto.Version),
		version:  "blocker-0.1",
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

func (n *Node) HandShake(ctx context.Context, clientVersion *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}
	peer, ok := peer.FromContext(ctx)
	if !ok {
		logs.Logger.Info("Unable to fetch peer from the context")
	}
	client, err := makeNodeClient(peer.Addr.String())
	if err != nil {
		return nil, err
	}
	logs.Logger.Infof("Recieved version from %s: %v ", peer, clientVersion)
	n.AddPeer(client, clientVersion)
	return ourVersion, nil
}

func (n *Node) AddPeer(c proto.NodeClient, peerVersion *proto.Version) {
	n.PeerLock.Lock()
	defer n.PeerLock.Unlock()
	logs.Logger.Infof("New peer connected (%s), height (%d)", peerVersion.Version, peerVersion.Height)
	n.Peers[c] = peerVersion
}

func (n *Node) RemovePeer(c proto.NodeClient) {
	n.PeerLock.Lock()
	defer n.PeerLock.Unlock()
	delete(n.Peers, c)
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	client, err := grpc.NewClient(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logs.Logger.Errorf("Error Making Transaction %+v", err)
		return nil, err
	}
	c := proto.NewNodeClient(client)
	return c, nil
}
