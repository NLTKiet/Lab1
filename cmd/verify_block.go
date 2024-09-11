package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var verifyBlockCmd = &cobra.Command{
	Use:   "verify_block [index]",
	Short: "Verify the integrity of a block in the blockchain",
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
		block := chain.Blocks[i]
		valid, err := block.VerifyBlock()
		if err != nil {
			fmt.Println("Error verifying block:", err)
			return
		}
		if valid {
			fmt.Println("Block is valid")
		} else {
			fmt.Println("Block is not valid")
		}
	},
}
