package cmd

import (
	"blockchain/blockchain"
	"fmt"

	"github.com/spf13/cobra"
)

var addBlockCmd = &cobra.Command{
	Use:   "add_block ...transaction_data",
	Short: "Add a new block to the blockchain",
	Args:  cobra.MinimumNArgs(1), // Require at least one argument for [data]
	Run: func(cmd *cobra.Command, args []string) {
		chain := cmd.Flag("chain").Value.(*blockchainFlag).Chain

		transactions := make([]*blockchain.Transaction, 0)
		for _, txData := range args {
			tx := blockchain.CreateTransaction(txData)
			transactions = append(transactions, tx)
		}

		block := blockchain.CreateBlock(transactions, chain.Blocks[len(chain.Blocks)-1].Hash)
		chain.Blocks = append(chain.Blocks, block)

		fmt.Println("Block added successfully!")
	},
}
