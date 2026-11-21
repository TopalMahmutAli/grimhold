package repl

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strings"

	"github.com/TopalMahmutAli/grimhold/internal/game"
	"github.com/TopalMahmutAli/grimhold/internal/open5e"
	"github.com/TopalMahmutAli/grimhold/internal/world"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func monsterAttack(monster *open5e.Monster, player *game.Character) error {
	attackName := "Basic Attack"
	attackBonus := 2
	damageDice := "1d4"
	damageBonus := 0

	for _, action := range monster.Actions {
		if action.DamageDice != "" {
			attackName = action.Name
			attackBonus = action.AttackBonus
			damageDice = action.DamageDice
			damageBonus = action.DamageBonus
			break
		}
	}

	fmt.Printf("\n%s attacks with %s!\n", monster.Name, attackName)

	roll := game.RollD20()

	fmt.Printf("%s rolls: %d + %d = %d\n",
		monster.Name, roll, attackBonus, roll+attackBonus)

	if !game.AttackHits(roll, attackBonus, player.ArmorClass) {
		fmt.Printf("%s misses!\n", monster.Name)
		return nil
	}

	damage, err := game.RollDamageDice(damageDice, damageBonus)
	if err != nil {
		return fmt.Errorf("monster damage: %w", err)
	}

	player.HP -= damage

	if player.HP < 0 {
		player.HP = 0
	}

	fmt.Printf("%s hits you for %d damage!\n", monster.Name, damage)
	fmt.Printf("Your HP: %d/%d\n", player.HP, player.MaxHP)

	return nil
}

func Start() {
	fmt.Println("Welcome to Grimhold!")

	w, err := world.Load("world.json")

	if err != nil {
		fmt.Println("Error loading world:", err)
		return
	}

	fmt.Printf("Loaded %d regions\n", len(w.Regions))

	currentRegionID := w.Start
	player := game.NewCharacter()
	hasRested := false
	var currentEncounter *open5e.Monster

	respawnIfDead := func() {
		if player.HP > 0 {
			return
		}

		fmt.Println("\nYou have been defeated!")
		fmt.Println("You wake up at the Crossroads with 1 HP.")

		player.HP = 1
		currentRegionID = w.Start
		currentEncounter = nil
		hasRested = false
	}

	slain := make(map[string]int)

	var currentPage open5e.MonsterPage
	pageNumber := 0

	commands := make(map[string]cliCommand)

	commands["help"] = cliCommand{
		name:        "help",
		description: "Show available commands",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("help does not accept arguments")
			}

			fmt.Println("Available commands:")

			names := make([]string, 0, len(commands))
			for name := range commands {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				cmd := commands[name]
				fmt.Printf(" %s - %s\n", cmd.name, cmd.description)
			}
			return nil
		},
	}
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Quit the game",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("exit does not accept arguments")
			}

			fmt.Println("Goodbye!")
			return nil
		},
	}

	commands["travel"] = cliCommand{
		name:        "travel",
		description: "Travel to a connected region",
		callback: func(args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: travel <region>")
			}

			if currentEncounter != nil {
				return fmt.Errorf("you must defeat or flee from the enemy first")
			}

			destinationID := args[0]

			for _, region := range w.Regions {
				if region.ID != currentRegionID {
					continue
				}

				for _, exit := range region.Exits {
					if exit == destinationID {

						for _, destination := range w.Regions {
							if destination.ID == destinationID {
								currentRegionID = destination.ID
								hasRested = false

								fmt.Println("You travel to", destination.Name)
								fmt.Println(destination.Description)

								return nil
							}
						}
						return fmt.Errorf("region not found: %s", destinationID)
					}
				}
				return fmt.Errorf(
					"cannot travel to %s from %s",
					destinationID,
					region.Name,
				)
			}
			return fmt.Errorf("current region not found: %s", currentRegionID)
		},
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "List all known regions",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("map does not accept arguments")
			}
			fmt.Println("World Map:")

			for _, region := range w.Regions {
				fmt.Printf(
					"%s (%s) - Danger: %d\n",
					region.Name,
					region.ID,
					region.Danger,
				)
				fmt.Println(" Exits:", strings.Join(region.Exits, ", "))
			}
			return nil
		},
	}

	commands["look"] = cliCommand{
		name:        "look",
		description: "Describe your surroundings",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("look does not accept arguments")
			}

			for _, region := range w.Regions {
				if region.ID == currentRegionID {
					fmt.Println("Region:", region.Name)
					fmt.Println("Description:", region.Description)
					fmt.Println("Danger:", region.Danger)
					fmt.Println("Exits:", strings.Join(region.Exits, ", "))
					return nil

				}
			}
			return fmt.Errorf("current region not found: %s", currentRegionID)
		},
	}

	commands["bestiary"] = cliCommand{
		name:        "bestiary",
		description: "List monsters from Open5e",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("bestiary does not accept arguments")
			}

			page, err := open5e.FetchMonsters(open5e.MonstersURL)
			if err != nil {
				return fmt.Errorf("load bestiary: %w", err)
			}

			currentPage = page
			pageNumber = 1

			fmt.Println("Monsters - Page 1:")

			for _, monster := range currentPage.Results {
				fmt.Println(monster.Name)
			}

			return nil
		},
	}

	commands["next"] = cliCommand{
		name:        "next",
		description: "Show the next page of monsters",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("next does not accept arguments")
			}

			if pageNumber == 0 {
				return fmt.Errorf("use bestiary first")
			}

			if currentPage.Next == nil {
				return fmt.Errorf("no more pages")
			}

			page, err := open5e.FetchMonsters(*currentPage.Next)
			if err != nil {
				return fmt.Errorf("load next page: %w", err)
			}

			currentPage = page
			pageNumber++

			fmt.Printf("Monsters - Page %d:\n", pageNumber)

			for _, monster := range currentPage.Results {
				fmt.Println(monster.Name)
			}

			return nil
		},
	}

	commands["prev"] = cliCommand{
		name:        "prev",
		description: "Show the previous page of monsters",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("prev does not accept arguments")
			}

			if pageNumber == 0 {
				return fmt.Errorf("use bestiary first")
			}

			if currentPage.Previous == nil {
				return fmt.Errorf("no previous page")
			}

			page, err := open5e.FetchMonsters(*currentPage.Previous)
			if err != nil {
				return fmt.Errorf("load previous page: %w", err)
			}

			currentPage = page
			pageNumber--

			fmt.Printf("Monsters - Page %d:\n", pageNumber)

			for _, monster := range currentPage.Results {
				fmt.Println(monster.Name)
			}

			return nil
		},
	}

	commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Show a monster's statistics",
		callback: func(args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: inspect <monster>")
			}

			monster, err := open5e.FetchMonster(args[0])
			if err != nil {
				return fmt.Errorf("inspect monster: %w", err)
			}

			fmt.Println("Name:", monster.Name)
			fmt.Println("Hit Points:", monster.HitPoints)
			fmt.Println("Armor Class:", monster.ArmorClass)
			fmt.Println("Challenge Rating:", monster.ChallengeRating)

			fmt.Println("Actions:")

			for _, action := range monster.Actions {
				fmt.Println(" -", action.Name)

				if action.DamageDice != "" {
					fmt.Printf(
						"   Attack: %+d | Damage: %s%+d\n",
						action.AttackBonus,
						action.DamageDice,
						action.DamageBonus,
					)
				}

				if action.Description != "" {
					fmt.Println("  ", action.Description)
				}
			}

			return nil
		},
	}

	commands["character"] = cliCommand{
		name:        "character",
		description: "Show your character sheet",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("character does not accept arguments")
			}

			fmt.Println("Character Sheet")
			fmt.Printf("HP: %d/%d\n", player.HP, player.MaxHP)
			fmt.Printf("Level: %d\n", player.Level)
			fmt.Printf("XP: %d\n", player.XP)
			fmt.Printf("Attack Bonus: %+d\n", player.AttackBonus)
			fmt.Printf("Armor Class: %d\n", player.ArmorClass)
			fmt.Printf("Weapon: %s\n", player.Weapon)
			fmt.Printf("Gold: %d\n", player.Gold)

			return nil
		},
	}

	commands["rest"] = cliCommand{
		name:        "rest",
		description: "Restore your health once per region visit",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("rest does not accept arguments")
			}
			if currentEncounter != nil {
				return fmt.Errorf("you cannot rest while facing an enemy")
			}

			if hasRested {
				return fmt.Errorf("you have already rested in this region")
			}

			player.HP = player.MaxHP

			hasRested = true

			fmt.Println("You rest and recover your HP.")
			fmt.Printf("HP: %d/%d\n", player.HP, player.MaxHP)

			return nil
		},
	}

	commands["explore"] = cliCommand{
		name:        "explore",
		description: "Explore the current region",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("explore does not accept arguments")
			}

			if currentEncounter != nil {
				return fmt.Errorf("you are already facing an enemy")
			}

			for _, region := range w.Regions {
				if region.ID != currentRegionID {
					continue
				}

				if len(region.Foes) == 0 {
					fmt.Println("You explore the area, but find no enemies.")
					return nil
				}

				index := rand.IntN(len(region.Foes))
				monsterSlug := region.Foes[index]

				monster, err := open5e.FetchMonster(monsterSlug)
				if err != nil {
					return fmt.Errorf(
						"fetch encounter monster %s: %w",
						monsterSlug,
						err,
					)
				}

				// Mémoriser le monstre rencontré
				currentEncounter = &monster

				fmt.Printf("A wild %s appears!\n", currentEncounter.Name)
				fmt.Printf("HP: %d\n", currentEncounter.HitPoints)
				fmt.Printf("Armor Class: %d\n", currentEncounter.ArmorClass)

				return nil
			}

			return fmt.Errorf(
				"current region not found: %s",
				currentRegionID,
			)
		},
	}

	commands["fight"] = cliCommand{
		name:        "fight",
		description: "Attack the current enemy",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("fight does not accept arguments")
			}

			if currentEncounter == nil {
				return fmt.Errorf("there is no enemy to fight")
			}

			roll := game.RollD20()
			total := roll + player.AttackBonus

			fmt.Printf("You roll a d20: %d + %d = %d\n",
				roll, player.AttackBonus, total)

			if game.AttackHits(roll, player.AttackBonus, currentEncounter.ArmorClass) {
				damage, err := game.RollDamage(1, 6, 0)
				if err != nil {
					return fmt.Errorf("player damage: %w", err)
				}

				currentEncounter.HitPoints -= damage

				fmt.Printf("You hit %s for %d damage!\n",
					currentEncounter.Name, damage)

				if currentEncounter.HitPoints <= 0 {
					fmt.Printf("You defeated %s!\n", currentEncounter.Name)

					xp := game.XPForChallengeRating(currentEncounter.ChallengeRating)

					gold := rand.IntN(6) + 1

					levelsGained := player.GainXP(xp)
					player.Gold += gold

					slain[currentEncounter.Name]++

					fmt.Printf("You gained %d XP and %d gold!\n", xp, gold)

					if levelsGained > 0 {
						fmt.Printf("LEVEL UP! You reached level %d!\n", player.Level)
						fmt.Printf("Max HP: %d\n", player.MaxHP)
						fmt.Printf("Attack Bonus: %+d\n", player.AttackBonus)
					}
					currentEncounter = nil
					return nil
				}

				fmt.Printf("%s has %d HP remaining.\n",
					currentEncounter.Name, currentEncounter.HitPoints)
			} else {
				fmt.Println("Your attack misses!")
			}
			if err := monsterAttack(currentEncounter, &player); err != nil {
				return fmt.Errorf("enemy counter attack: %w", err)
			}

			respawnIfDead()

			return nil
		},
	}

	commands["flee"] = cliCommand{
		name:        "flee",
		description: "Try to escape from an enemy",
		callback: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("flee does not accept arguments")
			}

			if currentEncounter == nil {
				return fmt.Errorf("there is no enemy to flee from")
			}

			roll := game.RollD20()
			bonus := 2
			total := roll + bonus

			fmt.Printf("You try to flee: %d + %d = %d\n",
				roll, bonus, total)

			if total >= 12 {
				fmt.Printf("You escaped from %s!\n", currentEncounter.Name)
				currentEncounter = nil
				return nil
			}

			fmt.Println("Escape failed! The enemy gets a free attack.")

			if err := monsterAttack(currentEncounter, &player); err != nil {
				return fmt.Errorf("enemy free attack: %w", err)
			}

			respawnIfDead()

			return nil
		},
	}

	commands["slain"] = cliCommand{
		name:        "slain bestiary",
		description: "Show monster you have defeated",
		callback: func(args []string) error {
			if len(args) != 1 || args[0] != "bestiary" {
				return fmt.Errorf("usage: slain bestiary")
			}
			if len(slain) == 0 {
				fmt.Println("You have not defeated any monsters yet.")
				return nil
			}
			fmt.Println("Slain Bestiary:")

			names := make([]string, 0, len(slain))
			for name := range slain {
				names = append(names, name)
			}
			sort.Strings(names)

			for _, name := range names {
				fmt.Printf("%s: %d kills\n", name, slain[name])
			}
			return nil
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("grimhold > ")

		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())

		if command == "" {
			continue
		}

		parts := strings.Fields(command)

		commandName := parts[0]
		args := parts[1:]

		cmd, exists := commands[commandName]

		if !exists {
			fmt.Println("Unknown command:", commandName)
			continue
		}

		err := cmd.callback(args)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if commandName == "exit" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Input error:", err)
	}
}
