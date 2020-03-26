package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

//整数转化16进制
func InttoHex(num int64) []byte {
	buff := new(bytes.Buffer)	//开辟内存，存储字节集
	err := binary.Write(buff,binary.BigEndian,num)//num转化字节集写入
	if err != nil{
		log.Panic(err)
	}

	return buff.Bytes()	//返回字节集合
}

//字节反序
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}