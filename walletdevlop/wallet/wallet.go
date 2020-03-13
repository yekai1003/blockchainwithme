package wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"

	"walletdevlop/hdkeystore"

	"github.com/howeyc/gopass"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

//BIP路径
const defaultDerivationPath = "m/44'/60'/0'/0/1"

//钱包结构体
type HDWallet struct {
	Address    common.Address
	HdKeyStore *hdkeystore.HDkeyStore
}

func create_mnemonic() (string, error) {
	//Entropy 生成，注意传入值y=32*X，并且128<=y<=256
	b, err := bip39.NewEntropy(128)
	if err != nil {
		fmt.Println("failed to NewEntropy:", err, b)
		return "", err
	}
	//生成助记词
	nm, err := bip39.NewMnemonic(b)
	if err != nil {
		fmt.Println("failed to NewMnemonic:", err)
		return "", err
	}
	fmt.Println(nm)
	return nm, nil
}

func NewKeyFromMnemonic(mn string) (*ecdsa.PrivateKey, error) {
	//1. 先推导路径
	path, err := accounts.ParseDerivationPath(defaultDerivationPath)
	if err != nil {
		panic(err)
	}
	//2. 获得种子
	seed, err := bip39.NewSeedWithErrorChecking(mn, defaultDerivationPath)
	//3. 获得主key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println("Failed to NewMaster", err)
		return nil, err
	}
	//4. 推导私钥
	privateKey, err := DerivePrivateKey(path, masterKey)

	return privateKey, nil
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

//钱包构造函数
func NewWallet(keypath string) (*HDWallet, error) {
	//1. 创建助记词
	mn, err := create_mnemonic()
	if err != nil {
		fmt.Println("Failed to NewWallet", err)
		return nil, err
	}
	//2. 推导私钥
	privateKey, err := NewKeyFromMnemonic(mn)
	if err != nil {
		fmt.Println("Failed to NewKeyFromMnemonic", err)
		return nil, err
	}
	//3. 获取地址
	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		fmt.Println("Failed to DerivePublicKey", err)
		return nil, err
	}
	//利用公钥推导私钥
	address := crypto.PubkeyToAddress(*publicKey)
	//4. 创建keystore
	hdks := hdkeystore.NewHDkeyStore(keypath, privateKey)
	//5. 创建钱包
	return &HDWallet{address, hdks}, nil
}

func LoadWallet(filename, datadir string) (*HDWallet, error) {
	hdks := hdkeystore.NewHDkeyStoreNoKey(datadir)
	//解决密码问题
	fmt.Println("Please input password for:", filename)
	pass, _ := gopass.GetPasswd()
	//filename也是账户地址
	fromaddr := common.HexToAddress(filename)
	_, err := hdks.GetKey(fromaddr, hdks.JoinPath(filename), string(pass))
	if err != nil {
		log.Panic("Failed to GetKey ", err)
	}
	return &HDWallet{fromaddr, hdks}, nil
}

func (w HDWallet) StoreKey(pass string) error {
	filename := w.HdKeyStore.JoinPath(w.Address.Hex())
	return w.HdKeyStore.StoreKey(filename, &w.HdKeyStore.Key, pass)
}
