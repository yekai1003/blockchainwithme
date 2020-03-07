package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func main() {
	//打开数据库
	db, _ := bolt.Open("my.db", 0600, nil)
	defer db.Close()
	//插入数据库数据
	db.Update(func(tx *bolt.Tx) error {
		//创建bucket
		bucket, _ := tx.CreateBucket([]byte("bucket1"))
		//设置key-val
		bucket.Put([]byte("name"), []byte("yekai"))
		return nil
	})
	//查询数据库数据
	db.View(func(tx *bolt.Tx) error {
		//获取bucket
		bucket := tx.Bucket([]byte("bucket1"))
		//获取key-val
		val := bucket.Get([]byte("name"))
		fmt.Println(string(val))
		return nil
	})
}
