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

type Bot struct {
	location                  Point
	radius                    int
	crossWith, sees, isSeenBy []*Bot
}

func NewBot(l Point, r int) *Bot {
	return &Bot{l, r, []*Bot{}, []*Bot{}, []*Bot{}}
}

func (b Bot) GetBox() Box {
	l, r := b.location, b.radius
	if r < 0 {
		panic("bad radius")
	}
	return Box{l.x - r, l.x + r, l.y - r, l.y + r, l.z - r, l.z + r}
}

type Box struct {
	xs, xe, ys, ye, zs, ze int
}

func (b Box) DimX() int {
	return b.xe - b.xs
}

func (b Box) DimY() int {
	return b.ye - b.ys
}
func (b Box) DimZ() int {
	return b.ze - b.zs
}

func (b Box) Volume() int {
	return b.DimX() * b.DimY() * b.DimZ()
}

func (b Box) crossWith(a Box) Box {
	return Box{max(b.xs, a.xs), min(b.xe, a.xe), max(b.ys, a.ys), min(b.ye, a.ye), max(b.zs, a.zs), min(b.ze, a.ze)}
}
func (b Box) isValid() bool {
	return b.xs <= b.xe && b.ys <= b.ye && b.zs <= b.ze
}

type BotSet struct {
	bots []*Bot
}

type CrossBoxState struct {
	box          Box
	crossedCount int
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	bots := []*Bot{}

	var x, y, z, r int
	read, error := fmt.Fscanf(stdin, "pos=<%d,%d,%d>, r=%d\n", &x, &y, &z, &r)
	if read == 4 && error == nil {
		bots = append(bots, NewBot(Point{x, y, z}, r))
	} else {
		fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
	}

	for error == nil {
		read, error = fmt.Fscanf(stdin, "pos=<%d,%d,%d>, r=%d\n", &x, &y, &z, &r)
		if read == 4 && error == nil {
			bots = append(bots, NewBot(Point{x, y, z}, r))
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
				bref.sees = append(bref.sees, b)
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
				bref.isSeenBy = append(bref.isSeenBy, b)
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
				bref.crossWith = append(bref.crossWith, b)
			}
		}
		if crossings > crossingCount {
			crossingCount = crossings
			crosser = bi
		}
		fmt.Printf("%d crosses with %d\n", bi, crossings)
	}

	fmt.Printf("best beacon %v at index %d best %d\n", bots[beacon].location, beacon, beaconVisible)
	fmt.Printf("best receiver %v at index %d best %d\n", bots[receiver].location, receiver, receiverVisible)
	fmt.Printf("best crosser %v at index %d best %d\n", bots[crosser].location, crosser, crossingCount)

	candidates := []int{}
	crossersOfCandidates := map[int][]int{}

	for bi, bref := range bots {
		crossings := 0
		crossers := []int{}
		for ci, b := range bots {
			if dist(bref.location, b.location) <= b.radius+bref.radius {
				crossings++
				crossers = append(crossers, ci)
			}
		}
		if crossings == crossingCount {
			candidates = append(candidates, bi)
			crossersOfCandidates[bi] = crossers
		}
	}

	fmt.Printf("candidates %v len:%d\n", candidates, len(candidates))

	zero := Point{0, 0, 0}

	//shortest := 4000000000
	thePoint := zero

	for _, pov := range bots {
		//pov:=bots[ci]
		sections := []*CrossBoxState{}
		for _, crosser := range pov.crossWith {
			cbox := crosser.GetBox()
			if !cbox.isValid() {
				panic("wt...")
			}
			crossedCount := 0
			for _, state := range sections {
				crossed := state.box.crossWith(cbox)
				if crossed.isValid() {
					crossedCount++
					state.box = crossed
					state.crossedCount++
				}
			}
			if crossedCount == 0 {
				newBox := pov.GetBox().crossWith(cbox)
				if !newBox.isValid() {
					panic("nope this is bad")
				}
				sections = append(sections, &CrossBoxState{newBox, 0})
			}
		}

		for _, s := range sections {
			fmt.Printf("Section size %d / dim %d  %v\n", s.crossedCount, s.box.Volume(), s.box)
		}

		fmt.Printf("\n")
		// fmt.Printf("search is %d\n", whereTooLook.DimX()*whereTooLook.DimY()*whereTooLook.DimZ())
	}

	// for _, ci := range candidates {
	// 	bot1 := bots[ci]

	// 	for _, cj := range candidates {
	// 		bot2 := bots[cj]

	// 		from := bot1.location
	// 		to := bot2.location

	// 	topLoop:
	// 		for x := from.x; abs(from.x-x) < bot1.radius && to.x != from.x; x += sign(to.x - from.x) {
	// 			for y := from.y; abs(from.x-x)+abs(from.y-y) < bot1.radius && to.y != from.y; y += sign(to.y - from.y) {
	// 				for z := from.z; abs(from.x-x)+abs(from.y-y)+abs(from.z-z) < bot1.radius && to.z != from.z; z += sign(to.z - from.z) {
	// 					p := Point{x, y, z}
	// 					if dist(p, to) <= bot2.radius {

	// 						visible := 0
	// 						for _, b := range bots {
	// 							if dist(b.location, p) <= b.radius {
	// 								visible++
	// 							}
	// 						}

	// 						fmt.Printf(" %d vs %d; cross at %v has %d visible\n", ci, cj, p, visible)
	// 						if visible == crossingCount {
	// 							break topLoop
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}

	// 	}
	// }

	// fmt.Printf("final candidate %v , dist to zero %d\n", thePoint, dist(thePoint, zero))

	// for x := bot.location.x - bot.radius; x <= bot.location.x+bot.radius; x++ {
	// 	eatenByX := abs(bot.location.x - x)
	// 	leftForY := bot.radius - eatenByX
	// 	for y := bot.location.y - leftForY; y <= bot.location.y+leftForY; y++ {
	// 		eatenByXY := eatenByX + abs(bot.location.y-y)
	// 		leftForZ := bot.radius - eatenByXY

	// 		c := Point{x, y, leftForZ}
	// 		visible := 0
	// 		for _, b := range bots {
	// 			if dist(c, b.location) <= b.radius {
	// 				visible++
	// 			}
	// 		}
	// 		if visible == crossingCount && dist(c, zero) < dist(thePoint, zero) {
	// 			thePoint = c
	// 		}
	// 	}
	// }

	fmt.Printf("final candidate %v , dist to zero %d\n", thePoint, dist(thePoint, zero))
}
