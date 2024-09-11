package cmd

import (
	"github.com/spf13/cobra"
)

var printChainCmd = &cobra.Command{
	Use:   "printchain",
	Short: "Print all blocks in the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		chain := cmd.Flag("chain").Value.(*blockchainFlag).Chain
		chain.ListBlocks()
	},
}
