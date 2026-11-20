package repl

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/TopalMahmutAli/grimhold/internal/open5e"
	"github.com/TopalMahmutAli/grimhold/internal/world"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
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

			// Récupérer le monstre depuis Open5e
			monster, err := open5e.FetchMonster(args[0])
			if err != nil {
				return fmt.Errorf("inspect monster: %w", err)
			}

			// Afficher les caractéristiques
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
