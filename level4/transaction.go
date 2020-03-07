package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const subsidy = 10

// 交易输入结构
type TXInput struct {
	Txid     []byte //引用交易ID
	VoutIdx  int    //使用的交易输出编号
	FromAddr string //输入方验签
}

// 交易输出结构
type TXOutput struct {
	Value  int    //输出金额
	ToAddr string //收方验签
}

// 交易结构
type Transaction struct {
	ID   []byte     //交易ID
	Vin  []TXInput  //交易输入项
	Vout []TXOutput //交易输出项
}

// 将交易信息转换为hash，并设为ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	enc.Encode(tx)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//创建CoinBase交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	//创建一个输入项
	txin := TXInput{[]byte{}, -1, data}
	//创建输出项
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}
