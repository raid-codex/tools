package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	sourceFile   = kingpin.Arg("source-file", "CSV File with champions list and grades").Required().String()
	targetFolder = kingpin.Arg("target-folder", "Folder in which to create the JSON files").Required().String()
)

func main() {
	kingpin.Parse()
	content, errContent := getSourceFileContent()
	if errContent != nil {
		fmt.Fprintf(os.Stderr, "cannot read file: %s\n", errContent)
		os.Exit(1)
	}
	errExport := exportContent(content)
	if errExport != nil {
		fmt.Fprintf(os.Stderr, "error while exporting: %s\n", errExport)
		os.Exit(1)
	}
}

type Champions struct {
	Champions []*Champion
}

type Champion struct {
	Faction string `json:"faction"`
	Name    string `json:"name"`
	Rarity  string `json:"rarity"`
	Element string `json:"element"`
	Type    string `json:"type"`
	Rating  Rating `json:"rating"`
}

func (c *Champion) sanitize() error {
	name, err := strconv.Unquote(fmt.Sprintf(`"%s"`, c.Name))
	if err != nil {
		return err
	} else if name != c.Name {
		return fmt.Errorf("please change name of %s", name)
	}
	re := regexp.MustCompile("^([a-zA-Z0-9']+).*")
	c.Name = re.ReplaceAllString(name, `$1`)
	return nil
}

func (c Champion) filename() string {
	filename := c.Name
	for _, part := range []string{
		`'`, `"`, ` `,
	} {
		filename = strings.Replace(filename, part, "_", -1)
	}
	filename = strings.ToLower(filename)
	return fmt.Sprintf("%s.json", filename)
}

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

func getSourceFileContent() (*Champions, error) {
	file, errFile := os.Open(*sourceFile)
	if errFile != nil {
		return nil, errFile
	}
	defer file.Close()
	champions := &Champions{}
	reader := csv.NewReader(file)
	content, errRead := reader.ReadAll()
	if errRead != nil {
		return nil, errRead
	}
	for idx, line := range content {
		if len(line) != 20 {
			return nil, fmt.Errorf("line %s has %d parts, not 20", strings.Join(line, ","), len(line))
		} else if line[0] == "Factions" {
			if strings.Join(line, ",") != csvSafeguardOrder {
				return nil, fmt.Errorf("invalid csv")
			}
			// this is a header line, skip it
			continue
		} else if line[0] == "" && idx > 0 {
			// in case faction is not mentionned on every line, use the one from previous line
			line[0] = content[idx-1][0]
		}
		champion := &Champion{
			Faction: line[0],
			Name:    line[1],
			Rarity:  line[2],
			Element: line[3],
			Type:    line[4],
			Rating: Rating{
				Overall:       line[5],
				Campaign:      line[6],
				ArenaOff:      line[7],
				ArenaDef:      line[8],
				ClanBossWoGS:  line[9],
				ClanBosswGS:   line[10],
				IceGuardian:   line[11],
				Dragon:        line[12],
				Spider:        line[13],
				FireKnight:    line[14],
				Minotaur:      line[15],
				ForceDungeon:  line[16],
				MagicDungeon:  line[17],
				SpiritDungeon: line[18],
				VoidDungeon:   line[19],
			},
		}
		errSanitize := champion.sanitize()
		if errSanitize != nil {
			return nil, errSanitize
		}
		champions.Champions = append(champions.Champions, champion)
	}
	return champions, nil
}

const (
	csvSafeguardOrder = `Factions,Champion,Rarity,Element,Typ,Overall,Campaign,Arena-Off,Arena-Deff,CB (- GS),CB (+GS),IceG,Dragon,Spider,FK,Mino,Force,Magic,Spirit,Void`
)

type indexChampion struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func exportContent(content *Champions) error {
	var champions []indexChampion
	for _, champion := range content.Champions {
		filename := champion.filename()
		errWrite := writeToFile(filename, champion)
		if errWrite != nil {
			return errWrite
		}
		champions = append(champions, indexChampion{
			Name: champion.Name,
			URL:  filename,
		})
	}
	errWrite := writeToFile("index.json", champions)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func writeToFile(filename string, val interface{}) error {
	f, errOpen := os.OpenFile(fmt.Sprintf("%s/%s", *targetFolder, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if errOpen != nil {
		return errOpen
	}
	defer f.Close()
	data, errJSON := json.MarshalIndent(val, "", "  ")
	if errJSON != nil {
		return errJSON
	}
	_, errWrite := f.Write(data)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
