package main

import (
	"time"
)

// 定义区块结构
type Block struct {
	Timestamp     int64  //时间戳
	Data          []byte //数据域
	PrevBlockHash []byte //前块hash值
	Hash          []byte //当前块hash值
	Nonce         int64  //随机值
}

// // 区块设置内部hash的方法
// func (b *Block) SetHash() {
// 	//将时间戳转换为[]byte
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	//将前块hash、交易信息、时间戳联合到一起
// 	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
// 	//计算本块hash值
// 	hash := sha256.Sum256(headers)
// 	//[32]byte -> []byte
// 	b.Hash = hash[:]
// }

// 创建Block，返回Block指针
func NewBlock(data string, prevBlockHash []byte) *Block {
	//先构造block
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	//需要先挖矿
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	//设置hash和nonce
	block.Hash = hash
	block.Nonce = nonce
	return block
}

// 创世块创建，返回创世块Block指针
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
