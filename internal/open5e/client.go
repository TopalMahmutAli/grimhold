package open5e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const MonstersURL = "https://api.open5e.com/v1/monsters/"

func FetchMonsters(url string) (MonsterPage, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return MonsterPage{}, fmt.Errorf("fetch monsters: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MonsterPage{}, fmt.Errorf(
			"Open5e returned HTTP %d",
			resp.StatusCode,
		)
	}

	var page MonsterPage

	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return MonsterPage{}, fmt.Errorf("decode monsters: %w", err)
	}
	return page, nil
}

func FetchMonster(slug string) (Monster, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	url := MonstersURL + slug + "/"

	resp, err := client.Get(url)
	if err != nil {
		return Monster{}, fmt.Errorf("fetch monster: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Monster{}, fmt.Errorf(
			"monster not found or API error: HTTP %d",
			resp.StatusCode,
		)
	}

	var monster Monster

	if err := json.NewDecoder(resp.Body).Decode(&monster); err != nil {
		return Monster{}, fmt.Errorf("decode monster: %w", err)
	}

	return monster, nil
}
