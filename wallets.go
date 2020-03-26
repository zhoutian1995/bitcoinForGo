package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"log"
)



type Wallets struct {
	Wallets map[string]*Wallet	//一个字符串对应一个钱包
}
//创建一个钱包，并抓取已经存在的钱包
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile()

	return &wallets, err
}

//创建一个钱包
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()	//创建一个钱包
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet	//保存钱包

	return address
}

//抓取所有钱包
func (ws *Wallets) GetAddresses() []string {
	var addresses []string	//所有钱包地址

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses	//返回所有钱包地址
}
//抓取一个钱包
func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

//从文件中读取钱包
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	//读取文件二进制并且解析
	var wallets Wallets	//钱包集合
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}
//钱包保存到文件
func (ws Wallets) SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())	//注册加密算法

	encoder := gob.NewEncoder(&content)	//解码
	err := encoder.Encode(ws)//解码
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}



