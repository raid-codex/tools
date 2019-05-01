package common

type Characteristics struct {
	HP             int64   `json:"hp"`
	Attack         int64   `json:"attack"`
	Defense        int64   `json:"defense"`
	Speed          int64   `json:"speed"`
	CriticalRate   float64 `json:"critical_rate"`
	CriticalDamage float64 `json:"critical_damage"`
	Resistance     int64   `json:"resistance"`
	Accuracy       int64   `json:"accuracy"`
}
