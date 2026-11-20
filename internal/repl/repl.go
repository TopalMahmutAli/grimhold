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

			page, err := open5e.FetchMonsters()
			if err != nil {
				return fmt.Errorf("load bestiary: %w", err)
			}

			for _, monster := range page.Results {
				fmt.Println(monster.Name)
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
