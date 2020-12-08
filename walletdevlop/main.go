package main

import (
	"github.com/yekai1003/blockchainwithme/walletdevlop/cli"
)

// func main1() {
// 	w, err := wallet.NewWallet("./keystore")
// 	if err != nil {
// 		fmt.Println("Failed to NewWallet", err)
// 		return
// 	}
// 	w.StoreKey("123")
// }

func main() {
	c := client.NewCmdClient("http://localhost:8545", "./keystore")
	c.Run()
}
