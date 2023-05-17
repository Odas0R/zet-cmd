package slugify

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func Slug(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	result = strings.ToLower(result)                  // Convert to lowercase
	result = strings.ReplaceAll(result, " ", "-")     // Replace spaces with hyphen
	result = strings.Map(replaceSpecialChars, result) // Replace special characters
	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func replaceSpecialChars(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
		return r
	}
	return '-'
}
