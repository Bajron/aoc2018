package main

import (
	"bufio"
	"fmt"
	"os"
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

	recipes := make([]byte, 0, 1024*1024)

	recipes = append(recipes, 3)
	recipes = append(recipes, 7)

	e1 := 0
	e2 := 1

	for i := 0; i < 1000000; i++ {
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

	PrintNum(recipes[0 : 0+10])
	fmt.Printf("\n")
	PrintNum(recipes[9 : 9+10])
	fmt.Printf("\n")
	PrintNum(recipes[5 : 5+10])
	fmt.Printf("\n")
	PrintNum(recipes[18 : 18+10])
	fmt.Printf("\n")
	PrintNum(recipes[2018 : 2018+10])
	fmt.Printf("\nmine: \n")
	PrintNum(recipes[77201 : 77201+10])
}
