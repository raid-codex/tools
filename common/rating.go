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
