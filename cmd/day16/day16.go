package main

import (
	"bufio"
	"fmt"
	"os"
)

type InOutTest struct {
	before, after [4]int
	op, A, B, C   int
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

	fmt.Printf("lines: %d\n", len(lines))

	ops := []opImpl{addr, addi, mulr, muli, bani, banr, bori, borr, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}

	behaveLike3OrMoreCount := 0

	for _, l := range lines {
		matches := 0
		for _, o := range ops {
			out := o(l.A, l.B, l.C, l.before)
			if l.after == out {
				matches++
			}
		}
		if matches >= 3 {
			behaveLike3OrMoreCount++
		}
	}

	fmt.Printf("behaveLike3OrMoreCount: %d\n", behaveLike3OrMoreCount)
}
