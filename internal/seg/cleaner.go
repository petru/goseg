package seg

import (
	"regexp"
	"strings"
)

func cleanText(text, docType string, lang *language) string {
	s := text
	s = removeAllNewlines(s, lang)
	s = applyRules(s, rule{reDoubleNLSpace, "\r"}, rule{reDoubleNL, "\r"})
	s = replaceNewlines(s, docType)
	s = applyRules(s,
		rule{reEscapedNL, "\n"},
		rule{reEscapedCR, "\r"},
		rule{reTypoEscapedNL, "\n"},
		rule{reTypoEscapedCR, "\r"},
	)
	s = applyRules(s, rule{reHTMLTag, ""}, rule{reEscapedHTML, ""})
	s = replacePunctuationInBrackets(s)
	s = applyRules(s, rule{reInlineFmt, ""})
	s = applyRules(s, rule{reQuotationsFirst, `"`}, rule{reQuotationsSecond, `"`})
	s = applyRules(s, rule{reTOC, "\r"}, rule{reConsecutivePeriods, " "}, rule{reConsecutiveSlash, ""})
	s = checkNoSpaceBetweenSentences(s, lang)
	s = applyRules(s, rule{reConsecutivePeriods, " "}, rule{reConsecutiveSlash, ""})
	if lang.replaceBackticks {
		s = strings.ReplaceAll(s, "`", "'")
	}
	return s
}

func removeAllNewlines(s string, lang *language) string {
	s = removeNewlineInMiddleOfSentence(s)
	// Remove a newline sitting between short letter-only fragments
	// (PDF/OCR artifacts like "W\nA\nRN" → "WARN"). Do not consume the
	// following letters so overlapping newlines are all removed.
	s = replaceFind(s, nil, reNL, reShortLetterNL, "")
	if lang.japaneseWordNewline {
		s = replaceFind(s, reJANoEnd, reNL, reNonSpaceStart, "")
	}
	return s
}

func removeNewlineInMiddleOfSentence(s string) string {
	parts := strings.Split(s, ".")
	for i, p := range parts {
		parts[i] = reNLMidSentence.ReplaceAllString(p, "$1$2")
	}
	return strings.Join(parts, ".")
}

func replaceNewlines(s, docType string) string {
	if docType == "pdf" {
		s = reNLFollowedByBullet.ReplaceAllString(s, "\r$1")
		s = rePDFNLMid.ReplaceAllString(s, "$1$2")
		s = rePDFNLNoSpace.ReplaceAllString(s, " $1")
		return s
	}
	s = reNLFollowedByPeriod.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "\n", "\r")
	return s
}

func replacePunctuationInBrackets(s string) string {
	return reBrackets.ReplaceAllStringFunc(s, func(m string) string {
		if strings.Contains(m, "?") {
			return strings.ReplaceAll(m, "?", subQuestion)
		}
		return m
	})
}

func checkNoSpaceBetweenSentences(s string, lang *language) string {
	words := strings.Split(s, " ")
	abbrs := lang.abbreviations
	if lang.cleanerSkipAbbr {
		abbrs = nil
	}
	for _, word := range words {
		s = searchConnected(word, s, reNoSpaceSent, "${1}. ${2}", abbrs)
		s = searchConnected(word, s, reNoSpaceSentDigit, "${1}. ${2}", abbrs)
	}
	return s
}

func searchConnected(word, txt string, re *regexp.Regexp, repl string, abbrs []string) string {
	if !re.MatchString(word) {
		return txt
	}
	if reURLEmail.MatchString(word) {
		return txt
	}
	lower := strings.ToLower(word)
	for _, a := range abbrs {
		if a != "" && strings.Contains(lower, strings.ToLower(a)) {
			return txt
		}
	}
	newWord := re.ReplaceAllString(word, repl)
	if newWord == word {
		return txt
	}
	return strings.ReplaceAll(txt, word, newWord)
}
