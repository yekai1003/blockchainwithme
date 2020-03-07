package main

func main1() {
	//1. 创世块初始化区块链
	bc := CreateBlockchain()
	defer bc.db.Close()
	//2. 获取叶开余额
	bc.getBalance("yekai")

}

func main() {
	//1. 创世块初始化区块链
	bc := CreateBlockchain()
	defer bc.db.Close()
	//2. 获取叶开余额
	bc.getBalance("yekai")
	//3. 发送5个给傅红雪
	bc.send("yekai", "fuhongxue", 8, "拿去生活吧")
	//4. 查询余额
	bc.getBalance("yekai")
	bc.getBalance("fuhongxue")
}
