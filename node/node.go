package node

import (
	"context"
	"sync"

	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/p2p"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type NodeOpts struct {
	ListenAddr string
	Transport  p2p.Transport
}
type Node struct {
	version string
	NodeOpts

	PeerLock sync.RWMutex
	Peers    map[proto.NodeClient]*proto.Version
	proto.UnimplementedNodeServer
}

func NewNode(listenAddr string) *Node {
	tcpOptions := p2p.TCPTransportOpts{
		ListenAddr: listenAddr,
	}
	tcp_transport := p2p.NewTCPTransport(tcpOptions)
	nodeOpts := NodeOpts{
		ListenAddr: listenAddr,
		Transport:  tcp_transport,
	}

	return &Node{
		PeerLock: sync.RWMutex{},
		Peers:    make(map[proto.NodeClient]*proto.Version),
		NodeOpts: nodeOpts,

		version: "blocker-0.1",
	}
}
func (n *Node) BootStrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		client, err := makeNodeClient(addr)
		if err != nil {
			return err
		}
		peerVersion, err := client.HandShake(context.Background(), n.getVersion())
		if err != nil {
			logs.Logger.Errorf("Error BootStraping Network")
			continue
		}
		n.AddPeer(client, peerVersion)
	}
	return nil
}

func (n *Node) StartServer() error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	listener, err := n.Transport.Listen()
	if err != nil {
		logs.Logger.Errorf("Error Listening to Port %+v", err)
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)
	if err := grpcServer.Serve(listener); err != nil {
		logs.Logger.Errorf("Error Spinning up grpc server %+v", err)
		return err
	}
	return nil
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
	client, err := makeNodeClient(clientVersion.ListenAddr)
	if err != nil {
		return nil, err
	}
	n.AddPeer(client, clientVersion)
	return n.getVersion(), nil
}

func (n *Node) AddPeer(c proto.NodeClient, peerVersion *proto.Version) {
	n.PeerLock.Lock()
	defer n.PeerLock.Unlock()
	logs.Logger.Infof("[%s]:: New peer connected (%s)", n.ListenAddr, peerVersion.ListenAddr)
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

func (n *Node) getVersion() *proto.Version {
	version := &proto.Version{
		Version:    n.version,
		Height:     100,
		ListenAddr: n.ListenAddr,
	}
	return version
}
