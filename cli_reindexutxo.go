package main

import "fmt"

func (cli *CLI) reindexUTXO() {
	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf(" 有 %d 交易在UTXO 集合中\n", count)
}
