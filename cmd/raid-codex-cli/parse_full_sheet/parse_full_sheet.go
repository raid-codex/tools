package parse_full_sheet

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/utils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Command struct {
	CSVFile            *string
	Type               *string
	ChampionsDirectory *string
}

// SRC: https://drive.google.com/file/d/1Xc66CarzqyoOPHcJqoFAAE_lqyyHlB7L/view

func New(cmd *kingpin.CmdClause) *Command {
	return &Command{
		CSVFile:            cmd.Flag("csv-file", "CSV File to parse").Required().String(),
		Type:               cmd.Flag("type", "Type of CSV sheet").Required().Enum("basic", "detail", "reviews"),
		ChampionsDirectory: cmd.Flag("champions-directory", "Champions directory").Required().String(),
	}
}

func (c *Command) Run() {
	champions, errChampions := c.fetchChampions()
	if errChampions != nil {
		utils.Exit(1, errChampions)
	}
	championsByName := map[string]*common.Champion{}
	for _, champion := range champions {
		championsByName[champion.Name] = champion
	}
	file, errFile := os.Open(*c.CSVFile)
	if errFile != nil {
		utils.Exit(1, errFile)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'
	content, errRead := reader.ReadAll()
	if errRead != nil {
		utils.Exit(1, errRead)
	}
	var err error
	switch *c.Type {
	case "basic":
		err = c.handleBasic(championsByName, content)
	case "detail":
		err = c.handleDetail(championsByName, content)
	case "reviews":
		err = c.handleReviews(championsByName, content)
	default:
		err = fmt.Errorf("%s not implemented", *c.Type)
	}
	if err != nil {
		utils.Exit(1, err)
	}
}

var (
	replaceChampionName = map[string]string{
		"BigUn":         "Big'Un",
		"Knight-Errant": "Knight Errant",
		"MaShalled":     "Ma'Shalled",
	}
)

func (c *Command) handleBasic(champions map[string]*common.Champion, content [][]string) error {
	for idx, line := range content {
		if len(line) != 11 {
			return fmt.Errorf("line %s has %d parts, not 11", strings.Join(line, ";"), len(line))
		} else if idx == 0 {
			if strings.Join(line, ";") != ";Name;Rarity;Type;Faction;Aura;Based On;Targets;Hits;Effect 1;Buff/Debuff Who" {
				return fmt.Errorf("invalid first line %s", strings.Join(line, ";"))
			}
		} else {
			if _, ok := replaceChampionName[line[1]]; ok {
				line[1] = replaceChampionName[line[1]]
			}
			champion := champions[line[1]]
			if champion == nil {
				panic(fmt.Sprintf("champion %s not found", line[1]))
			}
			if line[5] != "-" {
				champion.SetAura(line[5])
			}
			errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *c.ChampionsDirectory, champion.Filename()), champion)
			if errWrite != nil {
				utils.Exit(1, errWrite)
			}
		}
	}
	return nil
}

func (c *Command) handleReviews(champions map[string]*common.Champion, content [][]string) error {
	for idx, line := range content {
		if len(line) != 21 {
			return fmt.Errorf("line %s has %d parts, not 21", strings.Join(line, ";"), len(line))
		} else if idx == 0 {
			if strings.Join(line, ";") != "Got;;Name;Rarity;Type;Faction;Review;Camp;Arena Def;Arena Off;Min;Spid;Fire;Clan;Force;Dragon;Ice;Void;Spirit;Magic;" {
				return fmt.Errorf("invalid first line %s", strings.Join(line, ";"))
			}
		} else {
			if _, ok := replaceChampionName[line[2]]; ok {
				line[2] = replaceChampionName[line[2]]
			}
			champion := champions[line[2]]
			if champion == nil {
				panic(fmt.Sprintf("champion %s not found", line[2]))
			}
			review := &common.Review{}
			val, err := parseFloat(line[6])
			if err != nil {
				utils.Exit(1, err)
			}
			review.NumberOfReviews = int64(val)
			val, err = parseFloat(line[7])
			if err != nil {
				utils.Exit(1, err)
			}
			review.Campaign = val
			val, err = parseFloat(line[8])
			if err != nil {
				utils.Exit(1, err)
			}
			review.ArenaDef = val
			val, err = parseFloat(line[9])
			if err != nil {
				utils.Exit(1, err)
			}
			review.ArenaOff = val
			val, err = parseFloat(line[10])
			if err != nil {
				utils.Exit(1, err)
			}
			review.Minotaur = val
			val, err = parseFloat(line[11])
			if err != nil {
				utils.Exit(1, err)
			}
			review.Spider = val
			val, err = parseFloat(line[12])
			if err != nil {
				utils.Exit(1, err)
			}
			review.FireKnight = val
			val, err = parseFloat(line[13])
			if err != nil {
				utils.Exit(1, err)
			}
			review.ClanBoss = val
			val, err = parseFloat(line[14])
			if err != nil {
				utils.Exit(1, err)
			}
			review.ForceDungeon = val
			val, err = parseFloat(line[15])
			if err != nil {
				utils.Exit(1, err)
			}
			review.Dragon = val
			val, err = parseFloat(line[16])
			if err != nil {
				utils.Exit(1, err)
			}
			review.IceGuardian = val
			val, err = parseFloat(line[17])
			if err != nil {
				utils.Exit(1, err)
			}
			review.VoidDungeon = val
			val, err = parseFloat(line[18])
			if err != nil {
				utils.Exit(1, err)
			}
			review.SpiritDungeon = val
			val, err = parseFloat(line[19])
			if err != nil {
				utils.Exit(1, err)
			}
			review.MagicDungeon = val
			champion.Reviews = review
			errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *c.ChampionsDirectory, champion.Filename()), champion)
			if errWrite != nil {
				utils.Exit(1, errWrite)
			}
		}
	}
	return nil
}

func (c *Command) handleDetail(champions map[string]*common.Champion, content [][]string) error {
	for idx, line := range content {
		if len(line) != 39 {
			return fmt.Errorf("line %s has %d parts, not 39", strings.Join(line, ";"), len(line))
		} else if idx == 0 {
			if strings.Join(line, ";") != "Aura;title;skill_1;;Level;Based On;Targets;Hits;Who;Chance;%;Effect 1;B Who;Turns;Places If;Chance;%;Effect 2;B Who;Turns;Places If;Chance;%;Effect 3;B Who;Turns;Places If;Chance;%;Effect 4;B Who;Turns;Places If;Chance;%;Effect 5;B Who;Turns;Places If" {
				return fmt.Errorf("invalid first line %s", strings.Join(line, ";"))
			}
		} else {
			for i := range line {
				line[i] = strings.Trim(line[i], " !")
			}
			if _, ok := replaceChampionName[line[1]]; ok {
				line[1] = replaceChampionName[line[1]]
			}
			champion := champions[line[1]]
			if champion == nil {
				panic(fmt.Sprintf("champion %s not found", line[1]))
			}
			log.Printf("%d/%d ;; treating champion %s\n\tline: %s\n", idx, len(content), champion.Name, strings.Join(line, ";"))
			if line[0] != "" {
				champion.SetAura(line[0])
			}
			passive := false
			if strings.Contains(line[2], " [Passive]") {
				passive = true
				line[2] = strings.Replace(line[2], " [Passive]", "", -1)
			} else if strings.Contains(line[2], " [P]") {
				passive = true
				line[2] = strings.Replace(line[2], " [P]", "", -1)
			}
			skill, errSkill := champion.GetSkillByName(line[2])
			if errSkill != nil {
				skill = &common.Skill{Name: line[2]}
				champion.Skills = append(champion.Skills, skill)
			}
			skill.Passive = passive
			if line[3] != "-" {
				sd := &common.SkillData{Level: line[3]}
				basedOn := strings.Split(line[5], "/")
				if basedOn[0] != "-" {
					sd.BasedOn = basedOn
				}
				hits, errHits := parseInt(line[7])
				if errHits != nil {
					utils.Exit(1, errHits)
				}
				sd.Hits = hits
				sd.Target = &common.Target{Who: line[8], Targets: line[7]}
				for _, i := range []int64{9, 15, 21, 27} {
					if line[i+2] == "-" {
						continue
					}
					chance, errChance := parseFloat(line[i])
					if errChance != nil {
						log.Printf("invalid chance value\n")
						utils.Exit(1, errChance)
					}
					value, errValue := parseFloat(line[i+1])
					if errValue != nil {
						log.Printf("invalid value value\n")
						utils.Exit(1, errValue)
					}
					amount := int64(1)
					if strings.HasPrefix(line[i+1], "2x") {
						amount = 2
					} else if strings.HasPrefix(line[i+1], "3x") {
						amount = 3
					}
					effect := line[i+2]
					who := line[i+3]
					turns, errTurns := parseInt(line[i+4])
					if errTurns != nil {
						log.Printf("invalid turn value\n")
						utils.Exit(1, errTurns)
					}
					placesIf := line[i+5]
					sd.AddEffect(effect, who, turns, float64(chance)/100.0, placesIf, float64(value)/100.0, amount)
				}
				skill.SetSkillData(sd)
			}
			errSanitize := champion.Sanitize()
			if errSanitize != nil {
				utils.Exit(1, errSanitize)
			}
			errWrite := utils.WriteToFile(fmt.Sprintf("%s/%s", *c.ChampionsDirectory, champion.Filename()), champion)
			if errWrite != nil {
				utils.Exit(1, errWrite)
			}
		}
	}
	return nil
}

func parseInt(str string) (int64, error) {
	value, errValue := strconv.ParseInt(str, 10, 64)
	if errValue != nil && str != "-" {
		if strings.HasPrefix(str, "x") {
			return parseInt(str[1:])
		} else if strings.HasSuffix(str, "(2)") {
			return parseInt(str[:len(str)-3])
		} else if strings.HasSuffix(str, " turn") {
			return parseInt(str[:len(str)-5])
		}
		return 0, fmt.Errorf("cannot parse '%s': %s", str, errValue)
	}
	return value, nil
}

func parseFloat(str string) (float64, error) {
	value, errValue := strconv.ParseFloat(str, 64)
	if errValue != nil && str != "-" {
		if strings.HasPrefix(str, "x") {
			return parseFloat(str[1:])
		} else if strings.Contains(str, ",") {
			return parseFloat(strings.Replace(str, ",", ".", -1))
		} else if strings.HasSuffix(str, " turn") {
			return parseFloat(str[:len(str)-5])
		} else if strings.HasSuffix(str, " Turn") {
			return parseFloat(str[:len(str)-5])
		} else if strings.HasPrefix(str, "2x") || strings.HasPrefix(str, "3x") {
			return parseFloat(str[2:])
		} else if strings.HasSuffix(str, "+") {
			return parseFloat(str[:len(str)-1])
		}
		return 0, fmt.Errorf("cannot parse %s: %s", str, errValue)
	}
	return value, nil
}

func (c *Command) fetchChampions() (common.ChampionList, error) {
	dir, err := ioutil.ReadDir(*c.ChampionsDirectory)
	if err != nil {
		return nil, err
	}
	var champions common.ChampionList
	for _, file := range dir {
		if file.Name() == "index.json" {
			continue
		} else if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		var champion common.Champion
		err := func() error {
			f, errOpen := os.Open(fmt.Sprintf("%s/%s", *c.ChampionsDirectory, file.Name()))
			if errOpen != nil {
				return errOpen
			}
			defer f.Close()
			errJSON := json.NewDecoder(f).Decode(&champion)
			return errJSON
		}()
		if err != nil {
			return nil, err
		}
		champions = append(champions, &champion)
	}
	return champions, nil
}
