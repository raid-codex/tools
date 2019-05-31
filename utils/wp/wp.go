package wp

import (
	"bytes"
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common"
	"github.com/raid-codex/tools/common/paged"
	"github.com/sogko/go-wordpress"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

func GetWPClient() *wordpress.Client {
	// create wp-api client
	client := wordpress.NewClient(&wordpress.Options{
		BaseAPIURL: "https://raid-codex.com/wp-json/wp/v2",
		Username:   os.Getenv("WP_USER"),
		Password:   os.Getenv("WP_PASSWORD"),
	})
	return client
}

func getContent(page paged.Paged, templateFile, dataDirectory string) (string, error) {
	if templateFile == "" {
		return "", nil
	}
	data, errData := common.GetPageExtraData(dataDirectory)
	if errData != nil {
		return "", errData
	}
	inputFile, errInput := os.Open(templateFile)
	if errInput != nil {
		return "", errInput
	}
	defer inputFile.Close()
	buf := bytes.NewBufferString("")
	errTemplate := page.GetPageContent(inputFile, buf, data)
	if errTemplate != nil {
		return "", errTemplate
	}
	str := buf.String()
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	s, err := m.String("text/html", str)
	if err != nil {
		return "", err
	}
	return s, nil
}

func CreatePage(client *wordpress.Client, page paged.Paged, templateFile, dataDirectory string) error {
	content, err := getContent(page, templateFile, dataDirectory)
	if err != nil {
		return errors.Annotatef(err, "error while creating page")
	}
	_, _, body, err := client.Pages().Create(&wordpress.Page{
		Slug:     page.LinkName(),
		Title:    wordpress.Title{Raw: page.GetPageTitle()},
		Type:     "page",
		Template: page.GetPageTemplate(),
		Status:   "private",
		Parent:   page.GetParentPageID(),
		Content:  wordpress.Content{Raw: content},
		Excerpt:  wordpress.Excerpt{Raw: page.GetPageExcerpt()},
	})
	if err != nil {
		logWordpressError(body)
		return errors.Annotatef(err, "error while creating page")
	}
	return nil
}

func UpdatePage(client *wordpress.Client, wpPage *wordpress.Page, page paged.Paged, templateFile, dataDirectory string) error {
	content, err := getContent(page, templateFile, dataDirectory)
	if err != nil {
		return errors.Annotatef(err, "error while creating page")
	}
	_, _, body, err := client.Pages().Update(wpPage.ID, &wordpress.Page{
		Content:  wordpress.Content{Raw: content},
		Excerpt:  wordpress.Excerpt{Raw: page.GetPageExcerpt()},
		Template: page.GetPageTemplate(),
	})
	if err != nil {
		logWordpressError(body)
		return errors.Annotatef(err, "error while updating page")
	}
	return nil
}

func GetPageFromSlug(client *wordpress.Client, slug string) (*wordpress.Page, error) {
	pages, _, body, err := client.Pages().List(map[string]string{
		"slug":   slug,
		"status": "private,publish,draft",
	})
	if err != nil {
		logWordpressError(body)
		return nil, err
	} else if len(pages) == 0 {
		return nil, errors.NotFoundf("page")
	}
	return &pages[0], nil
}

func logWordpressError(body []byte) {
	fmt.Fprintf(os.Stderr, "%s\n", body)
}
