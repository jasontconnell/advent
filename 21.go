package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"regexp"
	"strconv"
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
	result = cnt <= 2
	return
}

func (player Player) String() string {
	return fmt.Sprintf("%v HP: %v DMG: %v DEF:%v GLDSPNT:%v", player.Name, player.HP, player.Damage, player.Defense, player.GoldSpent)
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
		reg := regexp.MustCompile(`^([a-zA-Z\+]+),([0-9]+),([0-9]+),([0-9]+),([0-9]+)$`)

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

	Simulate(player, boss, shop, false)

	fmt.Println("Time", time.Since(startTime))
}

func Simulate(player,boss Player, shop []ShopItem, alternate bool){
	if player.HP < 0 || boss.HP < 0 { return }

	if !player.HasWeapon(){
		for i := 0; i < len(shop); i++ {
			if shop[i].Type == WeaponType {
				player = CopyPlayer(player)
				Buy(&player, shop[i])
				Simulate(player, boss, shop, alternate)
			}
		}
	}

	// armor is optional
	for i := -1; i < len(shop); i++ {
		if i != -1 && shop[i].Type == ArmorType && !player.HasArmor() {
			player = CopyPlayer(player)
			Buy(&player, shop[i])
		}
	}

	if boss.HP > 0 && player.HP > 0 {
		fmt.Println(boss, player)
		if alternate {
			boss,player = Strike(boss, player)
		} else {
			player,boss = Strike(player, boss)
		}
		Simulate(player, boss, shop, !alternate)
	} else {
		fmt.Println(player, boss)
	}
}

func Buy(player *Player, item ShopItem) {
	player.Items = append(player.Items, item)
	player.Damage += item.Damage
	player.Defense += item.Defense
	player.GoldSpent += item.Gold
}

func CopyPlayer(player Player) Player {
	p := Player{ Name: player.Name, HP: player.HP }
	return p
}

func Strike(striker, strikee Player) (Player, Player) {
	dmg := striker.Damage - strikee.Defense
	if dmg < 1 { dmg = 1 }

	strikee.HP -= dmg
	return striker, strikee
}

// reg := regexp.MustCompile("-?[0-9]+")
/* 			
if groups := reg.FindStringSubmatch(txt); groups != nil && len(groups) > 1 {
				fmt.Println(groups[1:])
			}
			*/
