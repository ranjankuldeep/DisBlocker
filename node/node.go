package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/p2p"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"github.com/ranjankuldeep/DisBlocker/types.go"
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
	Mempool MemPool

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
		Mempool:  NewInMemPool(),
		version:  "blocker-0.1",
	}
}

func (n *Node) StartServer(bootStrapAddrs []string) error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	listener, err := n.Transport.Listen()
	if err != nil {
		logs.Logger.Errorf("Error Listening to Port %+v", err)
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)

	if len(bootStrapAddrs) > 0 {
		go n.bootStrapNetwork(bootStrapAddrs)
	}
	if err := grpcServer.Serve(listener); err != nil {
		logs.Logger.Errorf("Error Spinning up grpc server %+v", err)
		return err
	}
	return nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		logs.Logger.Errorf("No peer exist")
		return nil, nil
	}
	hash := hex.EncodeToString(types.HashTransaction(tx))

	// First add to your own mempool.
	if n.Mempool.Add(tx) {
		logs.Logger.Infof("Recived tx from: %s, hash %s", peer.Addr, hash)
		go func() {
			// Broadcast to other nodes in the network.
			if err := n.broadCast(tx); err != nil {
				logs.Logger.Errorf("Broadcast error %+v", err)
			}
		}()
	}
	return &proto.Ack{}, nil
}

func (n *Node) broadCast(msg any) error {
	for peer := range n.Peers {
		switch v := msg.(type) {
		case *proto.Transaction:
			_, err := peer.HandleTransaction(context.Background(), v)
			if err != nil {
				logs.Logger.Errorf("Error Broadcasting Transaction %+v", err)
				return err
			}
		}
	}
	return nil
}

func (n *Node) HandShake(ctx context.Context, clientVersion *proto.Version) (*proto.Version, error) {
	client, err := makeNodeClient(clientVersion.ListenAddr)
	if err != nil {
		return nil, err
	}
	n.addPeer(client, clientVersion)
	return n.getVersion(), nil
}

func (n *Node) addPeer(c proto.NodeClient, peerVersion *proto.Version) {
	n.PeerLock.Lock()
	defer n.PeerLock.Unlock()

	// Handle the logic where we decide to accept or drop the
	// incoming node connection.
	n.Peers[c] = peerVersion
	logs.Logger.Infof("[%s]:: New peer succesfully connected (%s)", n.ListenAddr, peerVersion.ListenAddr)

	if len(peerVersion.PeerList) > 0 {
		go n.bootStrapNetwork(peerVersion.PeerList)
	}

}

func (n *Node) bootStrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		if !n.canConnect(addr) {
			continue
		}
		logs.Logger.Infof("[%s]: Dialing Remote Node [%s]", n.ListenAddr, addr)
		peerClient, peerVersion, err := n.redialWithExponentialBackoff(addr)
		if err != nil {
			logs.Logger.Errorf("[%s]: Error Dialing Remote Node [%s]", n.ListenAddr, addr)
			continue
		}
		n.addPeer(peerClient, peerVersion)
	}
	return nil
}

func (n *Node) redialWithExponentialBackoff(addr string) (proto.NodeClient, *proto.Version, error) {
	var (
		maxRetries     = 10
		initialBackoff = time.Second
		maxBackoff     = time.Minute
	)

	backoff := initialBackoff
	retries := 0
	var lastErr error

	for retries < maxRetries {
		nodeClient, protoVersion, err, isDailed := n.dialRemote(addr)
		if isDailed {
			fmt.Println("Connected successfully")
			return nodeClient, protoVersion, err
		}

		fmt.Printf("Connection failed, retrying in %v...\n", backoff)
		time.Sleep(backoff)
		lastErr = err

		// Increase the backoff time
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
		retries++
	}
	fmt.Println("Failed to connect after multiple attempts")
	return nil, nil, lastErr
}

func (n *Node) dialRemote(addr string) (proto.NodeClient, *proto.Version, error, bool) {
	client, _ := makeNodeClient(addr)
	peerVersion, err := client.HandShake(context.Background(), n.getVersion())
	if err != nil {
		logs.Logger.Errorf("Error Dialing %+v", err)
		return nil, nil, err, false
	}
	return client, peerVersion, nil, true
}

func (n *Node) canConnect(addr string) bool {
	if addr == n.ListenAddr {
		return false
	}
	for _, connectedPeerAddr := range n.getPeerList() {
		if addr == connectedPeerAddr {
			return false
		}
	}
	return true
}

func (n *Node) removePeer(c proto.NodeClient) {
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
		PeerList:   n.getPeerList(),
	}
	return version
}

func (n *Node) getPeerList() []string {
	n.PeerLock.RLock()
	defer n.PeerLock.RUnlock()

	peers := []string{}
	for _, version := range n.Peers {
		peers = append(peers, version.ListenAddr)
	}
	return peers
}
