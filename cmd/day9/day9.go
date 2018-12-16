package main

import (
	"container/list"
	"fmt"
)

func main() {
	// 427 70723
	var score [427]int
	lastMarble := 70723 * 100
	marble := 0

	l := list.New()
	current := l.PushBack(marble)
	marble++

	for (marble - 1) != lastMarble {
		for p := 0; p < len(score); p++ {
			if marble%23 == 0 {
				score[p] += marble
				marble++

				toRemove := current.Prev()
				if toRemove == nil {
					toRemove = l.Back()
				}
				for k := 1; k < 7; k++ {
					toRemove = toRemove.Prev()
					if toRemove == nil {
						toRemove = l.Back()
					}
				}

				current = toRemove.Next()
				if current == nil {
					current = l.Front()
				}

				v, _ := toRemove.Value.(int)
				score[p] += v
				l.Remove(toRemove)
			} else {
				insertPoint := current.Next()
				if insertPoint == nil {
					insertPoint = l.Front()
				}
				current = l.InsertAfter(marble, insertPoint)
				marble++
			}

			if (marble - 1) == lastMarble {
				break
			}
		}
	}

	maxScore := 0
	winner := -1
	for pi, points := range score {
		if points > maxScore {
			maxScore = points
			winner = pi
		}
	}

	fmt.Printf("Winner %d with score %d\n", winner+1, maxScore)
}
