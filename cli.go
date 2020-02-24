package main

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"strconv"
)

//命令行接口
type CLI struct {
	blockchain *BlockChain
}

func (cli *CLI) createBlockChain(address string) {
	bc := CreateBlockchain(address)
	bc.db.Close()
	fmt.Println("创建成功")
}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockChain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value	//查找UTXO中所有没有消费的金额，再进行累加
	}

	fmt.Printf("查询的金额如下 %s : %d\n", address, balance)
}

//用法
func (cli *CLI) printUsage(){
	fmt.Println("用法如下：")
	fmt.Println("getbalance -address 你输入的地址 根据地址查询金额")
	fmt.Println("createblockchain -address 你输入的地址 根据地址创建区块链")
	fmt.Println("send -from From -to To -amount Amount 转账")
	fmt.Println("showBlock 显示区块")
}

func (cli *CLI) validateArgs(){
	if len(os.Args) < 2{
		cli.printUsage()	//显示区块
		os.Exit(1)
	}
}

func (cli *CLI)showBlockChain(){
	bc := NewBlockChain("")
	defer bc.db.Close()

	bci := bc.Iterator()	//迭代

	for {
		block := bci.Next()

		fmt.Printf("上一块哈希: %x\n", block.PrevBlockHash)
		fmt.Printf("当前Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Vaildate()))	//工作量证明
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {	//创世区块
			break
		}
	}
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockChain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaciton(from, to, amount, bc)	//转账
	bc.MineBlock([]*Transaction{tx})	//挖矿
	fmt.Println("交易成功")
}

func (cli *CLI)Run(){
	cli.validateArgs()	//校验
	//处理命令行参数
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	showChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "查询的金额")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "查询的地址")
	sendFrom := sendCmd.String("from", "", "谁给的")
	sendTo := sendCmd.String("to", "", "给谁的")
	sendAmount := sendCmd.Int("amount", 0, "金额")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "showBlock":
		err := showChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockChain(*createBlockchainAddress)//创建区块链
	}

	if showChainCmd.Parsed() {
		cli.showBlockChain()	//显示区块链
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}