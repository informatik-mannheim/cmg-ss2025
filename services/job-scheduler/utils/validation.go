package utils

import (
	"net/url"
)

func IsUrlValid(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsPortValid(port int) bool {
	return port > 0 && port <= 65535
}
