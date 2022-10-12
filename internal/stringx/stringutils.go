// Package stringx has string utilities
package stringx

import (
	"regexp"
	"strings"

	"github.com/spotlightpa/almanack/internal/syncx"
)

// First non-blank string
func First(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func mustBeLazy(s string) func() *regexp.Regexp {
	return syncx.Once(
		func() *regexp.Regexp {
			return regexp.MustCompile(s)
		})
}

var (
	articleRe   = mustBeLazy(`\b(the|an?)\b`)
	pennRe      = mustBeLazy(`\bpa\b`)
	possesiveRe = mustBeLazy(`\.?[’']s`)
	nonasciiRe  = mustBeLazy(`\W+`)
)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = articleRe().ReplaceAllString(s, " ")
	s = pennRe().ReplaceAllString(s, "pennsylvania")
	s = possesiveRe().ReplaceAllString(s, "s")
	s = nonasciiRe().ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	s = nonasciiRe().ReplaceAllString(s, "-")
	return s
}

func LastCut(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return "", s, false
}
