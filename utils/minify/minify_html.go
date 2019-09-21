package minify

import (
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

// HTML minifies the HTML string passed as parameter
func HTML(str string) (string, error) {
	str = strings.Replace(str, "\n", "", -1)
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	return m.String("text/html", str)
}
