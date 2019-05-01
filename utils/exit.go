package utils

import (
	"fmt"
	"os"
)

func Exit(code int, err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(code)
}
