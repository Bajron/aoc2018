package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var initial string

	fmt.Fscanf(stdin, "initial state: %s\n", &initial)
	stdin.ReadString('\n')

	positivePatterns := map[string]bool{}
	var pattern, result string

	r, err := fmt.Fscanf(stdin, "%5s => %s\n", &pattern, &result)
	if r >= 2 && err == nil {
		if result[0] == '#' {
			positivePatterns[pattern] = true
		}
	} else {
		fmt.Fprintf(stderr, "Read %d, error: %s\n", r, err)
	}
	for err == nil {
		r, err = fmt.Fscanf(stdin, "%5s => %s\n", &pattern, &result)
		if r >= 2 && err == nil {
			if result[0] == '#' {
				positivePatterns[pattern] = true
			}
		} else {
			fmt.Fprintf(stderr, "Read %d, error: %s\n", r, err)
		}
	}

	offset := 80
	lenTotal := offset + len(initial) + offset
	buf1, buf2 := make([]byte, lenTotal), make([]byte, lenTotal)

	for i := 0; i < lenTotal; i++ {
		buf1[i] = '.'
		buf2[i] = '.'
	}

	for i := offset; i < offset+len(initial); i++ {
		buf1[i] = initial[i-offset]
	}

	for generation := 0; generation < 50000000000; generation++ {
		for i := 0; i < lenTotal-5; i++ {
			if positivePatterns[string(buf1[i:i+5])] {
				buf2[i+2] = '#'
			} else {
				buf2[i+2] = '.'
			}
		}
		buf1, buf2 = buf2, buf1
	}

	checkSum := 0
	for i := 0; i < lenTotal; i++ {
		if buf1[i] == '#' {
			checkSum += (i - offset)
		}
	}

	fmt.Printf("input len: %d, positive patterns: %d\n", len(initial), len(positivePatterns))
	fmt.Printf("end: %s\n", string(buf1))
	fmt.Printf("checkSum: %d\n", checkSum)
}
