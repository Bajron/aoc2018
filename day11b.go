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

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	serial := 9221

	var grid [300][300]int
	var sums [300][300][]int

	for y := 0; y < 300; y++ {
		for x := 0; x < 300; x++ {
			grid[y][x] = Value(x+1, y+1, serial)
		}
	}

	var mx, my, ms int
	maxSum := 0
	for y := 0; y < 300; y++ {
		for x := 0; x < 300; x++ {
			sum := 0
			maxEdge := Min(300-y, 300-x)
			for edge := 1; edge < maxEdge; edge++ {
				for xx := x; xx < x+edge; xx++ {
					sum += grid[y+edge-1][xx]
				}

				for yy := y; yy < y+edge-1; yy++ {
					sum += grid[yy][x+edge-1]
				}

				sums[y][x] = append(sums[y][x], sum)

				if sum > maxSum {
					maxSum = sum
					mx = x
					my = y
					ms = edge
				}
			}
		}
	}

	fmt.Fprintf(stderr, "check %d == 4 \n", Value(3, 5, 8))
	fmt.Fprintf(stderr, "check %d == -5\n", Value(122, 79, 57))
	fmt.Fprintf(stderr, "check %d == 0\n", Value(217, 196, 39))
	fmt.Fprintf(stderr, "check %d == 4\n", Value(101, 153, 71))

	fmt.Printf("max sum: %d (%d,%d,%d)\n", maxSum, mx+1, my+1, ms)
}
