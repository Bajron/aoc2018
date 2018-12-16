package main

import (
	"bufio"
	"fmt"
	"os"
)

func opOnInt(initial int64, op rune, value int64) int64 {
	if op == '-' {
		return initial - value
	} else if op == '+' {
		return initial + value
	} else {
		return initial
	}
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var sum, number, lines int64 = 0, 0, 0
	var op rune = ' '

	read, err := fmt.Fscanf(stdin, "%c%d\n", &op, &number)
	if err == nil && read >= 2 {
		sum = opOnInt(sum, op, number)
		lines++
	} else {
		fmt.Fprintf(stderr, "%s got: %d\n", err, read)
	}

	for err == nil {
		read, err = fmt.Fscanf(stdin, "%c%d\n", &op, &number)
		if err == nil && read >= 2 {
			sum = opOnInt(sum, op, number)
			lines++
		} else {
			fmt.Fprintf(stderr, "%s got: %d\n", err, read)
		}
	}
	fmt.Printf("%d\n", sum)
	fmt.Fprintf(stderr, "Lines processed: %d\n", lines)
}
