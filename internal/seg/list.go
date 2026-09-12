package seg

import (
	"regexp"
	"strconv"
	"strings"
)

var romanNumerals = []string{
	"i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix", "x",
	"xi", "xii", "xiii", "xiv", "x", "xi", "xii", "xiii", "xv",
	"xvi", "xvii", "xviii", "xix", "xx",
}

func latinNumerals() []string {
	out := make([]string, 26)
	for i := 0; i < 26; i++ {
		out[i] = string(rune('a' + i))
	}
	return out
}

func addListLineBreaks(s string) string {
	s = formatAlphabeticalLists(s)
	s = formatRomanNumeralLists(s)
	s = formatNumberedListWithPeriods(s)
	s = formatNumberedListWithParens(s)
	return s
}

func replaceListParens(s string) string {
	return reRomanParens.ReplaceAllString(s, subLeftParen+"$1"+subRightParen+"$2")
}

func formatNumberedListWithParens(s string) string {
	s = replaceParensInNumberedList(s)
	s = addLineBreaksForNumberedListWithParens(s)
	s = strings.ReplaceAll(s, listParen, "")
	return s
}

func formatNumberedListWithPeriods(s string) string {
	s = replacePeriodsInNumberedList(s)
	s = addLineBreaksForNumberedListWithPeriods(s)
	s = strings.ReplaceAll(s, listPeriod, subPeriod)
	return s
}

func formatAlphabeticalLists(s string) string {
	s = addLineBreaksForAlphabeticalListWithPeriods(s, false)
	s = addLineBreaksForAlphabeticalListWithParens(s, false)
	return s
}

func formatRomanNumeralLists(s string) string {
	s = addLineBreaksForAlphabeticalListWithPeriods(s, true)
	s = addLineBreaksForAlphabeticalListWithParens(s, true)
	return s
}

func replacePeriodsInNumberedList(s string) string {
	nums := scanNumberedList(s, false)
	cons := consecutiveValues(nums)
	if len(cons) == 0 {
		return s
	}
	return replaceNumberedListMarker(s, cons, false)
}

func replaceParensInNumberedList(s string) string {
	for i := 0; i < 2; i++ {
		nums := scanNumberedList(s, true)
		cons := consecutiveValues(nums)
		if len(cons) == 0 {
			continue
		}
		s = replaceNumberedListMarker(s, cons, true)
	}
	return s
}

func addLineBreaksForNumberedListWithPeriods(s string) string {
	if strings.Contains(s, listPeriod) && !reListAlreadyNL.MatchString(s) && !reListForFalse.MatchString(s) {
		s = reSpaceList1.ReplaceAllString(s, "$1\r$3")
		s = reSpaceList1b.ReplaceAllString(s, "\r$2")
		s = reSpaceList2.ReplaceAllString(s, "$1\r$3")
		s = reSpaceList2b.ReplaceAllString(s, "\r$2")
	}
	return s
}

func addLineBreaksForNumberedListWithParens(s string) string {
	if strings.Contains(s, listParen) && !reListParenNL.MatchString(s) {
		s = reSpaceList3.ReplaceAllString(s, "$1\r$3")
		s = reSpaceList3b.ReplaceAllString(s, "\r$2")
	}
	return s
}

func consecutiveValues(nums []int) map[int]bool {
	m := map[int]bool{}
	if len(nums) == 0 {
		return m
	}
	for i, a := range nums {
		var prevPtr, nextPtr *int
		if i > 0 {
			p := nums[i-1]
			prevPtr = &p
		} else {
			p := nums[len(nums)-1]
			prevPtr = &p
		}
		if i+1 < len(nums) {
			n := nums[i+1]
			nextPtr = &n
		}
		ok := false
		if nextPtr != nil && a+1 == *nextPtr {
			ok = true
		}
		if prevPtr != nil && a-1 == *prevPtr {
			ok = true
		}
		if a == 0 && prevPtr != nil && *prevPtr == 9 {
			ok = true
		}
		if a == 9 && nextPtr != nil && *nextPtr == 0 {
			ok = true
		}
		if ok {
			m[a] = true
		}
	}
	return m
}

// scanNumberedList extracts 1-2 digit list markers matching the Ruby gem's
// numbered-list regexes. parens=true uses "N)" markers; otherwise "N." / "N.)".
func scanNumberedList(s string, parens bool) []int {
	var out []int
	n := len(s)
	for i := 0; i < n; i++ {
		if parens {
			if v, ok, adv := matchParenNumber(s, i); ok {
				out = append(out, v)
				i += adv - 1
			}
			continue
		}
		if v, ok, adv := matchPeriodNumber(s, i); ok {
			out = append(out, v)
			i += adv - 1
		}
	}
	return out
}

func matchPeriodNumber(s string, i int) (int, bool, int) {
	n := len(s)
	if i >= n || s[i] < '0' || s[i] > '9' {
		return 0, false, 0
	}
	okPrefix := false
	if i == 0 || isSpaceByte(s[i-1]) {
		okPrefix = true
	} else if i > 0 && s[i-1] == '-' {
		pre := i - 1
		if pre == 0 || isSpaceByte(s[pre-1]) || s[pre-1] == 's' {
			okPrefix = true
		}
	} else if i >= len("⁃") && s[i-len("⁃"):i] == "⁃" {
		pre := i - len("⁃")
		if pre == 0 || isSpaceByte(s[pre-1]) {
			okPrefix = true
		}
	}
	if !okPrefix {
		return 0, false, 0
	}
	j := i
	if j+1 < n && s[j+1] >= '0' && s[j+1] <= '9' {
		j++
	}
	num := s[i : j+1]
	k := j + 1
	if k >= n || s[k] != '.' {
		return 0, false, 0
	}
	k++
	if k >= n {
		return 0, false, 0
	}
	if s[k] != ' ' && s[k] != ')' {
		return 0, false, 0
	}
	v, err := strconv.Atoi(num)
	if err != nil {
		return 0, false, 0
	}
	return v, true, j + 1 - i
}

func matchParenNumber(s string, i int) (int, bool, int) {
	n := len(s)
	if i >= n || s[i] < '0' || s[i] > '9' {
		return 0, false, 0
	}
	j := i
	if j+1 < n && s[j+1] >= '0' && s[j+1] <= '9' {
		j++
	}
	k := j + 1
	if k+1 >= n || s[k] != ')' || s[k+1] != ' ' {
		return 0, false, 0
	}
	v, err := strconv.Atoi(s[i : j+1])
	if err != nil {
		return 0, false, 0
	}
	return v, true, (j + 1) - i
}

func isSpaceByte(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f'
}

func replaceNumberedListMarker(s string, cons map[int]bool, parens bool) string {
	var b strings.Builder
	b.Grow(len(s))
	n := len(s)
	for i := 0; i < n; {
		if parens {
			if v, ok, adv := matchParenNumber(s, i); ok && cons[v] {
				b.WriteString(strconv.Itoa(v))
				b.WriteString(listParen)
				i += adv
				continue
			}
		} else {
			if v, ok, _ := matchPeriodNumber(s, i); ok && cons[v] {
				// Replace the digits+period (REGEX_2 includes the period).
				num := strconv.Itoa(v)
				if strings.HasPrefix(s[i:], num+".") {
					b.WriteString(num)
					b.WriteString(listPeriod)
					i += len(num) + 1
					continue
				}
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

func addLineBreaksForAlphabeticalListWithPeriods(s string, roman bool) string {
	return iterateAlphabetArray(s, false, roman)
}

func addLineBreaksForAlphabeticalListWithParens(s string, roman bool) string {
	return iterateAlphabetArray(s, true, roman)
}

func iterateAlphabetArray(s string, parens, roman bool) string {
	var alphabet []string
	if roman {
		alphabet = romanNumerals
	} else {
		alphabet = latinNumerals()
	}
	alphaSet := toSet(alphabet)
	list := scanAlphabetItems(s, parens)
	filtered := make([]string, 0, len(list))
	for _, item := range list {
		if roman {
			ok := false
			for _, a := range alphabet {
				if strings.Contains(a, item) {
					ok = true
					break
				}
			}
			if ok {
				filtered = append(filtered, item)
			}
		} else if inSet(alphaSet, item) {
			filtered = append(filtered, item)
		}
	}
	idx := indexMap(alphabet)
	for i, a := range filtered {
		if i == len(filtered)-1 {
			if shouldReplaceLast(a, i, idx, filtered) {
				s = replaceCorrectAlphabetList(s, a, parens)
			}
		} else if shouldReplaceOther(a, i, idx, filtered) {
			s = replaceCorrectAlphabetList(s, a, parens)
		}
	}
	return s
}

func indexMap(alphabet []string) map[string]int {
	m := make(map[string]int, len(alphabet))
	for i, a := range alphabet {
		if _, ok := m[a]; !ok {
			m[a] = i
		}
	}
	return m
}

func shouldReplaceLast(a string, i int, idx map[string]int, list []string) bool {
	if i == 0 {
		return false
	}
	prev := list[i-1]
	ia, okA := idx[a]
	ip, okP := idx[prev]
	if !okA || !okP {
		return false
	}
	d := ia - ip
	if d < 0 {
		d = -d
	}
	return d == 1
}

func shouldReplaceOther(a string, i int, idx map[string]int, list []string) bool {
	if i != 0 && i+1 >= len(list) {
		return false
	}
	var prev string
	if i > 0 {
		prev = list[i-1]
	} else if len(list) > 0 {
		prev = list[len(list)-1]
	}
	next := list[i+1]
	ia, okA := idx[a]
	ip, okP := idx[prev]
	in, okN := idx[next]
	if !okA || !okP || !okN {
		return false
	}
	dPrev := ia - ip
	if dPrev < 0 {
		dPrev = -dPrev
	}
	if in-ia != 1 && dPrev != 1 {
		return false
	}
	return true
}

func scanAlphabetItems(s string, parens bool) []string {
	var out []string
	if parens {
		ms := reAlphaListParens.FindAllStringSubmatch(s, -1)
		for _, m := range ms {
			item := m[1]
			if item == "" {
				item = m[2]
			}
			out = append(out, strings.ToLower(item))
		}
		return out
	}
	ms := reAlphaListPeriod.FindAllStringSubmatch(s, -1)
	for _, m := range ms {
		out = append(out, strings.ToLower(m[1]))
	}
	return out
}

func replaceCorrectAlphabetList(s, a string, parens bool) string {
	if parens {
		return replaceAlphabetListParens(s, a)
	}
	return replaceAlphabetList(s, a)
}

func replaceAlphabetList(s, a string) string {
	re, err := regexp.Compile(`(?:^|([\s]))(` + regexp.QuoteMeta(a) + `)\.`)
	if err != nil {
		return s
	}
	return re.ReplaceAllStringFunc(s, func(m string) string {
		lead := ""
		rest := m
		if len(m) > 0 && isSpaceByte(m[0]) {
			lead = m[:1]
			rest = m[1:]
		}
		if strings.ToLower(strings.TrimSuffix(rest, ".")) != a {
			return m
		}
		return lead + "\r" + a + subPeriod
	})
}

func replaceAlphabetListParens(s, a string) string {
	// Match the full parenthetical letter group, not a substring of a
	// longer marker (so "i" does not fire inside "(ii)").
	return reAlphaListReplaceParens.ReplaceAllStringFunc(s, func(m string) string {
		if strings.HasPrefix(m, "(") {
			inner := strings.ToLower(m[1 : len(m)-1])
			if inner != a {
				return m
			}
			return "\r" + subLeftParen + m[1:]
		}
		lead := ""
		body := m
		if len(m) > 0 && isSpaceByte(m[0]) {
			lead = m[:1]
			body = m[1:]
		}
		inner := strings.ToLower(strings.TrimSuffix(body, ")"))
		if inner != a {
			return m
		}
		return lead + "\r" + body
	})
}
