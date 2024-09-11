package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type BlockChain struct {
	Blocks []*Block
}

func (c *BlockChain) AddBlock(transactions []*Transaction) {
	// Check if the blockchain is empty (only contains the genesis block)
	if len(c.Blocks) == 0 {
		genesisBlock := createGenesisBlock()
		c.Blocks = append(c.Blocks, genesisBlock)
	} else {
		prevBlock := c.Blocks[len(c.Blocks)-1]
		newBlock := CreateBlock(transactions, prevBlock.Hash)
		c.Blocks = append(c.Blocks, newBlock)
	}
}

// Create the first block in the chain
func createGenesisBlock() *Block {
	genesis := CreateBlock([]*Transaction{CreateTransaction("Genesis Block")}, nil)
	return genesis
}

func InitBlockChain() *BlockChain {
	return &BlockChain{Blocks: []*Block{createGenesisBlock()}}
}

func (c *BlockChain) ListBlocks() {
	fmt.Printf("Printing all %d blocks in the chain...\n", len(c.Blocks))
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Timestamp", "Hash", "PrevBlockHash", "Transaction Count"})
	for i, block := range c.Blocks {
		t.AppendRow(table.Row{i,
			block.Timestamp,
			hex.EncodeToString(block.Hash),
			hex.EncodeToString(block.PrevBlockHash),
			len(block.Transactions),
		})
	}
	t.Render()
}

// Find the block that contains a transaction with the given ID
// by iterating over the blocks in the chain
// and comparing the block's timestamp with the transaction's timestamp
// The first block which its timestamp is greater than the transaction's timestamp
// contains the transaction
func (c *BlockChain) FindBlockByTransactionId(id string) *Block {
	for _, block := range c.Blocks {
		transactionTimestamp, _ := SplitTransactionId(id)
		if block.Timestamp >= transactionTimestamp {
			return block
		}
	}
	return nil
}

func (c *BlockChain) VerifyTransactionByTxId(txId string) (bool, error) {
	// Find the block that contains the transaction
	block := c.FindBlockByTransactionId(txId)
	if block == nil {
		return false, fmt.Errorf("transaction %s not found in the blockchain", txId)
	}
	return block.VerifyBlockTransaction(txId)
}
