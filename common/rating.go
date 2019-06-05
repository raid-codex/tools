package common

type Rating struct {
	Overall       string `json:"overall"`
	Campaign      string `json:"campaign"`
	ArenaOff      string `json:"arena_offense"`
	ArenaDef      string `json:"arena_defense"`
	ClanBossWoGS  string `json:"clan_boss_without_giant_slayer"`
	ClanBosswGS   string `json:"clan_boss_with_giant_slayer"`
	IceGuardian   string `json:"ice_guardian"`
	Dragon        string `json:"dragon"`
	Spider        string `json:"spider"`
	FireKnight    string `json:"fire_knight"`
	Minotaur      string `json:"minotaur"`
	ForceDungeon  string `json:"force_dungeon"`
	MagicDungeon  string `json:"magic_dungeon"`
	SpiritDungeon string `json:"spirit_dungeon"`
	VoidDungeon   string `json:"void_dungeon"`
}

type Review struct {
	NumberOfReviews int64   `json:"amount"`
	Campaign        float64 `json:"campaign"`
	ArenaOff        float64 `json:"arena_offense"`
	ArenaDef        float64 `json:"arena_defense"`
	ClanBoss        float64 `json:"clan_boss"`
	IceGuardian     float64 `json:"ice_guardian"`
	Dragon          float64 `json:"dragon"`
	Spider          float64 `json:"spider"`
	FireKnight      float64 `json:"fire_knight"`
	Minotaur        float64 `json:"minotaur"`
	ForceDungeon    float64 `json:"force_dungeon"`
	MagicDungeon    float64 `json:"magic_dungeon"`
	SpiritDungeon   float64 `json:"spirit_dungeon"`
	VoidDungeon     float64 `json:"void_dungeon"`
}
