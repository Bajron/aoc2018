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
	fmt.Printf("ip: #%d\n", ip)

	ops := map[string]opImpl{"addr": addr, "addi": addi,
		"mulr": mulr, "muli": muli, "bani": bani, "banr": banr, "bori": bori, "borr": borr, "setr": setr, "seti": seti, "gtir": gtir, "gtri": gtri, "gtrr": gtrr,
		"eqir": eqir, "eqri": eqri, "eqrr": eqrr}

	// 2159153 at ip:28
	registers := [6]int{0, 0, 0, 0, 0, 0}

	stopWith := map[int]int{} // B/A value : ticks

	ticks := 0
	for 0 <= registers[ip] && registers[ip] < len(lines) {
		if registers[ip] == 28 {
			if _, ok := stopWith[registers[1]]; ok {
				// loop?
				break
			} else {
				fmt.Printf("%v\n", registers)
				stopWith[registers[1]] = ticks
			}
		}
		//fmt.Printf("%v\n", registers)

		registers = execute(lines[registers[ip]], registers, ops)
		registers[ip]++
		ticks++
	}

	bestRegisterSet := -1
	maxTicks := 0
	for k, v := range stopWith {
		if v > maxTicks {
			bestRegisterSet, maxTicks = k, v
		}
	}

	fmt.Printf("maxTicks: for B:%d ticks %d\n", bestRegisterSet, maxTicks)

	fmt.Printf("%v ticks %d\n", registers, ticks)
}
