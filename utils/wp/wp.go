package wp

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common/paged"
	"github.com/raid-codex/tools/utils/minify"
	"github.com/sogko/go-wordpress"
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

func getContent(page paged.Paged, templateFile, dataDirectory string, tmpl *template.Template) (string, error) {
	if templateFile == "" {
		return "", nil
	}
	data, errData := page.GetPageExtraData(dataDirectory)
	if errData != nil {
		return "", errData
	}
	buf := bytes.NewBufferString("")
	var errTemplate error
	if tmpl != nil {
		errTemplate = page.GetPageContent_Templates(tmpl, buf, data)
	} else {
		inputFile, errInput := os.Open(templateFile)
		if errInput != nil {
			return "", errInput
		}
		defer inputFile.Close()
		errTemplate = page.GetPageContent(inputFile, buf, data)
	}
	if errTemplate != nil {
		return "", errTemplate
	}
	return minify.HTML(buf.String())
}

func CreatePage(client *wordpress.Client, page paged.Paged, templateFile, dataDirectory string, tmpl *template.Template) error {
	content, err := getContent(page, templateFile, dataDirectory, tmpl)
	if err != nil {
		return errors.Annotatef(err, "error while creating page")
	}
	return CreatePage_Content(client, page, content)
}

func CreatePage_Content(client *wordpress.Client, page paged.Paged, content string) error {
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

func UpdatePage(client *wordpress.Client, wpPage *wordpress.Page, page paged.Paged, templateFile, dataDirectory string, tmpl *template.Template) error {
	content, err := getContent(page, templateFile, dataDirectory, tmpl)
	if err != nil {
		return errors.Annotatef(err, "error while creating page")
	}
	return UpdatePage_Content(client, wpPage, page, content)
}

func UpdatePage_Content(client *wordpress.Client, wpPage *wordpress.Page, page paged.Paged, content string) error {
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
