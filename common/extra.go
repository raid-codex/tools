package common

import (
	"encoding/json"
	"fmt"
	"os"
)

func fetchStatusEffects(dataDirectory string) (map[string]*StatusEffect, error) {
	var sl StatusEffectList

	file, errOpen := os.Open(fmt.Sprintf("%s/docs/status-effects/current/index.json", dataDirectory))
	if errOpen != nil {
		return nil, errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&sl)
	if errJSON != nil {
		return nil, errJSON
	}

	effects := map[string]*StatusEffect{}
	for _, effect := range sl {
		effects[effect.Slug] = effect
	}

	return effects, nil
}

func fetchChampions(dataDirectory string) (ChampionList, error) {
	var cl ChampionList

	file, errOpen := os.Open(fmt.Sprintf("%s/docs/champions/current/index.json", dataDirectory))
	if errOpen != nil {
		return nil, errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&cl)
	if errJSON != nil {
		return nil, errJSON
	}
	return cl, nil
}
