package main

import (
	"fmt"
)

//1LFHh6iUoj5gHNVVtChrP96JwA9fM2LBYg
func main1() {
	//1. 创建钱包地址
	wallet := NewWallet()
	from := fmt.Sprintf("%s", wallet.GetAddress())

	bc := CreateBlockchain(from)
	defer bc.db.Close()
	//2. 获取叶开余额
	bc.getBalance(from)
}

func main() {
	cli := NewCli()
	cli.Run()
}
