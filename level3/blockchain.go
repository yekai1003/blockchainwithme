package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

//db文件名
const dbFile = "blockchain.db"

//数据库bucket名
const blocksBucket = "blocks"

// 区块链：一个区块的指针数组
type Blockchain struct {
	//记录最新块hash值
	tip []byte
	//存放区块链的db
	db *bolt.DB
}

//迭代器
type BlockchainIterator struct {
	currentHash []byte   //当区块hash
	db          *bolt.DB //已经打开的数据库
}

//通过Blockchain构造迭代器
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// 创建区块链结构，初始化只有创世块
func NewBlockchain() *Blockchain {
	var tip []byte
	//1.打开数据库文件
	db, _ := bolt.Open(dbFile, 0600, nil)
	//2.更新数据库
	db.Update(func(tx *bolt.Tx) error {
		//2.1 获取bucket
		buck := tx.Bucket([]byte(blocksBucket))
		if buck == nil {
			//2.2.1 第一次使用，创建创世块
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()
			//2.2.2 区块数据编码
			block_data := genesis.Serialize()
			//2.2.3 创建新bucket，存入区块信息
			bucket, _ := tx.CreateBucket([]byte(blocksBucket))
			bucket.Put(genesis.Hash, block_data)
			bucket.Put([]byte("last"), genesis.Hash)
			tip = genesis.Hash

		} else {
			//2.3 不是第一次使用，之前有块
			tip = buck.Get([]byte("last"))
		}
		return nil
	})
	//3. 记录Blockchain 信息
	return &Blockchain{tip, db}
}

// 向区块链结构上增加一个区块
func (bc *Blockchain) AddBlock(data string) {
	var tip []byte
	//1. 获取tip值，此时不能再打开数据库文件，要用区块的结构
	bc.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))
		tip = buck.Get([]byte("last"))
		return nil
	})
	//2. 更新数据库
	bc.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))
		block := NewBlock(data, tip)
		//将新区块放入db
		buck.Put(block.Hash, block.Serialize())
		buck.Put([]byte("last"), block.Hash)
		//覆盖tip值
		bc.tip = block.Hash
		return nil
	})
}

//获取前一个区块hash，返回当区块数据
func (i *BlockchainIterator) PreBlock() (*Block, bool) {
	var block *Block
	//根据hash获取块数据
	i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		//解码当前块数据
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	//当前hash变更为前块hash
	i.currentHash = block.PrevBlockHash
	//返回区块
	return block, len(i.currentHash) > 0
}
