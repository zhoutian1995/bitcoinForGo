package main

func main(){
	block := NewBlockChain()//创建区块链
	defer block.db.Close()	//延迟关闭数据
	cil := CLI{block}//创建命令行
	cil.Run()//开启
}
