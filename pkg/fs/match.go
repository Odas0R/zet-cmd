package fs

import "strings"

func MatchSubstring(startS string, endS string, str string) string {
	s := strings.Index(str, startS)
	if s == -1 {
		return ""
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return ""
	}
	result := newS[:e]
	return result
}

func MatchAllSubstrings(startS string, endS string, str string) []string {
	var results []string

	s := strings.Index(str, startS)
	for s != -1 {
		newS := str[s+len(startS):]
		e := strings.Index(newS, endS)
		if e == -1 {
			break
		}
		results = append(results, newS[:e])

		// Look for the next link.
		str = newS[e+len(endS):]
		s = strings.Index(str, startS)
	}
	return results
}
