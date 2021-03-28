package parse_static_data

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Data struct {
	Strings  map[string]string `json:"StaticDataLocalization"`
	HeroData struct {
		HeroTypes                []HeroType
		HeroExperienceByKey      map[string]int64
		SacrificeExperienceByKey map[string]int64
		RankHeroCountByGrade     map[string]int64
		RankSilverByGrade        map[string]int64
		LevelUpPriceByGrade      map[string]struct {
			RawValues map[string]int64
		}
		FractionsByRace           map[string][]int64
		HeroIdsByRarities         map[string][]int64
		LevelUpMaterialsLimit     int64
		MultipleHeroesSummonCount int64
		MaxInventorySlotsCount    int64
		MaxStorageSlotsCount      int64
		InventorySlotsPrices      []struct {
			UserSlotsCount int64
			SilverPrice    struct {
				RawValues map[string]int64
			}
			GemsPrice struct {
				RawValues map[string]int64
			}
		}
		StorageSlotsPrices []struct {
			UserSlotsCount int64
			SilverPrice    struct {
				RawValues map[string]int64
			}
			GemsPrice struct {
				RawValues map[string]int64
			}
		}
		HeroesOnIntroFinish             []int64
		HeroRatingUpdateCooldownMinutes int64
		HeroPartsCountByHeroType        map[string]int64
	}
	SkillData struct {
		SkillTypes []SkillType
	}
}

type WeirdValue int64

func (v *WeirdValue) UnmarshalJSON(src []byte) error {
	val, err := strconv.ParseInt(string(src), 10, 64)
	if err != nil {
		return err
	}
	switch val {
	case 214748364:
		*v = 5
	case 4294967296:
		*v = 1
	default:
		*v = WeirdValue(val)
	}
	return nil
}

type Command struct {
	DataFile      *string
	DataDirectory *string
}

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		DataFile:      cmd.Flag("data-file", "Data File to parse").Required().String(),
		DataDirectory: cmd.Flag("data-directory", "Data directory").Required().String(),
	}
}

func (c *Command) Run() {
	content, errFile := ioutil.ReadFile(*c.DataFile)
	if errFile != nil {
		utils.Exit(1, errFile)
	}
	var data Data
	errDecode := json.Unmarshal(content, &data)
	if errDecode != nil {
		utils.Exit(1, errDecode)
	}
	skillsByID := map[int64]SkillType{}
	for _, skill := range data.SkillData.SkillTypes {
		skillsByID[skill.ID] = skill
	}
	championsByName := map[string]*HeroRepresentation{}
	for _, champion := range data.HeroData.HeroTypes {
		name := data.Strings[champion.Name.Key]
		rep, ok := championsByName[name]
		if !ok {
			rep = &HeroRepresentation{
				Name:   name,
				Awaken: map[uint8]HeroType{},
			}
			championsByName[name] = rep
		}
		giid, err := strconv.ParseInt(champion.AvatarName, 10, 64)
		if err != nil {
			utils.Exit(1, err)
		}
		rep.Awaken[uint8(champion.ID-giid)] = champion
	}
}

type HeroRepresentation struct {
	Name   string
	Awaken map[uint8]HeroType
}
