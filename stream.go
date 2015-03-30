package merkle

import (
	"hash"
	"log"
	"runtime"
)

// NewHash provides a hash.Hash to generate a merkle.Tree checksum, given a
// HashMaker for the checksums of the blocks written and the blockSize of each
// block per node in the tree.
func NewHash(hm HashMaker, merkleBlockLength int) hash.Hash {
	mh := new(merkleHash)
	mh.blockSize = merkleBlockLength
	mh.hm = hm
	mh.tree = &Tree{Nodes: []*Node{}, BlockLength: merkleBlockLength}
	mh.lastBlock = make([]byte, merkleBlockLength)
	return mh
}

// TODO make a similar hash.Hash, that accepts an argument of a merkle.Tree,
// that will validate nodes as the new bytes are written. If a new written
// block fails checksum, then return an error on the io.Writer

type merkleHash struct {
	blockSize       int
	tree            *Tree
	hm              HashMaker
	lastBlock       []byte // as needed, for Sum()
	lastBlockLen    int
	partialLastNode bool // true when Sum() has appended a Node for a partial block
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
	var (
		curBlock     = []byte{}
		offset   int = 0
	)
	if mh.partialLastNode {
		// if this is true, then we need to pop the last node
		mh.tree.Nodes = mh.tree.Nodes[:len(mh.tree.Nodes)-1]
		mh.partialLastNode = false
	}

	if mh.lastBlockLen > 0 {
		curBlock = append(curBlock[:], mh.lastBlock[:mh.lastBlockLen]...)
		mh.lastBlockLen = 0
	}

	if b != nil && len(b) > 0 {
		curBlock = append(curBlock, b...)
	}

	for i := 0; i < len(curBlock)/mh.blockSize; i++ {
		n, err := NewNodeHashBlock(mh.hm, curBlock[offset:(offset+mh.blockSize)])
		if err != nil {
			// XXX i hate to swallow an error here, but the `Sum() []byte` signature
			// :-\
			sBuf := make([]byte, 1024)
			runtime.Stack(sBuf, false)
			log.Printf("[ERROR]: %s %q", err, string(sBuf))
			return nil
		}
		mh.tree.Nodes = append(mh.tree.Nodes, n)
		offset = offset + mh.blockSize
	}

	// If there is remainder, we'll need to make a partial node and stash it
	if m := (len(curBlock) % mh.blockSize); m != 0 {
		mh.lastBlockLen = copy(mh.lastBlock, curBlock[offset:])

		n, err := NewNodeHashBlock(mh.hm, curBlock[offset:])
		if err != nil {
			sBuf := make([]byte, 1024)
			runtime.Stack(sBuf, false)
			log.Printf("[ERROR]: %s %q", err, string(sBuf))
			return nil
		}
		mh.tree.Nodes = append(mh.tree.Nodes, n)
		mh.partialLastNode = true
	}

	sum, err := mh.tree.Root().Checksum()
	if err != nil {
		// XXX i hate to swallow an error here, but the `Sum() []byte` signature
		// :-\
		sBuf := make([]byte, 1024)
		runtime.Stack(sBuf, false)
		log.Printf("[ERROR]: %s %q", err, string(sBuf))
		return nil
	}
	return sum
}

func (mh *merkleHash) Write(b []byte) (int, error) {
	// basically we need to:
	// * include prior partial lastBlock, if any
	// * chunk these writes into blockSize
	// * create Node of the sum
	// * add the Node to the tree
	// * stash remainder in the mh.lastBlock

	var (
		curBlock       = make([]byte, mh.blockSize)
		numBytes   int = 0
		numWritten int
		offset     int = 0
	)
	if mh.lastBlock != nil && mh.lastBlockLen > 0 {
		//                                         XXX off by one?
		numBytes = copy(curBlock[:], mh.lastBlock[:mh.lastBlockLen])
		// not adding to numWritten, since these blocks were accounted for in a
		// prior Write()

		// then we'll chunk the front of the incoming bytes
		offset = copy(curBlock[numBytes:], b[:(mh.blockSize-numBytes)])
		n, err := NewNodeHashBlock(mh.hm, curBlock)
		if err != nil {
			// XXX might need to stash again the prior lastBlock and first little chunk
			return numWritten, err
		}
		mh.tree.Nodes = append(mh.tree.Nodes, n)
		numWritten += offset
	}

	numBytes = (len(b) - offset)
	for i := 0; i < numBytes/mh.blockSize; i++ {
		//fmt.Printf("%s", b[offset:offset+mh.blockSize])
		numWritten += copy(curBlock, b[offset:offset+mh.blockSize])
		n, err := NewNodeHashBlock(mh.hm, curBlock)
		if err != nil {
			// XXX might need to stash again the prior lastBlock and first little chunk
			return numWritten, err
		}
		mh.tree.Nodes = append(mh.tree.Nodes, n)
		offset = offset + mh.blockSize
	}

	mh.lastBlockLen = numBytes % mh.blockSize
	//                                       XXX off by one?
	numWritten += copy(mh.lastBlock[:], b[(len(b)-mh.lastBlockLen):])

	return numWritten, nil
}

func (mh *merkleHash) Reset() {
	mh.tree = &Tree{}
	mh.lastBlock = nil
}

// likely not the best to pass this through and not use our own node block
// size, but let's revisit this.
func (mh *merkleHash) BlockSize() int { return mh.hm().BlockSize() }
func (mh *merkleHash) Size() int      { return mh.hm().Size() }
