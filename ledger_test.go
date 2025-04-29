package ledger

import (
	"bytes"
	"testing"
)

// Test Ledger.AddNode.
// Creates a Ledger with 1 Node.
// Passes if there are no errors creating the Node.
func TestNewLedger_ByteString(t *testing.T) {
	testLedger := NewLedger()
	data := []byte("hello world")
	if err := testLedger.AddNode(data); err != nil {
		t.Errorf("error adding node to ledger: %v", err)
		t.Logf("%+v", testLedger.nodes[0])
	}
}

// Tests Ledger.ValidateHash
// using a mock Node.
// Expects the function to return true.
func TestValidateHash(t *testing.T) {
	testData := []byte("sample payload")

	n, err := NewNode(0, nil, testData)
	if err != nil {
		t.Fatalf("error creating node: %s", err.Error())
	}

	valid, err := n.ValidateHash()
	if err != nil {
		t.Fatalf("error validating hash: %s", err.Error())
	}

	if !valid {
		t.Fatal("node hash is invalid")
	}
}

// Tests Ledger.ValidateHash
// using a mock Node.
// Alters Node before calling ValidateHash.
// Expects the function to return false.
func TestValidateHash_InvalidHash(t *testing.T) {
	testData := []byte("sample payload")

	n, err := NewNode(0, nil, testData)
	if err != nil {
		t.Fatalf("error creating node: %s", err.Error())
	}

	//
	n.hash = append(n.hash, '0')

	valid, err := n.ValidateHash()
	if err != nil {
		t.Fatalf("error validating hash: %s", err.Error())
	}

	if valid {
		t.Fatal("node hash is valid")
	}
}

// Tests Ledger.ValidateNode
// using a new Ledger with 1 mock Node.
// Expects the function to return true.
func TestValidateNode(t *testing.T) {
	testData := []byte("sample payload")

	l := NewLedger()
	err := l.AddNode(testData)
	if err != nil {
		t.Fatalf("error adding node: %s", err.Error())
	}

	valid, err := l.ValidateNode(0)
	if err != nil {
		t.Fatalf("error validating node: %s,", err.Error())
	}

	if !valid {
		t.Fatal("node 0 hash is invalid")
	}
}

// Tests Ledger.ValidateNode
// using a new Ledger with 1 mock Node.
// Alters Node before calling ValidateNode.
// Expects the function to return false.
func TestValidateNode_InvalidHash(t *testing.T) {
	testData := []byte("sample payload")

	l := NewLedger()
	err := l.AddNode(testData)
	if err != nil {
		t.Fatalf("error adding node: %s", err.Error())
	}

	// Change data to alter hash
	l.nodes[0].data = []byte("not sample payload")

	valid, err := l.ValidateNode(0)
	if err != nil {
		t.Fatalf("error validating node: %s,", err.Error())
	}

	if valid {
		t.Fatal("node 0 hash is valid")
	}
}

// TestNewNode checks basic creation of a Node
func TestNewNode(t *testing.T) {
	data := []byte("hello world")
	prevHash := []byte{}

	node, err := NewNode(1, prevHash, data)
	if err != nil {
		t.Fatalf("error creating node: %v", err)
	}

	if node.id != 1 {
		t.Errorf("expected id 1, got %d", node.id)
	}

	if node.data == nil || !bytes.Equal(node.data, data) {
		t.Errorf("data mismatch: expected %v, got %v", data, node.data)
	}

	if node.prevHash == nil || len(node.prevHash) != 0 {
		t.Errorf("expected empty prevHash, got %v", node.prevHash)
	}

	if node.hash == nil || len(node.hash) != 32 {
		t.Errorf("expected SHA-256 hash (32 bytes), got length %d", len(node.hash))
	}
}

// TestComputeHash checks that ComputeHash returns consistent results
func TestComputeHash(t *testing.T) {
	data := []byte("test data")
	// prevHash := []byte("previous hash")
	prevHash := []byte(nil)

	node, err := NewNode(42, prevHash, data)
	if err != nil {
		t.Fatalf("error creating node: %v", err)
	}

	rehash, err := node.ComputeHash()
	if err != nil {
		t.Fatalf("error recomputing hash: %v", err)
	}

	if !bytes.Equal(node.hash, rehash) {
		t.Errorf("rehash mismatch: expected %x, got %x", node.hash, rehash)
	}
}

// TestGetData checks that GetData returns the correct byte slice
func TestGetData(t *testing.T) {
	data := []byte("sample payload")
	node, err := NewNode(3, nil, data)
	if err != nil {
		t.Fatalf("error creating node: %v", err)
	}

	got := node.GetData()
	if !bytes.Equal(data, got) {
		t.Errorf("GetData() mismatch: expected %v, got %v", data, got)
	}
}

// Tests EncryptData and DecryptData
// Encodes a test struct, then decodes it.
// Passes if decoded struct matches the original.
func TestEncryptDecryptData(t *testing.T) {
	type testType struct {
		Id  byte
		Age uint8
	}

	testData := testType{
		Id:  'A',
		Age: 33,
	}

	encData, err := EncryptData(testData)
	if err != nil {
		t.Fatalf("error encrypting test data: %v", err)
	}

	decData, err := DecryptData[testType](encData)
	if err != nil {
		t.Fatalf("error decrypting test data: %v", err)
	}

	if testData != decData {
		t.Log("decrypted data does not match original data")
		t.Errorf("%+v\n%+v\n", testData, decData)
	}
}
