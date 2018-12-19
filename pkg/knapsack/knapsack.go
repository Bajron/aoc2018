package knapsack

type Item struct {
	weight int
	value  int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solve(items []Item, total int) (int, []int) {
	best := make([][]int, 0, len(items)+1)
	for i := 0; i < cap(best); i++ {
		best = append(best, make([]int, total+1))
	}

	for taken := 0; taken <= len(items); taken++ {
		best[taken][0] = 0
	}

	for weight := 0; weight <= total; weight++ {
		best[0][weight] = 0
	}

	for taken := 1; taken <= len(items); taken++ {
		item := items[taken-1]
		for weight := 1; weight < item.weight; weight++ {
			withoutIt := best[taken-1][weight]
			best[taken][weight] = withoutIt
		}
		for weight := item.weight; weight <= total; weight++ {
			withoutIt := best[taken-1][weight]
			withIt := best[taken-1][weight-item.weight] + item.value
			best[taken][weight] = max(withoutIt, withIt)
		}
	}

	chosen := []int{}

	for bw, bt := total, len(items); bw > 0 && bt > 0; bt-- {
		if best[bt][bw] == best[bt-1][bw] {
			continue
		}
		bw -= items[bt-1].weight
		chosen = append(chosen, bt-1)
	}

	return best[len(items)][total], chosen
}
