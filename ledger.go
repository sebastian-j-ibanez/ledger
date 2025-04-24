package ledger

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"time"
)

type Ledger struct {
	head *Node
	tail *Node
	size uint
}

type Node struct {
	id        int
	timestamp time.Time
	hash      hash.Hash
	prevHash  hash.Hash
	data      []byte
}

// Create a new Node.
// Expects an id and data from a byte slice.
func (n *Node) NewNode(id int, data []byte) Node {
	return Node{
		id:        id,
		timestamp: time.Now(),
		hash:      n.ComputeHash(),
		prevHash:  nil,
		data:      data,
	}
}

// Compute a SHA256 hash of the Node.
// Value returned as hash.Hash.
func (n *Node) ComputeHash() hash.Hash {
	hashString := n.GetHashString()
	hash := sha256.New()
	hash.Write([]byte(hashString))
	return hash
}

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

// Decodes a generic object from binary.
// Type T must match the type used to encrypt the object.
func DecryptData[T any](data []byte) (T, error) {
	var value T
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		err = errors.New("unable to decrypte data: " + err.Error())
	}
	return value, err
}

// Get the encoded byte slice representation of the Node data.
// Data is encoded in binary format, using encoding/binary.
func (n *Node) GetData() []byte {
	return n.data
}

// Get a formatted hash string using the Node information.
func (n *Node) GetHashString() string {
	return fmt.Sprintf("%d%s%s%s", n.id, n.timestamp.GoString(), n.prevHash, n.data)
}
