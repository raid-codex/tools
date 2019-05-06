package utils

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/common/paged"
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

func CreatePage(client *wordpress.Client, page paged.Paged) error {
	_, _, body, err := client.Pages().Create(&wordpress.Page{
		Slug:     page.LinkName(),
		Title:    wordpress.Title{Raw: page.GetPageTitle()},
		Type:     "page",
		Template: page.GetPageTemplate(),
		Status:   "private",
		Parent:   page.GetParentPageID(),
	})
	if err != nil {
		logWordpressError(body)
		return errors.Annotatef(err, "error while creating page")
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
