package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func distance(lhs, rhs string) (int, string) {
	if len(lhs) != len(rhs) {
		panic("bad input")
	}
	different := 0
	common := strings.Builder{}
	l, r := strings.NewReader(lhs), strings.NewReader(rhs)

	for {
		left, _, err_l := l.ReadRune()
		right, _, err_r := r.ReadRune()

		if err_l != nil || err_r != nil {
			break
		}

		if left == right {
			common.WriteRune(left)
		} else {
			different++
		}
	}
	return different, common.String()
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)

	ids := make([]string, 0, 2000)
	var line string

	read, err := fmt.Fscanf(stdin, "%s\n", &line)
	for err == nil && read == 1 {
		ids = append(ids, line)
		read, err = fmt.Fscanf(stdin, "%s\n", &line)
	}
	fmt.Fprintf(stderr, "Stopped after reading %d lines, with %s, read %d\n", len(ids), err, read)

	for _, id1 := range ids {
		for _, id2 := range ids {
			different, common := distance(id1, id2)
			if different == 1 {
				fmt.Printf("%s\n", common)
				os.Exit(0)
			}
		}
	}
}
