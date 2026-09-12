package seg

import (
	"regexp"
	"strings"
)

func betweenPunctuation(s string, lang *language) string {
	s = subBetweenSingleQuotes(s)
	s = replacePunctInMatches(s, reBetweenSlantedSingle, false)
	s = subBetweenDoubleQuotes(s, lang)
	s = replacePunctInMatches(s, reBetweenSquare, false)
	s = replacePunctInMatches(s, reBetweenParens, false)
	s = replacePunctInMatches(s, reBetweenQuoteArrow, false)
	s = replacePunctInMatches(s, reBetweenEmDashes, false)
	s = replacePunctInMatches(s, reBetweenQuoteSlanted, false)
	if lang.japaneseQuotes {
		s = replacePunctInMatches(s, reBetweenParensJA, false)
		s = replacePunctInMatches(s, reBetweenQuoteJA, false)
	}
	if lang.chineseQuotes {
		s = replacePunctInMatches(s, reBetweenAngleZH, false)
		s = replacePunctInMatches(s, reBetweenLBracket, false)
	}
	return s
}

func subBetweenSingleQuotes(s string) string {
	if reLeadingApostrophe.MatchString(s) && !reQuoteSpace.MatchString(s) {
		return s
	}
	return replacePunctInMatches(s, reBetweenSingleQuotes, true)
}

func subBetweenDoubleQuotes(s string, lang *language) string {
	if lang.germanQuotes {
		if strings.Contains(s, "„") {
			s = replacePunctInMatches(s, reDEDoubleQuote, false)
			s = replacePunctInMatches(s, reDESplitQuote, false)
			return s
		}
		if strings.Contains(s, ",,") {
			return replacePunctInMatches(s, reDEUnconvQuote, false)
		}
	}
	return replacePunctInMatches(s, reBetweenDoubleQuotes, false)
}

func replacePunctInMatches(s string, re *regexp.Regexp, single bool) string {
	locs := re.FindAllStringIndex(s, -1)
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
		b.WriteString(s[pos:loc[0]])
		b.WriteString(hidePunct(s[loc[0]:loc[1]], single))
		pos = loc[1]
	}
	b.WriteString(s[pos:])
	return b.String()
}

func hidePunct(m string, single bool) string {
	r := strings.NewReplacer(
		".", subPeriod,
		"。", subFullWidthPeriod,
		"．", subSpecialPeriod,
		"！", subFullWidthExcl,
		"!", subExcl,
		"?", subQuestion,
		"？", subFullWidthQuestion,
	)
	out := r.Replace(m)
	if !single {
		out = strings.ReplaceAll(out, "'", subSingleQuote)
	}
	return out
}

func applyExclamationWords(s string) string {
	return replacePunctInMatches(s, reExclWords, false)
}
