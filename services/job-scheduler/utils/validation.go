package utils

import (
	"net/url"
)

func IsUrlValid(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}
