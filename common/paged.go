package common

type Paged interface {
	GetPageSlug() string
	GetPageTitle() string
	GetPageTemplate() string
	GetParentPageID() int
	LinkName() string
}
