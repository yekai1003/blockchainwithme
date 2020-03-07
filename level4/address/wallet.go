package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

const ChecksumLen = 4

// 钱包结构
type Wallet struct {
	PrivateKey ecdsa.PrivateKey //私钥
	PublicKey  []byte           //公钥
}

// 创建钱包
func NewWallet() *Wallet {
	//随机生成秘钥对
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

//wallet.go
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//获得曲线

	curve := elliptic.P256()
	//生成私钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//利用私钥推导出公钥
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

//计算公钥hash
func HashPubKey(pubKey []byte) []byte {
	//1. 先hash一次
	publicSHA256 := sha256.Sum256(pubKey)
	//2. 计算ripemd160
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

//计算校验和，输入为0x00+公钥hash
func checksum(payload []byte) []byte {

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:ChecksumLen]
}

//生成地址
func (w Wallet) GetAddress() []byte {
	//1. 计算公钥hash
	pubKeyHash := HashPubKey(w.PublicKey)
	//2. 计算校验和
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	//3. 计算base58编码
	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

//交易发送
func (w Wallet) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

//验证地址
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
