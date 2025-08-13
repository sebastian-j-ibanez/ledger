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
	n := newNode(0, nil, testData)
	if !n.ValidHash() {
		t.Fatal("node hash is invalid")
	}
}

// Tests Ledger.ValidateHash
// using a mock Node.
// Alters Node before calling ValidateHash.
// Expects the function to return false.
func TestValidateHash_InvalidHash(t *testing.T) {
	testData := []byte("sample payload")
	n := newNode(0, nil, testData)
	n.Hash = append(n.Hash, '0')
	if n.ValidHash() {
		t.Fatal("node hash is valid")
	}
}

// TestNewNode checks basic creation of a Node
func TestNewNode(t *testing.T) {
	data := []byte("hello world")
	prevHash := []byte{}

	node := newNode(1, prevHash, data)

	if node.Id != 1 {
		t.Errorf("expected id 1, got %d", node.Id)
	}

	if node.Data == nil || !bytes.Equal(node.Data, data) {
		t.Errorf("data mismatch: expected %v, got %v", data, node.Data)
	}

	if node.PrevHash == nil || len(node.PrevHash) != 0 {
		t.Errorf("expected empty prevHash, got %v", node.PrevHash)
	}

	if node.Hash == nil || len(node.Hash) != 32 {
		t.Errorf("expected SHA-256 hash (32 bytes), got length %d", len(node.Hash))
	}
}

// TestComputeHash checks that ComputeHash returns consistent results
func TestComputeHash(t *testing.T) {
	data := []byte("test data")
	prevHash := []byte(nil)

	node := newNode(42, prevHash, data)
	rehash := node.computeHash()
	if !bytes.Equal(node.Hash, rehash) {
		t.Errorf("rehash mismatch: expected %x, got %x", node.Hash, rehash)
	}
}

// TestGetData checks that GetData returns the correct byte slice
func TestGetData(t *testing.T) {
	data := []byte("sample payload")
	node := newNode(3, nil, data)
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
