package paged

import "io"

type Paged interface {
	GetPageSlug() string
	GetPageTitle() string
	GetPageTemplate() string
	GetParentPageID() int
	LinkName() string
	GetPageContent(io.Reader, io.Writer) error
	GetPageExcerpt() string
}
