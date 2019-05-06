package seo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/juju/errors"
	"github.com/raid-codex/tools/utils"
)

type SEO struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

func (s *SEO) Apply(pageSlug string) error {
	payload := applyPayload{
		YoastTitle:           s.Title,
		YoastMetaDescription: s.Description,
		YoastFocusKeywords:   strings.Join(s.Keywords, " "),
		YoastOGDescription:   s.Description,
		YoastOGTitle:         s.Title,
	}
	client := utils.GetWPClient()
	page, errPage := utils.GetPageFromSlug(client, pageSlug)
	if errPage != nil {
		return errors.Annotate(errPage, "cannot load page")
	}
	var current applyPayload
	pageUrl := fmt.Sprintf("https://raid-codex.com/wp-json/wp/v2/pages/%d", page.ID)
	_, _, errGetCurrentSEOTags := client.Get(pageUrl, nil, &current)
	if errGetCurrentSEOTags != nil {
		return errors.Annotate(errGetCurrentSEOTags, "cannot fetch current SEO tags")
	}
	diff := current.Diff(payload)
	if diff != nil {
		for k, v := range diff {
			fmt.Printf("Diff with field %s:\n%s\n", k, v)
		}
	} else {
		// no diff, don't call API
		fmt.Printf("no diff, skipping\n")
		return nil
	}
	res := map[string]interface{}{}
	_, _, errUpdate := client.Update(pageUrl, payload, &res)
	if errUpdate != nil {
		return errors.Annotate(errUpdate, "cannot update SEO tags")
	}
	return nil
}

type applyPayload struct {
	YoastTitle           string `json:"_yoast_wpseo_title,omitempty"`
	YoastMetaDescription string `json:"_yoast_wpseo_metadesc,omitempty"`
	YoastFocusKeywords   string `json:"_yoast_wpseo_focuskw,omitempty"`
	YoastOGDescription   string `json:"_yoast_wpseo_opengraph-description,omitempty"`
	YoastOGTitle         string `json:"_yoast_wpseo_opengraph-title,omitempty"`
}

func (a applyPayload) Diff(b applyPayload) map[string]string {
	diff := map[string]string{}
	t := reflect.TypeOf(a)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	for i := 0; i < t.NumField(); i++ {
		if av.Field(i).String() != bv.Field(i).String() {
			diff[t.Field(i).Name] = fmt.Sprintf("\tcurrent: %s\n\tapplied: %s", av.Field(i).String(), bv.Field(i).String())
		}
	}
	if len(diff) == 0 {
		return nil
	}
	return diff
}
