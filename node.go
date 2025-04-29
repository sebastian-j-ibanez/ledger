package ledger

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"slices"
	"time"
)

// Node represents an element in the Ledger.
// Nodes should typically be constructed via Ledger.AddNode().
type Node struct {
	id        uint64
	timestamp time.Time
	hash      []byte
	prevHash  []byte
	data      []byte
}

// Create a new Node.
// Returns an error if there is an issue computing the Node hash.
func NewNode(id uint64, prevHash []byte, data []byte) (Node, error) {
	n := Node{
		id:        id,
		timestamp: time.Now(),
		hash:      nil,
		prevHash:  prevHash,
		data:      data,
	}

	var err error
	n.hash, err = n.ComputeHash()
	if err != nil {
		msg := fmt.Sprintf("unable to create node %d: %s", n.id, err.Error())
		return Node{}, errors.New(msg)
	}

	return n, nil
}

// Compute a SHA256 hash of the Node.
// Returns error if there is an issue serializing Node fields.
func (n *Node) ComputeHash() ([]byte, error) {
	hash := sha256.New()
	var buf bytes.Buffer

	// Write node fields to binary buffer.
	if err := binary.Write(&buf, binary.LittleEndian, n.id); err != nil {
		msg := fmt.Sprintf("unable to encode id as binary: %s", err.Error())
		return nil, errors.New(msg)
	}

	if err := binary.Write(&buf, binary.LittleEndian, n.timestamp.UnixNano()); err != nil {
		msg := fmt.Sprintf("unable to encode timestamp as binary: %s", err.Error())
		return nil, errors.New(msg)
	}

	if err := binary.Write(&buf, binary.LittleEndian, n.prevHash); err != nil {
		msg := fmt.Sprintf("unable to encode previous hash as binary: %s", err.Error())
		return nil, errors.New(msg)
	}

	if _, err := buf.Write(n.data); err != nil {
		msg := fmt.Sprintf("unable to encode data as binary: %s", err.Error())
		return nil, errors.New(msg)
	}

	// Make hash from the binary buffer.
	hash.Write(buf.Bytes())
	return hash.Sum(nil), nil
}

// Validate a Node's hash.
// Returns an error if there is an issue computing the hash.
func (n *Node) ValidateHash() (bool, error) {
	currentHash, err := n.ComputeHash()
	if err != nil {
		msg := fmt.Sprintf("unable to compute hash: %s", err.Error())
		return false, errors.New(msg)
	}
	valid := slices.Compare(n.hash, currentHash) == 0

	return valid, nil
}

// Get the encoded byte slice representation of the Node data.
// Data is encoded in binary format, using encoding/binary.
func (n *Node) GetData() []byte {
	return n.data
}

// Return Node as readable string.
func (n Node) String() string {
	return fmt.Sprintf("id: %d\ntimestamp: %v\nhash: %v\nprev hash: %v\ndata: %v\n",
		n.id, n.timestamp, bytesToDecimalString(n.hash), bytesToDecimalString(n.prevHash), n.data)
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
