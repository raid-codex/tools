package common

import (
	"encoding/json"
	"fmt"
	"os"
)

func GetPageExtraData(dataDirectory string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	if dataDirectory == "" {
		return data, nil
	}

	errStatusEffects := fetchStatusEffects(dataDirectory, data)
	if errStatusEffects != nil {
		return nil, errStatusEffects
	}

	return data, nil
}

func fetchStatusEffects(dataDirectory string, data map[string]interface{}) error {
	var sl StatusEffectList

	file, errOpen := os.Open(fmt.Sprintf("%s/docs/status-effects/current/index.json", dataDirectory))
	if errOpen != nil {
		return errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&sl)
	if errJSON != nil {
		return errJSON
	}

	effects := map[string]*StatusEffect{}
	for _, effect := range sl {
		effects[effect.Slug] = effect
	}
	data["AllEffects"] = effects

	return nil
}
