package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	from, to rune
}

type Graph struct {
	from, to map[rune][]rune
	nodes    map[rune]bool
}

func NewGraph(edges []Edge) *Graph {
	from, to := make(map[rune][]rune), make(map[rune][]rune)
	nodes := make(map[rune]bool)
	for _, e := range edges {
		from[e.from] = append(from[e.from], e.to)
		to[e.to] = append(to[e.to], e.from)
		nodes[e.from] = true
		nodes[e.to] = true
	}
	return &Graph{from, to, nodes}
}

type NodeList []rune

func (nl NodeList) Len() int {
	return len(nl)
}

func (nl NodeList) Swap(i, j int) {
	nl[i], nl[j] = nl[j], nl[i]
}

func (nl NodeList) Less(i, j int) bool {
	return nl[i] < nl[j]
}

func timeToRun(node rune) int {
	return 60 + int(node-rune('A')) + 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func countTrue(runeSet map[rune]bool) int {
	count := 0
	for _, v := range runeSet {
		if v {
			count++
		}
	}
	return count
}

func (g *Graph) Walk() (string, int) {
	var trace []rune

	completed := make(map[rune]bool)
	possible := make(map[rune]bool)

	for n := range g.nodes {
		if _, is := g.to[n]; !is {
			possible[n] = true
		}
	}

	timeSum := 0

	workerCount := 1 + 4
	running := make(map[rune]int)

	for countTrue(possible) > 0 {
		var possibleSequence NodeList
		for n, isPossible := range possible {
			if !isPossible {
				continue
			}
			if _, isRunning := running[n]; isRunning {
				continue
			}
			possibleSequence = append(possibleSequence, n)
		}
		sort.Sort(possibleSequence)

		runningAlready := len(running)
		//fmt.Printf("Possible %d [%v],  running %d [%v]\n", len(possibleSequence), possibleSequence, runningAlready, running)
		for i := 0; i < min(len(possibleSequence), workerCount-runningAlready); i++ {
			n := possibleSequence[i]
			running[n] = timeToRun(n)
		}

		fmt.Printf("  running %v\n", running)

		var completes rune
		completesTime := 10000
		for n, t := range running {
			if completesTime > t || completesTime == t && completes > n {
				completesTime = t
				completes = n
			}
		}

		for ri := range running {
			running[ri] -= completesTime
			fmt.Printf("  node %c will run for %d\n", ri, running[ri])
		}
		fmt.Printf(" ^^ %d total\n", len(running))

		timeSum += completesTime
		trace = append(trace, completes)

		completed[completes] = true
		possible[completes] = false

		delete(possible, completes)
		delete(running, completes)

		for n := range g.nodes {
			if completed[n] {
				continue
			}

			if _, isRunning := running[n]; isRunning {
				continue
			}

			v, hasDependency := g.to[n]
			if !hasDependency {
				possible[n] = true
			} else {
				requiredCount := len(v)
				for _, required := range v {
					if completed[required] {
						requiredCount--
					}
				}
				//fmt.Printf("node %d req %d\n", n, requiredCount)
				if requiredCount == 0 {
					possible[n] = true
				}
			}
		}
	}
	return string(trace), timeSum
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var edges []Edge

	var from, to rune
	read, error := fmt.Fscanf(stdin, "Step %c must be finished before step %c can begin.\n", &from, &to)
	if read == 2 && error == nil {
		edges = append(edges, Edge{from, to})
	} else {
		fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
	}

	for error == nil {
		read, error = fmt.Fscanf(stdin, "Step %c must be finished before step %c can begin.\n", &from, &to)
		if read == 2 && error == nil {
			edges = append(edges, Edge{from, to})
		} else {
			fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
		}
	}
	fmt.Fprintf(stderr, "Read %d edges\n", len(edges))

	g := NewGraph(edges)
	trace, time := g.Walk()

	fmt.Printf("#trace %s %d\n", trace, time)
}
