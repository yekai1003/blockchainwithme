package main

import (
	"fmt"
)

func main() {
	//创世块初始化区块链
	bc := CreateBlockchain()
	//创建2个块记录2笔交易
	// bc.AddBlock("Send 1 BTC to Yekai")
	// bc.AddBlock("Send 2 more BTC to Fuhongxue")

	bci := bc.Iterator()
	for {
		block, next := bci.PreBlock()
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Transactions[0].Vin[0].FromAddr)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %t\n", pow.Validate())
		fmt.Println()
		if !next {
			//next为假代表已经到创世块了
			break
		}
	}

}
