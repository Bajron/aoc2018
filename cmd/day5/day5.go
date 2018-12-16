package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)
	input, _ := stdin.ReadString('\n')
	input = input[:len(input)-1]

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

	fmt.Printf("len: %d %s\n", sp, string(s[:sp]))
}
