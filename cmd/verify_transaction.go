package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verifyTransactionCmd = &cobra.Command{
	Use:   "verify_transaction [id]",
	Short: "Verify a transaction in a blockchain",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txId := args[0]

		chain := cmd.Flag("chain").Value.(*blockchainFlag).Chain

		valid, err := chain.VerifyTransactionByTxId(txId)
		if err != nil {
			fmt.Println("Error verifying transaction:", err)
			return
		}
		if valid {
			fmt.Println("Transaction is valid")
		} else {
			fmt.Println("Transaction is not valid")
		}
	},
}
