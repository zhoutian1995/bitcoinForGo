package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

//命令行接口
type CLI struct {
	blockchain *BlockChain
}
//用法
func (cli *CLI) printUsage(){
	fmt.Println("用法如下：")
	fmt.Println("addBlock 向区块链增加块")
	fmt.Println("showBlock 显示区块")
}

func (cli *CLI) validateArgs(){
	if len(os.Args) < 2{
		cli.printUsage()	//显示区块
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string){
	cli.blockchain.AddBlock(data)	//增加一个区鲁哀
	fmt.Println("区块增加成功")
}

func (cli *CLI)showBlockChain(){
	bci := cli.blockchain.Iterator()//创建循环迭代器
	for{
		block := bci.Next()//取得下一个区块
		fmt.Printf("prev hash = %x\n",block.PrevBlockHash)
		fmt.Printf("data = %s\n",block.Data)
		fmt.Printf("this hash= %x\n",block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("pow %s",strconv.FormatBool(pow.Vaildate()))
		fmt.Println()
		
		if len(block.PrevBlockHash) == 0{//遇到创世区块
			break
		}
	}
}


func (cli *CLI)Run(){
	cli.validateArgs()//校验
	
	//处理命令行参数
	addblockcmd := flag.NewFlagSet("addblck",flag.ExitOnError)
	showchaincmd := flag.NewFlagSet("showchain",flag.ExitOnError)
	addbBlockData := addblockcmd.String("data","","Block data")
	switch os.Args[1] {
	case "addblock":
		err := addblockcmd.Parse(os.Args[2:])//解析参数
		if err != nil{
			fmt.Printf("1223")
		}
	case "showchain":
		err := showchaincmd.Parse(os.Args[2:])//解析参数
		if err != nil{
			fmt.Printf("1223")
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addblockcmd.Parsed(){
		if *addbBlockData == ""{
			addblockcmd.Usage()
			os.Exit(1)
		}else{
			cli.addBlock(*addbBlockData)//增加区块
		}
	}

	if showchaincmd.Parsed() {
		cli.showBlockChain()//显示区块链
	}
}