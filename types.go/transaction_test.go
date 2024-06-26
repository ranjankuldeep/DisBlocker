package types

import (
	"testing"

	"github.com/ranjankuldeep/DisBlocker/crypto"
	"github.com/ranjankuldeep/DisBlocker/proto"
	"github.com/ranjankuldeep/DisBlocker/util"
)

func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddress := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevHash:     util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}
	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}
	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}
	tx := &proto.Transaction{
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}
	sign := SignTransaction(fromPrivKey, tx)
	input.Signature = sign.Bytes()
}
