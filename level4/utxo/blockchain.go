package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

//db文件名
const dbFile = "blockchain.db"

//数据库bucket名
const blocksBucket = "blocks"

//定义矿工地址
const miner = "yekai"

//创世块留言
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

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

// 判断区块链是否已经存在
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// 创建区块链结构，初始化只有创世块
func CreateBlockchain() *Blockchain {
	//1. 只能第一次创建
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	//2. 没有则创建文件
	db, _ := bolt.Open(dbFile, 0600, nil)
	//3.更新数据库
	db.Update(func(tx *bolt.Tx) error {
		//2.1 获取bucket
		buck := tx.Bucket([]byte(blocksBucket))
		if buck == nil {
			//2.2.1 第一次使用，创建创世块
			fmt.Println("No existing blockchain found. Creating a new one...")
			cbtx := NewCoinbaseTX(miner, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
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
func (bc *Blockchain) MinedBlock(transactions []*Transaction, data string) {
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
		//创建CoinBase交易
		cbtx := NewCoinbaseTX(miner, data)
		transactions = append(transactions, cbtx)
		block := NewBlock(transactions, tip)
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

//查找账户可解锁的全部交易
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	//已经花出的UTXO，构建tx->VOutIdx的map
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block, next := bci.PreBlock()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// 如果已经被花出了，直接跳过此交易
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				//可以被address解锁，就代表属于address的utxo在此交易中
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			//用来维护spentTXOs，已经被引用过了，代表被使用
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VoutIdx)
					}
				}
			}
		}

		if !next {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	//先找所有交易
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			//可解锁代表是用户的资产
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

//获取部分满足交易的UTXO
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	//获取可使用交易
	unspentTXs := bc.FindUnspentTransactions(address)
	//记录余额
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				//utxo足够了就跳出循环 break可以跳出多重循环
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *Blockchain) getBalance(address string) {

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

//交易发送
func (bc *Blockchain) send(from, to string, amount int, data string) {
	//创建普通交易
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MinedBlock([]*Transaction{tx}, data)
	fmt.Println("Success!")
}
