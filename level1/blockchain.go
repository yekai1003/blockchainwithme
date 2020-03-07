package main

// 区块链：一个区块的指针数组
type Blockchain struct {
	blocks []*Block
}

// 向区块链结构上增加一个区块
func (bc *Blockchain) AddBlock(data string) {
	//获取前块信息
	prevBlock := bc.blocks[len(bc.blocks)-1]
	//利用前块生成新块
	newBlock := NewBlock(data, prevBlock.Hash)
	//添加到区块链结构中
	bc.blocks = append(bc.blocks, newBlock)
}

// 创建区块链结构，初始化只有创世块
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
