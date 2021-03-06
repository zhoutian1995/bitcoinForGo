package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

//定义区块
type Block struct {
	Timestamp int64	//时间线，从1970年1月1日0.00.00
	//Data []byte		//交易数据
	Transactions []*Transaction //交易集合
	PrevBlockHash []byte	//上一块数据的哈希
	Hash	[]byte	//当前模块的哈希
	Nonce int //工作量证明
}

//创建一个区块
func NewBlock(transcations []*Transaction ,prevBlockHash []byte)  *Block {
	//block是一个指针，取得一个对象初始化之后的地址
	block := &Block{time.Now().Unix(),
		transcations,
		prevBlockHash,
				[]byte{},0}
	pow := NewProofOfWork(block)	//挖矿这个区块
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}


//创建创世区块，意味着没有前一块
func NewGenesisBlock(coinbase *Transaction)  *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//对交易实现哈希计算
func  (block *Block)HashTransactions()[]byte{
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}



//把对象转化为二进制字节集，可以写入文件
func (block *Block) Serialize() []byte{
	var result bytes.Buffer	//开辟内存存放字节集合
	encoder := gob.NewEncoder(&result)	//编码对象创建
	err := encoder.Encode(block)	//编码操作
	if err != nil{
		log.Panic(err)	//处理错误
	}
	return result.Bytes()
}
//读取文件读到二进制字节并集合转化为对象
func DeserializeBlock(data []byte) *Block{
	var block Block	//创建存储用于字节转化的对象
	decoder := gob.NewDecoder(bytes.NewReader(data))//解码
	err := decoder.Decode(&block)	//尝试解码
	if err != nil{
		fmt.Println("12345")
	}

	return &block
}

