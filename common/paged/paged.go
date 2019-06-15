package paged

import (
	"html/template"
	"io"
)

type Paged interface {
	GetPageSlug() string
	GetPageTitle() string
	GetPageTemplate() string
	GetParentPageID() int
	LinkName() string
	GetPageContent(io.Reader, io.Writer, map[string]interface{}) error
	GetPageContent_Templates(*template.Template, io.Writer, map[string]interface{}) error
	GetPageExcerpt() string
	GetPageExtraData(string) (map[string]interface{}, error)
}
