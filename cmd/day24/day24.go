package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Army struct {
	id           int
	units        int
	hppu         int
	weak, immune string
	damage       int
	attack       string
	initiative   int
	fraction     int
	attackedBy   *Army
	attacks      *Army
}

func (a Army) effectivePower() int {
	return a.units * a.damage
}

func (a Army) wouldDeal(b Army) int {
	base := a.effectivePower()
	if len(a.attack) > 0 {
		if strings.Contains(b.immune, a.attack) {
			return 0
		}
		if strings.Contains(b.weak, a.attack) {
			return 2 * base
		}
	}
	return base
}

func (a *Army) takeDamage(d int) {
	kills := d / a.hppu
	a.units -= kills
}

func readArmy(fraction int, re *regexp.Regexp, reader *bufio.Reader, log io.Writer) []Army {
	id := 1
	ret := []Army{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil || len(line) <= 1 {
			fmt.Fprintf(log, "error read line %s, line: %s\n", err, line)
			break
		}
		matches := re.FindAllStringSubmatch(line, 1)
		fmt.Printf("found %d matching to %s", len(matches), line)

		if len(matches) < 1 {
			continue
		}

		var army Army
		army.id = id
		id++
		army.units, err = strconv.Atoi(matches[0][1])
		army.hppu, err = strconv.Atoi(matches[0][2])

		if matches[0][5] == "weak" {
			army.weak = matches[0][6]
		}
		if matches[0][5] == "immune" {
			army.immune = matches[0][6]
		}
		if matches[0][9] == "weak" {
			army.weak = matches[0][10]
		}
		if matches[0][9] == "immune" {
			army.immune = matches[0][10]
		}
		army.damage, err = strconv.Atoi(matches[0][11])
		army.attack = matches[0][12]
		army.initiative, err = strconv.Atoi(matches[0][13])
		army.fraction = fraction
		army.attackedBy = nil
		army.attacks = nil
		ret = append(ret, army)
	}
	return ret
}

type By func(p1, p2 *Army) bool

type unitsSorter struct {
	units []*Army
	by    func(p1, p2 *Army) bool
}

func (s *unitsSorter) Len() int {
	return len(s.units)
}

func (s *unitsSorter) Swap(i, j int) {
	s.units[i], s.units[j] = s.units[j], s.units[i]
}

func (s *unitsSorter) Less(i, j int) bool {
	return s.by(s.units[i], s.units[j])
}

func (by By) Sort(units []*Army) {
	ls := &unitsSorter{
		units: units,
		by:    by,
	}
	sort.Sort(ls)
}

func main() {
	stderr := bufio.NewWriter(os.Stderr)
	stdout := bufio.NewWriter(os.Stdout)
	defer stderr.Flush()
	defer stdout.Flush()

	stdin := bufio.NewReader(os.Stdin)

	re, err := regexp.Compile(
		"(\\d+) units each with (\\d+) hit points " +
			"(\\(((immune|weak) to ([^;\\)]+))?(; )?((immune|weak) to ([^;\\)]+))?\\) )?with an attack that does (\\d+) (\\w+) damage at initiative (\\d+)")
	if err != nil {
		fmt.Fprintf(stderr, "error regex compile %s\n", err)
	}
	read, err := fmt.Fscanf(stdin, "Immune System:\n")
	if err != nil || read != 0 {
		fmt.Fprintf(stderr, "error fscanf %s\n", err)
	}
	immune := readArmy(0, re, stdin, stderr)

	read, err = fmt.Fscanf(stdin, "Infection:\n")
	if err != nil || read != 0 {
		fmt.Fprintf(stderr, "error fscanf %s\n", err)
	}
	infection := readArmy(1, re, stdin, stderr)

	fmt.Printf("%v %v\n", immune, infection)

	units := []*Army{}
	for k := range immune {
		units = append(units, &immune[k])
	}
	for k := range infection {
		units = append(units, &infection[k])
	}

	fmt.Printf("Total armies: %d\n", len(units))

	targetingOrder := func(lhs, rhs *Army) bool {
		return lhs.effectivePower() > rhs.effectivePower() || (lhs.effectivePower() == rhs.effectivePower() && lhs.initiative > rhs.initiative)
	}
	initiativeOrder := func(lhs, rhs *Army) bool {
		return lhs.initiative > rhs.initiative
	}

	immuneCount, infectionCount := len(immune), len(infection)

	for immuneCount > 0 && infectionCount > 0 {
		immuneCount, infectionCount = 0, 0

		for i := range units {
			u := units[i]
			u.attacks = nil
			u.attackedBy = nil

			if u.units <= 0 {
				continue
			}
			if u.fraction == 0 {
				immuneCount++
				fmt.Printf("Immune group %d has %d units\n", u.id, u.units)
			}
			if u.fraction == 1 {
				infectionCount++
				fmt.Printf("Infection group %d has %d units\n", u.id, u.units)
			}
		}

		By(targetingOrder).Sort(units)
		for i := range units {
			attacker := units[i]
			if attacker.units <= 0 {
				continue
			}

			fmt.Printf("Attacker %v\n", *attacker)

			choice := -1
			deals := 0
			for j := range units {
				c := units[j]
				if i == j || attacker.fraction == c.fraction || c.units <= 0 || c.attackedBy != nil {
					continue
				}
				d := attacker.wouldDeal(*c)
				fmt.Printf("?? Fraction %d, id %d would deal %d to f:%d id:%d\n",
					attacker.fraction, attacker.id, d, c.fraction, c.id)

				if d == 0 {
					continue
				}

				if choice < 0 || d > deals ||
					(d == deals && c.effectivePower() > units[choice].effectivePower()) ||
					(d == deals && c.effectivePower() == units[choice].effectivePower() && c.initiative > units[choice].initiative) {
					deals = d
					choice = j
				}
			}
			if choice >= 0 {
				attacker.attacks = units[choice]
				units[choice].attackedBy = attacker

				fmt.Printf("Fraction %d, id %d would deal %d to f:%d id:%d\n",
					attacker.fraction, attacker.id, deals, units[choice].fraction, units[choice].id)
			}
		}
		fmt.Printf(" -- attack phase --\n")
		By(initiativeOrder).Sort(units)
		for i := range units {
			attacker := units[i]
			if attacker.units <= 0 || attacker.attacks == nil {
				continue
			}

			victim := attacker.attacks
			dmg := attacker.wouldDeal(*victim)
			victim.takeDamage(dmg)
		}

		fmt.Printf("\n")
	}

	fmt.Printf("immune: %d infection: %d\n", immuneCount, infectionCount)

	immuneUnits, infectionUnits := 0, 0
	for i := range units {
		u := units[i]
		if u.units <= 0 {
			continue
		}
		if u.fraction == 0 {
			immuneUnits += u.units
		}
		if u.fraction == 1 {
			infectionUnits += u.units
		}
	}
	fmt.Printf("immune: %d infection: %d units\n", immuneUnits, infectionUnits)
}
