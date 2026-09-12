package seg

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func replaceAbbreviations(s string, lang *language) string {
	s = rePossessiveAbbr.ReplaceAllString(s, subPeriod+"$1")
	if !lang.germanLowerLetterAbbr {
		s = reKommandit.ReplaceAllString(s, "${1}"+subPeriod+"${3}")
	}
	if lang.germanLowerLetterAbbr {
		s = applySingleLetterAbbr(s)
		s = reSingleLowerLetter.ReplaceAllString(s, "${1}"+subPeriod+"${3}")
		s = reSingleLowerLetterLine.ReplaceAllString(s, "${1}"+subPeriod+"${3}")
	} else {
		s = applySingleLetterAbbr(s)
	}

	s = searchAbbreviations(s, lang)
	s = replaceMultiPeriodAbbreviations(s, lang)
	s = applyAmPmRules(s)
	s = replaceAbbreviationAsSentenceBoundary(s, lang)
	if lang.cyrillicSingleLetter {
		s = reCyrLetterStart.ReplaceAllString(s, "${1}"+subPeriod+"${3}")
		s = reCyrLetter.ReplaceAllString(s, "${1}"+subPeriod+"${3}")
	}
	return s
}

func applySingleLetterAbbr(s string) string {
	// Ruby: (?<=^[A-Z])\.(?=,?\s) and (?<=\s[A-Z])\.(?=,?\s)
	// Only the period is replaced, so "A. L. Burt" hides both initials.
	var b strings.Builder
	b.Grow(len(s))
	pos := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '.' || i == 0 {
			continue
		}
		c := s[i-1]
		if c < 'A' || c > 'Z' {
			continue
		}
		if i >= 2 && !isWSByte(s[i-2]) && s[i-2] != '\n' {
			continue
		}
		rest := s[i+1:]
		if rest == "" {
			continue
		}
		j := 0
		if rest[0] == ',' {
			j = 1
			if j >= len(rest) {
				continue
			}
		}
		if !isWSByte(rest[j]) {
			continue
		}
		b.WriteString(s[pos:i])
		b.WriteString(subPeriod)
		pos = i + 1
	}
	if pos == 0 {
		return s
	}
	b.WriteString(s[pos:])
	return b.String()
}

func applyAmPmRules(s string) string {
	s = reAmPmUpperP.ReplaceAllString(s, "${1}.${2}")
	s = reAmPmUpperA.ReplaceAllString(s, "${1}.${2}")
	s = reAmPmLowerP.ReplaceAllString(s, "${1}.${2}")
	s = reAmPmLowerA.ReplaceAllString(s, "${1}.${2}")
	return s
}

func searchAbbreviations(s string, lang *language) string {
	original := s
	downcased := strings.ToLower(s)
	for _, ca := range lang.compiled {
		if !strings.Contains(downcased, ca.stripped) {
			continue
		}
		if lang.mode == abbrAlwaysBeforeSpace {
			if ca.search.MatchString(original) {
				s = hideAbbrPeriods(s, ca, reSpaceStart)
			}
			continue
		}
		if lang.mode == abbrAlwaysPeriod {
			if ca.search.MatchString(original) {
				s = hideAbbrPeriods(s, ca)
			}
			continue
		}
		if ca.pre {
			if ca.search.MatchString(original) {
				s = hideAbbrPeriods(s, ca, reSpaceStart, reAheadColonDigits)
			}
			continue
		}
		hits := ca.search.FindAllStringIndex(original, -1)
		if len(hits) == 0 {
			continue
		}
		s = applyAbbrHits(s, ca, lang, hits)
	}
	return s
}

func applyAbbrHits(s string, ca compiledAbbr, lang *language, hits [][]int) string {
	if !abbrShouldApply(s, ca, hits) {
		return s
	}
	if ca.num {
		return hideAbbrPeriods(s, ca, reAheadSpaceDigit, reAheadSpaceParen)
	}
	if lang.mode == abbrPeriodAnywhere {
		return hideAbbrPeriods(s, ca)
	}
	return hideAbbrPeriods(s, ca, reAheadAbbrPeriod, reAheadAbbrPeriod2, reAheadComma)
}

func abbrShouldApply(s string, ca compiledAbbr, hits [][]int) bool {
	chars := ca.nextChar.FindAllStringSubmatch(s, -1)
	for i := range hits {
		var ch string
		if i < len(chars) && len(chars[i]) > 1 {
			ch = chars[i][1]
		}
		if ch == "" {
			return true
		}
		r, _ := utf8.DecodeRuneInString(ch)
		if !unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

var (
	reAheadSpaceDigit  = mustRE(`^\s\d`)
	reAheadSpaceParen  = mustRE(`^\s+\(`)
	reAheadColonDigits = mustRE(`^:\d+`)
	reAheadAbbrPeriod  = mustRE(`^(?:[.:?-]|\s(?:[a-z]|I\s|I'm|I'll|\d|\())`)
	reAheadAbbrPeriod2 = mustRE(`^(?:[.:?]|\s(?:[a-z]|I\s|I'm|I'll|\d))`)
	reAheadComma       = mustRE(`^,`)
)

func hideAbbrPeriods(s string, ca compiledAbbr, aheads ...*regexp.Regexp) string {
	if ca.ascii {
		return hideAbbrPeriodsASCII(s, ca.stripped, aheads)
	}
	return hideAbbrPeriodsRE(s, ca.hide, aheads)
}

func hideAbbrPeriodsRE(s string, hideRE *regexp.Regexp, aheads []*regexp.Regexp) string {
	if hideRE == nil {
		return s
	}
	locs := hideRE.FindAllStringSubmatchIndex(s, -1)
	if len(locs) == 0 {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 32)
	pos := 0
	changed := false
	for _, loc := range locs {
		if len(loc) < 4 || loc[2] < pos {
			continue
		}
		if !aheadOK(s[loc[3]:], aheads) {
			continue
		}
		b.WriteString(s[pos:loc[2]])
		b.WriteString(subPeriod)
		pos = loc[3]
		changed = true
	}
	if !changed {
		return s
	}
	b.WriteString(s[pos:])
	return b.String()
}

func hideAbbrPeriodsASCII(s, abbr string, aheads []*regexp.Regexp) string {
	n := len(abbr)
	var b strings.Builder
	b.Grow(len(s) + 32)
	pos := 0
	start := 0
	changed := false
	for start+n < len(s) {
		i := indexASCIIFold(s, abbr, start)
		if i < 0 {
			break
		}
		if i > 0 && !isWSByte(s[i-1]) {
			start = i + 1
			continue
		}
		per := i + n
		if per >= len(s) || s[per] != '.' {
			start = i + 1
			continue
		}
		if !aheadOK(s[per+1:], aheads) {
			start = per + 1
			continue
		}
		b.WriteString(s[pos:per])
		b.WriteString(subPeriod)
		pos = per + 1
		changed = true
		start = per + 1
	}
	if !changed {
		return s
	}
	b.WriteString(s[pos:])
	return b.String()
}

func aheadOK(after string, aheads []*regexp.Regexp) bool {
	if len(aheads) == 0 {
		return true
	}
	for _, re := range aheads {
		if re != nil && re.MatchString(after) {
			return true
		}
	}
	return false
}

func indexASCIIFold(s, lowerAbbr string, start int) int {
	n := len(lowerAbbr)
	if n == 0 || start+n > len(s) {
		return -1
	}
	first := lowerAbbr[0]
	upper := first
	if first >= 'a' && first <= 'z' {
		upper = first - 'a' + 'A'
	}
	for i := start; i+n <= len(s); i++ {
		c := s[i]
		if c != first && c != upper {
			continue
		}
		if asciiFoldEqual(s[i:i+n], lowerAbbr) {
			return i
		}
	}
	return -1
}

func asciiFoldEqual(s, lower string) bool {
	if len(s) != len(lower) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c != lower[i] {
			return false
		}
	}
	return true
}

func isWSByte(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f' || b == '\v'
}

func replaceMultiPeriodAbbreviations(s string, lang *language) string {
	if lang.multiPeriodKazakh {
		// RE2 \b is ASCII-only; use a Unicode-aware boundary so Cyrillic
		// abbreviations like "Б.з.б." and "т. б." are found.
		for _, m := range reKazakhCyrillicPeriods.FindAllStringSubmatch(s, -1) {
			if len(m) > 1 {
				s = strings.ReplaceAll(s, m[1], strings.ReplaceAll(m[1], ".", subPeriod))
			}
		}
		for _, m := range reKazakhLatinPeriods.FindAllStringSubmatch(s, -1) {
			if len(m) > 1 {
				s = strings.ReplaceAll(s, m[1], strings.ReplaceAll(m[1], ".", subPeriod))
			}
		}
		return s
	}
	mpa := reMultiPeriodAbbr.FindAllString(s, -1)
	for _, r := range mpa {
		s = strings.ReplaceAll(s, r, strings.ReplaceAll(r, ".", subPeriod))
	}
	return s
}

func replaceAbbreviationAsSentenceBoundary(s string, lang *language) string {
	starters := lang.sentenceStarters
	if len(starters) == 0 {
		return s
	}
	if lang.danishAmPm {
		for _, word := range starters {
			esc := regexp.QuoteMeta(word)
			s = regexp.MustCompile(`U∯S∯\s`+esc+`\s`).ReplaceAllString(s, "U∯S. "+word+" ")
			s = regexp.MustCompile(`U\.S∯\s`+esc+`\s`).ReplaceAllString(s, "U.S. "+word+" ")
			s = regexp.MustCompile(`U∯K∯\s`+esc+`\s`).ReplaceAllString(s, "U∯K. "+word+" ")
			s = regexp.MustCompile(`U\.K∯\s`+esc+`\s`).ReplaceAllString(s, "U.K. "+word+" ")
			s = regexp.MustCompile(`E∯U∯\s`+esc+`\s`).ReplaceAllString(s, "E∯U. "+word+" ")
			s = regexp.MustCompile(`E\.U∯\s`+esc+`\s`).ReplaceAllString(s, "E.U. "+word+" ")
			s = regexp.MustCompile(`U∯S∯A∯\s`+esc+`\s`).ReplaceAllString(s, "U∯S∯A. "+word+" ")
			s = regexp.MustCompile(`U\.S\.A∯\s`+esc+`\s`).ReplaceAllString(s, "U.S.A. "+word+" ")
			s = regexp.MustCompile(`I∯\s`+esc+`\s`).ReplaceAllString(s, "I. "+word+" ")
			s = regexp.MustCompile(`s.u∯\s`+esc+`\s`).ReplaceAllString(s, "s.u. "+word+" ")
			s = regexp.MustCompile(`S.U∯\s`+esc+`\s`).ReplaceAllString(s, "S.U. "+word+" ")
		}
		return s
	}
	for _, re := range lang.starterRes {
		s = re.ReplaceAllString(s, "${1}.${2}")
	}
	return s
}
