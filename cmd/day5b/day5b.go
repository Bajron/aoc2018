package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func react(input string) int {
	var sp int = 0
	var s [50100]rune

	for _, c := range input {
		s[sp] = c
		sp++

		for sp > 1 {
			cur, prev := s[sp-1:sp], s[sp-2:sp-1]
			if cur[0] != prev[0] && unicode.ToLower(cur[0]) == unicode.ToLower(prev[0]) {
				sp -= 2
			} else {
				break
			}
		}
	}
	return sp
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)
	input, _ := stdin.ReadString('\n')
	input = input[:len(input)-1]

	lInput := strings.ToLower(input)

	letters := "abcdefghijklmnopqrstuvwxyz"

	bestLetter := letters[0]
	bestLen := len(input)

	for il := range letters {
		l := letters[il]
		ii := make([]byte, 0, len(input))
		for i := range input {
			if l == lInput[i] {
				continue
			}
			ii = append(ii, input[i])
		}
		cur := react(string(ii))
		if cur < bestLen {
			bestLen = cur
			bestLetter = l
		}
	}

	fmt.Printf("len: %d %c\n", bestLen, bestLetter)
}
