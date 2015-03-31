package merkle

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestMerkleHashWriter(t *testing.T) {
	msg := "the quick brown fox jumps over the lazy dog"
	expectedSum := "48940c1c72636648ad40aa59c162f2208e835b38"

	h := NewHash(DefaultHashMaker, 10)
	i, err := io.Copy(h, bytes.NewBufferString(msg))
	if err != nil {
		t.Fatal(err)
	}
	if i != int64(len(msg)) {
		t.Fatalf("expected to write %d, only wrote %d", len(msg), i)
	}

	// We're left with a partial lastBlock
	expectedNum := 4
	if len(h.Nodes()) != expectedNum {
		t.Errorf("expected %d nodes, got %d", expectedNum, len(h.Nodes()))
	}

	// Next test Sum()
	gotSum := fmt.Sprintf("%x", h.Sum(nil))
	if expectedSum != gotSum {
		t.Errorf("expected initial checksum %q; got %q", expectedSum, gotSum)
	}

	// count blocks again, we should get 5 nodes now
	expectedNum = 5
	if len(h.Nodes()) != expectedNum {
		t.Errorf("expected %d nodes, got %d", expectedNum, len(h.Nodes()))
	}

	// Test Sum() again, ensure same sum
	gotSum = fmt.Sprintf("%x", h.Sum(nil))
	if expectedSum != gotSum {
		t.Errorf("expected checksum %q; got %q", expectedSum, gotSum)
	}

	// test that Reset() nulls us out
	h.Reset()
	gotSum = fmt.Sprintf("%x", h.Sum(nil))
	if expectedSum == gotSum {
		t.Errorf("expected reset checksum to not equal %q; got %q", expectedSum, gotSum)
	}

	// write our msg again and get the same sum
	i, err = io.Copy(h, bytes.NewBufferString(msg))
	if err != nil {
		t.Fatal(err)
	}
	if i != int64(len(msg)) {
		t.Fatalf("expected to write %d, only wrote %d", len(msg), i)
	}
	// Test Sum(), ensure same sum
	gotSum = fmt.Sprintf("%x", h.Sum(nil))
	if expectedSum != gotSum {
		t.Errorf("expected checksum %q; got %q", expectedSum, gotSum)
	}

	// Write more. This should pop the last node, and use the lastBlock.
	i, err = io.Copy(h, bytes.NewBufferString(msg))
	if err != nil {
		t.Fatal(err)
	}
	if i != int64(len(msg)) {
		t.Fatalf("expected to write %d, only wrote %d", len(msg), i)
	}
	expectedNum = 9
	if len(h.Nodes()) != expectedNum {
		t.Errorf("expected %d nodes, got %d", expectedNum, len(h.Nodes()))
	}
	gotSum = fmt.Sprintf("%x", h.Sum(nil))
	if expectedSum == gotSum {
		t.Errorf("expected reset checksum to not equal %q; got %q", expectedSum, gotSum)
	}
	if len(h.Nodes()) != expectedNum {
		t.Errorf("expected %d nodes, got %d", expectedNum, len(h.Nodes()))
	}

}
