package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
	//"strings"
	//"rand"
)
var input = "15.txt"

type Ingredient struct {
	Name string
	Capacity, Durability, Flavor, Texture, Calories int
}

type IngredientAmount struct {
	Ingredient Ingredient
	Amount int
}

type Cookie struct {
	Score int
	Ingredients []IngredientAmount
}

func (cookie Cookie) String() string {
	s := ""

	for _,ing := range cookie.Ingredients {
		s += " " + ing.Ingredient.Name + ": " + strconv.Itoa(ing.Amount)
	}
	return s
}

func main() {
	startTime := time.Now()
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		reg := regexp.MustCompile(`^(\w+): (\w+) ([\-0-9]+), (\w+) ([\-0-9]+), (\w+) ([\-0-9]+), (\w+) ([\-0-9]+), (\w+) ([\-0-9]+)$`)

		ingmap := make(map[string]Ingredient)
		names := []string{}

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				cap,_ := strconv.Atoi(groups[3])
				dur,_ := strconv.Atoi(groups[5])
				fla,_ := strconv.Atoi(groups[7])
				tex,_ := strconv.Atoi(groups[9])
				cal,_ := strconv.Atoi(groups[11])

				ing := Ingredient{ Name: groups[1], Capacity: cap, Durability: dur, Flavor: fla, Texture: tex, Calories: cal }

				names = append(names, groups[1])
				ingmap[ing.Name] = ing
			}
		}

		allcookies := []Cookie{}
		max := 0
		calmax := 0

		perms := Perms(100, 4)

		fmt.Println(len(perms))
		return

		ints := make([]int, len(names))
		for i := 0; i < len(names); i++ {
			ints[i] = 0
		}
		iterations := 0

		for {
			ints = NextAmount(ints, 100)
			if ints == nil {
				break
			}

			cookie := Cookie{ Score: 0 }
			for i,n := range names{
				ing := ingmap[n]
				amt := ints[i]
				ingamt := IngredientAmount{ Amount: amt, Ingredient: ing }
				cookie.Ingredients = append(cookie.Ingredients, ingamt)
			}

			cap := 0
			dur := 0
			fla := 0
			tex := 0
			cal := 0
			for _,c := range cookie.Ingredients {
				ing := c.Ingredient
				amt := c.Amount

				cap += ing.Capacity * amt
				dur += ing.Durability * amt
				fla += ing.Flavor * amt
				tex += ing.Texture * amt
				cal += ing.Calories * amt
			}

			if cap < 0 {
				cap = 0
			}

			if dur < 0 {
				dur = 0
			}

			if fla < 0 {
				fla = 0
			}

			if tex < 0 {
				tex = 0
			}

			if cal < 0 {
				cal = 0
			}

			cookie.Score = cap * dur * fla * tex

			if cookie.Score > max {
				max = cookie.Score
				fmt.Println("new max score is ", max)
				fmt.Println(cookie)
			}

			if cal == 500 && cookie.Score > calmax {
				calmax = cookie.Score
				fmt.Println("-----new calorie max is ", calmax)
				fmt.Println(cookie)
			}

			iterations++

			if iterations % 1000000 == 0 {
				fmt.Println(iterations)
			}
		}

		fmt.Println("all cookie combos", len(allcookies))
		fmt.Println("max score", max)
	}

	fmt.Println("Time", time.Since(startTime))
}

func SumCheck(ints []int, total int) bool {
	sum := 0
	for _,i := range ints {
		sum += i
	}
	return sum == total
}

func NextAmount(ints []int, total int) []int {
	next, cont := Inc(&ints, len(ints)-1, total)
	for cont && !SumCheck(next, total){
		next,cont = Inc(&ints, len(ints)-1, total)
	}
	return next
}

func Perms(total, num int) [][]int {
	ret := [][]int{}
	cp := []int{}

	if num == 1 {
		cp = append(cp, total)
		ret = append(ret, cp)
		fmt.Println("ret", ret)
		return ret
	}

	for i := 0; i < total; i++ {
		if len(cp) < num {
			cp = append(cp, i)
			ret = append(ret, cp)
			perms := Perms(total - i, num-1)
			for _,p := range perms {
				p = append([]int{ total }, p...)
				ret = append(ret, p)
			}
		}
	}
	return ret
}

func Inc(ints *[]int, pos, wrap int) ([]int, bool) { // array, more
	if pos < 0 {
		return nil, false
	}
	(*ints)[pos]++
	if (*ints)[pos] == wrap {
		(*ints)[pos] = 0
		Inc(ints, pos-1, wrap)
	}
	return *ints, true
}
