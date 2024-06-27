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
	blockNode := node.NewNode()
	listenAddr := ":3000"

	nodeServer := node.NewNodeServer(*blockNode, listenAddr)
	logs.Logger.Infof("Server Running on Port %s", listenAddr)

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
	version := &proto.Version{
		Version: "blocker-0.1",
		Height:  100,
	}

	_, err = c.HandShake(context.TODO(), version)
	if err != nil {
		logs.Logger.Errorf("Error Making Transaction")
		return err
	}
	return nil
}
