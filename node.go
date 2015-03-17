package merkle

import (
	"crypto"
	_ "crypto/sha1"
	"fmt"
)

var (
	DefaultHash = crypto.SHA1
)

func NewNode() *Node {
	return &Node{hash: DefaultHash}
}

// Node is a fundamental part of the tree.
type Node struct {
	hash                crypto.Hash
	checksum            []byte
	Parent, Left, Right *Node
}

// If it is a leaf (no children) Node, then the Checksum is of the block of a
// payload. Otherwise, the Checksum is of it's two children's Checksum.
func (n Node) Checksum() ([]byte, error) {
	if n.checksum != nil {
		return n.checksum, nil
	}
	if n.Left != nil && n.Right != nil {

		// we'll ask our children for their sum and wait til they return
		var (
			lSumChan = make(chan childSumResponse)
			rSumChan = make(chan childSumResponse)
		)
		go func() {
			c, err := n.Left.Checksum()
			lSumChan <- childSumResponse{checksum: c, err: err}
		}()
		go func() {
			c, err := n.Right.Checksum()
			rSumChan <- childSumResponse{checksum: c, err: err}
		}()

		h := n.hash.New()

		// First left
		lSum := <-lSumChan
		if lSum.err != nil {
			return nil, lSum.err
		}
		if _, err := h.Write(lSum.checksum); err != nil {
			return nil, err
		}

		// then right
		rSum := <-rSumChan
		if rSum.err != nil {
			return nil, rSum.err
		}
		if _, err := h.Write(rSum.checksum); err != nil {
			return nil, err
		}

		return h.Sum(nil), nil
	}
	return nil, ErrNoChecksumAvailable{node: &n}
}

type ErrNoChecksumAvailable struct {
	node *Node
}

func (err ErrNoChecksumAvailable) Error() string {
	return fmt.Sprintf("no block or children available to derive checksum from: %#v", *err.node)
}

type childSumResponse struct {
	checksum []byte
	err      error
}
