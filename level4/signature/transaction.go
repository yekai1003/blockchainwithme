package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

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
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	//创建输出项
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetID()

	return &tx
}

//创建普通交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain, wallet *Wallet) *Transaction {
	//1.需要组合输入项和输出项
	var inputs []TXInput
	var outputs []TXOutput
	//2. 查询最小UTXO
	acc, validOutputs := bc.FindSpendableOutputs(HashPubKey(wallet.PublicKey), amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	//3. 构建输入项
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	//4. 构建输出项
	outputs = append(outputs, *NewTXOutput(amount, to))
	// 需要找零
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	//5. 交易生成
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

//交易签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//1. CoinBase交易无需签名
	if tx.IsCoinbase() {
		return
	}
	//2. 修剪交易
	txCopy := tx.TrimmedCopy()
	//3. 循环向输入项签名
	for inID, vin := range txCopy.Vin {
		//找到输入项引用的交易
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.VoutIdx].PubKeyHash
		txCopy.SetID()
		//txid生成后再把PubKey置空
		txCopy.Vin[inID].PubKey = nil
		//使用ecsda签名获得r和s
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		//形成签名数据
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}

// 交易修剪
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	//将原交易内的签名和公钥都置空
	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.VoutIdx, nil, nil})
	}
	//复制原输入项
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	//复制一份交易
	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}
