package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SaveState struct {
	Location  string         `json:"location"`
	Character Character      `json:"character"`
	Slain     map[string]int `json:"slain"`
	HasRested bool           `json:"has_rested"`
}

func Save(path string, state SaveState) error {
	data, err := json.MarshalIndent(state, "", " ")
	if err != nil {
		return fmt.Errorf("encode save: %w", err)
	}

	file, err := os.CreateTemp(filepath.Dir(path), "grimhold-*.temp")
	if err != nil {
		return fmt.Errorf("create temporary save : %w", err)
	}

	defer os.Remove(file.Name())
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write save: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("syn save: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close save: %w", err)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("replace save: %w", err)
	}
	return nil
}

func Load(path string) (SaveState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SaveState{}, fmt.Errorf("decode save: %w", err)
	}

	var state SaveState

	if err := json.Unmarshal(data, &state); err != nil {
		return SaveState{}, fmt.Errorf("decode save: %w", err)
	}

	if state.Character.Level < 1 ||
		state.Character.MaxHP < 1 ||
		state.Character.HP < 1 ||
		state.Character.HP > state.Character.MaxHP ||
		state.Character.XP < 0 ||
		state.Character.Gold < 0 {
		return SaveState{}, fmt.Errorf("invalid character in save")
	}

	if state.Slain == nil {
		state.Slain = make(map[string]int)
	}
	return state, nil
}
