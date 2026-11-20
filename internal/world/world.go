package world

import (
	"encoding/json"
	"fmt"
	"os"
)

type Region struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Danger      int      `json:"danger"`
	Foes        []string `json:"foes"`
	Exits       []string `json:"exits"`
}

type World struct {
	Start   string   `json:"start"`
	Regions []Region `json:"regions"`
}

func Load(path string) (World, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return World{}, fmt.Errorf("cannot read world: %w", err)
	}

	var w World

	if err := json.Unmarshal(data, &w); err != nil {
		return World{}, fmt.Errorf("cannot decode world: %w", err)
	}
	return w, nil
}
