package main

import (
	"context"
	"fmt"
	"math/big"
	_ "math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
0x
000000000000000000000000
382ab9a91b2986107d4e97f59edbfa9b69045e1e
0x000000000000000000000000746324a75d8ca24dfe61a83400e62dc78ac6d8ec
0x
00000000000000000000000000000000
000000000000000000000000000003e6
*/

func main() {
	//1. 连接到geth
	conn, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("failed to connet to geth", err)
		return
	}
	defer conn.Close()
	//2. 查询日志
	// 先设置过滤条件，设为空
	query := ethereum.FilterQuery{
		Addresses: []common.Address{},
		Topics:    [][]common.Hash{{}},
	}

	//查询全部日志
	logs, err := conn.FilterLogs(context.Background(), query)
	if err != nil {
		fmt.Println("failed to FilterLogs", err)
		return
	}
	//3.//合约地址处理
	cAddress := common.HexToAddress("0x4b6388442c218751604CC3aec7512efE850C7D15")
	topicHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	for _, v := range logs {
		if cAddress == v.Address {
			if len(v.Topics) == 3 {
				if v.Topics[0] == topicHash {
					from := v.Topics[1].Bytes()[len(v.Topics[1].Bytes())-20:]
					to := v.Topics[2].Bytes()[len(v.Topics[2].Bytes())-20:]
					val := big.NewInt(0)
					val.SetBytes(v.Data)
					fmt.Printf("BlockNumber : %d\n", v.BlockNumber)
					fmt.Printf("from : 0x%x\n", from)
					fmt.Printf("to : 0x%x\n", to)
					fmt.Printf("from : %d\n\n", val.Int64())
				}

			}

		}
	}

}
