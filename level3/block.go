package main

import (
	"bytes"
	"encoding/gob"
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

// 序列化区块
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	//编码器
	encoder := gob.NewEncoder(&result)
	//编码
	encoder.Encode(b)
	return result.Bytes()
}

//区块数据还原为Block
func DeserializeBlock(d []byte) *Block {
	var block Block
	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(d))
	//解析区块数据
	decoder.Decode(&block)
	return &block
}

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
