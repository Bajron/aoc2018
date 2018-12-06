package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
    x,y int
}

type Grain struct {
    distance int
    owner int // -1 not set, -2 tied
}

type Arena struct {
    x0, y0, x1, y1 int
    grains []Grain
}

func NewArena(x0, y0, x1, y1) *Arena {
    w := x1 - x0 + 1
    h := y1 - y0 + 1
    arena := make([]Grain, w*h) 
    for _, g := range arena {
        g = -1
    }
    return &Arena{x0,y0, x1, y1, arena}
}

func (arena Arena) translate(pos Point) (location int) {
    return (pos.x - x0) + (pox.y - y0) * (y1 - y0 + 1)
}

func (arena Arena) at(pos Point) Grain {
    return arena.grains[arena.translate(pos)]
}

func (arena *Arena) set(pos Point, owner, distance int) {
    arena.grains[arena.translate(pos)] = Grain{distance, owner}
    return
}


func (arena *Arena) put(pos Point, id int, distance int) (claimed int) {
    grain := arena.at(pos)
    if grain.owner == -1 {
        
    }
}

type FloodState {
    id int
    distance int
    border []Point
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
    
    for ;error == nil; {
        read, error = fmt.Fscanf(stdin, "%d, %d\n", &x, &y)
        if read == 2 && error == nil {
            points = append(points, Point{x, y})
        } else {
            fmt.Fprintf(stderr, "Read %d, error: %s\n", read, error)
        }
    }
    fmt.Fprintf(stderr, "Read %d points\n", len(points))
    
    var x0,y0, x1, y1 int = 1000000, 1000000, 0, 0
    for _, p := points {
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
    
	//fmt.Printf("len: %d %s\n", sp, string(s[:sp]))
}
