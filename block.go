package main

import (
	"time"
)

//定义区块
type Block struct {
	Timestamp int64	//时间线，从1970年1月1日0.00.00
	Data []byte		//交易数据
	PrevBlockHash []byte	//上一块数据的哈希
	Hash	[]byte	//当前模块的哈希
	Nonce int //工作量证明
}
/*
//设置结构体对象的哈希
func (block *Block)setHash(){
	//处理当前的时间，转换成十进制的字符串,再转化为字节集
	timestamp := []byte(strconv.FormatInt((block.Timestamp),10))
	//叠加要哈希的数据
	headers := bytes.Join([][]byte{block.PrevBlockHash,block.Data,timestamp},[]byte{})
	//计算出哈希地址
	hash := sha256.Sum256(headers)
	//设置哈希
	block.Hash = hash[:]
}
*/
//创建一个区块
func NewBlock(data string,prevBlockHash []byte)  *Block {
	//block是一个指针，取得一个对象初始化之后的地址
	block := &Block{time.Now().Unix(),[]byte(data),prevBlockHash,[]byte{},1}
	pow := NewProofOfWork(block)	//挖矿这个区块
	nonce,hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	//block.setHash()	//设置哈希
	return block
}

//创建创世区块，意味着没有前一块
func NewGenesisBlock()  *Block {
	return NewBlock("first block",[]byte{})
}