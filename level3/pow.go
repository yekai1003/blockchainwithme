package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = int64(math.MaxInt64)
)

const targetBits = 6

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// 将int64写入[]byte
func Int2Hex(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}

//准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			Int2Hex(pow.block.Timestamp),
			Int2Hex(int64(targetBits)),
			Int2Hex(nonce),
		},
		[]byte{},
	)

	return data
}

// 挖矿运行
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	fmt.Printf("Mining the block containing %s, maxNonce=%d\n", pow.block.Data, maxNonce)
	for nonce < maxNonce {
		//数据准备
		data := pow.prepareData(nonce)
		//计算hash
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		//按字节比较，hashInt.Cmp小于0代表找到目标Nonce
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// 校验区块正确性
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
