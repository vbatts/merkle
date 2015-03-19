package merkle

import "hash"

// NewHash provides a hash.Hash to generate a merkle.Tree checksum, given a
// HashMaker for the checksums of the blocks written and the blockSize of each
// block per node in the tree.
func NewHash(hm HashMaker, merkleBlockSize int) hash.Hash {
	mh := new(merkleHash)
	mh.blockSize = merkleBlockSize
	mh.hm = hm
	return mh
}

// TODO make a similar hash.Hash, that accepts an argument of a merkle.Tree,
// that will validate nodes as the new bytes are written. If a new written
// block fails checksum, then return an error on the io.Writer

// TODO satisfy the hash.Hash interface
type merkleHash struct {
	blockSize int
	tree      Tree
	hm        HashMaker
}

// XXX this will be tricky, as the last block can be less than the BlockSize.
// if they get the sum, it will be mh.tree.Root().Checksum() at that point.
//
// But if they continue writing, it would mean a continuation of the bytes in
// the last block. So popping the last node, and having a buffer for the bytes
// in that last partial block.
//
// if that last block was complete, then no worries. start the next node.
func (mh *merkleHash) Sum(b []byte) []byte {
	return nil
}
func (mh *merkleHash) Write(b []byte) (int, error) {
	return 0, nil
}

func (mh *merkleHash) Reset() {
	mh.Tree = Tree{}
}

func (mh *merkleHash) BlockSize() int { return mh.hm().BlockSize() }
func (mh *merkleHash) Size() int      { return mh.hm().Size() }
