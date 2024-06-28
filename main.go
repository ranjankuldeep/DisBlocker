package main

import (
	"context"
	"time"

	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/node"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(100 * time.Millisecond)
	makeNode(":4000", []string{":3000"})
	time.Sleep(4 * time.Second)
	makeNode(":5001", []string{":4000"})

	// go func() {
	// 	for {
	// 		makeTransaction()
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()
	select {}
}

func makeTransaction() error {
	client, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logs.Logger.Errorf("Error Making Transaction %+v", err)
		return err
	}
	c := proto.NewNodeClient(client)
	version := &proto.Version{
		Version:    "blocker-0.1",
		Height:     100,
		ListenAddr: ":4000",
	}

	_, err = c.HandShake(context.TODO(), version)
	if err != nil {
		logs.Logger.Errorf("Error Making Transaction")
		return err
	}
	return nil
}

func makeNode(listenAddr string, bootStrapAddrs []string) *node.Node {
	blockNode := node.NewNode(listenAddr)
	logs.Logger.Infof("Server Running on Port %s", listenAddr)
	go blockNode.StartServer(bootStrapAddrs)
	return blockNode
}
