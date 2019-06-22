package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func ImageFallback(urls ...string) (string, error) {
	for _, url := range urls {
		if strings.HasPrefix(url, "data:image/") {
			return url, nil
		}
		r, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == 404 {
			continue
		}
		// it's an image, return it
		return url, nil
	}
	return "", fmt.Errorf("no url provided")
}
