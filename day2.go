package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	stderr := bufio.NewWriter(os.Stderr)
    defer stderr.Flush()
    
    stdin := bufio.NewReader(os.Stdin)
    
    ids := make([]string, 0, 2000)
    var line string
    
    read, err := fmt.Fscanf(stdin, "%s\n", &line)
    for ;err == nil && read == 1; {
        ids = append(ids, line)
        read, err = fmt.Fscanf(stdin, "%s\n", &line)
    }
    fmt.Fprintf(stderr, "Stopped after reading %d lines, with %s, read %d\n", len(ids), err, read)
    
    var count2, count3 int64 = 0, 0
    for _, id := range ids {
        counts := make(map[rune]int)
        for _, c := range id {
            counts[c]++
        }
        has2, has3 := false, false
        for l := range counts {
            if counts[l] == 2 {
                has2 = true
            } else if counts[l] == 3 {
                has3 = true
            }
        }
        if has2 {
            count2++
        }
        if has3 {
            count3++
        }
    }
    
    fmt.Printf("%d\n", count2 * count3)
}
