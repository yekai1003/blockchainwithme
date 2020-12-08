package main

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

type Key struct {
	// Version 4 "random" for unique id not derived from key data
	Id uuid.UUID
	// 地址
	Address common.Address
	// 私钥
	PrivateKey *ecdsa.PrivateKey
}

type Transaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

// keyStore接口
type keyStore interface {
	// 解析文件为key
	GetKey(addr common.Address, filename string, auth string) (*Key, error)
	// 存储key
	StoreKey(filename string, k *Key, auth string) error
	// 在文件前加上keystore路径
	JoinPath(filename string) string
}

func create_mnemonic() {
	//Entropy 生成，注意传入值y=32*X，并且128<=y<=256
	b, err := bip39.NewEntropy(128)
	if err != nil {
		log.Panic("failed to NewEntropy:", err, b)
	}

	fmt.Println(b)

	//生成助记词
	nm, err := bip39.NewMnemonic(b)
	if err != nil {
		log.Panic("failed to NewMnemonic:", err)
	}
	fmt.Println(nm)
}

func DeriveAddressFromMnemonic() {
	//1. 先推导路径
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/1")
	if err != nil {
		panic(err)
	}
	//2. 获得种子
	nm := "cargo emotion slot dentist client hint will penalty wrestle divide inform ranch"

	seed, err := bip39.NewSeedWithErrorChecking(nm, "")
	//3. 获得主key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Failed to NewMaster", err)
		return
	}
	//4. 推导私钥
	privateKey, err := DerivePrivateKey(path, masterKey)
	//5. 推导公钥
	publicKey, err := DerivePublicKey(privateKey)

	//6. 利用公钥推导私钥
	address := crypto.PubkeyToAddress(*publicKey)

	fmt.Println(address.Hex())
}

func DerivePrivateKey(path accounts.DerivationPath, masterKey *hdkeychain.ExtendedKey) (*ecdsa.PrivateKey, error) {
	var err error
	key := masterKey
	for _, n := range path {
		//按照路径迭代获得最终key
		key, err = key.Child(n)
		if err != nil {
			return nil, err
		}
	}
	//将key转换为ecdsa私钥
	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}

	return privateKeyECDSA, nil
}

func DerivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}
	return publicKeyECDSA, nil
}

func main() {
	//create_mnemonic()
	DeriveAddressFromMnemonic()
}
