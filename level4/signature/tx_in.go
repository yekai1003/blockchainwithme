package main

import "bytes"

// 交易输入结构
type TXInput struct {
	Txid      []byte //引用交易ID
	VoutIdx   int    //使用的交易输出编号
	Signature []byte //签名信息
	PubKey    []byte //公钥
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
