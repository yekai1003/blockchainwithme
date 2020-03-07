package main

import (
	"crypto/sha256"
	"fmt"
)

func main1() {
	val := Base58Decode([]byte("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
	fmt.Printf("%x，%d\n", val, len(val))
	fmt.Printf("%x, %x\n", val[1:21], val[21:])
	hash1 := sha256.Sum256(val[:21])
	hash2 := sha256.Sum256(hash1[:])
	fmt.Printf("%x\n", hash2[:4])
}

func main() {
	//1. 构造钱包
	wallet := NewWallet()
	//2. 生成地址
	address := wallet.GetAddress()
	fmt.Printf("%s\n", address)
}
