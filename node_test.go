package merkle

import (
	"fmt"
	"strings"
	"testing"
)

var words string = `Who were expelled from the academy for crazy & publishing obscene odes on the windows of the skull`

func TestNodeSums(t *testing.T) {
	var (
		nodes            []*Node
		h                = DefaultHash.New()
		expectedChecksum = "819fe8fed7a46900bd0613344c5ba2be336c74db"
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
	if len(nodes) != 1 {
		t.Errorf("%d nodes", len(nodes))
	}
	c, err := nodes[0].Checksum()
	if err != nil {
		t.Error(err)
	}
	if gotChecksum := fmt.Sprintf("%x", c); gotChecksum != expectedChecksum {
		t.Errorf("expected checksum %q, got %q", expectedChecksum, gotChecksum)
	}
}
