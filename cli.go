package main

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

//命令行接口
type CLI struct {
}

//用法
func (cli *CLI) printUsage(){
	fmt.Println("用法如下：")
	fmt.Println("  createblockchain -address ADDRESS #创建一个区块链并将创世区块奖励发送给ADDRESS")
	fmt.Println("  createwallet #生成新的密钥对并将其保存到钱包文件中")
	fmt.Println("  getbalance -address ADDRESS #取得ADDRESS的余额")
	fmt.Println("  listaddresses #列出钱包文件的所有地址")
	fmt.Println("  printchain #打印区块链的所有块")
	fmt.Println("  reindexutxo #重建UTXO集合")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT #从FROM地址向TO发送AMOUNT个硬币")
}

func (cli *CLI) validateArgs(){
	if len(os.Args) < 2{
		cli.printUsage()	//显示区块
		os.Exit(1)
	}
}

func (cli *CLI)Run(){
	cli.validateArgs()	//校验
	//处理命令行参数
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

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
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
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
		cli.createBlockchain(*createBlockchainAddress)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}