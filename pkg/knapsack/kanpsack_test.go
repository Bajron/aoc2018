package knapsack

import (
	"fmt"
	"testing"
)

func TestGetModes(t *testing.T) {
	cormenInput := []Item{
		Item{20, 100}, Item{10, 60}, Item{30, 120},
	}
	optimal := 220
	s, choice := solve(cormenInput, 50)
	if s != optimal {
		t.Error("Bad solution")
	}
	fmt.Printf("%v\n", choice)
}
