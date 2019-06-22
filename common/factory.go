package common

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	factory           *Factory
	ErrNotInitialized = fmt.Errorf("factory not initialized")
)

type Factory struct {
	champions     ChampionList
	statusEffects StatusEffectList
	factions      FactionList
}

func InitFactory(dataDirectory string) error {
	factory = &Factory{}
	return factory.init(dataDirectory)
}

func (f *Factory) init(dataDirectory string) error {
	if err := f.fetchChampions(fmt.Sprintf("%s/docs/champions/current/", dataDirectory)); err != nil {
		return err
	}
	if err := f.fetchFactions(fmt.Sprintf("%s/docs/factions/current/", dataDirectory)); err != nil {
		return err
	}
	if err := f.fetchStatusEffects(fmt.Sprintf("%s/docs/status-effects/current/", dataDirectory)); err != nil {
		return err
	}
	return nil
}

func (f *Factory) fetchChampions(dir string) error {
	file, errOpen := os.Open(fmt.Sprintf("%s/index.json", dir))
	if errOpen != nil {
		return errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&f.champions)
	if errJSON != nil {
		return errJSON
	}
	return nil
}

func (f *Factory) fetchFactions(dir string) error {
	file, errOpen := os.Open(fmt.Sprintf("%s/index.json", dir))
	if errOpen != nil {
		return errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&f.factions)
	if errJSON != nil {
		return errJSON
	}
	return nil
}
func (f *Factory) fetchStatusEffects(dir string) error {
	file, errOpen := os.Open(fmt.Sprintf("%s/index.json", dir))
	if errOpen != nil {
		return errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&f.statusEffects)
	if errJSON != nil {
		return errJSON
	}
	return nil
}

func GetChampions(filters ...ChampionFilter) (ChampionList, error) {
	if factory == nil {
		return nil, ErrNotInitialized
	}
	cl := make(ChampionList, 0)
CHAMPIONS:
	for _, champion := range factory.champions {
		for _, filter := range filters {
			if !filter(champion) {
				continue CHAMPIONS
			}
		}
		cl = append(cl, champion)
	}
	return cl, nil
}