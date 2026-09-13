package seg

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func betweenPunctuation(s string, lang *language) string {
	s = subBetweenSingleQuotes(s)
	s = replacePunctInMatches(s, reBetweenSlantedSingle, false, true)
	s = subBetweenDoubleQuotes(s, lang)
	s = replacePunctInMatches(s, reBetweenSquare, false, false)
	s = replacePunctInMatches(s, reBetweenParens, false, false)
	s = replacePunctInMatches(s, reBetweenQuoteArrow, false, true)
	s = replacePunctInMatches(s, reBetweenEmDashes, false, false)
	s = replacePunctInMatches(s, reBetweenQuoteSlanted, false, true)
	if lang.japaneseQuotes {
		s = replacePunctInMatches(s, reBetweenParensJA, false, false)
		s = replacePunctInMatches(s, reBetweenQuoteJA, false, false)
	}
	if lang.chineseQuotes {
		s = replacePunctInMatches(s, reBetweenAngleZH, false, false)
		s = replacePunctInMatches(s, reBetweenLBracket, false, false)
	}
	return s
}

func subBetweenSingleQuotes(s string) string {
	if reLeadingApostrophe.MatchString(s) && !reQuoteSpace.MatchString(s) {
		return s
	}
	return replacePunctInMatches(s, reBetweenSingleQuotes, true, true)
}

func subBetweenDoubleQuotes(s string, lang *language) string {
	if lang.germanQuotes {
		if strings.Contains(s, "„") {
			s = replacePunctInMatches(s, reDEDoubleQuote, false, true)
			s = replacePunctInMatches(s, reDESplitQuote, false, true)
			return s
		}
		if strings.Contains(s, ",,") {
			return replacePunctInMatches(s, reDEUnconvQuote, false, true)
		}
	}
	return replacePunctInMatches(s, reBetweenDoubleQuotes, false, true)
}

// replacePunctInMatches hides non-boundary punctuation inside spans matched by
// re. If keepInner is set (dialogue quotes), a . ! or ? that is followed by a
// new sentence stays visible so the quote can be segmented internally.
func replacePunctInMatches(s string, re *regexp.Regexp, single, keepInner bool) string {
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
		b.WriteString(hidePunct(s[loc[0]:loc[1]], single, keepInner))
		pos = loc[1]
	}
	b.WriteString(s[pos:])
	return b.String()
}

func hidePunct(m string, single, keepInner bool) string {
	if !keepInner {
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
	var b strings.Builder
	b.Grow(len(m) + 8)
	for i := 0; i < len(m); {
		r, w := utf8.DecodeRuneInString(m[i:])
		if hid, ok := hiddenSentencePunct(r); ok {
			if followedByUpper(m[i+w:]) {
				b.WriteString(m[i : i+w])
			} else {
				b.WriteString(hid)
			}
			i += w
			continue
		}
		if r == '\'' && !single {
			b.WriteString(subSingleQuote)
			i += w
			continue
		}
		b.WriteString(m[i : i+w])
		i += w
	}
	return b.String()
}

func hiddenSentencePunct(r rune) (string, bool) {
	switch r {
	case '.':
		return subPeriod, true
	case '。':
		return subFullWidthPeriod, true
	case '．':
		return subSpecialPeriod, true
	case '！':
		return subFullWidthExcl, true
	case '!':
		return subExcl, true
	case '?':
		return subQuestion, true
	case '？':
		return subFullWidthQuestion, true
	default:
		return "", false
	}
}

func isSentenceEndPunct(r rune) bool {
	_, ok := hiddenSentencePunct(r)
	return ok
}

func followedByUpper(s string) bool {
	skipped := false
	for i := 0; i < len(s); {
		r, w := utf8.DecodeRuneInString(s[i:])
		if unicode.IsLetter(r) {
			return skipped && unicode.IsUpper(r)
		}
		skipped = true
		i += w
	}
	return false
}

func hasInnerSentenceBoundary(s string) bool {
	for i := 0; i < len(s); {
		r, w := utf8.DecodeRuneInString(s[i:])
		if isSentenceEndPunct(r) && followedByUpper(s[i+w:]) {
			return true
		}
		i += w
	}
	return false
}

func applyExclamationWords(s string) string {
	return replacePunctInMatches(s, reExclWords, false, false)
}
