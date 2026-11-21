package game

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

func RollD20() int {
	return rand.IntN(20) + 1
}

func AttackHits(roll, bonus, armorClass int) bool {
	return roll+bonus >= armorClass
}

func RollDamage(count, sides, bonus int) (int, error) {
	if count <= 0 || sides <= 0 {
		return 0, fmt.Errorf("invalid dice: %dd%d", count, sides)
	}

	total := bonus

	for i := 0; i < count; i++ {
		total += rand.IntN(sides) + 1
	}

	return total, nil
}

func RollDamageDice(dice string, bonus int) (int, error) {
	countText, sidesText, ok := strings.Cut(strings.TrimSpace(dice), "d")
	if !ok {
		return 0, fmt.Errorf("invalid damage dice: %s", dice)
	}
	count, err := strconv.Atoi(countText)
	if err != nil {
		return 0, fmt.Errorf("invalid dice count: %w", err)
	}
	sides, err := strconv.Atoi(sidesText)
	if err != nil {
		return 0, fmt.Errorf("invalid dice sides: %w", err)
	}

	return RollDamage(count, sides, bonus)
}
