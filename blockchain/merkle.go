package blockchain

import (
	"bytes"
	"crypto/sha256"
	"errors"
)

type NodeData interface {
	CalculateHash() ([]byte, error)
	GetId() string
}

type MerkleTree struct {
	RootNode *Node
	Leafs    []*Node
}

type Node struct {
	Parent       *Node
	Left         *Node
	Right        *Node
	IsLeaf       bool
	IsDuplicated bool
	Data         NodeData
	Hash         []byte
}

func newMerkleTree(data []NodeData) (*MerkleTree, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot construct merkle tree with no data")
	}
	var leafs []*Node
	for _, d := range data {
		hash, err := d.CalculateHash()
		if err != nil {
			return nil, err
		}
		leafs = append(leafs, &Node{Data: d, Hash: hash, IsLeaf: true})
	}

	// If the number of leafs is odd, duplicate the last leaf
	if len(leafs)%2 != 0 {
		leaf := leafs[len(leafs)-1]
		leafs = append(leafs, &Node{Data: leaf.Data, Hash: leaf.Hash, IsLeaf: true, IsDuplicated: true})
	}

	root := buildIntermediate(leafs)

	return &MerkleTree{RootNode: root, Leafs: leafs}, nil
}

func calculateCombinedHash(left, right []byte) []byte {
	hash := sha256.Sum256(append(left, right...))

	return hash[:]
}

func buildIntermediate(prevLayer []*Node) *Node {
	var currentLayer []*Node
	for i := 0; i < len(prevLayer); i += 2 {
		// If there is only one node left, duplicate it
		var left, right int = i, i + 1
		if right == len(prevLayer) {
			right = i
		}

		hash := calculateCombinedHash(prevLayer[left].Hash, prevLayer[right].Hash)
		parent := &Node{Left: prevLayer[left], Right: prevLayer[right], Hash: hash[:]}
		prevLayer[left].Parent = parent
		prevLayer[right].Parent = parent
		currentLayer = append(currentLayer, parent)

		// If there are only two nodes in the current layer, return the parent (root node)
		if len(prevLayer) == 2 {
			return parent
		}
	}

	return buildIntermediate(currentLayer)
}

func (n *Node) CalculateNodeHash() []byte {
	if n.IsLeaf {
		return n.Hash
	}

	hash := calculateCombinedHash(n.Left.CalculateNodeHash(), n.Right.CalculateNodeHash())

	return hash[:]
}

func (t *MerkleTree) VerifyNodeDataByLeafIndex(index int) (bool, error) {
	if index < 0 || index >= len(t.Leafs) {
		return false, errors.New("index out of range")
	}

	currentParent := t.Leafs[index].Parent
	for currentParent != nil {
		rightHash := currentParent.Right.CalculateNodeHash()
		leftHash := currentParent.Left.CalculateNodeHash()
		combinedHash := calculateCombinedHash(leftHash, rightHash)

		if !bytes.Equal(currentParent.Hash, combinedHash) {
			return false, nil
		}
		currentParent = currentParent.Parent
	}
	return true, nil
}
