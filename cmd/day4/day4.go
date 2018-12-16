package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

type Log struct {
	date    time.Time
	message string
}
type By func(p1, p2 *Log) bool

type logSorter struct {
	logs []Log
	by   func(p1, p2 *Log) bool
}

func (s *logSorter) Len() int {
	return len(s.logs)
}

func (s *logSorter) Swap(i, j int) {
	s.logs[i], s.logs[j] = s.logs[j], s.logs[i]
}

func (s *logSorter) Less(i, j int) bool {
	return s.by(&s.logs[i], &s.logs[j])
}

func (by By) Sort(logs []Log) {
	ls := &logSorter{
		logs: logs,
		by:   by,
	}
	sort.Sort(ls)
}

type GuardInfo struct {
	id      int
	total   int
	minutes [60]int
}

func (g GuardInfo) getMaxMinute() (maxMinute, maxTimes int) {
	maxMinute = 0
	maxTimes = 0
	for m, times := range g.minutes {
		if times > maxTimes {
			maxMinute = m
			maxTimes = times
		}
	}
	return
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	defer stderr.Flush()

	stdin := bufio.NewReader(os.Stdin)

	var year, month, day, hour, min int
	var message string
	var logs []Log

	read, err := fmt.Fscanf(stdin, "[%d-%d-%d %d:%d] ", &year, &month, &day, &hour, &min)
	if read >= 5 && err == nil {
		message, err = stdin.ReadString('\n')
		logs = append(logs, Log{time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC), message})
	} else {
		fmt.Fprintf(stderr, "%s got: %d\n", err, read)
	}

	for err == nil {
		read, err = fmt.Fscanf(stdin, "[%d-%d-%d %d:%d] ", &year, &month, &day, &hour, &min)
		if read >= 5 && err == nil {
			message, err = stdin.ReadString('\n')
			logs = append(logs, Log{time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC), message})
		} else {
			fmt.Fprintf(stderr, "%s got: %d\n", err, read)
		}
	}

	fmt.Fprintf(stderr, "Lines total: %d\n", len(logs))

	date := func(p1, p2 *Log) bool {
		return p1.date.Before(p2.date)
	}
	By(date).Sort(logs)

	guards := make(map[int]*GuardInfo)
	current := 0
	state := 1
	sleepMinute := 0

	for _, l := range logs {
		if l.message[0] == 'G' {
			if state == 0 {
				panic("Guard did not wake up at the end of shift")
			}
			fmt.Sscanf(l.message, "Guard #%d", &current)
			if guards[current] == nil {
				guards[current] = &GuardInfo{current, 0, [60]int{}}
			}
			state = 1
		} else if l.message[0] == 'f' {
			state = 0
			sleepMinute = l.date.Minute()
		} else if l.message[0] == 'w' {
			state = 1

			slept := l.date.Minute() - sleepMinute
			guards[current].total += slept
			for m := sleepMinute; m < l.date.Minute(); m++ {
				guards[current].minutes[m]++
			}

			// fmt.Fprintf(stderr, "#%d slept %d\n", current, slept)
		}
	}

	var maxMinutesGuard *GuardInfo
	maxMinutes := 0
	for _, g := range guards {
		if g.total > maxMinutes {
			maxMinutes = g.total
			maxMinutesGuard = g
		}
	}

	maxMinute, _ := maxMinutesGuard.getMaxMinute()

	fmt.Printf("#%d %d/%d %d\n", maxMinutesGuard.id, maxMinutesGuard.total, maxMinute, maxMinute*maxMinutesGuard.id)
}
