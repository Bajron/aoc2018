package main

import (
	"bufio"
	"fmt"
	"os"
)

type Rope []*string

type Node struct {
	text             string
	kind             byte
	children         []*Node
	generatedStrings []Rope
}

func MaxLen(node *Node) int {
	if len(node.children) == 0 {
		return len(node.text)
	}
	if node.text == "()" {
		maxLen := 0
		for _, ch := range node.children {
			ml := MaxLen(ch)
			if ml > maxLen {
				maxLen = ml
			}
		}
		return maxLen
	}

	l := 0
	for _, ch := range node.children {
		l += MaxLen(ch)
	}
	return l
}

func Dump(node *Node, tab string) {
	fmt.Printf("%s %s %s\n", tab, string(node.kind), node.text)

	for _, ch := range node.children {
		Dump(ch, tab+" ")
	}
	return
}

type Point struct {
	x, y int
}
type Door struct {
	from, to Point
}

func Move(p Point, direction byte) Point {
	if direction == 'W' {
		return Point{p.x - 1, p.y}
	}
	if direction == 'E' {
		return Point{p.x + 1, p.y}
	}
	if direction == 'N' {
		return Point{p.x, p.y + 1}
	}
	if direction == 'S' {
		return Point{p.x, p.y - 1}
	}
	panic("bad direction " + string(direction))
}

// func GenerateStrings(node *Node) []string {
// 	if node.generatedStrings != nil {
// 		return node.generatedStrings
// 	}

// 	if len(node.children) == 0 {
// 		fmt.Printf("generating leaf %s\n", node.text)
// 		node.generatedStrings = []string{node.text}
// 		return node.generatedStrings
// 	}

// 	if node.text == "()" {
// 		generated := make([]string, 0, 1000)
// 		for _, ch := range node.children {
// 			strs := GenerateStrings(ch)
// 			for _, s := range strs {
// 				generated = append(generated, s)
// 			}
// 		}
// 		node.generatedStrings = generated
// 		fmt.Printf("produced () %d\n", len(node.generatedStrings))
// 		return node.generatedStrings
// 	}

// 	starts := []string{""}
// 	nextStarts := []string{}
// 	for _, ch := range node.children {
// 		nextStarts := nextStarts[:0]
// 		tmp := GenerateStrings(ch)
// 		for _, s := range starts {
// 			for _, d := range tmp {
// 				nextStarts = append(nextStarts, s+d)
// 			}
// 		}
// 		starts, nextStarts = nextStarts, starts
// 	}
// 	node.generatedStrings = starts
// 	fmt.Printf("produced , %d\n", len(node.generatedStrings))
// 	return node.generatedStrings
// }

func GenerateSize(node *Node) int {
	if len(node.children) == 0 {
		return 1
	}

	if node.text == "()" {
		ret := 0
		for _, ch := range node.children {
			ret += GenerateSize(ch)
		}
		return ret
	}

	ret := 1
	for _, ch := range node.children {
		ret *= GenerateSize(ch)
	}
	return ret
}

func FindDoors(node *Node, pos Point, doors *map[int]map[int]bool) []Point {

	if len(node.children) == 0 {
		prev := pos
		for i := range node.text {
			pos = Move(pos, node.text[i])
			pr := Gnode(prev)
			po := Gnode(pos)
			//*doors = append(*doors, Door{prev, pos})

			if _, ok := (*doors)[pr]; !ok {
				(*doors)[pr] = map[int]bool{}
			}
			if _, ok := (*doors)[po]; !ok {
				(*doors)[po] = map[int]bool{}
			}

			(*doors)[pr][po] = true
			(*doors)[po][pr] = true
			prev = pos
		}

		//fmt.Printf(" doors now %d \n", len(*doors))
		return []Point{prev}
	}

	if node.text == "()" {
		finishes := make([]Point, 0, 1000)
		for _, ch := range node.children {
			f := FindDoors(ch, pos, doors)
			for _, d := range f {
				finishes = append(finishes, d)
			}
		}
		return finishes
	}

	starts := []Point{pos}
	for _, ch := range node.children {
		nextStarts := map[Point]bool{}
		for _, s := range starts {
			tmp := FindDoors(ch, s, doors)
			for _, d := range tmp {
				if d != pos {
					nextStarts[d] = true
				}
			}
		}
		starts = starts[:0]
		for p := range nextStarts {
			starts = append(starts, p)
		}
	}
	return starts
}

func Gnode(p Point) int {
	return p.x + 1000 + (p.y+1000)*10000
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	root := Node{"", '$', []*Node{}, nil}
	stack := []*Node{}
	sp := -1
	text := ""

	b, err := stdin.ReadByte()
	if b != '^' {
		panic("bad input")
	}

	stack = append(stack, &root)
	sp++

	for err == nil {
		//fmt.Printf("byte %s stack %v \n", string(b), stack)
		if b != '(' && b != '|' && b != ')' && b != '$' && b != '^' {
			text += string(b)
		}

		if b == '(' {
			if len(text) > 0 {
				stack[sp].children = append(stack[sp].children, &Node{text, 'a', []*Node{}, nil})
				text = ""
			}
			stack = append(stack, &Node{"()", '(', []*Node{}, nil})
			sp++

			stack = append(stack, &Node{"", ',', []*Node{}, nil})
			sp++
		}
		if b == '|' {
			if len(text) > 0 {
				stack[sp].children = append(stack[sp].children, &Node{text, 'a', []*Node{}, nil})
				text = ""
			}

			finished := stack[sp]
			//fmt.Printf("Finishing on |, %s appending to %s\n", string(finished.kind), string(stack[sp-1].kind))
			stack = stack[:sp]
			sp--
			stack[sp].children = append(stack[sp].children, finished)

			stack = append(stack, &Node{"", ',', []*Node{}, nil})
			sp++
		}
		if b == ')' || b == '$' {
			if len(text) > 0 {
				stack[sp].children = append(stack[sp].children, &Node{text, 'a', []*Node{}, nil})
				text = ""
			}

			if b != '$' {
				finished := stack[sp]
				//fmt.Printf("Finishing on %s, %s appending to %s\n", string(b), string(finished.kind), string(stack[sp-1].kind))
				stack = stack[:sp]
				sp--
				stack[sp].children = append(stack[sp].children, finished)

				finished = stack[sp]
				//fmt.Printf("Finishing on %s, %s appending to %s\n", string(b), string(finished.kind), string(stack[sp-1].kind))
				stack = stack[:sp]
				sp--
				stack[sp].children = append(stack[sp].children, finished)
			}
		}
		b, err = stdin.ReadByte()
	}
	fmt.Fprintf(stderr, "read error was %s\n", err)
	fmt.Printf("root has %d children\n", len(root.children))
	fmt.Printf("maxlen: %d\n", MaxLen(&root))

	//Dump(&root, "")

	doors := map[int]map[int]bool{}

	//strs := GenerateStrings(&root)
	fmt.Printf("generated set count %d\n", GenerateSize(&root))
	// for i, s := range strs {
	// 	fmt.Printf("%d: %s\n", i, s)
	// }

	//doors := []Door{}
	FindDoors(&root, Point{0, 0}, &doors)
	//fmt.Printf("doors %d / %v\n", len(doors), doors)

	gnodes := map[int]bool{}
	edges := doors
	for from := range edges {
		gnodes[from] = true
	}

	startingGnode := Gnode(Point{0, 0})

	fmt.Printf("edges %d, gnodes %d\n", len(edges), len(gnodes))
	fmt.Printf("starting room neighbors %d %v\n", len(edges[startingGnode]), edges[startingGnode])

	best := map[int]int{}
	for g := range gnodes {
		best[g] = 1000000
	}
	wave := []int{startingGnode}
	distance := 0
	for len(wave) > 0 {
		newWave := []int{}
		for _, w := range wave {
			if distance < best[w] {
				best[w] = distance

				for neighbor := range edges[w] {
					newWave = append(newWave, neighbor)
				}
			}
		}
		distance++
		wave = newWave
	}

	atLeast1000 := 0
	rm, dm := -1, 0
	for room, maxDistance := range best {
		if maxDistance > dm {
			dm = maxDistance
			rm = room
		}
		if maxDistance >= 1000 {
			atLeast1000++
		}
	}

	fmt.Printf("best %v\n", best)
	fmt.Printf("furthest room %d, with distance %d\n", rm, dm)
	fmt.Printf("at least 1000 %d\n", atLeast1000)
}
