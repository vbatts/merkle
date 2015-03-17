package merkle

import (
	"log"
	"strings"
	"testing"
)

var words string = `Who were expelled from the academy for crazy & publishing obscene odes on the windows of the skull`

func TestNodeSums(t *testing.T) {
	var (
		nodes []*Node
		h     = DefaultHash.New()
	)
	for _, word := range strings.Split(words, " ") {
		h.Reset()
		if _, err := h.Write([]byte(word)); err != nil {
			t.Errorf("on word %q, encountered %s", word, err)
		}
		sum := h.Sum(nil)
		nodes = append(nodes, &Node{checksum: sum})
	}

	for {
		nodes = levelUp(nodes)
		if len(nodes) == 1 {
			break
		}
	}
	for i := range nodes {
		c, err := nodes[i].Checksum()
		if err != nil {
			t.Error(err)
		}
		t.Logf("checksum %x", c)
	}
	if len(nodes) > 0 {
		t.Errorf("%d nodes; %d characters", len(nodes), len(words))
	}
}

func levelUp(nodes []*Node) []*Node {
	var (
		newNodes []*Node
		last     = len(nodes) - 1
	)

	for i := range nodes {
		if i%2 == 0 {
			if i == last {
				// TODO rebalance the last parent
				log.Println("WHOOP")
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
