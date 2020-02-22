package main

import (
	"fmt"
	"strconv"
)

func main(){
	fmt.Println("digging bitcoin...")
	bc := NewBlockchain()	//创建一个区块链
	bc.AddBlock("add block A")
	bc.AddBlock("add block B")
	bc.AddBlock("add block C")

	for i,block := range bc.blocks{
		fmt.Printf("%d prev hash = %x\n",i,block.PrevBlockHash)
		fmt.Printf("%d data = %s\n",i,block.Data)
		fmt.Printf("%d this hash= %x\n",i,block.Hash)
		pow := NewProofOfWork(block)//校验工作量
		fmt.Printf("pow %s\n",strconv.FormatBool(pow.Vaildata()))
		fmt.Println()
	}
}
