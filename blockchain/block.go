package blockchain

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	MerkleTree    *MerkleTree
}

type Transaction struct {
	Id   string
	Data []byte
}

func Serialize[D any](data D) []byte {
	var serializedBuffer bytes.Buffer
	encoder := gob.NewEncoder(&serializedBuffer)
	encoder.Encode(data)
	return serializedBuffer.Bytes()
}

func Deserialize[T any](buffer []byte) *T {
	var data T
	decoder := gob.NewDecoder(bytes.NewReader(buffer))
	decoder.Decode(&data)
	return &data
}

func (t *Transaction) CalculateHash() ([]byte, error) {
	serializedTransaction := Serialize(t)

	h := sha256.New()
	if _, err := h.Write(serializedTransaction); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (t *Transaction) GetId() string {
	return t.Id
}

func (b *Block) HashTransactions() ([]byte, error) {
	var transactions []NodeData
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx)
	}

	mTree, err := newMerkleTree(transactions)
	if err != nil {
		return nil, err
	}

	// This is a temporary solution to store the Merkle tree in the block
	b.MerkleTree = mTree

	return mTree.RootNode.Hash, nil
}

func (b *Block) calculateHash() ([]byte, error) {
	hashedTransaction, err := b.HashTransactions()
	if err != nil {
		return nil, err
	}

	record := []byte(strconv.Itoa(int(b.Timestamp)))
	record = append(record, hashedTransaction...)
	record = append(record, b.PrevBlockHash...)
	hash := sha256.Sum256(record)
	return hash[:], nil
}

func (b *Block) SetHash() error {
	hash, err := b.calculateHash()
	if err != nil {
		return err
	}
	b.Hash = hash
	return nil
}

func (b *Block) VerifyBlock() (bool, error) {
	hash, err := b.calculateHash()
	if err != nil {
		return false, err
	}

	return bytes.Equal(b.Hash, hash[:]), nil
}

func CreateBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().UnixMicro(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}

	block.SetHash()
	return block
}

func CreateTransaction(data string) *Transaction {
	tx := &Transaction{
		Id:   strconv.Itoa(int(time.Now().UnixMicro())) + "+" + randomHex(8),
		Data: []byte(data),
	}

	// Simulate some processing time
	time.Sleep(100 * time.Microsecond)

	return tx
}

func (b *Block) PrintTransactions() string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Id", "Data"})
	for _, tx := range b.Transactions {
		t.AppendRow(table.Row{tx.Id, string(tx.Data)})
	}
	return t.Render()
}

func (b *Block) AsString() string {
	return "Timestamp: " + strconv.FormatInt(b.Timestamp, 10) + "\n" +
		"Contains: " + strconv.Itoa(len(b.Transactions)) + " transactions\n" +
		"Prev. hash: " + hex.EncodeToString(b.PrevBlockHash) + "\n" +
		"Hash: " + hex.EncodeToString(b.Hash) + "\n"
}

func (b *Block) VerifyBlockTransaction(txId string) (bool, error) {
	// Find the transaction index in the block
	txIndex := b.FindTransactionIndexById(txId)
	if txIndex == -1 {
		return false, nil
	}

	valid, err := b.MerkleTree.VerifyNodeDataByLeafIndex(0)

	if err != nil {
		return false, err
	}

	return valid, nil
}

func (b *Block) FindTransactionIndexById(id string) int {
	for index, tx := range b.Transactions {
		if tx.GetId() == id {
			return index
		}
	}
	return -1
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func SplitTransactionId(id string) (int64, string) {
	splitedId := strings.Split(id, "+")
	timestamp, err := strconv.ParseInt(splitedId[0], 10, 64)
	if err != nil {
		panic(err)
	}
	return timestamp, id[11:]
}
