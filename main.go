package main

import (
	"context"
	"time"

	"github.com/ranjankuldeep/DisBlocker/crypto"
	"github.com/ranjankuldeep/DisBlocker/logs"
	"github.com/ranjankuldeep/DisBlocker/node"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"github.com/ranjankuldeep/DisBlocker/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(2 * time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(2 * time.Second)
	makeNode(":6000", []string{":4000"})
	time.Sleep(4 * time.Second)
	makeTransaction()
	select {}
}

func makeTransaction() error {
	client, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logs.Logger.Errorf("Error Making Transaction %+v", err)
		return err
	}
	c := proto.NewNodeClient(client)
	privKey := crypto.GeneratePrivateKey()
	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevHash:     util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}
	_, err = c.HandleTransaction(context.Background(), tx)
	if err != nil {
		logs.Logger.Errorf("Error Handling Transaction %+v", err)
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
