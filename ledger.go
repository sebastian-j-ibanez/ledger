// Ledger is a library for a blockchain ledger data structure.
//
// The Ledger is a tool that provides data validation and integrity.
// Each Node of the Ledger contains a unique hash.
// The hash can be validated to ensure data integrity.
// This library does not contain any external dependencies.
package ledger

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
		prevHash = l.nodes[len(l.nodes)-1].Hash
	}

	node := newNode(id, prevHash, data)
	l.nodes = append(l.nodes, node)
	return nil
}

// Check that each node has a valid hash.
func (l *Ledger) ValidateLedger() (bool, error) {
	for _, node := range l.nodes {
		if !node.ValidHash() {
			return false, nil
		}
	}

	return true, nil
}

// Get all Nodes.
func (l *Ledger) GetNodes() []Node {
	return l.nodes
}

// Get Node at ith index.
func (l *Ledger) GetNode(index int) *Node {
	if len(l.nodes) <= index && index >= 0 {
		return &l.nodes[index]
	}
	return nil
}
