package main

import (
	"fmt"
	"log"
)

func (cli *CLI) getBalance(address string) {
	if !ValidateAddress(address) {
		log.Panic("地址错误")
	}
	bc := NewBlockchain()	//根据地址创建
	UTXOSet := UTXOSet{bc}	//延迟关闭数据库
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))	//提取公钥
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]//除去版本号
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)//查找交易金额

	for _, out := range UTXOs {
		balance += out.Value	//取出金额
	}

	fmt.Printf("查询地址%s金额如下 %d\n", address, balance)
}
