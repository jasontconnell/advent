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

		max := 0
		calmax := 0

		perms := Perms(100, 4)

		for _,p := range perms {
			cookie := Cookie{ Score: 0 }
			for i,n := range names {
				ing := ingmap[n]
				amt := p[i]
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
			}

			if cal == 500 && cookie.Score > calmax {
				calmax = cookie.Score
			}
		}

		fmt.Println("all cookie combos", len(perms))
		fmt.Println("max score", max)
		fmt.Println("max calorie score", calmax)
	}

	fmt.Println("Time", time.Since(startTime))
}

func Perms(total, num int) [][]int {
	ret := [][]int{}

	if num == 2 {
		for i := 0; i < total/2 + 1; i++ {
			ret = append(ret, []int{ total-i, i })
			if i != total - i {
				ret = append(ret, []int{ i, total-i })
			}
		}
	} else {
		for i := 0; i <= total; i++ {
			perms := Perms(total-i, num-1)
			for _, p := range perms {
				q := append([]int{ i }, p...)
				ret = append(ret, q)
			}
		}
	}
	return ret
}
