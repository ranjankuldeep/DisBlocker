package main

import (
	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/node"
	"github.com/ranjankuldeep/DisBlocker/p2p"
)

func main() {
	blockNode := node.NewNode()
	tcpOptions := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
	}
	tcp_transport := p2p.NewTCPTransport(tcpOptions)

	nodeOpts := node.NodeServerOpts{
		Transport: tcp_transport,
	}
	nodeServer := node.NewNodeServer(*blockNode, &nodeOpts)

	logs.Logger.Infof("Server Running on Port %s", tcpOptions.ListenAddr)
	if err := nodeServer.StartServer(); err != nil {
		logs.Logger.Errorf("Error Starting NOde server %+v", err)
		panic(err)
	}
}
