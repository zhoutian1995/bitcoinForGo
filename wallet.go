package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)	//钱包版本
const walletFile = "wallet.dat"	//钱包文件
const addressChecksumLen = 4	//检测地址长度


type Wallet struct {
	PrivateKey ecdsa.PrivateKey	//私钥  钱包的权限
	PublicKey  []byte	//公钥  收款地址
}
//创建一个钱包

func NewWallet() *Wallet {
	private, public := newKeyPair()	//创建公钥私钥
	wallet := Wallet{private, public}//创建钱包

	return &wallet
}
//创建公钥私钥
//抓取钱包的地址
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)	//取得哈希值

	versionedPayload := append([]byte{version}, pubKeyHash...)	//
	checksum := checksum(versionedPayload)//检测版本和公钥

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address	//返回钱包地址
}


//公钥哈希处理
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)	//处理公钥

	RIPEMD160Hasher := ripemd160.New()		//创建一个哈希
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])	//写入处理
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)	//叠加运算

	return publicRIPEMD160
}

//校验钱包的地址
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))	//解码
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]	//取得版本
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

//公钥的校验
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()	//创建加密算法
	private, err := ecdsa.GenerateKey(curve, rand.Reader)	//生产私钥
	if err != nil {
		log.Panic(err)
	}
	//生成公钥
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
