package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	//"strings"
	"sort"
)
var input = "9.txt"

type Route struct {
	Source string
	Destination string
	Miles int
}

func (route Route) String() string{
	return route.Source + " to " + route.Destination + " is " + strconv.Itoa(route.Miles) + " miles"
}

func (route Route) Key() string {
	return route.Source + "-" + route.Destination
}

type FullRoute struct {
	Cities []string
	Miles int
}

type FullRouteList []FullRoute

type FullRouteListSorter struct {
	Entries FullRouteList
}
func (p FullRouteListSorter) Len() int {
	return len(p.Entries)
}
func (p FullRouteListSorter) Less(i, j int) bool {
	return p.Entries[i].Miles < p.Entries[j].Miles
}
func (p FullRouteListSorter) Swap(i, j int) {
	p.Entries[i], p.Entries[j] = p.Entries[j], p.Entries[i]
}


func (fullRoute FullRoute) String() string {
	s := ""
	for _,c := range fullRoute.Cities {
		s += c + "-> "
	}
	s += "=" + strconv.Itoa(fullRoute.Miles)
	return s
}

func main() {
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)

		reg,_ := regexp.Compile(`^(.*?) to (.*?) = ([0-9]+)$`)

		cities := []string{}
		routes := make(map[string]Route)

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				miles,err := strconv.Atoi(groups[3])
				if err != nil {
					panic(err)
				}
				route := Route{Source: groups[1], Destination: groups[2], Miles: miles }
				routes[route.Key()] = route

				if !Contains(cities, route.Source) {
					cities = append(cities, route.Source)
				}

				if !Contains(cities, route.Destination) {
					cities = append(cities, route.Destination)
				}
			}
		}

		permutations := Permutate(cities)
		list := FullRouteList{}

		for _,perm := range permutations {
			fullRoute := FullRoute{}
			fullRoute.Cities = append(fullRoute.Cities, perm[0])

			for i := 1; i < len(perm); i++ {
				key := perm[i-1] + "-" + perm[i]
				key2 := perm[i] + "-" + perm[i-1]
				var route Route
				if r, ok := routes[key]; ok {
					route = r
				} else if r, ok := routes[key2]; ok {
					route = r
				}

				fullRoute.Cities = append(fullRoute.Cities, perm[i])
				fullRoute.Miles += route.Miles
			}

			list = append(list, fullRoute)
		}

		fmt.Println(len(list))

		sorter := FullRouteListSorter{ Entries: list }
		sort.Sort(sorter)

		fmt.Println(sorter.Entries[0])
		fmt.Println(sorter.Entries[len(sorter.Entries)-1])
	}
}

//tmp = append([]string{ str[0] }, Permutate(str[1:])...)

func Permutate(str []string) [][]string {
	var ret [][]string

	if len(str) == 2 {
		ret = append(ret, []string{ str[0], str[1] })
		ret = append(ret, []string{ str[1], str[0] })
	} else {

		for i := 0; i < len(str); i++ {
			strc := make([]string, len(str))
			copy(strc, str)

			t := strc[i]
			sh := append(strc[:i], strc[i+1:]...)
			perm := Permutate(sh)
			
			for _,p := range perm {
				p = append([]string{ t }, p...)
				ret = append(ret, p)
			}
		}
	}

	return ret
}

func Contains(ss []string, s string) bool {
	for _,t := range ss {
		if t == s {
			return true
		}
	}
	return false
}