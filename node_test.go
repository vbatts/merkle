package merkle

import (
	"fmt"
	"strings"
	"testing"
)

func TestNodeSums(t *testing.T) {
	var (
		nodes            []*Node
		h                = DefaultHashMaker()
		words            = `Who were expelled from the academy for crazy & publishing obscene odes on the windows of the skull`
		expectedChecksum = "819fe8fed7a46900bd0613344c5ba2be336c74db"
	)
	for _, word := range strings.Split(words, " ") {
		h.Reset()
		if _, err := h.Write([]byte(word)); err != nil {
			t.Errorf("on word %q, encountered %s", word, err)
		}
		sum := h.Sum(nil)
		nodes = append(nodes, &Node{checksum: sum, hash: DefaultHashMaker})
	}

	newNodes := nodes
	for {
		newNodes = levelUp(newNodes)
		if len(newNodes) == 1 {
			break
		}
	}
	if len(newNodes) != 1 {
		t.Errorf("%d nodes", len(newNodes))
	}
	c, err := newNodes[0].Checksum()
	if err != nil {
		t.Error(err)
	}
	gotChecksum := fmt.Sprintf("%x", c)
	if gotChecksum != expectedChecksum {
		t.Errorf("expected checksum %q, got %q", expectedChecksum, gotChecksum)
	}

	tree := Tree{Nodes: nodes}
	c, err = tree.Root().Checksum()
	if err != nil {
		t.Error(err)
	}
	rootChecksum := fmt.Sprintf("%x", c)
	if rootChecksum != gotChecksum {
		t.Errorf("expected checksum %q, got %q", gotChecksum, rootChecksum)
	}

	expectedPieces := `7d531617dd394cef59d3cf58fc32b3bc458f6744a315dee0bd22f45265f67268f091869cca3cbf4ac267872aa7424b933c7e2b4de64e7c91b710686b0b1e95cfd9775191a7224d0a218ae79187e80c1dbbccdf2efb33b52e6c9d0a14dd70b2d415fbea6ecb2766cf39b9ee567af0081faffc4bb74c2b1fba43eef9a62abb8b1e1654f8a890aae054abffa82b33b501a5f87749b22562d3a7d38f8db6ccb80fe97c4d33785daa5c2370201ffa236b427aa37c99963fea93d27d200a96fc9e41ada467fda07ed68560efc7daae2005c903a8cb459ff1d51aee2988a3b3b04666d10863651a70ac9859cbeb83e919460bd3db3d405b10675998c030223177d42e71b4e7a312bbccdf2efb33b52e6c9d0a14dd70b2d415fbea6eab378b80a8a4aafabac7db7ae169f25796e65994de04fa0e29f9b35e24905d2e512bedc9bb6e09e4bbccdf2efb33b52e6c9d0a14dd70b2d415fbea6e15e9abb2e818480bc62afceb1b7f438663f7f08f`
	gotPieces := fmt.Sprintf("%x", tree.Pieces())
	if gotPieces != expectedPieces {
		t.Errorf("expected pieces %q, got %q", expectedPieces, gotPieces)
	}
}
