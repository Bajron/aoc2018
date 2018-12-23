package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y, z int
}

type Bot struct {
	location Point
	radius   int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func dist(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)
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

	bots := []Bot{}

	var x, y, z, r int
	read, error := fmt.Fscanf(stdin, "pos=<%d,%d,%d>, r=%d\n", &x, &y, &z, &r)
	if read == 4 && error == nil {
		bots = append(bots, Bot{Point{x, y, z}, r})
	} else {
		fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
	}

	for error == nil {
		read, error = fmt.Fscanf(stdin, "pos=<%d,%d,%d>, r=%d\n", &x, &y, &z, &r)
		if read == 4 && error == nil {
			bots = append(bots, Bot{Point{x, y, z}, r})
		} else {
			fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
		}
	}
	fmt.Printf("bots %d\n", len(bots))

	maxR := 0
	maxI := -1
	for i, b := range bots {
		if b.radius > maxR {
			maxI = i
			maxR = b.radius
		}
	}

	inRange := 0
	mb := bots[maxI]
	for _, b := range bots {
		if dist(mb.location, b.location) <= mb.radius {
			inRange++
		}
	}

	fmt.Printf("in range %d\n", inRange)

	for bi, bref := range bots {
		visible := 0
		for _, b := range bots {
			if dist(bref.location, b.location) <= b.radius {
				inRange++
			}
		}
	}

	minx, maxx, miny, maxy, minz, maxz := 2000000000, -2000000000, 2000000000, -2000000000, 2000000000, -2000000000
	for _, b := range bots {
		if b.location.x < minx {
			minx = b.location.x
		}
		if b.location.x > maxx {
			maxx = b.location.x
		}
		if b.location.y < miny {
			miny = b.location.y
		}
		if b.location.y > maxy {
			maxy = b.location.y
		}
		if b.location.z < minz {
			minz = b.location.z
		}
		if b.location.z > maxz {
			maxz = b.location.z
		}
	}
	grain := 200
	dx, dy, dz := (maxx-minx)/grain, (maxy-miny)/grain, (maxz-minz)/grain

	if dx <= 0 {
		dx = 1
	}
	if dy <= 0 {
		dy = 1
	}
	if dz <= 0 {
		dz = 1
	}

	zero := Point{0, 0, 0}
	best := 0
	candidate := Point{0, 0, 0}
	for x := minx; x < maxx; x += dx {
		for y := miny; y < maxy; y += dy {
			for z := minz; z < maxz; z += dz {
				visible := 0
				c := Point{x, y, z}
				for _, b := range bots {
					if dist(c, b.location) <= b.radius {
						visible++
					}
				}
				if visible > best || (visible == best && dist(c, zero) < dist(candidate, zero)) {
					candidate = c
					best = visible
				}
			}
		}
	}
	fmt.Printf("candidate %v best %d\n", candidate, best)

	fine := Point{0, 0, 0}
	fineBest := 0

	for {
		prevCandidate := candidate
		for x := candidate.x - grain; x < candidate.x+grain; x++ {
			for y := candidate.y - grain; y < candidate.y+grain; y++ {
				for z := candidate.z - grain; z < candidate.z+grain; z++ {
					visible := 0
					c := Point{x, y, z}
					for _, b := range bots {
						if dist(c, b.location) <= b.radius {
							visible++
						}
					}
					if visible > fineBest || (visible == fineBest && dist(c, zero) < dist(fine, zero)) {
						fine = c
						fineBest = visible
					}
				}
			}
		}

		for {
			visible := 0
			c := Point{fine.x + sign(zero.x-fine.x), fine.y + sign(zero.y-fine.y), fine.z + sign(zero.z-fine.z)}
			for _, b := range bots {
				if dist(c, b.location) <= b.radius {
					visible++
				}
			}

			if visible != fineBest {
				break
			}
			fine = c
		}

		if fine == prevCandidate {
			break
		}
		candidate = fine
		fmt.Printf("fine candidate %v best %d\n", fine, fineBest)
	}

	fmt.Printf("final candidate %v best %d, dist to zero %d\n", fine, fineBest, dist(fine, zero))
}
