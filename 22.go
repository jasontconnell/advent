package main

import (
	"fmt"
	"time"
	"math/rand"
	"runtime"
	"sync"
)

type Player struct {
	Name string
	HP int
	Damage int
	Defense int
	Mana int
	HardMode bool
}

func (p Player) String() string {
	return fmt.Sprintf("%v HP: %v Defense: %v Mana: %v Hard Mode: %v", p.Name, p.HP, p.Defense, p.Mana, p.HardMode)
}

type Spell struct {
	Name string
	Cost int
	Effects []Effect
}

func (s Spell) String() string {
	return fmt.Sprintf("%v Cost: %v  Effects: %v", s.Name, s.Cost, s.Effects)
}

type EffectProperty int
const (
	Damage EffectProperty = iota
	HP EffectProperty = iota
	Defense EffectProperty = iota
	Mana EffectProperty = iota
)

type Effect struct {
	Countdown int
	Prop EffectProperty
	Value int
	OneTime bool
	Instant bool
}

func (eff Effect) String() string {
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

	return fmt.Sprintf("%v Countdown: %v Value: %v One Time: %v Instant: %v", prop, eff.Countdown, eff.Value, eff.OneTime, eff.Instant)
}

type Result struct {
	Mana int
	Story Story
}

type Story []string

func (story *Story) AddLine(s string) {
	(*story) = append(*story, s)
}

func (story Story) String() string {
	s := ""
	for _, line := range story {
		s += fmt.Sprintf("%v\n", line)
	}
	return s
}

func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	startTime := time.Now()

	rand.Seed(startTime.Unix())

	spells := []Spell{}
	mm := Spell{ Name: "Magic Missile", Cost: 53, Effects: []Effect{ Effect{ Instant: true, Prop: Damage, Value: 4 } } }
	drain := Spell{ Name: "Drain", Cost: 73, Effects: []Effect{ Effect{ Instant: true, Prop: HP, Value: 2 }, Effect{ Instant: true, Prop: Damage, Value: 2 } } }
	shield := Spell{ Name: "Shield", Cost: 113, Effects: []Effect{ Effect{ Countdown: 6, Prop: Defense, Value: 7, OneTime: true } } }
	poison := Spell{ Name: "Poison", Cost: 173, Effects: []Effect{ Effect{ Countdown: 6, Prop: Damage, Value: 3 } } }
	recharge := Spell{ Name: "Recharge", Cost: 229, Effects: []Effect{ Effect{ Countdown: 5, Prop: Mana, Value: 101 } } }

	spells = append(spells, mm)
	spells = append(spells, drain)
	spells = append(spells, shield)
	spells = append(spells, poison)
	spells = append(spells, recharge)

	min := 1400
	hardMode := true

	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	results := make(chan Result)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(threadID int) {
			RunSims(threadID, spells, min, hardMode, results)
			wg.Done()
		}(i)
	}

	go func(){
		for r := range results {
			if r.Mana < min {
				min = r.Mana
				fmt.Println(r.Mana, "\n", r.Story)
			}
		}
		close(results)
	}()
	wg.Wait()

	fmt.Println("Time", time.Since(startTime))
}

func RunSims(threadID int, spells []Spell, min int, hardMode bool, res chan Result)  {
	minmana := min
	for i := 0; i < 1000; i++ {
		player := Player{ Name: "Jason", HP: 50, Mana: 500, HardMode: hardMode }
		boss := Player{ Name: "Boss", HP: 58, Damage: 9 }

		story,mana,win := RunSim(&player, &boss, spells, minmana)

		if win && mana < minmana {
			minmana = mana
			result := Result{ Mana: mana, Story: story }
			res <- result
		}
	}
}

func RunSim(player, boss *Player, availableSpells []Spell, minmana int) (Story, int, bool) {
	effects := []Effect{}
	story := Story{}
	alternate := false
	mana := 0
	turn := 0
	win := false
	for mana < minmana {
		story.AddLine(fmt.Sprintf("\nStarting turn %v\n----------", turn))

		if player.HardMode {
			player.HP--
			story.AddLine(fmt.Sprintf("Player hard mode. Player HP: %v", player.HP))
		}

		// process pre-turn effects
		ApplyEffects(&story, player, boss, effects)
		// count down effects counters
		effects = CountdownEffects(&story, player, boss, effects)

		if player.HP <= 0 || boss.HP <= 0 {
			// dead
			if player.HP <= 0 {
				story.AddLine(fmt.Sprintf("Player has died"))
			} else {
				win = true
				story.AddLine(fmt.Sprintf("Boss has died"))
			}
			break
		}

		// attack
		if alternate {
			damage := Attack(boss.Damage, player.Defense)
			player.HP -= damage
			story.AddLine(fmt.Sprintf("Boss attacks for %v, player HP: %v", damage, player.HP))
		} else {
			spell := GetRandSpell(*boss, *player, availableSpells, effects)
			if spell.Cost > 0 {
				effs := Cast(&story, player, spell)
				ApplyOneTime(&story, player, boss, effs)

				effects = append(effects, effs...)
				mana += spell.Cost
			} else {
				break // if can't cast any spell you lose
			}
		}

		alternate = !alternate
		turn++
	}

	return story, mana, win
}

func CountdownEffects(story *Story, player, boss *Player, effects []Effect) []Effect {
	remaining := []Effect{}
	for i := len(effects)-1; i >= 0; i-- {
		if !effects[i].Instant {
			effects[i].Countdown--
			story.AddLine(fmt.Sprintf("Countdown updated: %v", effects[i]))
		}

		if effects[i].Countdown == 0 {
			RollbackEffect(player, boss, effects[i])
			story.AddLine(fmt.Sprintf("Effect %v has ended", effects[i]))
		} else {
			remaining = append(remaining, effects[i])
		}
	}
	return remaining
}

func RollbackEffect(player, boss *Player, eff Effect){
	if eff.OneTime {
		switch eff.Prop {
			case Defense: 
				player.Defense -= eff.Value
			break
		}
	}
}

func ApplyOneTime(story *Story, player, boss *Player, effects []Effect){
	for _,eff := range effects {
		if eff.OneTime {
			switch eff.Prop {
				case Defense: 
					player.Defense += eff.Value
					story.AddLine(fmt.Sprintf("Effect Defense has started. Player Defense: %v", player.Defense))
				break
			}
		}
	}
}

func ApplyEffects(story *Story, player, boss *Player, effects []Effect){
	for _,eff := range effects {
		switch eff.Prop {
			case Damage: 
				boss.HP -= eff.Value
				story.AddLine(fmt.Sprintf("Applied effect Damage, Boss HP: %v", boss.HP))
			break
			case HP: 
				player.HP += eff.Value
				story.AddLine(fmt.Sprintf("Applied effect HP, Player HP: %v", player.HP))
			break
			case Mana: 
				player.Mana += eff.Value
				story.AddLine(fmt.Sprintf("Applied effect Mana, Player Mana: %v", player.Mana))
			break
		}
	}
}

func GetRandSpell(boss, player Player, spells []Spell, effects []Effect) Spell {
	available := []Spell{}
	for _,sp := range spells {
		canCast := true
		if sp.Cost <= player.Mana {
			for _,eff := range effects {
				for _, spellEffect := range sp.Effects {
					if eff.Prop == spellEffect.Prop && eff.Instant == spellEffect.Instant {
						canCast = false
						break
					}
				}
			}
		} else {
			canCast = false
		}

		if canCast {
			available = append(available, sp)
		}
	}

	if len(available) > 0 {
		index := rand.Intn(len(available))
		return available[index]
	} else {
		return Spell{}
	}
}

func Attack(damage, defense int) (dmg int) {
	dmg = damage - defense
	if dmg < 1 { dmg = 1 }
	return dmg
}

func Cast(story *Story, caster *Player, spell Spell) []Effect {
	caster.Mana -= spell.Cost
	story.AddLine(fmt.Sprintf("Player casting %v for %v mana, %v mana remaining", spell.Name, spell.Cost, caster.Mana))
	effects := []Effect{}
	for _,eff := range spell.Effects {
		effects = append(effects, Effect{ Countdown: eff.Countdown, Prop: eff.Prop, Value: eff.Value, OneTime: eff.OneTime, Instant: eff.Instant })
	}
	return effects
}