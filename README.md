
# Grimhold

Grimhold is a text-based dungeon crawler written in Go.
Explore regions, fight monsters, gain experience and gold,
and save your progress.

Monster data is fetched from the Open5e API.

## Requirements

- Go 1.22 or higher
- Internet connection

## Build and run

Build the game:

```bash
go build -o grimhold ./cmd/grimhold
```

Run the game:

```bash
./grimhold
```

Or run directly:

```bash
go run ./cmd/grimhold
```

## Commands

| Command | Description |
|---|---|
| help | Show available commands |
| map | List all regions |
| look | Describe the current region |
| travel <region> | Travel to a connected region |
| explore | Encounter a random monster |
| fight | Attack the enemy |
| flee | Try to escape |
| character | Show character statistics |
| rest | Restore HP once per region visit |
| bestiary | List monsters |
| next / prev | Navigate bestiary pages |
| inspect <monster> | Show monster statistics |
| slain bestiary | Show defeated monsters |
| save | Save the game |
| load | Load the saved game |
| exit | Quit the game |

## Save and load

Use `save` to save your progress to `save.json`.
Use `load` to restore your previous game.

Saving and loading are not allowed during combat.

## Example session

```text
grimhold > travel ashen-marsh
You travel to Ashen Marsh
Fog, iron, and things that breathe under the water.

grimhold > explore
A wild Goblin appears!
HP: 7
Armor Class: 15

grimhold > flee
You try to flee: 14 + 2 = 16
You escaped from Goblin!

grimhold > save
Game saved successfully!

grimhold > exit
Goodbye!
```

## Project Architecture

The project is organized into packages:

```text
grimhold/
├── cmd/grimhold/ 
│   └── main.go
├── internal/
│   ├── repl/
│   ├── game/
│   ├── world/
│   └── open5e/
├── world.json
├── go.mod
└── README.md
```

- **cmd/grimhold**: Entry point of the application.
- **internal/repl**: Handles commands and user input.
- **internal/game**: Manages the character, combat, XP and saves.
- **internal/world**: Loads the world and manages regions.
- **internal/open5e**: Fetches monster data from the Open5e API.
- **world.json**: Contains regions, enemies and connections.

The project uses only the Go standard library.

## API

https://api.open5e.com/v1/monsters/
