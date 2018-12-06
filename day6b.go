package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type Grain struct {
	distance int
	owner    int // -1 not set, -2 tied
}

type Arena struct {
	x0, y0, x1, y1 int
	grains         []Grain
}

func NewArena(x0, y0, x1, y1 int) *Arena {
	w := x1 - x0 + 1
	h := y1 - y0 + 1
	arena := make([]Grain, w*h)
	for gi := range arena {
		arena[gi] = Grain{w * h, -1}
	}
	return &Arena{x0, y0, x1, y1, arena}
}

func (arena Arena) translate(pos Point) (location int) {
	return (pos.x - arena.x0) + (pos.y-arena.y0)*(arena.x1-arena.x0+1)
}

func (arena Arena) at(pos Point) Grain {
	return arena.grains[arena.translate(pos)]
}

func (arena *Arena) set(pos Point, owner, distance int) {
	arena.grains[arena.translate(pos)] = Grain{distance, owner}
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (arena *Arena) FillAll(points []Point) {
	for pi, p := range points {
		for x := arena.x0; x <= arena.x1; x++ {
			for y := arena.y0; y <= arena.y1; y++ {
				distance := abs(p.x-x) + abs(p.y-y)
				arena.tryPut(Point{x, y}, pi, distance)
			}
		}
	}
	return
}

func (arena *Arena) SumAll(points []Point) {
	for gi := range arena.grains {
		arena.grains[gi].distance = 0
	}

	for _, p := range points {
		for x := arena.x0; x <= arena.x1; x++ {
			for y := arena.y0; y <= arena.y1; y++ {
				distance := abs(p.x-x) + abs(p.y-y)
				arena.grains[arena.translate(Point{x, y})].distance += distance
			}
		}
	}
	return
}

func (arena Arena) CountFields() map[int]int {
	ret := make(map[int]int)
	inf := arena.findInfinite()
	for id := range inf {
		ret[id] = -1
	}
	for x := arena.x0; x <= arena.x1; x++ {
		for y := arena.y0; y <= arena.y1; y++ {
			grain := arena.at(Point{x, y})
			if !inf[grain.owner] {
				ret[grain.owner]++
			}
		}
	}
	return ret
}

func (arena Arena) findInfinite() map[int]bool {
	ret := make(map[int]bool)
	for _, x := range []int{arena.x0, arena.x1} {
		for y := arena.y0; y <= arena.y1; y++ {
			ret[arena.at(Point{x, y}).owner] = true
		}
	}

	for x := arena.x0; x <= arena.x1; x++ {
		for _, y := range []int{arena.y0, arena.y1} {
			ret[arena.at(Point{x, y}).owner] = true
		}
	}
	return ret
}

func (arena *Arena) tryPut(pos Point, id int, distance int) {
	grain := arena.at(pos)
	if grain.owner == -1 || grain.distance > distance {
		arena.set(pos, id, distance)
	} else if grain.distance == distance && grain.owner != id {
		arena.set(pos, -2, distance)
	}
	return
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var points []Point

	var x, y int
	read, error := fmt.Fscanf(stdin, "%d, %d\n", &x, &y)
	if read == 2 && error == nil {
		points = append(points, Point{x, y})
	} else {
		fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
	}

	for error == nil {
		read, error = fmt.Fscanf(stdin, "%d, %d\n", &x, &y)
		if read == 2 && error == nil {
			points = append(points, Point{x, y})
		} else {
			fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
		}
	}
	fmt.Fprintf(stderr, "Read %d points\n", len(points))

	var x0, y0, x1, y1 int = 1000000, 1000000, 0, 0
	for _, p := range points {
		if p.x < x0 {
			x0 = p.x
		}
		if p.y < y0 {
			y0 = p.y
		}
		if p.x > x1 {
			x1 = p.x
		}
		if p.y > y1 {
			y1 = p.y
		}
	}

	arena := NewArena(x0, y0, x1, y1)
	arena.SumAll(points)

	count := 0
	for _, g := range arena.grains {
		if g.distance < 10000 {
			count++
		}
	}

	fmt.Printf("#region %d\n", count)
}
