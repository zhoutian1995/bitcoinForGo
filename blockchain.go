package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const dbFile = "blockchain.db"	//数据库文件名
const blocksBucket = "blocks"	//名称
const genesisCoinbaseData = "genesisCoinbase"

type BlockChain struct{
	tip []byte//二进制数据
	db *bolt.DB	//数据库
}

func CreateBlockchain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("数据库已经存在，无需创建")
		os.Exit(1)
	}

	var tip []byte	//存储区块链的二进制数据
	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)	//创建创世区块的事物交易
	genesis := NewGenesisBlock(cbtx)	//创建创区块的快


	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
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

//新加一个区块
func NewBlockchain() *BlockChain {
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
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}



//找到包含未花费输出的交易,获取没有使用输出的交易列表
func (blockchain *BlockChain) FindUnspentTransactions(pubkeyhash []byte) []Transaction {
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

				if out.IsLockedWithKey(pubkeyhash) {
					unspentTXs = append(unspentTXs, *tx)	//没有花费的钱，加入列表
				}
			}

			if tx.IsCoinBase() == false {	//不是挖矿的奖励
				for _, in := range tx.Vin {
					if in.UsesKey(pubkeyhash) {
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

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, nil
}

//获取所有没有使用的交易
func (bc *BlockChain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinBase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

//迭代器
func (block *BlockChain) Iterator() *BlockChainIterator{
	bcit := &BlockChainIterator{block.tip,block.db}

	return bcit	//根据区块链创建区块链迭代器
}

//挖矿带来的交易
func (blockchain *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte//最后的哈希

	for _, tx := range transactions {
		if blockchain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := blockchain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))//查看数据
		lastHash = b.Get([]byte("l"))//取出最后个区块的哈希

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)//创建一个新的区块

	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))//取出索引
		err := b.Put(newBlock.Hash, newBlock.Serialize())//存入数据库
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		blockchain.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}









//交易签名
func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

//交易确认
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)	//查找交易
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
//判断数据库是否存在
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}