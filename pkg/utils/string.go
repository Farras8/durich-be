package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GenerateSlug(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	input, _, _ = transform.String(t, input)

	reg := regexp.MustCompile("[^a-z0-9]+")
	input = reg.ReplaceAllString(input, "-")

	input = strings.Trim(input, "-")

	reg = regexp.MustCompile("-+")
	input = reg.ReplaceAllString(input, "-")

	if len(input) > 200 {
		input = input[:200]
		input = strings.TrimRight(input, "-")
	}

	return input
}
