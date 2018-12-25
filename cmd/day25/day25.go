package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point [4]int

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dist(a, b Point) int {
	d := 0
	for i := 0; i < 4; i++ {
		d += abs(a[i] - b[i])
	}
	return d
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	points := []Point{}

	var c0, c1, c2, c3 int
	read, error := fmt.Fscanf(stdin, "%d,%d,%d,%d\n", &c0, &c1, &c2, &c3)
	if read == 4 && error == nil {
		points = append(points, Point{c0, c1, c2, c3})
	} else {
		fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
	}

	for error == nil {
		read, error = fmt.Fscanf(stdin, "%d,%d,%d,%d\n", &c0, &c1, &c2, &c3)
		if read == 4 && error == nil {
			points = append(points, Point{c0, c1, c2, c3})
		} else {
			fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
		}
	}
	fmt.Printf("points %d\n", len(points))

	isClose := map[int]map[int]bool{}
	for i := range points {
		for j := range points {
			if i == j {
				continue
			}
			d := dist(points[i], points[j])
			if d <= 3 {
				if _, ok := isClose[i]; !ok {
					isClose[i] = map[int]bool{}
				}
				isClose[i][j] = true

				if _, ok := isClose[j]; !ok {
					isClose[j] = map[int]bool{}
				}
				isClose[j][i] = true
			}
		}
	}

	inConstelation := map[int]int{}
	constelation := 0

	for i := range points {
		if _, ok := inConstelation[i]; ok {
			continue
		}

		wave := []int{i}
		for len(wave) > 0 {
			nextWave := []int{}

			for _, w := range wave {
				if vc, ok := inConstelation[w]; ok {
					if vc != constelation {
						panic("this flooding went wrong")
					}
					continue
				}
				inConstelation[w] = constelation
				for k := range isClose[w] {
					nextWave = append(nextWave, k)
				}
			}
			wave = nextWave
		}
		constelation++
	}

	assignedConstelations := map[int]bool{}
	for _, v := range inConstelation {
		assignedConstelations[v] = true
	}

	fmt.Printf("constelation count %d  -- %v\n", len(assignedConstelations), assignedConstelations)
}
