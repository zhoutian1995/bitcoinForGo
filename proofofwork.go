package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64	//最大的64位整数
)

const targetBits = 24 	//对比位数，实际上是难度

type ProofofWork struct {
	block *Block 	//区块
	target *big.Int 	//存储计算哈希值的整数
}

//创建一个工作量证明的挖矿对象
func NewProofOfWork(block *Block) *ProofofWork{
	target := big.NewInt(1)	//初始目标整数
	target.Lsh(target,uint(256 - targetBits))	//数据转换,向左移动targetBits位数 10000000
	pow := &ProofofWork{block,target}
	return  pow
}
//准备数据进行计算
func (pow * ProofofWork)prepareData(nonce int) []byte{
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,	//上一块哈希
			pow.block.Data,	//当前数据
			InttoHex(pow.block.Timestamp),	//时间十六进制
			InttoHex(int64(targetBits)),	//位数十六进制
			InttoHex(int64(nonce)),	//保存工作量的nonce
		},[]byte{},
		)
	return data
}
//挖矿执行
func  (pow * ProofofWork) Run()(int,[]byte){
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("当前挖矿计算的区块数据%s",pow.block.Data)
	for nonce < maxNonce{
		data := pow.prepareData(nonce)	//准备好的数据
		hash = sha256.Sum256(data)	//计算出哈希
		fmt.Printf("\r%x",hash)	//打印哈希
		hashInt.SetBytes(hash[:])	//获取要对比的数据
		if hashInt.Cmp(pow.target) == -1{	//挖矿校验	比100000000小就是代表左边有几个0
			break
		}else {
			nonce++
		}
	}
	fmt.Println("\n\n")
	return nonce,hash[:]
}
//校验挖矿是不是真的成功
func (pow * ProofofWork) Vaildata() bool{
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)	//准备好数据
	hash := sha256.Sum256(data)	//计算出哈希
	hashInt.SetBytes(hash[:])
	isValid := (hashInt.Cmp(pow.target) == -1)//校验数据 比100000000小就是代表左边有几个0

	return isValid
}

//把对象转化为二进制字节集，可以写入文件
func (block *Block) serialize() []byte{
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
		log.Panic(err)
	}
	return &block
}



