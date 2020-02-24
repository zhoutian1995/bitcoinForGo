package main

import (
	"crypto/sha256"
	"encoding/gob"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10	//矿工挖矿给予的奖励

//输入
type TXInput struct {
	Txid []byte	//保存了交易的id
	Vout int	//保存该交易一个output索引
	ScriptSig string	//仅保存了一个任意的用户定义的钱包
}

//输出
type TXOutput struct {
	Value int	//一定量的比特币
	ScriptPubkey string	//一个锁定脚本
}

//是否可以解锁输出
func (out *TXOutput)CanBeunlockedWith(unlockingData string) bool {
	return out.ScriptPubkey == unlockingData 	//判断是否可以解锁
}

//交易，编号，输入，输出
type Transaction struct {
	ID []byte
	Vin []TXInput
	Vout []TXOutput
}

//检查交易事物是否为coinbase，就是挖矿得来的奖励币
func (tx *Transaction)IsCoinBase ()bool{
	//有输入,无交易内容，无输出
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

//从二进制数据中，设置交易的ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer//设置交易
	var hash[32] byte
	enc := gob.NewEncoder(&encoded)//解码对象
	err := enc.Encode(tx)	//解码
	if err != nil{
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())	//计算哈希
	tx.ID = hash[:]
}
//检查地址是否启动事物
func (input *TXInput) CanUnlockOutPutWith(unlockingData string) bool{
	return input.ScriptSig == unlockingData
}
//挖矿交易
func NewCoinBaseTX(to , data string) *Transaction{
	if data == ""{
		data = fmt.Sprintf("挖矿奖励给%s",to)	//挖矿信息赋值
	}

	txin := TXInput{[]byte{},-1,data}		//设置输入
	txout := TXOutput{subsidy,to}	//设置输出
	tx := Transaction{nil,[]TXInput{txin},[]TXOutput{txout}}		//设置交易信息

	return &tx	//返回交易信息
}
//转账交易
func NewUTXOTransaciton (from, to string, amount int, bc *BlockChain) *Transaction{
	var inputs[] TXInput	//输入
	var outputs [] TXOutput	//输出
	acc,vaildOutputs := bc.FindSpendableOutputs(from,amount)

	if acc < amount{
		log.Panic("交易金额不足")
	}

	for txid,outs := range vaildOutputs{	//循环遍历有效输出
		txID,err := hex.DecodeString(txid)	//解码
		if err != nil{
			log.Panic(err)	//处理错误
		}
		for _,out:=range outs{
			input := TXInput{txID,out,from}	//输入交易
			inputs = append(inputs,input)	//输出交易
		}
	}
	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
