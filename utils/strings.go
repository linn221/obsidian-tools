package utils

import (
	"regexp"
)

// for expressions such as `hello, (\d*)` with one capture group
func MatchRegex(re *regexp.Regexp, s string, f func(match string)) {
	matches := re.FindStringSubmatch(s)
	if len(matches) >= 2 {
		f(matches[1])
	}
}

func IsRegexMatch(re *regexp.Regexp, s string) bool {
	return re.Match([]byte(s))
}
