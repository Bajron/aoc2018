package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	meta     []int
	children []*Node
}

func ReadNode(stdin *bufio.Reader) *Node {
	var node Node
	var childrenCount int
	var metaCount int

	fmt.Fscanf(stdin, "%d %d", &childrenCount, &metaCount)

	for i := 0; i < childrenCount; i++ {
		ch := ReadNode(stdin)
		node.children = append(node.children, ch)
	}

	var m int
	for i := 0; i < metaCount; i++ {
		fmt.Fscan(stdin, &m)
		node.meta = append(node.meta, m)
	}

	return &node
}

func WalkCountMeta(node *Node) int {
	sum := 0
	for _, m := range node.meta {
		sum += m
	}
	for _, ch := range node.children {
		sum += WalkCountMeta(ch)
	}
	return sum
}

func WalkCountValue(node *Node) int {
	value := 0
	if len(node.children) == 0 {
		for _, m := range node.meta {
			value += m
		}
	} else {
		for _, m := range node.meta {
			i := m - 1
			if 0 <= i && i < len(node.children) {
				value += WalkCountValue(node.children[i])
			}
		}
	}
	return value
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)

	root := ReadNode(stdin)

	ms := WalkCountMeta(root)
	fmt.Printf("Sum: %d\n", ms)
	v := WalkCountValue(root)
	fmt.Printf("Value: %d\n", v)
}
