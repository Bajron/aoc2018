package main

import (
	"bufio"
	"fmt"
	"os"
    "time"
    "sort"
)

type Log struct {
	date time.Time
    message string
}
type By func(p1, p2 *Log) bool

type logSorter struct {
	logs []Log
	by func(p1, p2 *Log) bool
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
		by: by,
	}
	sort.Sort(ls)
}

type GuardInfo struct {
    total int
    minutes [60]int
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
    
    for ;err==nil; {
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
    
    guards := make(map[int]GuardInfo)
    current := 0 
    state := 1
    sleepMinute := 0
    
    for _, l := range logs {
        if l.message[0] == 'G' {
            fmt.Sscanf(l.message, "Guard #%d", &current)
            state = 1
        } else if l.message[0] == 'f' {
            state = 0
            sleepMinute = l.date.Minute()
        } else if l.message[0] == 'w' {
            state = 1
        }
        
        if state == 0 {
            guards[current].total++
            guards[current].minutes[sleepMinute]++
            sleepMinute++
        }
    }
    
    fmt.Printf("%s %s", logs[0].date.String(), logs[0].message)
}
