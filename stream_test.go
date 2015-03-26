package merkle

import (
	"bytes"
	"io"
	"testing"
)

func TestMerkleHashWriter(t *testing.T) {
	msg := "the quick brown fox jumps over the lazy dog"
	h := NewHash(DefaultHashMaker, 10)
	i, err := io.Copy(h, bytes.NewBufferString(msg))
	if err != nil {
		t.Fatal(err)
	}
	if i != int64(len(msg)) {
		t.Fatalf("expected to write %d, only wrote %d", len(msg), i)
	}

	var (
		mh *merkleHash
		ok bool
	)
	if mh, ok = h.(*merkleHash); !ok {
		t.Fatalf("expected to get merkleHash, but got %#t", h)
	}

	// We're left with a partial lastBlock
	expectedNum := 4
	if len(mh.tree.Nodes) != expectedNum {
		t.Errorf("expected %d nodes, got %d", expectedNum, len(mh.tree.Nodes))
	}

	// Next test Sum()

	// count blocks again, we should get 5 nodes now

	// Test Sum() again, ensure same sum

	// Write more. This should pop the last node, and use the lastBlock.

}
