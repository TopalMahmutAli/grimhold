package open5e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const monstersURL = "https://api.open5e.com/v1/monsters/"

func FetchMonsters() (MonsterPage, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(monstersURL)
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
