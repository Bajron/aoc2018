package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cave struct {
	elfPower int
	caveMap  [][]rune
	state    [][]rune
	m        [][]int
	dudes    Dudes
	girth    int
}

const (
	FullRound = iota
	Interrupted
	ElfDead
	Unknown
)

type Dude struct {
	id       int
	hp       int
	pos      Point
	fraction rune
}

func (d Dude) getSurround() []Point {
	ret := make([]Point, 4)
	ret[0] = Point{d.pos.x, d.pos.y - 1}
	ret[1] = Point{d.pos.x - 1, d.pos.y}
	ret[2] = Point{d.pos.x + 1, d.pos.y}
	ret[3] = Point{d.pos.x, d.pos.y + 1}
	return ret
}

func (p Point) getSurround() []Point {
	ret := make([]Point, 4)
	ret[0] = Point{p.x, p.y - 1}
	ret[1] = Point{p.x - 1, p.y}
	ret[2] = Point{p.x + 1, p.y}
	ret[3] = Point{p.x, p.y + 1}
	return ret
}

type Dudes []*Dude

type Point struct {
	x, y int
}

type Points []Point

func NewCave(lines [][]rune, elfPower int) *Cave {
	dudeId := 0
	var dudes []*Dude
	caveState := make([][]rune, len(lines))
	caveMap := make([][]rune, len(lines))
	m := make([][]int, len(lines))

	for r := range lines {
		l := make([]rune, len(lines[r]))
		copy(l, lines[r])
		caveMap[r] = l

		l = make([]rune, len(lines[r]))
		copy(l, lines[r])
		caveState[r] = l

		m[r] = make([]int, len(lines[r]))
	}
	for y, row := range lines {
		for x, b := range row {
			if isDude(b) {
				dudes = append(dudes, &Dude{dudeId, 200, Point{x, y}, b})
				dudeId++
				caveMap[y][x] = '.'
			}
		}
	}
	return &Cave{elfPower, caveMap, caveState, m, dudes, len(caveState) * len(caveState[0])}
}

func (cave Cave) Print() {
	for _, row := range cave.state {
		// Right now newlines from the input are also there :D
		fmt.Printf("%s", string(row))
	}
}

func (cave Cave) dumpFlood() {
	for y := range cave.m {
		for x := range cave.m[y] {
			fmt.Printf("%5d ", cave.m[y][x])
		}
		fmt.Printf("\n")
	}
}

func (cave Cave) dumpDudes() {
	for _, d := range cave.dudes {
		fmt.Printf("%d %c %d,%d hp:%d \n", d.id, d.fraction, d.pos.x, d.pos.y, d.hp)
	}
}

func (cave *Cave) sweepDeadDudes() {
	aliveDudes := Dudes{}
	for _, d := range cave.dudes {
		if d != nil {
			aliveDudes = append(aliveDudes, d)
		}
	}
	cave.dudes = aliveDudes
}

func (cave *Cave) Tick() int {
	sort.Sort(cave.dudes)
	defer cave.sweepDeadDudes()

	for _, d := range cave.dudes {
		if d == nil {
			continue
		}
		fmt.Printf("%d %c %d,%d (hp:%d) starts ", d.id, d.fraction, d.pos.x, d.pos.y, d.hp)
		enemyAround := false
		for _, p := range d.getSurround() {
			if cave.isEnemy(p, d.fraction) {
				enemyAround = true
				break
			}
		}
		if !enemyAround {
			// Move
			targets := Points{}
			enemyCount := 0

			for _, dd := range cave.dudes {
				if dd == nil {
					continue
				}
				if dd.fraction == d.fraction {
					continue
				}
				enemyCount++
				for _, p := range dd.getSurround() {
					if cave.state[p.y][p.x] != '.' {
						continue
					}

					targets = append(targets, p)
				}
			}

			if enemyCount == 0 {
				return Interrupted
			}

			sort.Sort(targets)
			fmt.Printf("sees %d targets ", len(targets))

			var newPos Point
			cost := cave.girth

			for _, t := range targets {
				for y := range cave.m {
					for x := range cave.m[y] {
						cave.m[y][x] = cave.girth
					}
				}

				cave.flood(t, d.pos)

				for _, p := range d.getSurround() {
					if cave.m[p.y][p.x] < cost {
						cost = cave.m[p.y][p.x]
						newPos = p
					}
				}
			}

			if cost != cave.girth {
				fmt.Printf("goes to %d, %d (%d)", newPos.x, newPos.y, cost)

				cave.state[d.pos.y][d.pos.x] = cave.caveMap[d.pos.y][d.pos.x]
				d.pos = newPos
				cave.state[d.pos.y][d.pos.x] = d.fraction
			} else {
				fmt.Printf("has nowhere to go")
			}
		}

		targetHp := 400
		var toAttackIndex int
		var toAttack *Dude = nil
		for _, p := range d.getSurround() {
			if cave.isEnemy(p, d.fraction) {
				ci, candidate := cave.dudes.Find(p)
				if candidate.hp < targetHp {
					targetHp = candidate.hp
					toAttack = candidate
					toAttackIndex = ci
				}
			}
		}

		if toAttack != nil {
			fmt.Printf("attack %d at %d,%d", toAttack.id, toAttack.pos.x, toAttack.pos.y)
			if d.fraction == 'E' {
				toAttack.hp -= cave.elfPower
			} else {
				toAttack.hp -= 3
			}
			if toAttack.hp <= 0 {
				cave.dudes[toAttackIndex] = nil
				cave.state[toAttack.pos.y][toAttack.pos.x] = cave.caveMap[toAttack.pos.y][toAttack.pos.x]

				if toAttack.fraction == 'E' {
					return ElfDead
				}
			}
		}

		fmt.Println()

		// cave.Print()
	}

	return FullRound
}

func (cave *Cave) flood(from, to Point) {
	distance := 1
	border := []Point{from}

	for len(border) > 0 {
		wave := []Point{}
		for _, p := range border {
			if p == to {
				if distance < cave.m[p.y][p.x] {
					cave.m[p.y][p.x] = distance
				}
				continue
			}

			if cave.state[p.y][p.x] != '.' {
				continue
			}

			if cave.m[p.y][p.x] <= distance {
				continue
			}

			cave.m[p.y][p.x] = distance

			wave = append(wave, p.getSurround()...)
		}

		distance++
		border = wave
	}
}

func (cave *Cave) isEnemy(p Point, r rune) bool {
	tile := cave.state[p.y][p.x]
	return isDude(tile) && tile != r
}

func isDude(r rune) bool {
	return r == 'G' || r == 'E'
}

func (d Dudes) Len() int {
	return len(d)
}

func (d Dudes) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Dudes) Less(i, j int) bool {
	return d[i].pos.y < d[j].pos.y || (d[i].pos.y == d[j].pos.y && d[i].pos.x < d[j].pos.x)
}

func (dudes Dudes) Find(p Point) (int, *Dude) {
	for i, d := range dudes {
		if d != nil && d.pos == p {
			return i, d
		}
	}
	return -1, nil
}

func (p Points) Len() int {
	return len(p)
}

func (p Points) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Points) Less(i, j int) bool {
	return p[i].y < p[j].y || (p[i].y == p[j].y && p[i].x < p[j].x)
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var lines [][]rune

	line, err := stdin.ReadString('\n')
	if err == nil {
		lines = append(lines, []rune(line))
	} else {
		fmt.Fprintf(stderr, "Read error: %s\n", err)
	}
	for err == nil {
		line, err = stdin.ReadString('\n')
		if err == nil {
			lines = append(lines, []rune(line))
		} else {
			fmt.Fprintf(stderr, "Read error: %s\n", err)
		}
	}

	fmt.Printf("lines: %d\n", len(lines))

	elfPower := 2
	var cave *Cave = nil
	rounds := 0
	lastOutcome := ElfDead

	for lastOutcome == ElfDead {
		elfPower++
		cave = NewCave(lines, elfPower)

		rounds = 0
		lastOutcome = Unknown

		fmt.Printf("elfPower %d\n", elfPower)

		for !(lastOutcome == ElfDead || lastOutcome == Interrupted) {
			fmt.Printf("Starting round %d\n", rounds+1)
			cave.Print()

			lastOutcome = cave.Tick()
			if lastOutcome == FullRound {
				rounds++
			}
		}
	}
	remainingHp := 0
	for _, d := range cave.dudes {
		remainingHp += d.hp
	}

	fmt.Printf("\n\n elfPower: %d, rounds: %d, hp: %d; checksum: %d\n\n", elfPower, rounds, remainingHp, rounds*remainingHp)
}
