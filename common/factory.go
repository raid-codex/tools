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
	fusions       FusionList
	masteries     MasteryList
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
	if err := f.fetchFusions(fmt.Sprintf("%s/docs/fusions/current/", dataDirectory)); err != nil {
		return err
	}
	if err := f.fetchMasteries(fmt.Sprintf("%s/docs/masteries/current/", dataDirectory)); err != nil {
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

func (f *Factory) fetchMasteries(dir string) error {
	file, errOpen := os.Open(fmt.Sprintf("%s/index.json", dir))
	if errOpen != nil {
		return errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&f.masteries)
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
func (f *Factory) fetchFusions(dir string) error {
	file, errOpen := os.Open(fmt.Sprintf("%s/index.json", dir))
	if errOpen != nil {
		return errOpen
	}
	defer file.Close()
	errJSON := json.NewDecoder(file).Decode(&f.fusions)
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
	cl.Sort()
	return cl, nil
}

func GetStatuseffects(filters ...StatusEffectFilter) (StatusEffectList, error) {
	if factory == nil {
		return nil, ErrNotInitialized
	}
	sel := make(StatusEffectList, 0)
STATUSEFFECTS:
	for _, statusEffect := range factory.statusEffects {
		for _, filter := range filters {
			if !filter(statusEffect) {
				continue STATUSEFFECTS
			}
		}
		sel = append(sel, statusEffect)
	}
	sel.Sort()
	return sel, nil
}

func GetFactions(filters ...FactionFilter) (FactionList, error) {
	if factory == nil {
		return nil, ErrNotInitialized
	}
	cl := make(FactionList, 0)
FACTIONS:
	for _, faction := range factory.factions {
		for _, filter := range filters {
			if !filter(faction) {
				continue FACTIONS
			}
		}
		cl = append(cl, faction)
	}
	cl.Sort()
	return cl, nil
}

func GetFusions(filters ...FusionFilter) (FusionList, error) {
	if factory == nil {
		return nil, ErrNotInitialized
	}
	cl := make(FusionList, 0)
FUSION:
	for _, fusion := range factory.fusions {
		for _, filter := range filters {
			if !filter(fusion) {
				continue FUSION
			}
		}
		cl = append(cl, fusion)
	}
	cl.Sort()
	return cl, nil
}

func GetMasteries(filters ...MasteryFilter) (MasteryList, error) {
	if factory == nil {
		return nil, ErrNotInitialized
	}
	ml := make(MasteryList, 0)
MASTERY:
	for _, mastery := range factory.masteries {
		for _, filter := range filters {
			if !filter(mastery) {
				continue MASTERY
			}
		}
		ml = append(ml, mastery)
	}
	ml.Sort()
	return ml, nil
}
