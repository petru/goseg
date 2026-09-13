package seg

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func process(text string, lang *language) []string {
	s := addListLineBreaks(text)
	s = replaceAbbreviations(s, lang)
	s = replaceNumbers(s, lang)
	s = replaceContinuousPunctuation(s)
	s = replacePeriodsBeforeNumericReferences(s)
	s = reWordDotWord.ReplaceAllString(s, "${1}"+subMultiPeriod+"${3}")
	s = replaceLiteralFind(s, reGeoBehind, ".", reGeoAhead, subPeriod)
	s = replaceLiteralFind(s, reSpaceEnd, ".", reFileExtAhead, subPeriod)
	return splitIntoSegments(s, lang)
}

func replaceNumbers(s string, lang *language) string {
	s = replaceLiteralFind(s, nil, ".", reDigitStart, subPeriod)
	s = replaceLiteralFind(s, reDigitEnd, ".", reNonSpaceStart, subPeriod)
	s = replaceLiteralFind(s, reCRDigitEnd, ".", reSpaceSStart, subPeriod)
	s = replaceLiteralFind(s, reStartOneDigit, ".", reSpaceSStart, subPeriod)
	s = replaceLiteralFind(s, reStartTwoDigits, ".", reSpaceSStart, subPeriod)
	if lang.extraNumberPeriodSpace {
		s = reDENumberPeriod.ReplaceAllString(s, "${1}"+subPeriod+"${3}")
		s = reDENegNumber.ReplaceAllString(s, "${1}"+subPeriod+"${3}")
	}
	for _, ahead := range lang.monthAheads {
		s = replaceLiteralFind(s, reDigitEnd, ".", ahead, subPeriod)
	}
	return s
}

func replaceContinuousPunctuation(s string) string {
	return replaceFindFunc(s, reNonSpaceEnd, reContinuousPunct, reSpaceOrEnd, func(m string) string {
		m = strings.ReplaceAll(m, "!", subExcl)
		m = strings.ReplaceAll(m, "?", subQuestion)
		return m
	})
}

func replacePeriodsBeforeNumericReferences(s string) string {
	return replaceFindFunc(s, reNonDigitSpaceEnd, reNumberedRefMid, reUpperStart, func(m string) string {
		// m is (.|∯)(ref)(space); replacement is ∯ + ref + \r + space
		if m == "" {
			return m
		}
		rest := m
		if rest[0] == '.' || strings.HasPrefix(rest, subPeriod) {
			if strings.HasPrefix(rest, subPeriod) {
				rest = rest[len(subPeriod):]
			} else {
				rest = rest[1:]
			}
		}
		if len(rest) == 0 {
			return subPeriod + "\r"
		}
		space := ""
		body := rest
		if rest[len(rest)-1] == ' ' {
			space = " "
			body = rest[:len(rest)-1]
		}
		return subPeriod + body + "\r" + space
	})
}

func splitIntoSegments(s string, lang *language) []string {
	s = checkParensBetweenQuotes(s)
	parts := strings.Split(s, "\r")
	var segs []string
	for _, part := range parts {
		part = applyRules(part, rule{reSingleNL, subNewline})
		part = applyEllipsisRules(part)
		for _, p := range checkForPunctuation(part, lang) {
			p = applyRules(p, subSymbolRules...)
			for _, q := range postProcessSegments(p) {
				if q == "" {
					continue
				}
				q = strings.ReplaceAll(q, subSingleQuote, "'")
				segs = append(segs, q)
			}
		}
	}
	if segs == nil {
		return []string{}
	}
	return segs
}

func applyEllipsisRules(s string) string {
	s = reEllipsis3space.ReplaceAllString(s, ellipsis3space)
	s = reEllipsis4space.ReplaceAllString(s, "${1}"+ellipsis4space)
	s = replaceFind(s, reFourDotBehind, reThreeDotFind, reFourDotAhead, ellipsis3)
	s = replaceFind(s, nil, reThreeDotFind, reThreeDotAhead, ellipsis2+".")
	s = strings.ReplaceAll(s, "...", ellipsis3)
	return s
}

func checkParensBetweenQuotes(s string) string {
	return reParensBetweenQuotes.ReplaceAllStringFunc(s, func(m string) string {
		m = replaceFind(m, nil, mustRE(`\s`), mustRE(`^\(`), "\r")
		m = replaceFind(m, mustRE(`\)\z`), mustRE(`\s`), nil, "\r")
		return m
	})
}

func checkForPunctuation(s string, lang *language) []string {
	if lang.hasPunct(s) {
		return processText(s, lang)
	}
	return []string{s}
}

func processText(s string, lang *language) []string {
	if !lang.endsWithPunct(s) {
		s += tmpEnding
	}
	s = applyExclamationWords(s)
	s = betweenPunctuation(s, lang)
	s = applyRules(s, doublePunctRules...)
	s = reQInQuote.ReplaceAllString(s, subQuestion+"$1")
	s = reEInQuote.ReplaceAllString(s, subExcl+"$1")
	s = reEComma.ReplaceAllString(s, subExcl+"$1")
	s = reEMid.ReplaceAllString(s, subExcl+"$1")
	s = replaceListParens(s)
	if lang.kazakhExtras {
		s = reKazakhQDash.ReplaceAllString(s, "${1}"+subQuestion+"${2}")
		s = reKazakhEDash.ReplaceAllString(s, "${1}"+subExcl+"${2}")
	}
	return sentenceBoundaryPunctuation(s, lang)
}

func sentenceBoundaryPunctuation(s string, lang *language) []string {
	if lang.colonBetweenNumbers {
		s = reColonBetween.ReplaceAllString(s, "${1}"+subColon+"${2}")
	}
	if lang.nonBoundaryComma {
		s = reArabicCommaNSB.ReplaceAllString(s, subArabicComma+"$1")
	}
	if lang.simpleBoundary {
		return scanSimpleBoundary(s, lang.boundaryPuncts)
	}
	return scanCommonBoundary(s)
}

func scanSimpleBoundary(s string, puncts []rune) []string {
	punctSet := make(map[rune]bool, len(puncts))
	for _, p := range puncts {
		punctSet[p] = true
	}
	var out []string
	start := 0
	i := 0
	for i < len(s) {
		r, w := utf8.DecodeRuneInString(s[i:])
		if r == '\n' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + w
			i += w
			continue
		}
		if punctSet[r] {
			out = append(out, s[start:i+w])
			start = i + w
			i += w
			continue
		}
		i += w
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func scanCommonBoundary(s string) []string {
	var out []string
	pos := 0
	n := len(s)
	for pos < n {
		if m, adv := matchSpecialBoundary(s, pos); adv > 0 {
			out = append(out, m)
			pos += adv
			continue
		}
		rest := s[pos:]
		loc := reCommonBoundary.FindStringIndex(rest)
		if loc == nil {
			break
		}
		out = append(out, rest[loc[0]:loc[1]])
		pos += loc[1]
	}
	return out
}

func matchSpecialBoundary(s string, pos int) (string, int) {
	rest := s[pos:]
	try := func(re *regexp.Regexp, ahead *regexp.Regexp) (string, int) {
		loc := re.FindStringIndex(rest)
		if loc == nil || loc[0] != 0 {
			return "", 0
		}
		after := rest[loc[1]:]
		if !matchesAtStart(after, ahead) {
			return "", 0
		}
		return rest[loc[0]:loc[1]], loc[1]
	}
	if m, adv := try(reFF08, reFF08Ahead); adv > 0 {
		return m, adv
	}
	if m, adv := try(re300c, re300cAhead); adv > 0 {
		return m, adv
	}
	if m, adv := try(reParenSent, reParenSentAhead); adv > 0 {
		return m, adv
	}
	if m, adv := try(reSQuoteSent, reQuoteSentAhead); adv > 0 && !hasInnerSentenceBoundary(m) {
		return m, adv
	}
	if m, adv := try(reDQuoteSent, reQuoteSentAhead); adv > 0 && !hasInnerSentenceBoundary(m) {
		return m, adv
	}
	if m, adv := try(reSmartQuoteSent, reQuoteSentAhead); adv > 0 && !hasInnerSentenceBoundary(m) {
		return m, adv
	}
	return "", 0
}

func postProcessSegments(txt string) []string {
	n := utf8.RuneCountInString(txt)
	if n < 2 && reLetterOnly.MatchString(txt) {
		return []string{txt}
	}
	if consecutiveUnderscore(txt) || n < 2 {
		return nil
	}
	txt = applyRules(txt, reinsertEllipsisRules...)
	txt = reExtraSpace.ReplaceAllString(txt, " ")
	if reQuoteAtEnd.MatchString(txt) {
		return splitQuoteAtEnd(txt)
	}
	txt = strings.ReplaceAll(txt, "\n", "")
	return []string{strings.TrimSpace(txt)}
}

func splitQuoteAtEnd(txt string) []string {
	locs := reSplitQuote.FindAllStringSubmatchIndex(txt, -1)
	if len(locs) == 0 {
		return []string{strings.TrimSpace(strings.ReplaceAll(txt, "\n", ""))}
	}
	var out []string
	last := 0
	for _, loc := range locs {
		// groups: 1 = punct+quote, 2 = next sentence start (capital or opening quote).
		spaceStart := loc[3] // end of group 1
		spaceEnd := loc[4]   // start of group 2
		out = append(out, txt[last:spaceStart])
		last = spaceEnd
	}
	out = append(out, txt[last:])
	for i := range out {
		out[i] = strings.TrimSpace(strings.ReplaceAll(out[i], "\n", ""))
	}
	return out
}

func consecutiveUnderscore(txt string) bool {
	return reUnderscores.ReplaceAllString(txt, "") == ""
}
