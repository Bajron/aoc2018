package main

import (
	"bufio"
	"fmt"
	"os"
)

type InOutTest struct {
    before, after [4]int
    op, A, B, C int
}

type ProgramLine struct {
    op, A, B, C int
}

type opImpl func(A, B, C int, input [4]int) [4]int 

func addr(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] + input[B]
    return ret
}

func addi(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] + B
    return ret
}

func mulr(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] * input[B]
    return ret
}

func muli(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] * B
    return ret
}

func banr(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] & input[B]
    return ret
}

func bani(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] & B
    return ret
}

func borr(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] | input[B]
    return ret
}

func bori(A, B, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A] | B
    return ret
}

func setr(A, _, C int, input [4]int) [4]int {
    ret := input
    ret[C] = input[A]
    return ret
}

func seti(A, _, C int, input [4]int) [4]int {
    ret := input
    ret[C] = A
    return ret
}

func gtrr(A, B, C int, input [4]int) [4]int {
    ret := input
    if input[A] > input[B] {
        ret[C] = 1
    } else {
        ret[C] = 0
    }
    return ret
}

func gtir(A, B, C int, input [4]int) [4]int {
    ret := input
    if A > input[B] {
        ret[C] = 1
    } else {
        ret[C] = 0
    }
    return ret
}

func gtri(A, B, C int, input [4]int) [4]int {
    ret := input
    if input[A] > B {
        ret[C] = 1
    } else {
        ret[C] = 0
    }
    return ret
}

func eqrr(A, B, C int, input [4]int) [4]int {
    ret := input
    if input[A] == input[B] {
        ret[C] = 1
    } else {
        ret[C] = 0
    }
    return ret
}

func eqir(A, B, C int, input [4]int) [4]int {
    ret := input
    if A == input[B] {
        ret[C] = 1
    } else {
        ret[C] = 0
    }
    return ret
}

func eqri(A, B, C int, input [4]int) [4]int {
    ret := input
    if input[A] == B {
        ret[C] = 1
    } else {
        ret[C] = 0
    }
    return ret
}


func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var lines []InOutTest
    var tmp InOutTest

	read, err := fmt.Fscanf(stdin, "Before: [%d, %d, %d, %d]\n%d %d %d %d\nAfter:  [%d, %d, %d, %d]\n\n",
        &tmp.before[0], &tmp.before[1], &tmp.before[2], &tmp.before[3],
        &tmp.op, &tmp.A, &tmp.B, &tmp.C,
        &tmp.after[0], &tmp.after[1], &tmp.after[2], &tmp.after[3])
    
	if read == 12 && err == nil {
		lines = append(lines, tmp)
	} else {
		fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
	}
	for err == nil {
		read, err = fmt.Fscanf(stdin, "Before: [%d, %d, %d, %d]\n%d %d %d %d\nAfter:  [%d, %d, %d, %d]\n\n",
            &tmp.before[0], &tmp.before[1], &tmp.before[2], &tmp.before[3],
            &tmp.op, &tmp.A, &tmp.B, &tmp.C,
            &tmp.after[0], &tmp.after[1], &tmp.after[2], &tmp.after[3])
    
        if read == 12 && err == nil {
            lines = append(lines, tmp)
        } else {
            fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
        }
	}

    stdin.ReadString('\n')
    stdin.ReadString('\n')
    
    var t ProgramLine
    var program []ProgramLine
    
    read, err = fmt.Fscanf(stdin, "%d %d %d %d\n", &t.op, &t.A, &t.B, &t.C)
    if read == 4 && err == nil {
		program = append(program, t)
	} else {
		fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
	}
    for err == nil {
		read, err = fmt.Fscanf(stdin, "%d %d %d %d\n", &t.op, &t.A, &t.B, &t.C)
    
        if read == 4 && err == nil {
            program = append(program, t)
        } else {
            fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
        }
	}
    
    
	fmt.Printf("lines: %d\n", len(lines))
    fmt.Printf("program lines: %d\n", len(program))

    
    ops := []opImpl{ addr, addi, mulr, muli, bani, banr, bori, borr, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}
    opMap := map[int]int{} // from op to index in ops
    opStat := [16]map[int]bool{} // "op could be at"
    
    for i := range opStat {
        opStat[i] = make(map[int]bool)
    }
    
    lastMatched := -1
    lastOperation := -1
    
    for _, l := range lines {
        matches := 0
        for oi, o := range ops {
            out := o(l.A, l.B, l.C, l.before)
            
            if l.after == out {
                opStat[l.op][oi] = true
                matches++
                lastMatched = oi
                lastOperation = l.op
            }
        }
        if matches == 1 {
            opMap[lastOperation] = lastMatched
        }
    }
    fmt.Printf("%v\n", opMap)
    fmt.Printf("%v\n", opStat)
    
    nextLines := []InOutTest{}
    for _, l := range lines {
        if _,ok := opMap[l.op]; ok {
            continue
        }
        nextLines = append(nextLines, l)
    }
    lines, nextLines = nextLines, lines
    fmt.Printf("lines: %d\n", len(lines))
    for _, l := range lines {
        matches := 0
        for oi, o := range ops {
            out := o(l.A, l.B, l.C, l.before)
            if l.after == out {
                matches++
                lastMatched = oi
                lastOperation = l.op
            }
        }
        if matches == 1 {
            opMap[lastOperation] = lastMatched
        } else {
            fmt.Printf("matches %d\n", matches)
        }
    }
    fmt.Printf("%v\n", opMap)
}
