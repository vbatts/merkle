package merkle

import (
	"crypto"
	_ "crypto/sha1" // to satisfy our DefaultHash
	"fmt"
	"hash"
)

var (
	// DefaultHash is for checksum of blocks and nodes
	DefaultHash = crypto.SHA1
)

// HashMaker produces a new has for use in making checksums
type HashMaker func() hash.Hash

// NewNode returns a new Node with the DefaultHash for checksums
func NewNode() *Node {
	return NewNodeHash(DefaultHash.New)
}

// NewNodeHash returns a new Node using the provided crypto.Hash for checksums
func NewNodeHash(h HashMaker) *Node {
	return &Node{hash: h}
}

// Node is a fundamental part of the tree.
type Node struct {
	hash                HashMaker
	checksum            []byte
	Parent, Left, Right *Node

	pos int // XXX maybe keep their order when it is a direct block's hash
}

// IsLeaf indicates this node is for specific block (and has no children)
func (n Node) IsLeaf() bool {
	return len(n.checksum) != 0 && (n.Left == nil && n.Right == nil)
}

// Checksum returns the checksum of the block, or the checksum of this nodes
// children (left.checksum + right.checksum)
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

		h := n.hash()

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

// ErrNoChecksumAvailable is for nodes that do not have the means to provide
// their checksum
type ErrNoChecksumAvailable struct {
	node *Node
}

// Error shows the message with information on the node
func (err ErrNoChecksumAvailable) Error() string {
	return fmt.Sprintf("no block or children available to derive checksum from: %#v", *err.node)
}

type childSumResponse struct {
	checksum []byte
	err      error
}
