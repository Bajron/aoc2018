package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintNum(bytes []byte) {
	for _, b := range bytes {
		fmt.Printf("%d", int(b))
	}
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	// stdin := bufio.NewReader(os.Stdin)

	recipes := make([]byte, 0, 1024*1024*1024)

	recipes = append(recipes, 3)
	recipes = append(recipes, 7)

	e1 := 0
	e2 := 1

	for i := 0; i < 1000000000; i++ {
		created := recipes[e1] + recipes[e2]
		c1 := created / 10
		c2 := created % 10

		if c1 > 0 {
			recipes = append(recipes, c1)
		}
		recipes = append(recipes, c2)

		e1 = (e1 + int(recipes[e1]) + 1) % len(recipes)
		e2 = (e2 + int(recipes[e2]) + 1) % len(recipes)
	}

	for i := range recipes {
		recipes[i] += '0'
	}
	s := string(recipes)

	fmt.Printf("9: %d\n", strings.Index(s, "51589"))
	fmt.Printf("5: %d\n", strings.Index(s, "01245"))
	fmt.Printf("18: %d\n", strings.Index(s, "92510"))
	fmt.Printf("2018: %d\n", strings.Index(s, "59414"))
	fmt.Printf("mine: %d\n", strings.Index(s, "077201"))
}
