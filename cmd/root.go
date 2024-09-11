package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"blockchain/blockchain"
)

var rootCmd = &cobra.Command{
	Use:   "blockchain-cli",
	Short: "CLI tool for managing a blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Blockchain CLI. Use --help to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize blockchain
	chain := blockchain.InitBlockChain()

	// Create the blockchainFlag instance
	chainFlag := NewBlockchainFlag(chain)

	rootCmd.PersistentFlags().Var(
		chainFlag,
		"chain",
		"Blockchain to use",
	)

	// Add commands
	rootCmd.AddCommand(addBlockCmd)
	rootCmd.AddCommand(printChainCmd)
	rootCmd.AddCommand(verifyBlockCmd)
	rootCmd.AddCommand(verifyTransactionCmd)
	rootCmd.AddCommand(listTransactionCmd)

	// Set blockchain instance to commands
	addBlockCmd.PersistentFlags().Var(
		chainFlag,
		"chain",
		"Blockchain to use",
	)
	printChainCmd.PersistentFlags().Var(
		chainFlag,
		"chain",
		"Blockchain to use",
	)
	verifyBlockCmd.PersistentFlags().Var(
		chainFlag,
		"chain",
		"Blockchain to use",
	)
	verifyTransactionCmd.PersistentFlags().Var(
		chainFlag,
		"chain",
		"Blockchain to use",
	)
	listTransactionCmd.PersistentFlags().Var(
		chainFlag,
		"chain",
		"Blockchain to use",
	)

	// Call displayMenu after all commands are added
	displayMenu()
}

// Custom flag type to pass blockchain instance
type blockchainFlag struct {
	Chain *blockchain.BlockChain
}

// NewBlockchainFlag creates a new blockchainFlag instance
func NewBlockchainFlag(chain *blockchain.BlockChain) *blockchainFlag {
	return &blockchainFlag{Chain: chain}
}

// Set sets the value of the blockchainFlag based on the string value
func (b *blockchainFlag) Set(value string) error {
	// Implement logic to set the blockchain.Chain based on the value provided
	b.Chain = initializeBlockchain(value)
	return nil
}

func initializeBlockchain(_ string) *blockchain.BlockChain {
	// Example logic to initialize blockchain from value (adjust as needed)
	chain := blockchain.InitBlockChain()
	return chain
}

// String returns a string representation of the blockchainFlag
func (b *blockchainFlag) String() string {
	if b.Chain == nil {
		return ""
	}
	return "blockchain.Chain"
}

// Type returns the data type of the flag's value
func (b *blockchainFlag) Type() string {
	return "Blockchain"
}

func displayMenu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Menu:")
		fmt.Println("1. Add Block")
		fmt.Println("2. Print Chain")
		fmt.Println("3. Print Transactions of a Block")
		fmt.Println("4. Verify Block")
		fmt.Println("5. Verify Transaction")
		fmt.Println("0. Exit")
		fmt.Print("Enter your choice: ")

		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Remove newline character from the end of input
		choiceStr = strings.TrimSpace(choiceStr)

		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("You chose: Add Block")
			addBlock(reader)
		case 2:
			fmt.Println("You chose: Print Chain")
			printChain()
		case 3:
			fmt.Println("You chose: Print Transactions of a Block")
			printTransactions(reader)
		case 4:
			fmt.Println("You chose: Verify Block")
			verifyBlock(reader)
		case 5:
			fmt.Println("You chose: Verify Transaction")
			fmt.Print("Enter the transaction id: ")
			txId, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			txId = strings.TrimSpace(txId)
			// Simulate command-line argument for verifyTransactionCmd
			os.Args = []string{"app", "verify_transaction", txId}
			err = rootCmd.Execute()
			if err != nil {
				fmt.Println("Error executing command:", err)
			}
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 5.")
		}
	}
}

func addBlock(reader *bufio.Reader) {
	fmt.Print("Enter transactions's data for the new block, enter empty to start sealing the block\n")
	transactionsData := make([]string, 0)
	for {
		fmt.Print("Enter transaction data: ")
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		data = strings.TrimSpace(data)
		if data == "" {
			break
		}
		transactionsData = append(transactionsData, data)
	}

	// Run the addBlock command with provided data
	addBlockCmd.Run(addBlockCmd, transactionsData)
}

// Function to handle printing the blockchain
func printChain() {
	printChainCmd.Run(printChainCmd, nil)
}

func printTransactions(reader *bufio.Reader) {
	fmt.Print("Enter the block index: ")
	indexStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	indexStr = strings.TrimSpace(indexStr)

	// Simulate command-line argument for listTransactionCmd
	os.Args = []string{"app", "list_transaction", indexStr}
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}

// Function to handle verifying a block
func verifyBlock(reader *bufio.Reader) {
	fmt.Print("Enter the block index: ")
	indexStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	indexStr = strings.TrimSpace(indexStr)

	// Simulate command-line argument for verifyBlockCmd
	os.Args = []string{"app", "verify_block", indexStr}
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}
