package hdkeystore

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type HDkeyStore struct {
	keysDirPath string       //文件所在路径
	scryptN     int          //生成加密文件的参数N
	scryptP     int          //生成加密文件的参数P
	Key         keystore.Key //keystore对应的key
}
type UUID []byte

//全局加密随机阅读器
var rander = rand.Reader

//生成UUID
func NewRandom() UUID {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	//版本4规范处理与变形
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid
}

//给出一个生成HDkeyStore对象的方法,通过privatekey生成
func NewHDkeyStore(path string, privateKey *ecdsa.PrivateKey) *HDkeyStore {
	//获得UUID
	uuid := []byte(NewRandom())
	if privateKey == nil {
		return &HDkeyStore{
			keysDirPath: path,
			scryptN:     keystore.LightScryptN,
			scryptP:     keystore.LightScryptP,
			Key:         keystore.Key{},
		}
	}
	key := keystore.Key{
		Id:         uuid,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}
	return &HDkeyStore{
		keysDirPath: path,
		scryptN:     keystore.LightScryptN,
		scryptP:     keystore.LightScryptP,
		Key:         key,
	}
}

func NewHDkeyStoreNoKey(path string) *HDkeyStore {
	return &HDkeyStore{
		keysDirPath: path,
		scryptN:     keystore.LightScryptN,
		scryptP:     keystore.LightScryptP,
		Key:         keystore.Key{},
	}
}

func (ks *HDkeyStore) GetKey(addr common.Address, filename, auth string) (*keystore.Key, error) {
	// 读取json文件
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 解析key
	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	// 验证
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	ks.Key = *key
	return key, nil
}

//存储key为keystore文件
func (ks HDkeyStore) StoreKey(filename string, key *keystore.Key, auth string) error {
	//编码key为json
	keyjson, err := keystore.EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	if err != nil {
		return err
	}
	//写入文件
	return WriteKeyFile(filename, keyjson)
}

func (ks HDkeyStore) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}

	return filepath.Join(ks.keysDirPath, filename)
}

func WriteKeyFile(file string, content []byte) error {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()

	return os.Rename(f.Name(), file)
}

func (ks HDkeyStore) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {

	// Sign the transaction and verify the sender to avoid hardware fault surprises
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, ks.Key.PrivateKey)
	if err != nil {
		return nil, err
	}

	//验证 remix->msg
	msg, err := signedTx.AsMessage(types.HomesteadSigner{})
	if err != nil {
		return nil, err
	}

	sender := msg.From()
	if sender != address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", address.Hex(), sender.Hex())
	}

	return signedTx, nil
}

func (ks HDkeyStore) NewTransactOpts() *bind.TransactOpts {
	return bind.NewKeyedTransactor(ks.Key.PrivateKey)
}
