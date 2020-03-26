package main

import "bytes"

//输入
type TXInput struct {
	Txid []byte	//保存了交易的id
	Vout int	//保存该交易一个output索引
	Signature []byte//签名
	PubKey	[]byte	//公钥
}


//key检测一下地址与交易
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}