package main

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "blockchain.db"	//数据库文件名
const blockBucket = "blocks"	//名称
const genesisCoinbaseData = "genesisCoinbase"

type BlockChain struct{
	tip []byte//二进制数据
	db *bolt.DB	//数据库
}

type BlockChainIterator struct {
	currentHash []byte	//当前要找的哈希
	db *bolt.DB	//数据库
}

//挖矿带来的交易
func (blockchain *BlockChain) MineBlock(transactions []*Transaction) {
	var lastHash []byte//最后的哈希

	err := blockchain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))//查看数据
		lastHash = b.Get([]byte("l"))	//取出最后个区块的哈希

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)	//创建一个新的区块

	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))//取出索引
		err := b.Put(newBlock.Hash, newBlock.serialize())//存入数据库
		if err != nil {
			log.Panic(err)	//处理错误
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		blockchain.tip = newBlock.Hash

		return nil
	})
}


//找到包含未花费输出的交易
func (blockchain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction	//交易事物
	spentTXOs := make(map[string][]int)	//开辟内存
	bci := blockchain.Iterator()	//迭代器

	for {
		block := bci.Next()	//循环下一个

		for _, tx := range block.Transactions {	//循环每个交易
			txID := hex.EncodeToString(tx.ID)	//获取交易编号

		Outputs:
			for outIdx, out := range tx.Vout {	//循环输出交易，只有当输出不等于输入(真的花了钱)，才把这个交易事物加入数组列表中

				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs	//循环到不等为止
						}
					}
				}

				if out.CanBeunlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)	//没有花费的钱，加入列表
				}
			}

			if tx.IsCoinBase() == false {	//不是挖矿的奖励
				for _, in := range tx.Vin {
					if in.CanUnlockOutPutWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)//花费的钱，加入列表
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}
//获取所有没有使用的交易
func (blockchain *BlockChain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := blockchain.FindUnspentTransactions(address)//查找所有的交易

	for _, tx := range unspentTransactions {	//循环所有的事物交易
		for _, out := range tx.Vout {
			if out.CanBeunlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}
//对所有的未花费交易进行迭代，并对它的值进行累加。
// 当累加值大于或等于我们想要传送的值时，它就会停止并返回累加值，
// 同时返回的还有通过交易 ID 进行分组的输出索引。
func (blockchain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)	//输出
	unspentTXs := blockchain.FindUnspentTransactions(address)	//根据地址查找所有交易
	accumulated := 0	//累计

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)	//获取编号

		for outIdx, out := range tx.Vout {
			if out.CanBeunlockedWith(address) && accumulated < amount {
				accumulated += out.Value	//统计金额
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

//判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

//迭代器
func (block *BlockChain) Iterator() *BlockChainIterator{
	bcit := &BlockChainIterator{block.tip,block.db}

	return bcit	//根据区块链创建区块链迭代器
}

//取得下一个区块
func (it *BlockChainIterator) Next() *Block{
	var block *Block

	err := it.db.View(func (tx *bolt.Tx ) error{
		bucket := tx.Bucket([]byte(blockBucket))
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
//新加一个区块
func NewBlockChain(address string) *BlockChain {
	if dbExists() == false {
		fmt.Println("不存在区块，需要新建一个")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

func CreateBlockchain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("数据库已经存在，无需创建")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinBaseTX(address, genesisCoinbaseData)	//创建创世区块的事物交易
		genesis := NewGenesisBlock(cbtx)	//创建创区块的快

		b, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}
