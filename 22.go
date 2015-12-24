package main

import (
	"fmt"
	"time"
	"math/rand"
)

type Player struct {
	Name string
	HP int
	Damage int
	Defense int
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

type BattleEffect struct {
	Countdown int
	Prop EffectProperty
	Value int
}
func (eff BattleEffect) String() string {
	prop := ""
	if eff.Prop == Damage { 
		prop = "Damage" 
	} else if eff.Prop == HP { 
		prop = "HP" 
	} else if eff.Prop == Defense { 
		prop = "Defense" 
	} else if eff.Prop == Mana { 
		prop = "Mana" 
	}

	return fmt.Sprintf("%v Countdown: %v Value: %v", prop, eff.Countdown, eff.Value)
}

func main() {
	startTime := time.Now()

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

	minmana := 1000000
	var winningCasts []Spell

	for i := 0; i < 6000000; i++ {
		//fmt.Println("running sim #", i)
		player := Player{ Name: "Jason", HP: 50, Mana: 500, Spells: spells }
		boss := Player{ Name: "Boss", HP: 58, Damage: 9 }

		casts := RunSim(&player, &boss, spells)
		mana := 0

		for _, c := range casts {
			mana += c.Cost
		}

		if (player.HP > 0 && boss.HP <= 0) && mana < minmana {
			minmana = mana
			winningCasts = casts
			fmt.Println("Got new winning combo", player, boss)
		}
	}

	fmt.Println("min mana =", minmana)
	fmt.Println("winning combo", winningCasts)

	fmt.Println("Time", time.Since(startTime))
}

func RunSim(player, boss *Player, availableSpells []Spell) []Spell {
	effects := []BattleEffect{}
	casts := []Spell{}
	alternate := false
	for {
		// process pre-turn effects
		ApplyEffects(player, boss, effects)
		// attack
		if alternate {
			Attack(boss, player)
		} else {
			spell := GetRandSpell(player.Mana, availableSpells)
			if spell.Name == "" { break }

			effs := Cast(player, spell)
			casts = append(casts, spell)

			for i := len(effs)-1; i >= 0; i--  {
				if effs[i].Countdown == 0 {
					ApplyEffects(player, boss, []BattleEffect{effs[i]})
					effs = append(effs[:i], effs[i+1:]...)
				} else {
					effects = append(effects, effs[i])
				}
			}
		}


		// determine if one is dead
		if player.HP <= 0 || boss.HP <= 0 {
			// dead
			break
		}

		player.Defense = 0

		// count down effects counters
		for i := len(effects)-1; i >= 0; i-- {
			effects[i].Countdown--

			if effects[i].Countdown <= 0 {
				effects = append(effects[:i], effects[i+1:]...)  // get rid of current effect
			}
		}
		alternate = !alternate
	}

	return casts
}

func ApplyEffects(player, boss *Player, effects []BattleEffect){
	for _,eff := range effects {
		switch eff.Prop {
			case Damage: boss.HP -= eff.Value
			break
			case HP: player.HP += eff.Value
			break
			case Defense: player.Defense += eff.Value
			break
			case Mana: player.Mana += eff.Value
			break
		}
	}
}

func GetRandSpell(mana int, spells []Spell) Spell {
	cancast := []Spell{}
	for _,sp := range spells {
		if sp.Cost <= mana {
			cancast = append(cancast, sp)
		}
	}
	if len(cancast) > 0 {
		n := rand.Intn(len(cancast))
		return cancast[n]
	} else {
		return Spell{}
	}
}

func Attack(striker, strikee *Player) (dmg, def, hp int) {
	dmg = striker.Damage - strikee.Defense
	if dmg < 1 { dmg = 1 }
	strikee.HP -= dmg
	return dmg, strikee.Defense, strikee.HP
}

func Cast(caster *Player, spell Spell) []BattleEffect {
	caster.Mana -= spell.Cost
	effects := []BattleEffect{}
	for _,eff := range spell.Effects {
		effects = append(effects, BattleEffect{ Countdown: eff.Length, Prop: eff.Prop, Value: eff.Value })
	}
	return effects
}