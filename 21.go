package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
	"sort"
	//"strings"
)
var input = "21.txt"
var shopfile = "21.shop.txt"

type ShopItem struct {
	Name string
	Gold, Damage, Defense int
	Type int /// 0 weapon, 1 armor, 2 ring
}

var WeaponType int = 0
var ArmorType int = 1
var RingType int = 2

type Weapon struct {
	Name string
	Damage int
}

type Armor struct {
	Name string
	Defense int
}

type Ring struct {
	Name string
	Damage int
	Defense int
}

type Player struct {
	Name string
	HP int
	Damage int
	Defense int
	GoldSpent int
	Items []ShopItem
}

func (player Player) String() string {
	s := fmt.Sprintf("%v HP: %v DMG: %v DEF:%v GLDSPNT:%v\n", player.Name, player.HP, player.Damage, player.Defense, player.GoldSpent)
	for _,item := range player.Items {
		s += item.Name + "\n"
	}
	return s
}

type Strike struct {
	Striker, Strikee string
	DMG, DEF int
	HP int
}

type Variation struct {
	Player Player
	Boss Player
	BossDefeated bool
	GoldSpent int
	Strikes []Strike
}

type VariationList []Variation

type VariationListSorter struct {
	Entries VariationList
}
func (p VariationListSorter) Len() int {
	return len(p.Entries)
}
func (p VariationListSorter) Less(i, j int) bool {
	return p.Entries[i].GoldSpent < p.Entries[j].GoldSpent
}
func (p VariationListSorter) Swap(i, j int) {
	p.Entries[i], p.Entries[j] = p.Entries[j], p.Entries[i]
}


func (variation Variation) String() string {
	s := fmt.Sprintf("%v GLDSPNT:%v\n", variation.BossDefeated, variation.GoldSpent)
	s = s + fmt.Sprintf("Player %v\n", variation.Player)
	s = s + fmt.Sprintf("Boss   %v\n", variation.Boss)

	for _,strike := range variation.Strikes {
		s += fmt.Sprintf("%v hit for %v damage, %v had %v defense, %v hit points left\n", strike.Striker, strike.DMG, strike.Strikee, strike.DEF, strike.HP)
	}

	return s
}

func (player Player) HasWeapon() (result bool){
	for _,item := range player.Items {
		if item.Type == WeaponType {
			result = true
		}
	}
	return
}

func (player Player) HasArmor() (result bool){
	for _,item := range player.Items {
		if item.Type == ArmorType {
			result = true
		}
	}
	return
}

func (player Player) HasMaxRings() (result bool){
	cnt := 0
	for _,item := range player.Items {
		if item.Type == RingType {
			cnt++
		}
	}
	result = cnt >= 2
	return
}


func main() {
	startTime := time.Now()
	boss := Player{ Name: "Boss" }
	if f, err := os.Open(input); err == nil {
		scanner := bufio.NewScanner(f)
		reg := regexp.MustCompile(`^([a-zA-Z ]+): ([0-9]+)$`)

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				va,_ := strconv.Atoi(groups[2])

				if groups[1] == "Hit Points" {
					boss.HP = va
				} else if groups[1] == "Damage" {
					boss.Damage = va
				} else if groups[1] == "Armor" {
					boss.Defense = va
				}
			}
		}
	}
	shop := []ShopItem{}
	if f, err := os.Open(shopfile); err == nil {
		scanner := bufio.NewScanner(f)
		reg := regexp.MustCompile(`^([a-zA-Z0-9\+]+),([0-9]+),([0-9]+),([0-9]+),([0-9]+)$`)

		for scanner.Scan() {
			var txt = scanner.Text()
			if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				t,_ := strconv.Atoi(groups[5])
				gld,_ := strconv.Atoi(groups[2])
				dmg,_ := strconv.Atoi(groups[3])
				def,_ := strconv.Atoi(groups[4])
				shop = append(shop, ShopItem{ Name: groups[1], Type: t, Gold: gld, Damage: dmg, Defense: def})
			}
		}
	}

	player := Player{ Name: "Jason", Damage: 0, Defense: 0, HP: 100 }
	variationList := Variations(player, boss, shop)

	wins := VariationList{}
	losses := VariationList{}
	for _,v := range variationList {
		Simulate(&v)
		if v.BossDefeated {
			wins = append(wins, v)
		} else {
			losses = append(losses, v)
		}
	}

	sorter := VariationListSorter{ Entries: wins }
	sort.Sort(sort.Reverse(sorter))

	losssorter := VariationListSorter{ Entries: losses }
	sort.Sort(losssorter)


	fmt.Println(sorter.Entries[len(sorter.Entries)-1])
	fmt.Println(losssorter.Entries[len(losssorter.Entries)-1])

	fmt.Println("Time", time.Since(startTime))
}

func Variations(player, boss Player, shop []ShopItem) VariationList {
	variations := VariationList{}

	for i := 0; i < len(shop); i++ {
		p,b := Copy(player), Copy(boss)
		if shop[i].Type == WeaponType && !p.HasWeapon() {
			if Buy(&p, shop[i]) && Valid(p) {
				v := Variation{ Player: p, Boss: b }
				variations = append(variations, v)
				subvars := Variations(Copy(p), Copy(b), shop)
				variations = append(variations, subvars...)
			}
		}

		if shop[i].Type == ArmorType && !p.HasArmor() {
			if Buy(&p, shop[i]) && Valid(p) {
				v := Variation{ Player: p, Boss: b }
				variations = append(variations, v)
				subvars := Variations(Copy(p), Copy(b), shop)
				variations = append(variations, subvars...)
			}
		}

		if shop[i].Type == RingType && !p.HasMaxRings() {
			if Buy(&p, shop[i]) && Valid(p) {
				v := Variation{ Player: p, Boss: b }
				variations = append(variations, v)
				subvars := Variations(Copy(p), Copy(b), shop)
				variations = append(variations, subvars...)
			}
		}
	}

	return variations
}

func Valid(p Player) bool {
	valid := false
	for _,item := range p.Items {
		if item.Type == WeaponType {
			valid = true
			break
		}
	}
	return valid
}

func Simulate(variation *Variation){
	alternate := false
	for variation.Boss.HP > 0 && variation.Player.HP > 0 {
		dmg, def, hp := 0,0,0
		striker,strikee := "",""
		if alternate {
			striker = variation.Boss.Name
			strikee = variation.Player.Name
			dmg, def, hp = DoStrike(&variation.Boss, &variation.Player)
		} else {
			striker = variation.Player.Name
			strikee = variation.Boss.Name
			dmg, def, hp = DoStrike(&variation.Player, &variation.Boss)
		}
		variation.Strikes = append(variation.Strikes, Strike{ Striker: striker, Strikee: strikee, DMG: dmg, DEF: def, HP: hp })

		alternate = !alternate
	}

	variation.BossDefeated = variation.Boss.HP <= 0
	variation.GoldSpent = variation.Player.GoldSpent
}

func Buy(player *Player, item ShopItem) (success bool) {
	exists := false
	for _,i := range player.Items {
		if i.Name == item.Name {
			// already owned
			exists = true
			break
		}
	}

	if !exists {
		player.Items = append(player.Items, item)
		player.Damage += item.Damage
		player.Defense += item.Defense
		player.GoldSpent += item.Gold
	}
	return !exists
}

func Copy(player Player) Player {
	p := Player{ Name: player.Name, HP: player.HP, Defense: player.Defense, Damage: player.Damage, GoldSpent: player.GoldSpent, Items: []ShopItem{} }
	for _,item := range player.Items {
		p.Items = append(p.Items, item)
	}
	return p
}

func DoStrike(striker, strikee *Player) (dmg, def, hp int) {
	dmg = striker.Damage - strikee.Defense
	if dmg < 1 { dmg = 1 }
	strikee.HP -= dmg
	return dmg, strikee.Defense, strikee.HP
}

// reg := regexp.MustCompile("-?[0-9]+")
/* 			
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
			*/
