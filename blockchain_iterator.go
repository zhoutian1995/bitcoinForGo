package main

import (
	"github.com/boltdb/bolt"
	"github.com/labstack/gommon/log"
)

type BlockChainIterator struct {
	currentHash []byte	//当前要找的哈希
	db *bolt.DB	//数据库
}

//取得下一个区块
func (it *BlockChainIterator) Next() *Block{
	var block *Block

	err := it.db.View(func (tx *bolt.Tx ) error{
		bucket := tx.Bucket([]byte(blocksBucket))
		encodeBlock := bucket.Get(it.currentHash)	//抓取二进制数据
		block = DeserializeBlock(encodeBlock)	//解码

		return nil
	})

	if err != nil{
		log.Panic(err)
	}

	it.currentHash = block.PrevBlockHash	//当前要找的哈希更替

	return block
}