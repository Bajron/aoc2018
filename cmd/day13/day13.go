package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Arena struct {
	roads [][]rune
	carts Carts
}

type Cart struct {
	id               int
	currentPosition  Point
	currentDirection rune
	lastTurn         Turn
}

type Carts []*Cart
type Turn int

const (
	Left = iota
	Straight
	Right
	TurnCount
)

type Point struct {
	x, y int
}

func NewArena(roads [][]rune) *Arena {
	cartId := 0
	var carts []*Cart
	for y, row := range roads {
		for x, b := range row {
			if isCart(b) {
				carts = append(carts, &Cart{cartId, Point{x, y}, b, Right})
				cartId++
				roads[y][x] = cartToRoad(b)
			}
		}
	}
	return &Arena{roads, carts}
}

func (arena *Arena) Tick() bool {
	sort.Sort(arena.carts)

	occupied := map[Point]int{}
	for _, c := range arena.carts {
		occupied[c.currentPosition] = c.id
	}

	for _, c := range arena.carts {
		delete(occupied, c.currentPosition)
		switch c.currentDirection {
		case '>':
			c.currentPosition.x++
		case '<':
			c.currentPosition.x--
		case '^':
			c.currentPosition.y--
		case 'v':
			c.currentPosition.y++
		}
		if _, ok := occupied[c.currentPosition]; ok {
			fmt.Printf("Crash at %d,%d\n", c.currentPosition.x, c.currentPosition.y)
			return true
		} else {
			occupied[c.currentPosition] = c.id
		}

		c.UpdateDirection(arena.roads[c.currentPosition.y][c.currentPosition.x])
	}
	return false
}

func (c *Cart) UpdateDirection(arrivedAt rune) {
	if arrivedAt == '+' {
		c.lastTurn++
		c.lastTurn = c.lastTurn % TurnCount
		c.currentDirection = calculateDirection(c.currentDirection, c.lastTurn)

		return
	}

	switch c.currentDirection {
	case '>':
		switch arrivedAt {
		case '\\':
			c.currentDirection = 'v'
		case '/':
			c.currentDirection = '^'
		}
	case '<':
		switch arrivedAt {
		case '\\':
			c.currentDirection = '^'
		case '/':
			c.currentDirection = 'v'
		}
	case '^':
		switch arrivedAt {
		case '\\':
			c.currentDirection = '<'
		case '/':
			c.currentDirection = '>'
		}
	case 'v':
		switch arrivedAt {
		case '\\':
			c.currentDirection = '>'
		case '/':
			c.currentDirection = '<'
		}
	}
	return
}

func calculateDirection(direction rune, turn Turn) rune {
	switch direction {
	case '>':
		switch turn {
		case Left:
			return '^'
		case Right:
			return 'v'
		default:
			return direction
		}
	case '<':
		switch turn {
		case Left:
			return 'v'
		case Right:
			return '^'
		default:
			return direction
		}
	case '^':
		switch turn {
		case Left:
			return '<'
		case Right:
			return '>'
		default:
			return direction
		}
	case 'v':
		switch turn {
		case Left:
			return '>'
		case Right:
			return '<'
		default:
			return direction
		}
	}
	return direction
}

func isCart(b rune) bool {
	return b == '^' || b == 'v' || b == '<' || b == '>'
}

func cartToRoad(b rune) rune {
	if b == '^' || b == 'v' {
		return '|'
	}
	if b == '<' || b == '>' {
		return '-'
	}
	return b
}

func (c Carts) Len() int {
	return len(c)
}

func (c Carts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Carts) Less(i, j int) bool {
	return c[i].currentPosition.y < c[j].currentPosition.y || (c[i].currentPosition.y == c[j].currentPosition.y && c[i].currentPosition.x < c[j].currentPosition.x)
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var roads [][]rune

	line, err := stdin.ReadString('\n')
	if err == nil {
		roads = append(roads, []rune(line))
	} else {
		fmt.Fprintf(stderr, "Read error: %s\n", err)
	}
	for err == nil {
		line, err = stdin.ReadString('\n')
		if err == nil {
			roads = append(roads, []rune(line))
		} else {
			fmt.Fprintf(stderr, "Read error: %s\n", err)
		}
	}

	fmt.Printf("lines: %d\n", len(roads))

	arena := NewArena(roads)
	crashed := false
	for !crashed {
		crashed = arena.Tick()
	}
}
