package merkle

// Tree is the information on the structure of a set of nodes
//
// TODO more docs here
type Tree struct {
	Nodes       []*Node `json:"pieces"`
	BlockLength int     `json:"piece length"`
}

// Pieces returns the concatenation of hash values of all blocks
//
// TODO integrate with hash size
func (t *Tree) Pieces() []byte {
	if len(t.Nodes) == 0 {
		return nil
	}
	pieces := []byte{}
	for _, n := range t.Nodes {
		if n.checksum == nil || len(n.checksum) == 0 {
			continue
		}
		pieces = append(pieces, n.checksum...)
	}
	return pieces
}

// Root generates a hash tree bash on the current nodes, and returns the root
// of the tree
func (t *Tree) Root() *Node {
	newNodes := t.Nodes
	for {
		newNodes = levelUp(newNodes)
		if len(newNodes) == 1 {
			break
		}
	}
	return newNodes[0]
}

func levelUp(nodes []*Node) []*Node {
	var (
		newNodes []*Node
		last     = len(nodes) - 1
	)

	for i := range nodes {
		if i%2 == 0 {
			if i == last {
				// last nodes on uneven node counts get pushed up, to be in the next
				// level up
				newNodes = append(newNodes, nodes[i])
				continue
			}
			//n := NewNodeHash(nodes[i].hash) // use the node's hash type
			n := NewNode()
			n.Left = nodes[i]
			n.Left.Parent = n
			newNodes = append(newNodes, n)
		} else {
			n := newNodes[len(newNodes)-1]
			n.Right = nodes[i]
			n.Right.Parent = n
		}
	}
	return newNodes
}
