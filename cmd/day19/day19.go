package main

import (
	"bufio"
	"fmt"
	"os"
)

type Line struct {
	op      string
	A, B, C int
}

type opImpl func(A, B, C int, input [6]int) [6]int

func addr(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] + input[B]
	return ret
}

func addi(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] + B
	return ret
}

func mulr(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] * input[B]
	return ret
}

func muli(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] * B
	return ret
}

func banr(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] & input[B]
	return ret
}

func bani(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] & B
	return ret
}

func borr(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] | input[B]
	return ret
}

func bori(A, B, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A] | B
	return ret
}

func setr(A, _, C int, input [6]int) [6]int {
	ret := input
	ret[C] = input[A]
	return ret
}

func seti(A, _, C int, input [6]int) [6]int {
	ret := input
	ret[C] = A
	return ret
}

func gtrr(A, B, C int, input [6]int) [6]int {
	ret := input
	if input[A] > input[B] {
		ret[C] = 1
	} else {
		ret[C] = 0
	}
	return ret
}

func gtir(A, B, C int, input [6]int) [6]int {
	ret := input
	if A > input[B] {
		ret[C] = 1
	} else {
		ret[C] = 0
	}
	return ret
}

func gtri(A, B, C int, input [6]int) [6]int {
	ret := input
	if input[A] > B {
		ret[C] = 1
	} else {
		ret[C] = 0
	}
	return ret
}

func eqrr(A, B, C int, input [6]int) [6]int {
	ret := input
	if input[A] == input[B] {
		ret[C] = 1
	} else {
		ret[C] = 0
	}
	return ret
}

func eqir(A, B, C int, input [6]int) [6]int {
	ret := input
	if A == input[B] {
		ret[C] = 1
	} else {
		ret[C] = 0
	}
	return ret
}

func eqri(A, B, C int, input [6]int) [6]int {
	ret := input
	if input[A] == B {
		ret[C] = 1
	} else {
		ret[C] = 0
	}
	return ret
}

func execute(l Line, input [6]int, ops map[string]opImpl) [6]int {
	return ops[l.op](l.A, l.B, l.C, input)
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var lines []Line
	var tmp Line
	ip := 0

	read, err := fmt.Fscanf(stdin, "#ip %d\n", &ip)
	if read != 1 || err != nil {
		fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
	}

	read, err = fmt.Fscanf(stdin, "%s %d %d %d\n",
		&tmp.op, &tmp.A, &tmp.B, &tmp.C)

	if read == 4 && err == nil {
		lines = append(lines, tmp)
	} else {
		fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
	}
	for err == nil {
		read, err = fmt.Fscanf(stdin, "%s %d %d %d\n",
			&tmp.op, &tmp.A, &tmp.B, &tmp.C)

		if read == 4 && err == nil {
			lines = append(lines, tmp)
		} else {
			fmt.Fprintf(stderr, "Read error: %s (read %d)\n", err, read)
		}
	}

	fmt.Printf("lines: %d\n", len(lines))

	ops := map[string]opImpl{"addr": addr, "addi": addi,
		"mulr": mulr, "muli": muli, "bani": bani, "banr": banr, "bori": bori, "borr": borr, "setr": setr, "seti": seti, "gtir": gtir, "gtri": gtri, "gtrr": gtrr,
		"eqir": eqir, "eqri": eqri, "eqrr": eqrr}

	// seems to be calculating occurences of x*y that x*y == 10551345

	// factors 3 5 7 317 317

	counter := 0
	for i := 1; i <= 10551345; i++ {
		if 10551345%i == 0 {
			counter += i
		}
	}
	fmt.Printf("counter %d\n", counter)

	// registers := [6]int{1,0,0,0,0,0}
	//registers := [6]int{0,0,0,0,0,0}
	//registers := [6]int{0, 1, 10551300, 0, 8, 10551345}
	registers := [6]int{1, 3, 3500000, 0, 10, 10551345}

	for 0 <= registers[ip] && registers[ip] < len(lines) {
		registers = execute(lines[registers[ip]], registers, ops)
		registers[ip]++
		//fmt.Printf("%v\n", registers)
	}

	fmt.Printf("%v\n", registers)
}
