// Ledger is a library for a blockchain ledger data structure.
//
// The Ledger is a tool that provides data validation and integrity.
// Each Node of the Ledger contains a unique hash.
// The hash can be validated to ensure data integrity.
// This library does not contain any external dependencies.
package ledger

import (
	"errors"
	"fmt"
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

// Validate a Node's hash.
// Return error if id is not valid.
func (l *Ledger) ValidateNode(id uint64) (bool, error) {
	// THIS WILL NOT WORK IF ELEMENTS ARE REMOVED (sharding)
	if id >= uint64(len(l.nodes)) {
		return false, errors.New("id is out of bounds")
	}

	node := l.nodes[id]
	valid, err := node.ValidateHash()
	if err != nil {
		msg := fmt.Sprintf("node %d: unable to validate hash: %s", node.id, err.Error())
		return false, errors.New(msg)
	}

	return valid, nil
}

// Recompute hash of each Node in the Ledger.
// Return error if there is issue computing a hash.
func (l *Ledger) RecomputeHashes() error {
	for _, node := range l.nodes {
		var err error
		if node.hash, err = node.ComputeHash(); err != nil {
			msg := fmt.Sprintf("error recomputing hashes: node %d: %s", node.id, err.Error())
			return errors.New(msg)
		}
	}

	return nil
}
