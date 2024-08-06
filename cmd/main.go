package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/blocksignalio/core"
)

func main() {
	const (
		dao  = "0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413"
		weth = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
		pepe = "0x6982508145454Ce325dDbE47a25d4ec3d2311933"
	)

	db, err := core.Open()
	if err != nil {
		panic(err)
	}
	core.BackfillLogs(context.Background(), db, "0x57a9cbED053f37EB67d6f5932b1F2f9Afbe347F3")
	// fmt.Println(core.SelectLogs(db, "0x57a9cbED053F37EB67D6F5932B1F2F9AFBE347F3", "", 0, 0))

	if true {
		return
	}

	var fromBlock uint64 = 17899693
	contract := weth
	toBlock, logs, err := core.QueryLogs(context.TODO(), fromBlock, contract)
	if err != nil {
		panic(err)
	}
	for i, log := range logs[:3] {
		fmt.Printf("logs[%d]:\n", i)
		fmt.Println("\taddress:", hex.EncodeToString(log.Address[:]))
		fmt.Println("\tblockHash:", hex.EncodeToString(log.BlockHash[:]))
		fmt.Println("\tblockNumber:", log.BlockNumber)
		fmt.Println("\tdata:", hex.EncodeToString(log.Data))
		fmt.Println("\tindex:", log.Index)
		fmt.Println("\tremoved:", log.Removed)
		for j, topic := range log.Topics {
			fmt.Printf("\ttopic[%d]: %s\n", j, hex.EncodeToString(topic[:]))
		}
		fmt.Println("\ttxHash:", hex.EncodeToString(log.TxHash[:]))
		fmt.Println("\ttxIndex:", log.TxIndex)
	}
	fmt.Println(fromBlock, toBlock, toBlock-fromBlock, len(logs))

	// ch := make(chan<- types.Log)
	// sub, err := client.SubscribeFilterLogs(context.TODO(), query, ch)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(sub)

	// balance, err := client.BalanceAt(context.TODO(), target, big.NewInt(19984478))
	// if err != nil {
	// 	panic(err)
	// }

	// eth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	// fmt.Println(new(big.Int).Div(balance, eth))
}
