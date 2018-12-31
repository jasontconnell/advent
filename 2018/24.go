package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var input = "24.txt"

type AttackType int

const (
	None AttackType = iota
	Bludgeoning
	Radiation
	Fire
	Cold
	Slashing
)

type group struct {
	id         int
	team       int
	units      int
	hp         int
	atktype    AttackType
	dmg        int
	weak       []AttackType
	immune     []AttackType
	initiative int
	name       string
}

func (g group) String() string {
	return g.name + fmt.Sprintf(" hp: %d units: %d atk: %s dmg: %d effpwr: %d init: %d weak: %s immune: %s", g.hp, g.units, attackstr(g.atktype), g.dmg, g.effectivePower(), g.initiative, resistancestr(g.weak), resistancestr(g.immune))
}

func (g group) effectivePower() int {
	return g.dmg * g.units
}

func (g group) effectiveDamage(a group) int {
	weak := false
	for _, at := range a.weak {
		if at == g.atktype {
			weak = true
		}
	}

	immune := false
	for _, at := range a.immune {
		if at == g.atktype {
			immune = true
		}
	}

	mult := 1
	if immune {
		mult = 0
	} else if weak {
		mult = 2
	}
	return g.units * g.dmg * mult
}

func reverse(groups []group) []group {
	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[j], groups[i] = groups[i], groups[j]
	}
	return groups
}

func sortGroupsEffPower(groups []group) []group {
	sfunc := func(i, j int) bool {
		epi, epj := groups[i].effectivePower(), groups[j].effectivePower()
		ini, inj := -groups[i].initiative, -groups[j].initiative

		if epi == epj {
			return ini < inj
		}
		return epi < epj
	}

	sort.Slice(groups, sfunc)
	return reverse(groups)
}

func sortGroupsEffDamage(attacker group, groups []group) []group {
	sfunc := func(i, j int) bool {
		edi, edj := attacker.effectiveDamage(groups[i]), attacker.effectiveDamage(groups[j])
		epi, epj := groups[i].effectivePower(), groups[j].effectivePower()
		ini, inj := -groups[i].initiative, -groups[j].initiative

		if edi == edj {
			if epi == epj {
				return ini < inj
			} else {
				return epi < epj
			}
		}
		return edi < edj
	}
	sort.Slice(groups, sfunc)

	return reverse(groups)
}

func sortGroupsInitiative(groups []group) []group {
	sfunc := func(i, j int) bool {
		ini, inj := groups[i].initiative, groups[j].initiative
		return ini < inj
	}
	sort.Slice(groups, sfunc)
	return reverse(groups)
}

func main() {
	startTime := time.Now()

	f, err := os.Open(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		var txt = scanner.Text()
		lines = append(lines, txt)
	}

	groups := parseGroups(lines)
	clone := append([]group{}, groups...)

	p1 := solve(groups)
	fmt.Println("Part 1:", p1)

	p2 := solvePart2(clone)
	fmt.Println("Part 2:", p2)

	fmt.Println("Time", time.Since(startTime))
}

func solve(groups []group) int {
	gclone := append([]group{}, groups...)
	gclone = sim(gclone)

	return livingUnits(gclone)
}

func solvePart2(groups []group) int {
	boost := 1
	team := 1
	ans := 0
	done := false
	for !done {
		fmt.Println("Trying boost", boost)
		grps := append([]group{}, groups...)
		grps = boostTeam(grps, team, boost)
		grps = sim(grps)
		if teamDead(grps) && winningTeam(grps) == team {
			ans = livingUnits(grps)
			fmt.Println("boost is", boost, grps)
			done = true
		}
		boost++
	}
	return ans
}

func winningTeam(groups []group) int {
	team := 0
	for _, g := range groups {
		if g.units > 0 {
			team = g.team
		}
	}
	return team
}

func boostTeam(groups []group, team, boost int) []group {
	for i := 0; i < len(groups); i++ {
		if groups[i].team == team {
			groups[i].dmg += boost
		}
	}
	return groups
}

func sim(groups []group) []group {
	gclone := append([]group{}, groups...)
	done := false
	var total int
	for !done {
		gclone = sortGroupsEffPower(gclone)
		emap := enemiesMap(gclone)

		atkmap := chooseAttacks(gclone, emap)
		gclone, total = attack(gclone, atkmap)

		done = teamDead(gclone) || len(atkmap) == 0 || total == 0
	}

	return gclone
}

func livingUnits(groups []group) int {
	ans := 0
	for _, g := range groups {
		if g.units > 0 {
			ans += g.units
		}
	}
	return ans
}

func teamDead(groups []group) bool {
	team1, team2 := true, true

	for _, g := range groups {
		if g.team == 1 && g.units > 0 {
			team1 = false
		}

		if g.team == 2 && g.units > 0 {
			team2 = false
		}
	}
	return team1 || team2
}

func attack(groups []group, m map[int]int) ([]group, int) {
	gclone := append([]group{}, groups...)
	gclone = sortGroupsInitiative(gclone)
	total := 0
	gmap := make(map[int]group)
	for _, g := range groups {
		gmap[g.id] = g
	}

	for _, g := range gclone {
		if g.units == 0 {
			continue
		}
		if agid, ok := m[g.id]; ok {
			ag, ok := gmap[agid] // get attacked group
			if !ok {
				break
			}

			if ag.units == 0 {
				// fmt.Println(g, "Can't attack", ag)
				continue
			}

			ed := g.effectiveDamage(ag)

			killed := ed / ag.hp
			if killed >= ag.units {
				killed = ag.units
			}

			// fmt.Println(g, ed, "attacked", ag, ag.units * ag.hp, "and killed", killed, g.atktype, ag.weak)
			total += killed
			ag.units = ag.units - killed
			gclone = updateUnits(ag, gclone)
		}
	}

	return gclone, total
}

func updateUnits(g group, groups []group) []group {
	for i := 0; i < len(groups); i++ {
		if groups[i].id == g.id {
			groups[i].units = g.units
		}
	}
	return groups
}

func chooseAttacks(groups []group, m map[int]map[int]int) map[int]int {
	atks := make(map[int]int)
	gmap := make(map[int]group)
	for _, g := range groups {
		gmap[g.id] = g
	}

	for _, g := range groups {
		enlist := enemies(g, groups)
		enlist = sortGroupsEffDamage(g, enlist)
		if g.effectivePower() == 0 {
			continue
		}
		// fmt.Println(g, "can attack", enlist)
		tm := m[g.team] // target map
		chosen := false
		for _, en := range enlist {
			if chosen {
				break
			}

			if f, ok := tm[en.id]; !ok || f != -1 {
				continue // already marked for attack
			}

			if g.effectiveDamage(en) == 0 {
				continue // can't do any damage ... because this is sorted by effective damage first, probably an end state
			}

			if g.effectiveDamage(en) < en.hp && len(atks) == 1 {
				// fmt.Println(g, " !!! ", en )
			}

			if en.units == 0 {
				continue // can't attack empty group
			}

			// fmt.Println(g.name, g.id, "chooses to attack", en.name, en.id)
			tm[en.id] = g.id // mark this enemy as being attacked by this g
			atks[g.id] = en.id
			chosen = true
		}
		m[g.team] = tm
	}
	return atks
}

func enemies(g group, groups []group) []group {
	enlist := []group{}
	for _, eg := range groups {
		if eg.team != g.team && g.units > 0 {
			enlist = append(enlist, eg)
		}
	}
	return enlist
}

// returns a map of that team's enemies, and whether that group is targeted
func enemiesMap(groups []group) map[int]map[int]int {
	m := make(map[int]map[int]int)
	for _, g := range groups {
		if g.units == 0 {
			continue
		}
		opposing := 1
		if g.team == 1 {
			opposing = 2
		}

		if _, ok := m[opposing]; !ok {
			m[opposing] = make(map[int]int)
		}

		m[opposing][g.id] = -1
	}
	return m
}

func parseGroups(lines []string) []group {
	groups := []group{}
	team := 1
	teamName := "Immune System"
	gc := 1
	for _, line := range lines {
		if line == "Immune System:" || line == "" {
			continue
		} else if line == "Infection:" {
			gc = 1
			team = 2
			teamName = "Infection"
			continue
		}

		g := parseGroup(line, gc, team, teamName)
		if g != nil {
			groups = append(groups, *g)
		}
		gc++
	}
	return groups
}

var reg *regexp.Regexp = regexp.MustCompile(`(\d+) units each with (\d+) hit points( ?\(?.*?\)?)? with an attack that does (\d+) (\w+) damage at initiative (\d+)`)

func parseGroup(line string, grp, team int, teamName string) *group {
	if groups := reg.FindStringSubmatch(line); groups != nil && len(groups) > 1 {
		groupName := "group " + strconv.Itoa(grp)
		c, _ := strconv.Atoi(groups[1])
		hp, _ := strconv.Atoi(groups[2])
		dmg, _ := strconv.Atoi(groups[4])
		atktype := getAttackType(groups[5])
		init, _ := strconv.Atoi(groups[6])

		id := grp + (team * 100)

		g := &group{id: id, name: teamName + " " + groupName, team: team, units: c, hp: hp, dmg: dmg, atktype: atktype, initiative: init}

		if len(groups[3]) > 0 {
			sub := strings.TrimSpace(groups[3])
			sub = string(sub[1 : len(sub)-1])
			sp := strings.Split(sub, ";")
			for _, s := range sp {
				s = strings.TrimSpace(s)
				if strings.HasPrefix(s, "weak to ") {
					g.weak = getAttackTypes(string(s[len("weak to "):]))
				} else if strings.HasPrefix(s, "immune to ") {
					g.immune = getAttackTypes(string(s[len("immune to "):]))
				}
			}
		}
		return g
	}
	return nil
}

func getAttackTypes(txt string) []AttackType {
	atks := []AttackType{}
	sp := strings.Split(txt, ",")
	for _, s := range sp {
		s = strings.TrimSpace(s)
		atks = append(atks, getAttackType(s))
	}
	return atks
}

func getAttackType(txt string) AttackType {
	atk := None
	switch txt {
	case "fire":
		atk = Fire
	case "bludgeoning":
		atk = Bludgeoning
	case "slashing":
		atk = Slashing
	case "radiation":
		atk = Radiation
	case "cold":
		atk = Cold
	}
	return atk
}

func resistancestr(atks []AttackType) string {
	s := ""
	for i, a := range atks {
		s = s + attackstr(a)
		if i < len(atks)-1 {
			s = s + ","
		}
	}
	return s
}

func attackstr(atk AttackType) string {
	var val string
	switch atk {
	case Fire:
		val = "Fire"
	case Bludgeoning:
		val = "Bludgeoning"
	case Slashing:
		val = "Slashing"
	case Cold:
		val = "Cold"
	case Radiation:
		val = "Radiation"
	}
	return val
}
