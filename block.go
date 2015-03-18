package merkle

const (
	// MaxBlockSize reasonable max byte size for blocks that are checksummed for
	// a Node
	MaxBlockSize = 1024 * 16
)

// DetermineBlockSize returns a reasonable block size to use, based on the
// provided size
func DetermineBlockSize(blockSize int) int {
	var b = blockSize
	for b > MaxBlockSize {
		b /= 2
	}
	if b == 0 || (blockSize%b != 0) {
		return 0
	}
	return b
}
