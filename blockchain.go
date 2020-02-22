package main

type BlockChain struct {
	blocks []*Block
}

//增加一个区块
func (blocks *BlockChain)AddBlock (data string){
	prevBlcok := blocks.blocks[len(blocks.blocks) - 1]	//取出最后一个区块
	newBlock := NewBlock(data,prevBlcok.Hash)	//创建一个区块
	blocks.blocks = append(blocks.blocks,newBlock)
}

//创建一个区块链
func NewBlockchain() *BlockChain{
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}




