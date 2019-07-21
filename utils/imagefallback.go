package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	cache = map[string]bool{}
)

func ImageFallback(urls ...string) (string, error) {
	for _, url := range urls {
		if v, ok := cache[url]; ok {
			if v {
				return url, nil
			} else {
				continue
			}
		}
		if err := imageFallback(url); err == nil {
			cache[url] = true
			return url, nil
		} else {
			cache[url] = false
			log.Printf("url: %s -> %v\n", url, err)
		}
	}
	return "", fmt.Errorf("no url provided")
}

func imageFallback(url string) error {
	if strings.HasPrefix(url, "data:image/") {
		return nil
	}
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	} else if resp.StatusCode == 404 {
		return err
	}
	return nil
}
