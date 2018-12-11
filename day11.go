package main

import (
	"bufio"
	"fmt"
	"os"
)

func Value(x, y int, serial int) int {
	rackId := x + 10
	v := rackId*y + serial
	v = v * rackId
	v = (v / 100) % 10
	v -= 5

	return v
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	serial := 9221

	var grid, sums [300][300]int

	for y := 0; y < 300; y++ {
		for x := 0; x < 300; x++ {
			grid[y][x] = Value(x+1, y+1, serial)
		}
	}

	var mx, my int
	maxSum := 0
	for y := 0; y < 300-3; y++ {
		for x := 0; x < 300-3; x++ {
			sum := 0
			for yy := y; yy < y+3; yy++ {
				for xx := x; xx < x+3; xx++ {
					sum += grid[yy][xx]
				}
			}
			sums[y][x] = sum

			if sum > maxSum {
				maxSum = sum
				mx = x
				my = y
			}
		}
	}

	fmt.Fprintf(stderr, "check %d == 4 \n", Value(3, 5, 8))
	fmt.Fprintf(stderr, "check %d == -5\n", Value(122, 79, 57))
	fmt.Fprintf(stderr, "check %d == 0\n", Value(217, 196, 39))
	fmt.Fprintf(stderr, "check %d == 4\n", Value(101, 153, 71))

	fmt.Printf("max sum: %d (%d, %d)\n", maxSum, mx+1, my+1)
}
