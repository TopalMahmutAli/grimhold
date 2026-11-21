package game

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/TopalMahmutAli/grimhold/internal/open5e"
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

func MonsterAttack(monster *open5e.Monster, player *Character) error {
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

	fmt.Printf("\n%s attacks with %s!\n",
		monster.Name, attackName)

	roll := RollD20()

	fmt.Printf("%s rolls: %d + %d = %d\n",
		monster.Name, roll, attackBonus, roll+attackBonus)

	if !AttackHits(roll, attackBonus, player.ArmorClass) {
		fmt.Printf("%s misses!\n", monster.Name)
		return nil
	}

	damage, err := RollDamageDice(damageDice, damageBonus)
	if err != nil {
		return fmt.Errorf("monster damage: %w", err)
	}

	player.HP = max(0, player.HP-max(0, damage))

	fmt.Printf("%s hits you for %d damage!\n",
		monster.Name, max(0, damage))

	fmt.Printf("Your HP: %d/%d\n", player.HP, player.MaxHP)

	return nil
}
