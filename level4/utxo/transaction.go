package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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

//是否为CoinBase交易
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 &&
		len(tx.Vin[0].Txid) == 0 &&
		tx.Vin[0].VoutIdx == -1
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

// 判断该输入是否可以被某账户使用
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.FromAddr == unlockingData
}

// 判断某输出是否可以被账户使用
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ToAddr == unlockingData
}

//创建普通交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	//1.需要组合输入项和输出项
	var inputs []TXInput
	var outputs []TXOutput
	//2. 查询最小UTXO
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	//3. 构建输入项
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	//4. 构建输出项
	outputs = append(outputs, TXOutput{amount, to})
	// 需要找零
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}
	//5. 交易生成
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
