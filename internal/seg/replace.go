package seg

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// matchesAtEnd reports whether re matches a suffix of s.
// Behind-patterns are written to end with \z, so a short suffix is enough.
func matchesAtEnd(s string, re *regexp.Regexp) bool {
	if re == nil {
		return true
	}
	const window = 96
	if len(s) > window {
		s = s[len(s)-window:]
	}
	return re.MatchString(s)
}

func matchesAtStart(s string, re *regexp.Regexp) bool {
	if re == nil {
		return true
	}
	const window = 96
	if len(s) > window {
		s = s[:window]
	}
	loc := re.FindStringIndex(s)
	return loc != nil && loc[0] == 0
}

// replaceFind replaces matches of find whose prefix ends with behind and
// whose suffix starts with ahead. behind/ahead may be nil (unconditional).
// Lookaround is not consumed, so overlapping candidates work like Ruby gsub.
func replaceFind(s string, behind, find, ahead *regexp.Regexp, repl string) string {
	return replaceFindFunc(s, behind, find, ahead, func(string) string { return repl })
}

func replaceFindFunc(s string, behind, find, ahead *regexp.Regexp, repl func(match string) string) string {
	if s == "" || find == nil {
		return s
	}
	locs := find.FindAllStringIndex(s, -1)
	if len(locs) == 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	pos := 0
	for _, loc := range locs {
		if loc[0] < pos {
			continue
		}
		prefix := s[:loc[0]]
		suffix := s[loc[1]:]
		if !matchesAtEnd(prefix, behind) || !matchesAtStart(suffix, ahead) {
			continue
		}
		b.WriteString(s[pos:loc[0]])
		b.WriteString(repl(s[loc[0]:loc[1]]))
		pos = loc[1]
	}
	b.WriteString(s[pos:])
	return b.String()
}

// replaceLiteralFind is replaceFind for a literal target string.
func replaceLiteralFind(s string, behind *regexp.Regexp, target string, ahead *regexp.Regexp, repl string) string {
	if target == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	pos := 0
	for pos <= len(s)-len(target) {
		j := strings.Index(s[pos:], target)
		if j < 0 {
			break
		}
		i := pos + j
		prefix := s[:i]
		suffix := s[i+len(target):]
		if matchesAtEnd(prefix, behind) && matchesAtStart(suffix, ahead) {
			b.WriteString(s[pos:i])
			b.WriteString(repl)
			pos = i + len(target)
			continue
		}
		// Advance one rune so we can find overlapping candidates.
		_, w := utf8.DecodeRuneInString(s[i:])
		if w < 1 {
			w = 1
		}
		b.WriteString(s[pos : i+w])
		pos = i + w
	}
	b.WriteString(s[pos:])
	return b.String()
}
