package merkle

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestMerkleHashWriterLargeChunk(t *testing.T) {
	// make a large enough test file of increments, corresponding to our blockSize
	bs := 512 * 1024
	fh, err := ioutil.TempFile("", "merkleChunks.")
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	defer os.Remove(fh.Name())

	// slow, i know ... FIXME
	for i := 0; i < 5; i++ {
		b := []byte{byte(i)}
		for j := 0; j < bs; j++ {
			fh.Write(b)
		}
	}
	if err := fh.Sync(); err != nil {
		t.Fatal(err)
	}
	if _, err := fh.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	expectedSums := []string{
		"6a521e1d2a632c26e53b83d2cc4b0edecfc1e68c", // 0's
		"316c136d75ffdeb6ac5f1262c45dd8c6ec50fd85", // 1's
		"a56e9c245b9c50d61a91c6c4299813b5e6313722", // 2's
		"58bed752c036310cc48d9dd0d25c4ee9ad0d7ff1", // 3's
		"bf382d8394213b897424803c27f3e2ec2223e5fd", // 4's
	}

	h := NewHash(DefaultHashMaker, bs)
	if _, err = io.Copy(h, fh); err != nil {
		t.Fatal(err)
	}
	h.Sum(nil)
	for i, node := range h.Nodes() {
		c, err := node.Checksum()
		if err != nil {
			t.Fatal(err)
		}
		if cs := fmt.Sprintf("%x", c); cs != expectedSums[i] {
			t.Errorf("expected sum %q; got %q", expectedSums[i], cs)
		}
	}
}

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

var bench = NewHash(DefaultHashMaker, 8192)
var buf = make([]byte, 8192)

func benchmarkSize(b *testing.B, size int) {
	b.SetBytes(int64(size))
	sum := make([]byte, bench.Size())
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:size])
		bench.Sum(sum[:0])
	}
}

func BenchmarkHash8Bytes(b *testing.B) {
	benchmarkSize(b, 8)
}

func BenchmarkHash1K(b *testing.B) {
	benchmarkSize(b, 1024)
}

func BenchmarkHash8K(b *testing.B) {
	benchmarkSize(b, 8192)
}
