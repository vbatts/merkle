/*
What do you expect from a merkle tree API?
* streaming support
 - building a tree from an io.Reader
 - validating a tree from an io.Reader
* concurrency safe
 - any buffer or hash.Hash reuse
*/

package merkle
