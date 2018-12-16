package main

import (
	"bufio"
	"fmt"
	"os"
)

type Patch struct {
	id   int
	x, y int
	w, h int
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var bufPatch Patch
	var patches []Patch

	read, err := fmt.Fscanf(stdin, "#%d @ %d,%d: %dx%d\n", &bufPatch.id, &bufPatch.x, &bufPatch.y, &bufPatch.w, &bufPatch.h)
	if err == nil && read >= 2 {
		patches = append(patches, bufPatch)
	} else {
		fmt.Fprintf(stderr, "%s got: %d\n", err, read)
	}

	for err == nil {
		read, err = fmt.Fscanf(stdin, "#%d @ %d,%d: %dx%d\n", &bufPatch.id, &bufPatch.x, &bufPatch.y, &bufPatch.w, &bufPatch.h)
		if err == nil && read >= 2 {
			patches = append(patches, bufPatch)
		} else {
			fmt.Fprintf(stderr, "%s got: %d\n", err, read)
		}
	}

	fmt.Fprintf(stderr, "Lines total: %d\n", len(patches))

	heatMap := [1024][1024]int{}

	for ir := range heatMap {
		for ic := range heatMap[ir] {
			heatMap[ir][ic] = 0
		}
	}

	for _, p := range patches {
		for i := p.x; i < p.x+p.w; i++ {
			for j := p.y; j < p.y+p.h; j++ {
				heatMap[i][j]++
			}
		}
	}
	var collisionCount int = 0

	for _, r := range heatMap {
		for _, c := range r {
			if c > 1 {
				collisionCount++
			}
		}
	}

	fmt.Printf("%d\n", collisionCount)
}
