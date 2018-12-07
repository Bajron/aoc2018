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

func (g *Graph) Walk() string {
	var trace []rune

	completed := make(map[rune]bool)
	var possible NodeList

	for n := range g.nodes {
		if _, is := g.to[n]; !is {
			possible = append(possible, n)
		}
	}

	for len(possible) > 0 {
		sort.Sort(possible)
		trace = append(trace, possible[0])
		completed[possible[0]] = true

		possible = possible[0:0]

		for n := range g.nodes {
			if completed[n] {
				continue
			}

			v, is := g.to[n]
			if !is {
				possible = append(possible, n)
			} else {
				requiredCount := len(v)
				for _, required := range v {
					if completed[required] {
						requiredCount--
					}
				}
				//fmt.Printf("node %d req %d\n", n, requiredCount)
				if requiredCount == 0 {
					possible = append(possible, n)
				}
			}
		}
	}
	return string(trace)
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
	trace := g.Walk()

	fmt.Printf("#trace %s\n", trace)
}
