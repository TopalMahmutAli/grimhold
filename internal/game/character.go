package game

type Character struct {
	HP          int
	MaxHP       int
	Level       int
	XP          int
	AttackBonus int
	ArmorClass  int
	Weapon      string
	Gold        int
}

func NewCharacter() Character {
	return Character{
		HP:          20,
		MaxHP:       20,
		Level:       1,
		XP:          0,
		AttackBonus: 3,
		ArmorClass:  13,
		Weapon:      "Rusty Sword",
		Gold:        0,
	}
}
