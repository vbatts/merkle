package merkle

type Tree struct {
	Nodes []*Node
}

func levelUp(nodes []*Node) []*Node {
	var (
		newNodes []*Node
		last     = len(nodes) - 1
	)

	for i := range nodes {
		if i%2 == 0 {
			if i == last {
				// last nodes on uneven node counts get pushed up, to be in the next level up
				newNodes = append(newNodes, nodes[i])
				continue
			}
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
