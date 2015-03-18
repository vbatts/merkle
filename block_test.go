package merkle

import "testing"

func TestBlockSize(t *testing.T) {
	var testSet = [][2]int{
		{1024 * 1024, 16384},
		{1023 * 1023, 0}, // Not a evenly divisible
		{1023, 1023},     // less than the max
	}
	for _, item := range testSet {
		got := DetermineBlockSize(item[0])
		if got != item[1] {
			t.Errorf("expected %d, got %d", item[1], got)
		}
	}
}
