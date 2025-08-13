package ledger

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"slices"
	"time"
)

// Node represents an element in the Ledger.
// Nodes should typically be constructed via Ledger.AddNode().
type Node struct {
	Id        uint64
	Timestamp time.Time
	Hash      []byte
	PrevHash  []byte
	Data      []byte
}

// Create a new Node.
// Returns an error if there is an issue computing the Node hash.
func newNode(id uint64, prevHash []byte, data []byte) Node {
	n := Node{
		Id:        id,
		Timestamp: time.Now(),
		Hash:      nil,
		PrevHash:  prevHash,
		Data:      data,
	}
	n.Hash = n.computeHash()
	return n
}

// Compute a SHA256 hash of the Node.
// Returns error if there is an issue serializing Node fields.
func (n *Node) computeHash() []byte {
	// Encode node data to gob
	var encBuffer bytes.Buffer
	encoder := gob.NewEncoder(&encBuffer)
	encoder.Encode(n.Id)
	encoder.Encode(n.Timestamp.UnixNano())
	encoder.Encode(n.PrevHash)
	encoder.Encode(n.Data)

	// Make hash from the encoded gob data.
	hash := sha256.New()
	hash.Write(encBuffer.Bytes())
	return hash.Sum(nil)
}

// Validate a Node's hash.
// Returns an error if there is an issue computing the hash.
func (n *Node) ValidHash() bool {
	currentHash := n.computeHash()
	validHash := slices.Compare(n.Hash, currentHash) == 0
	return validHash
}

// Get the encoded byte slice representation of the Node data.
// Data is encoded in binary format, using encoding/binary.
func (n *Node) GetData() []byte {
	return n.Data
}

// Return Node as readable string.
func (n Node) String() string {
	return fmt.Sprintf("id: %d\ntimestamp: %v\nhash: %v\nprev hash: %v\ndata: %v\n",
		n.Id, n.Timestamp, bytesToDecimalString(n.Hash), bytesToDecimalString(n.PrevHash), n.Data)
}

// Convert a byte slice to a string of literal numeric values.
// Primarily used to print hashes.
func bytesToDecimalString(data []byte) string {
	var buf bytes.Buffer
	for _, elem := range data {
		fmt.Fprintf(&buf, "%d", elem)
	}
	return buf.String()
}
