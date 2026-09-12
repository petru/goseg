package seg

import (
	"regexp"
	"strings"
	"sync"
)

var langCache sync.Map // string -> *language

type compiledAbbr struct {
	stripped string
	search   *regexp.Regexp
	nextChar *regexp.Regexp
	hide     *regexp.Regexp // boundary + abbr + captured period
	pre      bool
	num      bool
	ascii    bool
}

type abbrMode int

const (
	abbrDefault abbrMode = iota
	abbrAlwaysBeforeSpace
	abbrAlwaysPeriod
	abbrPeriodAnywhere
)

type language struct {
	code string

	punctuations []string

	sentenceStarters []string
	months           []string

	abbreviations []string
	abbrSet       map[string]struct{}
	prepositive   []string
	preSet        map[string]struct{}
	numberAbbr    []string
	numSet        map[string]struct{}

	cleanerSkipAbbr        bool
	replaceBackticks       bool
	japaneseWordNewline    bool
	germanQuotes           bool
	japaneseQuotes         bool
	chineseQuotes          bool
	simpleBoundary         bool
	boundaryPuncts         []rune
	mode                   abbrMode
	colonBetweenNumbers    bool
	nonBoundaryComma       bool
	kazakhExtras           bool
	danishAmPm             bool
	germanLowerLetterAbbr  bool
	cyrillicSingleLetter   bool
	extraNumberPeriodSpace bool
	multiPeriodKazakh      bool

	compiled    []compiledAbbr
	starterRes  []*regexp.Regexp
	monthAheads []*regexp.Regexp
}

func newLanguage(code string) *language {
	if code == "" {
		code = "en"
	}
	if v, ok := langCache.Load(code); ok {
		return v.(*language)
	}
	lang := buildLanguage(code)
	actual, _ := langCache.LoadOrStore(code, lang)
	return actual.(*language)
}

func buildLanguage(code string) *language {
	lang := baseLanguage(code)
	abbr, pre, num := rawAbbr(code)
	if inheritCommonAbbr(code) {
		abbr, pre, num = rawAbbr("common")
	}
	lang.abbreviations = abbr
	lang.prepositive = pre
	lang.numberAbbr = num
	lang.abbrSet = toSet(abbr)
	lang.preSet = toSet(pre)
	lang.numSet = toSet(num)
	lang.compile()
	return lang
}

func (l *language) compile() {
	l.compiled = make([]compiledAbbr, 0, len(l.abbreviations))
	for _, a := range l.abbreviations {
		stripped := strings.TrimSpace(a)
		if stripped == "" {
			continue
		}
		esc := regexp.QuoteMeta(stripped)
		hide, err := regexp.Compile(`(?i)(?:^|[\s\r\n])` + stripped + `(\.)`)
		if err != nil {
			hide = regexp.MustCompile(`(?i)(?:^|[\s\r\n])` + esc + `(\.)`)
		}
		l.compiled = append(l.compiled, compiledAbbr{
			stripped: stripped,
			search:   regexp.MustCompile(`(?i)(?:^|[\s\r\n])` + esc),
			nextChar: regexp.MustCompile(`(?s)` + esc + ` (.)`),
			hide:     hide,
			pre:      inSet(l.preSet, stripped),
			num:      inSet(l.numSet, stripped),
			ascii:    isASCII(stripped) && !hasRegexMeta(stripped),
		})
	}
	if len(l.sentenceStarters) > 0 && !l.danishAmPm {
		l.starterRes = make([]*regexp.Regexp, len(l.sentenceStarters))
		for i, word := range l.sentenceStarters {
			esc := regexp.QuoteMeta(word)
			l.starterRes[i] = regexp.MustCompile(`(U∯S|U\.S|U∯K|E∯U|E\.U|U∯S∯A|U\.S\.A|I|i.v|I.V)∯(\s` + esc + `\s)`)
		}
	}
	for _, month := range l.months {
		l.monthAheads = append(l.monthAheads, regexp.MustCompile(`^\s*`+regexp.QuoteMeta(month)))
	}
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func hasRegexMeta(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.', '+', '*', '?', '(', ')', '[', ']', '{', '}', '|', '\\', '^', '$':
			return true
		}
	}
	return false
}

func inheritCommonAbbr(code string) bool {
	switch code {
	case "en", "ja", "el", "am", "hi", "hy", "fa", "my", "ur", "zh", "common":
		return true
	default:
		return false
	}
}

func toSet(xs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		m[x] = struct{}{}
	}
	return m
}

func inSet(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}

func baseLanguage(code string) *language {
	enStarters := []string{
		"A", "Being", "Did", "For", "He", "How", "However", "I", "In", "It",
		"Millions", "More", "She", "That", "The", "There", "They", "We",
		"What", "When", "Where", "Who", "Why",
	}
	commonPunct := []string{"。", "．", ".", "！", "!", "?", "？"}

	lang := &language{
		code:             code,
		punctuations:     commonPunct,
		sentenceStarters: enStarters,
	}

	switch code {
	case "en":
		lang.cleanerSkipAbbr = true
		lang.replaceBackticks = true
	case "da":
		lang.cleanerSkipAbbr = true
		lang.replaceBackticks = true
		lang.extraNumberPeriodSpace = true
		lang.danishAmPm = true
		lang.sentenceStarters = []string{
			"At", "De", "Dem", "Den", "Der", "Det", "Du", "En", "Et", "For", "Få",
			"Gjorde", "Han", "Hun", "Hvad", "Hvem", "Hvilke", "Hvor", "Hvordan",
			"Hvorfor", "Hvorledes", "Hvornår", "I", "Jeg", "Mange", "Vi", "Være",
		}
	case "de":
		lang.germanQuotes = true
		lang.germanLowerLetterAbbr = true
		lang.extraNumberPeriodSpace = true
		lang.mode = abbrAlwaysBeforeSpace
		lang.months = []string{
			"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli",
			"August", "September", "Oktober", "November", "Dezember",
		}
		lang.sentenceStarters = []string{
			"Am", "Auch", "Auf", "Bei", "Da", "Das", "Der", "Die", "Ein", "Eine",
			"Es", "Für", "Heute", "Ich", "Im", "In", "Ist", "Jetzt", "Mein", "Mit",
			"Nach", "So", "Und", "Warum", "Was", "Wenn", "Wer", "Wie", "Wir",
		}
	case "el":
		lang.simpleBoundary = true
		lang.punctuations = []string{".", "!", ";", "?"}
		lang.boundaryPuncts = []rune{'.', '!', ';', '?'}
		lang.sentenceStarters = nil
	case "hi":
		lang.simpleBoundary = true
		lang.punctuations = []string{"।", "|", ".", "!", "?"}
		lang.boundaryPuncts = []rune{'।', '|', '.', '!', '?'}
		lang.sentenceStarters = nil
	case "hy":
		lang.simpleBoundary = true
		lang.punctuations = []string{"։", "՜", ":"}
		lang.boundaryPuncts = []rune{'։', '՜', ':'}
		lang.sentenceStarters = nil
	case "my":
		lang.simpleBoundary = true
		lang.punctuations = []string{"။", "၏", "?", "!"}
		lang.boundaryPuncts = []rune{'။', '၏', '?', '!'}
		lang.sentenceStarters = nil
	case "am":
		lang.simpleBoundary = true
		lang.punctuations = []string{"።", "፧", "?", "!"}
		lang.boundaryPuncts = []rune{'።', '፧', '?', '!'}
		lang.sentenceStarters = nil
	case "ur":
		lang.simpleBoundary = true
		lang.punctuations = []string{"?", "!", "۔", "؟"}
		lang.boundaryPuncts = []rune{'?', '!', '۔', '؟'}
		lang.sentenceStarters = nil
	case "fa":
		lang.simpleBoundary = true
		lang.punctuations = []string{"?", "!", ":", ".", "؟"}
		lang.boundaryPuncts = []rune{'?', '!', ':', '.', '؟'}
		lang.colonBetweenNumbers = true
		lang.nonBoundaryComma = true
		lang.mode = abbrAlwaysPeriod
		lang.sentenceStarters = nil
	case "ar":
		lang.simpleBoundary = true
		lang.punctuations = []string{"?", "!", ":", ".", "؟", "،"}
		lang.boundaryPuncts = []rune{'?', '!', ':', '.', '؟', '،'}
		lang.colonBetweenNumbers = true
		lang.nonBoundaryComma = true
		lang.mode = abbrAlwaysPeriod
		lang.sentenceStarters = nil
	case "ja":
		lang.japaneseWordNewline = true
		lang.japaneseQuotes = true
		lang.sentenceStarters = nil
	case "zh":
		lang.chineseQuotes = true
		lang.sentenceStarters = nil
	case "ru":
		lang.mode = abbrPeriodAnywhere
		lang.sentenceStarters = nil
	case "bg":
		lang.mode = abbrPeriodAnywhere
		lang.sentenceStarters = nil
	case "kk":
		lang.kazakhExtras = true
		lang.cyrillicSingleLetter = true
		lang.multiPeriodKazakh = true
		lang.sentenceStarters = nil
	case "fr", "it", "es", "nl", "pl":
		lang.sentenceStarters = nil
	}
	return lang
}

func (l *language) hasPunct(s string) bool {
	for _, p := range l.punctuations {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

func (l *language) endsWithPunct(s string) bool {
	if s == "" {
		return false
	}
	for _, p := range l.punctuations {
		if len(s) >= len(p) && s[len(s)-len(p):] == p {
			return true
		}
	}
	return false
}
