package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	Rocky  = 0
	Wet    = 1
	Narrow = 2

	Neither  = 0
	Torch    = 1
	Climbing = 2

	MaxDim = 2000
)

type Point struct {
	x, y int
}

func (p Point) getSurround() []Point {
	ret := make([]Point, 0, 4)
	if p.y > 0 {
		ret = append(ret, Point{p.x, p.y - 1})
	}
	if p.x > 0 {
		ret = append(ret, Point{p.x - 1, p.y})
	}
	if p.x < MaxDim-1 {
		ret = append(ret, Point{p.x + 1, p.y})
	}
	if p.y < MaxDim-1 {
		ret = append(ret, Point{p.x, p.y + 1})
	}
	return ret
}

// Point Time Gear
type Ptg struct {
	p Point
	t int
	g int
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	// stdin := bufio.NewReader(os.Stdin)

	depth := 11109
	tx, ty := 9, 731

	// depth := 510
	// tx, ty := 10, 10

	M := 20183

	var cave [MaxDim][MaxDim]int

	cave[0][0] = 0
	for x := 1; x < MaxDim; x++ {
		cave[0][x] = (x * 16807) % M
	}
	for y := 1; y < MaxDim; y++ {
		cave[y][0] = (y * 48271) % M
	}

	for y := 1; y < MaxDim; y++ {
		for x := 1; x < MaxDim; x++ {
			if x == tx && y == ty {
				cave[y][x] = 0
				continue
			}
			cave[y][x] = (cave[y-1][x] + depth) * (cave[y][x-1] + depth) % M
		}
	}

	risk := 0
	for y := 0; y <= ty; y++ {
		for x := 0; x <= tx; x++ {
			risk += ((cave[y][x] + depth) % M) % 3
		}
	}
	fmt.Printf("risk %d\n", risk)
	var terrain [MaxDim][MaxDim]int
	var visited [MaxDim][MaxDim][3]int

	for y := 0; y < MaxDim; y++ {
		for x := 0; x < MaxDim; x++ {
			terrain[y][x] += ((cave[y][x] + depth) % M) % 3

			for wearing := 0; wearing < 3; wearing++ {
				visited[y][x][wearing] = 1000000000
			}
		}
	}

	bestNow := visited[ty][tx][Torch]
	bestNotChanged := 0
	wave := []Ptg{Ptg{Point{0, 0}, 0, Torch}}
	for len(wave) > 0 {
		nextWave := []Ptg{}

		for _, tmp := range wave {
			if tmp.g == terrain[tmp.p.y][tmp.p.x] {
				// cannot be here with this
				continue
			}

			if tmp.t >= visited[tmp.p.y][tmp.p.x][tmp.g] {
				// not better choice
				continue
			}

			visited[tmp.p.y][tmp.p.x][tmp.g] = tmp.t

			for _, p := range tmp.p.getSurround() {
				nextWave = append(nextWave, Ptg{p, tmp.t + 1, tmp.g})
			}
			nextWave = append(nextWave, Ptg{tmp.p, tmp.t + 7, (tmp.g + 1) % 3})
			nextWave = append(nextWave, Ptg{tmp.p, tmp.t + 7, (tmp.g + 2) % 3})
		}

		wave = nextWave
		if bestNow == visited[ty][tx][Torch] {
			bestNotChanged++
		} else {
			fmt.Printf("best now %d\n", visited[ty][tx][Torch])
			bestNow = visited[ty][tx][Torch]
			bestNotChanged = 0
		}

		if bestNotChanged > MaxDim*MaxDim*7*3 {
			fmt.Printf("stopping search with no change for %d\n", bestNotChanged)
			break
		}
	}

	fmt.Printf("best found %d\n", visited[ty][tx][Torch])
}
