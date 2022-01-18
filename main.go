package main

import (
	"context"
	"fmt"

	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
)

const alice = "tz1VSUr8wwNhLAzempoch5d6hLRiTh8Cjcjb"
const bob = "tz1aSkwEot3L2kmUvcoxzjMomb9mvBNuzFK6"

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := rpc.NewClient("http://localhost:20000", nil)
	if err != nil {
		panic(err)
	}

	block, err := client.GetTipHeader(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Height: %v\n", block.Level)

	players := map[string]string{alice: "alice", bob: "bob"}

	height := block.Level
	var lastBlock *rpc.Block
	var chainID string
	for i := int64(0); i <= height; i++ {
		block, err := client.GetBlockHeight(ctx, i)
		if err != nil {
			panic(err)
		}
		if block == nil {
			fmt.Printf("Nil block for level %v\n", i)
			break
		}

		if lastBlock == nil {
			chainID = block.ChainId.String()
		} else {
			if block.ChainId.String() != chainID {
				fmt.Print("Chain broken")
			}
			if block.Header.Predecessor.String() != lastBlock.Hash.String() {
				fmt.Printf("Expected hash chain\n")
			}
		}

		if val, ok := players[block.Metadata.Baker.String()]; ok {
			fmt.Printf("%v was the baker\n", val)
		}

		for _, update := range block.Metadata.BalanceUpdates {
			if val, ok := players[update.Address().String()]; ok {
				fmt.Printf("Block %d has a balance update for %v\n", i, val)
			}
		}

		for _, operationList := range block.Operations {
			for _, operation := range operationList {
				for _, opContents := range operation.Contents {
					if opContents.OpKind() == tezos.OpTypeOrigination {
						if typedOp, ok := opContents.(*rpc.OriginationOp); ok {
							if val, ok := players[typedOp.Source.String()]; ok {
								fmt.Printf("Operation %d made by %v\n", i, val)
							}
						}
					}
				}
			}
		}

		lastBlock = block
	}
}
