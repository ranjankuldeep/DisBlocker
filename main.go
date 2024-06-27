package main

import (
	"context"
	"time"

	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/node"
	"github.com/ranjankuldeep/DisBlocker/p2p"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	go func() {
		for {
			makeTransaction()
			time.Sleep(2 * time.Second)
		}
	}()

	if err := nodeServer.StartServer(); err != nil {
		logs.Logger.Errorf("Error Starting NOde server %+v", err)
		panic(err)
	}
}

func makeTransaction() error {
	client, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logs.Logger.Errorf("Error Making Transaction %+v", err)
		return err
	}
	c := proto.NewNodeClient(client)
	tx := &proto.Transaction{
		Version: 1,
	}
	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		logs.Logger.Errorf("Error Making Transaction")
		return err
	}
	return nil
}
