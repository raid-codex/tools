package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raid-codex/tools/common"
	"github.com/sogko/go-wordpress"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	championTypeCmd              = kingpin.Command("create-champion", "Create a champion page")
	championFile                 = championTypeCmd.Arg("champion-file", "Filename for the champion").Required().String()
	factionTypeCmd               = kingpin.Command("create-faction", "Create a faction page")
	factionFile                  = factionTypeCmd.Arg("faction-file", "Filename for the faction").Required().String()
	championImageCmd             = kingpin.Command("champion-image", "Upload a champion image")
	championImageFile            = championImageCmd.Arg("image-file", "Link to image").Required().String()
	championImageCmdChampionFile = championImageCmd.Arg("champion-file", "Which champion this is uploaded for").Required().String()
)

func main() {
	// create wp-api client
	client := wordpress.NewClient(&wordpress.Options{
		BaseAPIURL: "https://raid-codex.com/wp-json/wp/v2",
		Username:   os.Getenv("WP_USER"),
		Password:   os.Getenv("WP_PASSWORD"),
	})

	switch kingpin.Parse() {
	case "create-champion":
		processPage(client, *championFile, func(decoder func(v interface{}) error) (common.Paged, error) {
			champion := common.Champion{}
			return &champion, decoder(&champion)
		})
	case "create-faction":
		processPage(client, *factionFile, func(decoder func(v interface{}) error) (common.Paged, error) {
			faction := common.Faction{}
			return &faction, decoder(&faction)
		})
	case "champion-image":
		processMedia(client)
	default:
		exit(fmt.Errorf("	"))
	}
}

func processPage(client *wordpress.Client, filename string, decode func(func(v interface{}) error) (common.Paged, error)) {
	fmt.Printf("creating page for file %s\n", filename)
	file, errFile := os.Open(filename)
	if errFile != nil {
		exit(errFile)
	}
	defer file.Close()

	pageData, errJSON := decode(json.NewDecoder(file).Decode)
	if errJSON != nil {
		exit(errJSON)
	}
	slug := pageData.GetPageSlug()
	// does the page already exist ?
	if pageExists, err := pageWithSlugExists(client, slug); err != nil {
		exit(err)
	} else if pageExists != nil {
		fmt.Println("page already exists")
		_, _, body, err := client.Pages().Update(pageExists.ID, &wordpress.Page{
			Slug:   pageData.LinkName(),
			Parent: pageData.GetParentPageID(),
		})
		if err != nil {
			fmt.Printf("%s\n", body)
			exit(err)
		}
	} else {
		fmt.Println("will create page")
		// create page
		_, _, body, err := client.Pages().Create(&wordpress.Page{
			Slug:     pageData.LinkName(),
			Title:    wordpress.Title{Raw: pageData.GetPageTitle()},
			Type:     "page",
			Template: pageData.GetPageTemplate(),
			Status:   "private",
			Parent:   pageData.GetParentPageID(),
		})
		if err != nil {
			fmt.Printf("%s\n", body)
			exit(err)
		}

	}
}

func pageWithSlugExists(client *wordpress.Client, slug string) (*wordpress.Page, error) {
	pages, _, body, err := client.Pages().List(map[string]string{
		"slug":   slug,
		"status": "private,publish,draft",
	})
	if err != nil {
		fmt.Printf("%s\n", body)
		exit(err)
	} else if len(pages) > 0 {
		return &pages[0], nil
	}
	return nil, nil
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

func processMedia(client *wordpress.Client) {
	exit(fmt.Errorf("test"))
}
