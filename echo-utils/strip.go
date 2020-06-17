package utils

import (
	"net/url"
	"regexp"
	"strings"
)

func StripForCache(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func RemoveNonAlphanumeric(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")

	reg, err := regexp.Compile(`[^a-zA-Z0-9]+`)
	if err != nil {
		panic(err)
	}

	s = reg.ReplaceAllString(s, "")

	return s
}

func StripSubdomainPlus(s string) string {

	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) > 2 {
		if parts[0] != "" {
			s = strings.ReplaceAll(s, parts[0]+".", "")
		}
	}

	s = strings.ReplaceAll(s, "http://", "")
	s = strings.ReplaceAll(s, "https://", "")

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return s
}
