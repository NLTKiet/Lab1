package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var listTransactionCmd = &cobra.Command{
	Use:   "list_transaction [index]",
	Short: "List a transaction in a blockchain",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index := args[0]
		chain := cmd.Flag("chain").Value.(*blockchainFlag).Chain
		i, err := strconv.Atoi(index)
		if err != nil {
			fmt.Println("Invalid block index")
			return
		}
		if i < 0 || i >= len(chain.Blocks) {
			fmt.Println("Block index out of range")
			return
		}

		fmt.Println(chain.Blocks[i].PrintTransactions())
	},
}
