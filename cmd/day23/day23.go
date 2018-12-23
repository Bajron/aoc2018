package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct {
	x, y, z int
}

func (p Point) add(a Point) Point {
	return Point{p.x + a.x, p.y + a.y, p.z + a.z}
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

func sortUnique(c []int) []int {
	r := make([]int, 0, cap(c))
	sort.Ints(c)
	r = append(r, c[0])
	for i := 1; i < len(c); i++ {
		if r[len(r)-1] != c[i] {
			r = append(r, c[i])
		}
	}
	return r
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
	beacon := -1
	beaconVisible := 0
	for bi, bref := range bots {
		visible := 0
		for _, b := range bots {
			if dist(bref.location, b.location) <= bref.radius {
				visible++
			}
		}
		if visible > beaconVisible {
			beaconVisible = visible
			beacon = bi
		}
		fmt.Printf("%d sees %d\n", bi, visible)
	}

	receiver := -1
	receiverVisible := 0
	for bi, bref := range bots {
		visible := 0
		for _, b := range bots {
			if dist(bref.location, b.location) <= b.radius {
				visible++
			}
		}
		if visible > receiverVisible {
			receiverVisible = visible
			receiver = bi
		}
		fmt.Printf("%d is seen by %d\n", bi, visible)
	}

	crosser := -1
	crossingCount := 0
	for bi, bref := range bots {
		crossings := 0
		for _, b := range bots {
			if dist(bref.location, b.location) <= b.radius+bref.radius {
				crossings++
			}
		}
		if crossings > crossingCount {
			crossingCount = crossings
			crosser = bi
		}
		fmt.Printf("%d is crosses with by %d\n", bi, crossings)
	}

	fmt.Printf("best beacon %v at index %d best %d\n", bots[beacon], beacon, beaconVisible)
	fmt.Printf("best receiver %v at index %d best %d\n", bots[receiver], receiver, receiverVisible)
	fmt.Printf("best crosser %v at index %d best %d\n", bots[crosser], crosser, crossingCount)

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

	zero := Point{0, 0, 0}

	best := 0
	candidate := Point{0, 0, 0}

	interesting := []Point{}

	interestingX := []int{}
	interestingY := []int{}
	interestingZ := []int{}
	for _, b := range bots {
		interestingX = append(interestingX, b.location.x, b.location.x-b.radius, b.location.x+b.radius)
		interestingY = append(interestingY, b.location.y, b.location.y-b.radius, b.location.y+b.radius)
		interestingZ = append(interestingZ, b.location.z, b.location.z-b.radius, b.location.z+b.radius)

		interesting = append(interesting, b.location)
		interesting = append(interesting, b.location.add(Point{1, 0, 0}))
		interesting = append(interesting, b.location.add(Point{-1, 0, 0}))
		interesting = append(interesting, b.location.add(Point{0, 1, 0}))
		interesting = append(interesting, b.location.add(Point{0, -1, 0}))
		interesting = append(interesting, b.location.add(Point{0, 0, 1}))
		interesting = append(interesting, b.location.add(Point{0, 0, -1}))
	}
	interestingX = sortUnique(interestingX)
	interestingY = sortUnique(interestingY)
	interestingZ = sortUnique(interestingZ)

	fmt.Printf("eh %d %d %d\n", len(interestingX), len(interestingY), len(interestingZ))

	for _, c := range interesting {
		visible := 0
		for _, b := range bots {
			if dist(c, b.location) <= b.radius {
				visible++
			}
		}
		if visible > best || (visible == best && dist(c, zero) < dist(candidate, zero)) {
			candidate = c
			best = visible

			fmt.Printf("something %v  best %d \n", candidate, best)
		}

	}

	best = 0
	candidate = Point{0, 0, 0}

	grain := 100
	for i := 0; i < 20; i++ {
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
		minx, maxx = candidate.x-33*dx, candidate.x+33*dx
		miny, maxy = candidate.y-33*dy, candidate.y+33*dy
		minz, maxz = candidate.z-33*dz, candidate.z+33*dz

		fmt.Printf("%d: grid lookup %v  best %d (%d,%d,%d)\n", i, candidate, best, dx, dy, dz)
	}

	fmt.Printf("grid lookup %v  best %d\n", candidate, best)

	for j := 0; j < 10; j++ {
		minx, maxx = candidate.x-150000, candidate.x+150000
		miny, maxy = candidate.y-150000, candidate.y+150000
		minz, maxz = candidate.z-150000, candidate.z+150000

		for i := 0; i < 20; i++ {
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
			minx, maxx = candidate.x-33*dx, candidate.x+33*dx
			miny, maxy = candidate.y-33*dy, candidate.y+33*dy
			minz, maxz = candidate.z-33*dz, candidate.z+33*dz

			fmt.Printf(" %d %d: grid lookup %v  best %d d=%d (%d,%d,%d)\n", j, i, candidate, best, dist(zero, candidate), dx, dy, dz)
		}
	}

	for {
		nextCandidate := candidate
		for x := candidate.x - 100; x < candidate.x+100; x++ {
			for y := candidate.y - 100; y < candidate.y+100; y++ {
				for z := candidate.z - 100; z < candidate.z+100; z++ {
					visible := 0
					c := Point{x, y, z}
					for _, b := range bots {
						if dist(c, b.location) <= b.radius {
							visible++
						}
					}
					if visible > best || (visible == best && dist(c, zero) < dist(candidate, zero)) {
						nextCandidate = c
						best = visible
					}
				}
			}
		}
		if nextCandidate == candidate {
			break
		}
		fmt.Printf("grid lookup refined %v  best %d dist %d\n", candidate, best, dist(zero, candidate))
		candidate = nextCandidate
	}

	// bot := bots[receiver]
	// candidate = bot.location
	// best = 0
	// fmt.Printf("candidate %v  best %d\n", candidate, best)

	// for x := bot.location.x - bot.radius; x <= bot.location.x+bot.radius; x++ {
	// 	eatenByX := abs(bot.location.x - x)
	// 	leftForY := bot.radius - eatenByX
	// 	for y := bot.location.y - leftForY; y <= bot.location.y+leftForY; y++ {
	// 		eatenByXY := eatenByX + abs(bot.location.y-y)
	// 		leftForZ := bot.radius - eatenByXY
	// 		for z := bot.location.z - leftForZ; y <= bot.location.z+leftForZ; z++ {
	// 			c := Point{x, y, z}
	// 			visible := 0
	// 			for _, b := range bots {
	// 				if dist(c, b.location) <= b.radius {
	// 					visible++
	// 				}
	// 			}
	// 			if visible > best || (visible == best && dist(c, zero) < dist(candidate, zero)) {
	// 				candidate = c
	// 				best = visible
	// 			}
	// 		}
	// 	}
	// }

	//fmt.Printf("final candidate %v best %d, dist to zero %d\n", candidate, best, dist(candidate, zero))
}
