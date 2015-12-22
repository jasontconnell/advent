package main

import (
	"fmt"
	"time"
	//"regexp"
	//"strconv"
	//"strings"
)

type Player struct {
	Name string
	HP int
	Damage int
	Mana int
	Spells []Spell
}

type Spell struct {
	Name string
	Cost int
	Effects []Effect
}

type EffectProperty int
const (
	Damage EffectProperty = iota
	HP EffectProperty = iota
	Defense EffectProperty = iota
	Mana EffectProperty = iota
)

type Effect struct {
	Length int
	Prop EffectProperty
	Value int
}

func main() {
	startTime := time.Now()
	boss := Player{ Name: "Boss", HP: 58, Damage: 9 }

	spells := []Spell{}
	mm := Spell{ Name: "Magic Missile", Cost: 53, Effects: []Effect{ Effect{ Length: 0, Prop: HP, Value: 4 } } }
	drain := Spell{ Name: "Drain", Cost: 73, Effects: []Effect{ Effect{ Length: 0, Prop: HP, Value: 2 }, Effect{ Length: 0, Prop: Damage, Value: 2 } } }
	shield := Spell{ Name: "Shield", Cost: 113, Effects: []Effect{ Effect{ Length: 6, Prop: Defense, Value: 7 } } }
	poison := Spell{ Name: "Poison", Cost: 173, Effects: []Effect{ Effect{ Length: 6, Prop: Damage, Value: 3 } } }
	recharge := Spell{ Name: "Recharge", Cost: 229, Effects: []Effect{ Effect{ Length: 5, Prop: Mana, Value: 101 } } }

	spells = append(spells, mm)
	spells = append(spells, drain)
	spells = append(spells, shield)
	spells = append(spells, poison)
	spells = append(spells, recharge)

	player := Player{ Name: "Jason", HP: 50, Mana: 500, Spells: spells }

	fmt.Println(boss, player)

	fmt.Println("Time", time.Since(startTime))
}