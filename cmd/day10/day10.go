package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type Light struct {
	location, v Point
}

type Sky []*Light

func (sky Sky) Print(frame int, writer *bufio.Writer) {
	var buf [2001][3001]byte

	minx, miny := -1500, -1000
	maxx, maxy := 1500, 1000

	for y := maxy; y >= miny; y-- {
		for x := minx; x <= maxx; x++ {
			buf[y-miny][x-minx] = '.'
		}
	}

	visible := 0
	for _, p := range sky {
		l := p.location
		if minx <= l.x && l.x <= maxx && miny <= l.y && l.y <= maxy {
			buf[l.y-miny][l.x-minx] = '#'
			visible++
		}
	}

	checkIt := false
	aligned := 0
	maxAligned := 0
	for x := minx; x <= maxx; x++ {
		for y := maxy; y >= miny; y-- {
			if buf[y-miny][x-minx] == '#' {
				aligned++
			} else {
				aligned = 0
			}

			if aligned > maxAligned {
				maxAligned = aligned
			}

			if aligned > 3 {
				checkIt = true
			}
		}
	}

	if checkIt {
		fmt.Fprintf(writer, " ~~ check %d - %d %d ~~\n", frame, visible, maxAligned)
		for y := maxy; y >= miny; y-- {
			for x := minx; x <= maxx; x++ {
				fmt.Fprintf(writer, "%c", buf[y-miny][x-minx])
			}
			fmt.Fprintf(writer, "\n")
		}
	} else {
		fmt.Fprintf(writer, "Frame seems to be too chaotic (%d), visible %d, maxAligned %d\n", frame, visible, maxAligned)
	}
	writer.Flush()
}

// func (sky Sky) CheckAlignment() int {
// var inLine map[int][]int{}

// for _, p := range sky {
// l := p.location
// inLine[l.x] = append(inLine[l.x], y)
// }

// for

// }

func (sky *Sky) Tick() {
	for _, p := range []*Light(*sky) {
		p.location.x += p.v.x
		p.location.y += p.v.y
	}
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var x, y, vx, vy int
	var lights Sky

	read, err := fmt.Fscanf(stdin, "position=<%d,  %d> velocity=< %d, %d>", &x, &y, &vx, &vy)
	if read >= 4 && err == nil {
		_, err = stdin.ReadString('\n')
		lights = append(lights, &Light{Point{x, y}, Point{vx, vy}})
	} else {
		fmt.Fprintf(stderr, "%s got: %d\n", err, read)
	}

	for err == nil {
		read, err = fmt.Fscanf(stdin, "position=<%d, %d> velocity=<%d, %d>", &x, &y, &vx, &vy)
		if read >= 4 && err == nil {
			_, err = stdin.ReadString('\n')
			lights = append(lights, &Light{Point{x, y}, Point{vx, vy}})
		} else {
			fmt.Fprintf(stderr, "%s got: %d\n", err, read)
		}
	}

	fmt.Fprintf(stderr, "Lines total: %d\n", len(lights))

	for i := 0; i < 15000; i++ {
		if 10500 < i && i < 11000 {
			lights.Print(i, stdout)
		}
		lights.Tick()
	}

	//fmt.Printf("#%d %d/%d %d\n", maxMinutesGuard.id, maxMinutesGuard.total, maxMinute, maxMinute*maxMinutesGuard.id)
}
