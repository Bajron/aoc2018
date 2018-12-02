package main

import (
	"bufio"
	"fmt"
	"os"
)

type Instruction struct {
    operation rune
    value int64
}

func (i Instruction) apply(initial int64) int64 {
	if i.operation == '-' {
		return initial - i.value
	} else if i.operation == '+' {
		return initial + i.value
	} else {
		return initial
	}
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
    defer stderr.Flush()
    
    stdin := bufio.NewReader(os.Stdin)
    
	var sum, number int64 = 0, 0
	var op rune = ' '

    input := make([]Instruction, 0, 2000)
    
	read, err := fmt.Fscanf(stdin, "%c%d\n", &op, &number)
	if err == nil && read >= 2 {
        input = append(input, Instruction{op, number})
	} else {
		fmt.Fprintf(stderr, "%s got: %d\n", err, read)
	}

	for err == nil {
		read, err = fmt.Fscanf(stdin, "%c%d\n", &op, &number)
		if err == nil && read >= 2 {
            input = append(input, Instruction{op, number})
		} else {
			fmt.Fprintf(stderr, "%s got: %d\n", err, read)
		}
	}
	
    received := make(map[int64]bool)
    received[sum] = true
    
    for x := 0; x < len(input); x++ {
        for i := range input {
            sum = input[i].apply(sum)
            if received[sum] {
                fmt.Printf("%d\n", sum)
                os.Exit(0)
            }
            received[sum] = true
        }
    }
    
    fmt.Fprintf(stderr, "Lines processed: %d\n", len(input))
}
