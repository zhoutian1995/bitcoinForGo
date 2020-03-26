package main

import (
	"bytes"
	"math/big"
)
//字母表格，最终展示的字符
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)	//输入的数据存入二进制

	base := big.NewInt(int64(len(b58Alphabet)))	//创建了一个大数
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)//求余数
		result = append(result, b58Alphabet[mod.Int64()])
	}

	if input[0] == 0x00 {
		result = append(result, b58Alphabet[0])
	}

	ReverseBytes(result)	//字节集合翻转

	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)	//初始化为0

	for _, b := range input {
		charIndex := bytes.IndexByte(b58Alphabet, b)	//字母表格
		result.Mul(result, big.NewInt(58))	//乘法
		result.Add(result, big.NewInt(int64(charIndex)))//加法
	}

	decoded := result.Bytes()//解码

	if input[0] == b58Alphabet[0] {
		decoded = append([]byte{0x00}, decoded...)	//叠加
	}

	return decoded
}
