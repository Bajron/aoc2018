package main

import (
	"bufio"
	"fmt"
	"os"
)

type Entry struct {
	stable           rune
	stableCoordinate int
	vein             rune
	begin, end       int
}

type Arena struct {
	ground [2000][2000]rune
}

func NewArena() *Arena {
	var arena Arena
	for y := range arena.ground {
		for x := range arena.ground {
			arena.ground[y][x] = '.'
		}
	}
	return &arena
}

func (a *Arena) addVein(e Entry) {
	if e.stable == 'x' {
		for y := e.begin; y <= e.end; y++ {
			a.ground[y][e.stableCoordinate] = '#'
		}
	}
	if e.stable == 'y' {
		for x := e.begin; x <= e.end; x++ {
			a.ground[e.stableCoordinate][x] = '#'
		}
	}
}

func (a *Arena) flow(x, y int) {
	//fmt.Printf("flow %d, %d\n", x, y)
	a.ground[y][x] = '|'

	var lower int
	for lower = y + 1; lower < 2000 && !(a.ground[lower][x] == '#' || a.ground[lower][x] == '~'); lower++ {
		a.ground[lower][x] = '|'
	}

	if lower >= 2000 {
		return
	}

	lower--
	if a.ground[lower][x-1] != '#' {
		a.flowSide(x-1, lower, -1)
	}
	if a.ground[lower][x+1] != '#' {
		a.flowSide(x+1, lower, 1)
	}
}

func (a *Arena) flowSide(x, y, move int) {
	//fmt.Printf("flowSide %d, %d   %d\n", x, y, move)
	a.ground[y][x] = '|'

	if !(a.ground[y+1][x] == '#' || a.ground[y+1][x] == '~') {
		a.flow(x, y+1)
		return
	}

	next := x + move
	if next < 0 || next >= 2000 {
		return
	}

	for a.ground[y][next] != '#' {
		a.ground[y][next] = '|'
		if !(a.ground[y+1][next] == '#' || a.ground[y+1][next] == '~') {
			a.flow(next, y+1)
			break
		}
		next += move
	}
}

func (a Arena) countWater(ymin, ymax int) (int, int) {
	waterFlowing, waterStable := 0, 0
	for y := ymin; y <= ymax; y++ {
		for x := 0; x < 2000; x++ {
			if a.ground[y][x] == '|' {
				waterFlowing++
			}
			if a.ground[y][x] == '~' {
				waterStable++
			}
		}
	}
	return waterStable, waterFlowing
}

func (a *Arena) stabilize() int {
	stabilized := 0
	for y := 0; y < 2000; y++ {
		state := 0
		open := -1
		for x := 0; x < 2000; x++ {
			if state == 0 {
				if a.ground[y][x] == '#' {
					state = 1
					open = x
				}
				continue
			}
			if state == 1 || state == 2 {
				if a.ground[y][x] == '#' {
					for xx := x - 1; xx > open && a.ground[y][xx] == '|'; xx-- {
						a.ground[y][xx] = '~'
						stabilized++
					}
					state = 1
					open = x
					continue
				}

				if !(a.ground[y+1][x] == '#' || a.ground[y+1][x] == '~') {
					state = 0
				} else if a.ground[y][x] == '|' {
					state = 2
				} else {
					state = 0
				}
				continue
			}
		}
	}
	return stabilized
}

func (a Arena) getYRange() (int, int) {
	ymin, ymax := 2000, 0
	for y := 0; y < 2000; y++ {
		for x := 0; x < 2000; x++ {
			if a.ground[y][x] == '#' {
				if y < ymin {
					ymin = y
				}
				if y > ymax {
					ymax = y
				}
			}
		}
	}
	return ymin, ymax
}

func (a Arena) dump(ymin, ymax, xmin, xmax int) {
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			fmt.Printf("%c", a.ground[y][x])
		}
		fmt.Printf("\n")
	}
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var input []Entry

	var stable, vein rune
	var stableCoord, begin, end int

	read, err := fmt.Fscanf(stdin,
		"%c=%d, %c=%d..%d\n",
		&stable, &stableCoord, &vein, &begin, &end)

	if read == 5 && err == nil {
		input = append(input, Entry{stable, stableCoord, vein, begin, end})
	} else {
		fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
	}
	for err == nil {
		read, err = fmt.Fscanf(stdin,
			"%c=%d, %c=%d..%d\n",
			&stable, &stableCoord, &vein, &begin, &end)

		if read == 5 && err == nil {
			input = append(input, Entry{stable, stableCoord, vein, begin, end})
		} else {
			fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
		}
	}

	fmt.Printf("lines: %d\n", len(input))

	arena := NewArena()
	for _, i := range input {
		arena.addVein(i)
	}
	stabilized := -1
	ymin, ymax := arena.getYRange()

	for stabilized != 0 {
		arena.flow(500, 0)
		stabilized = arena.stabilize()
		// arena.dump(0, 13, 494, 507)
	}

	waterStable, waterFlowing := arena.countWater(ymin, ymax)
	// arena.dump(0, 2000-1, 0, 2000-1)
	fmt.Printf("waterStable: %d, waterFlowing: %d    sum: %d\n", waterStable, waterFlowing, waterFlowing+waterStable)
}
