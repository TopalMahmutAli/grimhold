package open5e

type Monster struct {
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	HitPoints       int    `json:"hit_points"`
	ArmorClass      int    `json:"armor_class"`
	ChallengeRating string `json:"challenge_rating"`
}

type MonsterPage struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Monster `json:"results"`
}
