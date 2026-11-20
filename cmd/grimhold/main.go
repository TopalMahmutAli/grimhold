package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	fmt.Println("Welcome to Grimhold!")

	commands := make(map[string]cliCommand)

	commands["help"] = cliCommand{
		name:        "help",
		description: "Show available commands",
		callback: func() error {
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
		callback: func() error {
			fmt.Println("Goodbye!")
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

		cmd, exists := commands[command]

		if !exists {
			fmt.Println("Unknown command:", command)
			continue
		}

		err := cmd.callback()
		if err != nil {
			fmt.Println("Error:", err)
		}

		if command == "exit" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Input error:", err)
	}
}
