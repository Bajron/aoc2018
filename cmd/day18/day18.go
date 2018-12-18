package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

func surround(x, y int) []Point {
	return []Point{
		Point{x - 1, y + 1}, Point{x, y + 1}, Point{x + 1, y + 1},
		Point{x - 1, y}, Point{x + 1, y},
		Point{x - 1, y - 1}, Point{x, y - 1}, Point{x + 1, y - 1}}
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	data := [][]rune{}
	filler := []rune{}

	read, err := stdin.ReadString('\n')
	if err == nil {
		filler = append(filler, 'x')
		for i := 0; i < len(read); i++ {
			filler = append(filler, 'x')
		}
		data = append(data, filler)

		runes := []rune(read)
		toAdd := make([]rune, len(runes)+1)
		toAdd[0] = 'x'
		for i := range []rune(runes) {
			toAdd[i+1] = runes[i]
		}
		toAdd[len(toAdd)-1] = 'x'
		data = append(data, toAdd)
	} else {
		fmt.Fprintf(stderr, "Read error: %s\n", err)
	}
	for err == nil {
		read, err = stdin.ReadString('\n')
		if err == nil {
			runes := []rune(read)
			toAdd := make([]rune, len(runes)+1)
			toAdd[0] = 'x'
			for i := range []rune(runes) {
				toAdd[i+1] = runes[i]
			}
			toAdd[len(toAdd)-1] = 'x'
			data = append(data, toAdd)
		} else {
			fmt.Fprintf(stderr, "Read error: %s\n", err)
		}
	}
	data = append(data, filler)

	fmt.Printf("lines: %d\n", len(data))

	buf := make([][]rune, len(data))
	for i := 0; i < len(data); i++ {
		buf[i] = make([]rune, len(data[i]))
		copy(buf[i], data[i])
	}
	changed := (len(data) - 2) * (len(data[0]) - 2)
	minute := 0

	history := map[string]int{}

	keep := make([]rune, len(data)*(len(data[0])-2))
	var cycleStart, cycleEnd int

	for minute = 0; minute < 1000000000 && changed > 0; minute++ {
		changed = 0

		for y := 1; y < len(data)-1; y++ {
			for x := 0; x < len(data[y]); x++ {
				keep[(y-1)*len(data[0])+x] = data[y][x]
			}
		}
		ks := string(keep)
		if v, ok := history[ks]; ok {
			cycleStart = v
			cycleEnd = minute
			fmt.Printf("already saw on %d; cycle %d --> %d\n", v, v, minute)
			break
		}

		history[string(keep)] = minute

		for y := 1; y < len(data)-1; y++ {
			for x := 1; x < len(data[y])-1; x++ {
				trees := 0
				lumber := 0
				free := 0

				for _, p := range surround(x, y) {
					if data[p.y][p.x] == '|' {
						trees++
					}
					if data[p.y][p.x] == '#' {
						lumber++
					}
					if data[p.y][p.x] == '.' {
						free++
					}
				}

				if data[y][x] == '.' {
					if trees >= 3 {
						buf[y][x] = '|'
						changed++
					} else {
						buf[y][x] = '.'
					}
				}
				if data[y][x] == '|' {
					if lumber >= 3 {
						buf[y][x] = '#'
						changed++
					} else {
						buf[y][x] = '|'
					}
				}
				if data[y][x] == '#' {
					if lumber >= 1 && trees >= 1 {
						buf[y][x] = '#'
					} else {
						buf[y][x] = '.'
						changed++
					}
				}
			}
		}
		fmt.Printf("changed %d in minute %d\n", changed, minute)
		data, buf = buf, data
	}

	targetCycle := cycleStart + (1000000000-cycleStart)%(cycleEnd-cycleStart)
	fmt.Printf("Looking for %d\n", targetCycle)

	for k, v := range history {
		if v == targetCycle {
			trees := 0
			lumber := 0
			free := 0

			for _, b := range k {
				if b == '|' {
					trees++
				}
				if b == '#' {
					lumber++
				}
				if b == '.' {
					free++
				}
			}

			fmt.Printf("HERE: t: %d l: %d f: %d   v: %d\n", trees, lumber, free, lumber*trees)
		}
	}

	// for y := 0; y < len(data); y++ {
	// 	for x := 0; x < len(data[y]); x++ {
	// 		fmt.Printf("%c", data[y][x])
	// 	}
	// 	fmt.Printf("\n")
	// }
	fmt.Printf("minute %d\n", minute)
	trees := 0
	lumber := 0
	free := 0
	for y := 1; y < len(data)-1; y++ {
		for x := 1; x < len(data[y])-1; x++ {
			if data[y][x] == '|' {
				trees++
			}
			if data[y][x] == '#' {
				lumber++
			}
			if data[y][x] == '.' {
				free++
			}
		}
	}
	fmt.Printf("t: %d l: %d f: %d   v: %d\n", trees, lumber, free, lumber*trees)
}
