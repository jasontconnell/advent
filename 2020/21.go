package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

var input = "21.txt"

type food struct {
	ingredients map[string]string
	allergens   map[string]string
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

	foods := getFoods(lines)
	p1, p2 := solve(foods)
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", strings.Join(p2, ","))
	fmt.Println("Time", time.Since(startTime))
}

func solve(foods []food) (int, []string) {
	possibilities := make(map[string]map[string]string)
	for _, f := range foods {
		for _, a := range f.allergens {
			if _, ok := possibilities[a]; !ok {
				possibilities[a] = make(map[string]string)
			}
			for ing, _ := range f.ingredients {
				possibilities[a][ing] = ing
			}
		}
	}

	done := false

	for !done {
		for a, m := range possibilities {
			for ing := range m {
				if !isPossible(foods, ing, a) {
					delete(m, ing)
				}
			}
		}

		done = true
		for a, m := range possibilities {
			if len(m) == 1 {
				matched := ""
				for _, v := range m {
					matched = v
				}
				remove(possibilities, a, matched)
			}
			done = done && len(m) == 1
		}
	}

	list := []string{}
	amap := make(map[string]string)
	for a, m := range possibilities {
		for v := range m {
			amap[v] = a
			list = append(list, v)
		}
	}

	count := 0
	for _, f := range foods {
		for _, ing := range f.ingredients {
			if _, ok := amap[ing]; !ok {
				count++
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		ai := amap[list[i]]
		aj := amap[list[j]]

		return ai < aj
	})

	return count, list
}

func remove(p map[string]map[string]string, notKey, value string) {
	for a, m := range p {
		if a == notKey {
			continue
		}

		delete(m, value)
	}
}

func isPossible(foods []food, ing string, allergen string) bool {
	ret := true
	for _, f := range foods {
		_, hasAllergen := f.allergens[allergen]
		_, hasIngredient := f.ingredients[ing]
		if hasAllergen && !hasIngredient {
			ret = false
		}
	}
	return ret
}

func getFoods(lines []string) []food {
	foods := []food{}
	//reg := regexp.MustCompile("([a-z]+)")
	for _, line := range lines {
		r := strings.Replace(line, "(", "", -1)
		r = strings.Replace(r, ")", "", -1)
		r = strings.Replace(r, ",", "", -1)
		flds := strings.Fields(r)

		contains := false
		f := food{}
		f.ingredients = make(map[string]string)
		f.allergens = make(map[string]string)
		for _, fld := range flds {
			if fld == "contains" {
				contains = true
				continue
			}
			if !contains {
				f.ingredients[fld] = fld
			} else {
				f.allergens[fld] = fld
			}
		}
		foods = append(foods, f)
	}
	return foods
}

// reg := regexp.MustCompile("-?[0-9]+")
/*
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
*/
