package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Slugify(input string) string {
	t := norm.NFD.String(input)

	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}

	s := strings.ToLower(b.String())

	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	reg2 := regexp.MustCompile(`-+`)
	s = reg2.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	if len(s) > 60 {
		s = s[:60]
		s = strings.Trim(s, "-")
	}

	return s
}
