// Ledger is a library for a blockchain ledger data structure.
//
// The Ledger is a tool that provides data validation and integrity.
// Each Node of the Ledger contains a unique hash.
// The hash can be validated to ensure data integrity.
// This library does not contain any external dependencies.
package ledger

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// The Ledger is a data structure that stores data with a hash.
// Each element in the Ledger can be cryptographically validated.
type Ledger struct {
	nodes []Node
}

// Create a new Ledger.
func NewLedger() *Ledger {
	return &Ledger{
		nodes: make([]Node, 0),
	}
}

// Add a new Node to the Ledger.
// Returns error if there is an issue generating the Node hash.
func (l *Ledger) AddNode(data []byte) error {
	id := uint64(len(l.nodes))

	prevHash := []byte(nil)
	if len(l.nodes) > 1 {
		prevHash = l.nodes[len(l.nodes)-1].hash
	}

	node, err := NewNode(id, prevHash, data)
	if err != nil {
		return err
	}

	l.nodes = append(l.nodes, node)
	return nil
}

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

// NOTE: REFACTOR TO USE GOB ENCODING
// Encodes a generic object to binary.
// Stores encrypted value in a byte slice.
// Object must have fixed size in memory.
func EncryptData[T any](data T) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		return nil, errors.New("unable to encrypt data")
	}

	return buf.Bytes(), nil
}

// NOTE: REFACTOR TO USE GOB ENCODING
// Decodes a generic object from binary.
// Type T must match the type used to encrypt the object.
func DecryptData[T any](data []byte) (T, error) {
	var value T
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		err = errors.New("unable to decrypt data: " + err.Error())
	}
	return value, err
}
