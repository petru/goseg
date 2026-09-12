package seg

import "regexp"

// Placeholder characters used to temporarily hide punctuation that is not a
// sentence boundary. They match pragmatic_segmenter so later reverse rules
// stay compatible with the original algorithm.
const (
	subPeriod            = "∯"
	subMultiPeriod       = "∮"
	subArabicComma       = "♬"
	subColon             = "♭"
	subFullWidthPeriod   = "&ᓰ&"
	subSpecialPeriod     = "&ᓱ&"
	subFullWidthExcl     = "&ᓳ&"
	subExcl              = "&ᓴ&"
	subQuestion          = "&ᓷ&"
	subFullWidthQuestion = "&ᓸ&"
	subMixedQE           = "☉"
	subMixedQQ           = "☇"
	subMixedEQ           = "☈"
	subMixedEE           = "☄"
	subLeftParen         = "&✂&"
	subRightParen        = "&⌬&"
	subSingleQuote       = "&⎋&"
	tmpEnding            = "ȸ"
	subNewline           = "ȹ"
	ellipsis3            = "ƪ"
	ellipsis3space       = "♟"
	ellipsis4space       = "♝"
	ellipsis2            = "☏"
	listPeriod           = "♨"
	listParen            = "☝"
)

type rule struct {
	re   *regexp.Regexp
	repl string
}

func applyRules(s string, rules ...rule) string {
	for _, r := range rules {
		if r.re == nil {
			continue
		}
		s = r.re.ReplaceAllString(s, r.repl)
	}
	return s
}

func mustRE(p string) *regexp.Regexp {
	return regexp.MustCompile(p)
}

var (
	rePossessiveAbbr = mustRE(`\.('s(?:\s|$))`)
	reKommandit      = mustRE(`(Co)(\.)(\sKG)`)

	reSpaceStart    = mustRE(`^\s`)
	reNL            = mustRE(`\n`)
	reShortLetterNL = mustRE(`^[a-zA-Z]{1,2}\n`)
	reJANoEnd       = mustRE(`の\z`)
	reDigitEnd      = mustRE(`\d\z`)
	reDigitStart    = mustRE(`^\d`)
	reNonSpaceStart = mustRE(`^\S`)
	reSpaceSStart   = mustRE(`^(?:\s\S|\))`)

	reCRDigitEnd     = mustRE(`\r\d\z`)
	reStartOneDigit  = mustRE(`(?:\A|\n)\d\z`)
	reStartTwoDigits = mustRE(`(?:\A|\n)\d\d\z`)

	reWordDotWord = mustRE(`(\w)(\.)(\w)`)

	reGeoBehind = mustRE(`[a-zA-z]°\z`)
	reGeoAhead  = mustRE(`^\s*\d`)

	reFileExtAhead = mustRE(`^(?:jpe?g|png|gif|tiff?|pdf|ps|docx?|xlsx?|svg|bmp|tga|exif|odt|html?|txt|rtf|bat|sxw|xml|zip|exe|msi|blend|wmv|mp[34]|pptx?|flac|rb|cpp|cs|js)\s`)
	reSpaceEnd     = mustRE(`\s\z`)

	reDoubleQE = mustRE(`\?!`)
	reDoubleEQ = mustRE(`!\?`)
	reDoubleQQ = mustRE(`\?\?`)
	reDoubleEE = mustRE(`!!`)

	reQInQuote = mustRE(`\?(['"])`)
	reEInQuote = mustRE(`!(['"])`)
	reEComma   = mustRE(`!(,\s[a-z])`)
	reEMid     = mustRE(`!(\s[a-z])`)

	reEllipsis3space = mustRE(`(?:\s\.){3}\s`)
	reEllipsis4space = mustRE(`([a-z])((?:\.\s){3}\.(?:\n|$))`)
	reFourDotBehind  = mustRE(`\S\z`)
	reFourDotAhead   = mustRE(`^\.\s[A-Z]`)
	reThreeDotAhead  = mustRE(`^\s+[A-Z]`)
	reThreeDotFind   = mustRE(`\.\.\.`)

	reSingleNL = mustRE(`\n`)

	reSubPeriod            = mustRE(`∯`)
	reSubArabicComma       = mustRE(`♬`)
	reSubColon             = mustRE(`♭`)
	reSubFullWidthPeriod   = mustRE(`&ᓰ&`)
	reSubSpecialPeriod     = mustRE(`&ᓱ&`)
	reSubFullWidthExcl     = mustRE(`&ᓳ&`)
	reSubExcl              = mustRE(`&ᓴ&`)
	reSubQuestion          = mustRE(`&ᓷ&`)
	reSubFullWidthQuestion = mustRE(`&ᓸ&`)
	reSubMixedQE           = mustRE(`☉`)
	reSubMixedQQ           = mustRE(`☇`)
	reSubMixedEQ           = mustRE(`☈`)
	reSubMixedEE           = mustRE(`☄`)
	reSubLeftParen         = mustRE(`&✂&`)
	reSubRightParen        = mustRE(`&⌬&`)
	reTmpEnding            = mustRE(`ȸ`)
	reSubNewline           = mustRE(`ȹ`)
	reEllipsis3            = mustRE(`ƪ`)
	reEllipsis3spSym       = mustRE(`♟`)
	reEllipsis4spSym       = mustRE(`♝`)
	reEllipsis2            = mustRE(`☏`)
	reSubMultiPeriod       = mustRE(`∮`)
	reExtraSpace           = mustRE(`\s{3,}`)

	reAmPmUpperP = mustRE(`(P∯M)∯(\s[A-Z])`)
	reAmPmUpperA = mustRE(`(A∯M)∯(\s[A-Z])`)
	reAmPmLowerP = mustRE(`(p∯m)∯(\s[A-Z])`)
	reAmPmLowerA = mustRE(`(a∯m)∯(\s[A-Z])`)

	reMultiPeriodAbbr       = mustRE(`(?i)\b[a-z](?:\.[a-z])+[.]`)
	reKazakhCyrillicPeriods = mustRE(`(?i)(?:^|[^\p{L}\p{N}_])(\p{Cyrillic}(?:\.\s?\p{Cyrillic})+[.])`)
	reKazakhLatinPeriods    = mustRE(`(?i)(?:^|[^\p{L}\p{N}_])([a-z](?:\.[a-z])+[.])`)

	reContinuousPunct = mustRE(`[!?]{3,}`)
	reNonSpaceEnd     = mustRE(`\S\z`)
	reSpaceOrEnd      = mustRE(`^(?:\s|$)`)

	reNumberedRefMid   = mustRE(`(?:\.|∯)(?:(?:\[(?:\d{1,3},?\s?-?\s?)?\b\d{1,3}\])+|(?:\d{1,3}\s?){0,3}\d{1,3})\s`)
	reNonDigitSpaceEnd = mustRE(`[^\d\s]\z`)
	reUpperStart       = mustRE(`^[A-Z]`)

	reParensBetweenQuotes = mustRE(`["”]\s\(.*\)\s["“]`)

	// Sentence-ending punct before a closer, then whitespace, then a new
	// sentence (capital) or another opening quote (consecutive dialogue).
	// Includes en/em dashes so interrupted speech (“life—”) can close a turn.
	reQuoteAtEnd = mustRE(`[!?.\x{2013}\x{2014}-]["'\x{201d}\x{201c}]\s+[A-Z"'\x{201c}\x{2018}\x{201e}]`)
	reSplitQuote = mustRE(`([!?.\x{2013}\x{2014}-]["'\x{201d}\x{201c}])\s+([A-Z"'\x{201c}\x{2018}\x{201e}])`)

	reLetterOnly     = mustRE(`\A[a-zA-Z]*\z`)
	reUnderscores    = mustRE(`_{3,}`)
	reColonBetween   = mustRE(`(\d):(\d)`)
	reArabicCommaNSB = mustRE(`،(\s\S+،)`)

	reDENumberPeriod = mustRE(`(\s(?:[1-9][0-9]|[0-9]))(\.)(\s)`)
	reDENegNumber    = mustRE(`(-(?:[1-9][0-9]|[0-9]))(\.)(\s)`)

	reSingleLowerLetter     = mustRE(`(\s[a-z])(\.)(\s)`)
	reSingleLowerLetterLine = mustRE(`(?:\A|\n)([a-z])(\.)(\s)`)

	reCyrLetterStart = mustRE(`(?:\A|\n)([А-ЯЁ])(\.)(\s)`)
	reCyrLetter      = mustRE(`(\s[А-ЯЁ])(\.)(\s)`)

	reKazakhQDash = mustRE(`(\p{Ll})\?(\s*[-—]\s*\p{Ll})`)
	reKazakhEDash = mustRE(`(\p{Ll})!(\s*[-—]\s*\p{Ll})`)

	reListAlreadyNL = mustRE(`♨.+\n.+♨|♨.+\r.+♨`)
	reListForFalse  = mustRE(`for\s\d{1,2}♨\s[a-z]`)
	reListParenNL   = mustRE(`☝.+\n.+☝|☝.+\r.+☝`)

	reSpaceList1  = mustRE(`(\S\S)(\s)(\S\s*\d{1,2}♨)`)
	reSpaceList1b = mustRE(`^(\s)(\S\s*\d{1,2}♨)`)
	reSpaceList2  = mustRE(`(\S\S)(\s)(\d{1,2}♨)`)
	reSpaceList2b = mustRE(`^(\s)(\d{1,2}♨)`)
	reSpaceList3  = mustRE(`(\S\S)(\s)(\d{1,2}☝)`)
	reSpaceList3b = mustRE(`^(\s)(\d{1,2}☝)`)

	reRomanParens = mustRE(`\(([mdclxvi]+)\)(\s[A-Z])`)

	reBetweenSingleQuotes  = mustRE(`(\s)'(?:[^']|'[a-zA-Z])*'`)
	reBetweenSlantedSingle = mustRE(`(\s)‘(?:[^’]|’[a-zA-Z])*’`)
	reBetweenDoubleQuotes  = mustRE(`"([^"\\]|\\.)*"`)
	reBetweenQuoteArrow    = mustRE(`«([^»\\]|\\.)*»`)
	reBetweenQuoteSlanted  = mustRE(`“([^”\\]|\\.)*”`)
	reBetweenSquare        = mustRE(`\[([^\]\\]|\\.)*\]`)
	reBetweenParens        = mustRE(`\(([^()\\]|\\.)*\)`)
	reLeadingApostrophe    = mustRE(`\s'(?:[^']|'[a-zA-Z])*'\S`)
	reQuoteSpace           = mustRE(`'\s`)
	reBetweenEmDashes      = mustRE(`--[^-]*--`)

	reBetweenQuoteJA  = mustRE("\u300c([^\u300c\u300d\\\\]|\\\\.)*\u300d")
	reBetweenParensJA = mustRE("\uff08([^\uff08\uff09\\\\]|\\\\.)*\uff09")
	reBetweenAngleZH  = mustRE(`《([^》\\]|\\.)*》`)
	reBetweenLBracket = mustRE(`「([^」\\]|\\.)*」`)

	reDEDoubleQuote = mustRE(`„([^“\\]|\\.)*“`)
	reDESplitQuote  = mustRE(`\A„([^“\\]|\\.)*“`)
	reDEUnconvQuote = mustRE(`,,([^“\\]|\\.)*“`)

	reExclWords = mustRE(`!Xũ|!Kung|ǃʼOǃKung|!Xuun|!Kung-Ekoka|ǃHu|ǃKhung|ǃKu|ǃung|ǃXo|ǃXû|ǃXung|ǃXũ|!Xun|Yahoo!|Y!J|Yum!`)

	reCommonBoundary = mustRE(`\S.*?[。．.！!?？ȸȹ☉☈☇☄]`)
	reFF08           = mustRE("\uff08[^\uff09]*\uff09")
	reFF08Ahead      = mustRE(`^\s?[A-Z]`)
	re300c           = mustRE("\u300c[^\u300d]*\u300d")
	re300cAhead      = mustRE(`^\s[A-Z]`)
	reParenSent      = mustRE(`\([^)]{2,}\)`)
	reParenSentAhead = mustRE(`^\s[A-Z]`)
	reSQuoteSent     = mustRE(`'[^']*[^,]'`)
	reDQuoteSent     = mustRE(`"[^"]*[^,]"`)
	reSmartQuoteSent = mustRE(`“[^”]*[^,]”`)
	reQuoteSentAhead = mustRE(`^\s[A-Z]`)

	reDoubleNLSpace      = mustRE(`\n \n`)
	reDoubleNL           = mustRE(`\n\n`)
	reNLFollowedByPeriod = mustRE(`\n(\.(?:\s|\n))`)
	reEscapedNL          = mustRE(`\\n`)
	reEscapedCR          = mustRE(`\\r`)
	reTypoEscapedNL      = mustRE(`\\ n`)
	reTypoEscapedCR      = mustRE(`\\ r`)
	reInlineFmt          = mustRE(`\{b\^&gt;\d*&lt;b\^\}|\{b\^>\d*<b\^\}`)
	reTOC                = mustRE(`\.{5,}\s*\d+-*\d*`)
	reConsecutivePeriods = mustRE(`\.{5,}`)
	reConsecutiveSlash   = mustRE(`/{3}`)
	reNoSpaceSent        = mustRE(`([a-z])\.([A-Z])`)
	reNoSpaceSentDigit   = mustRE(`(\d)\.([A-Z])`)
	reNLMidSentence      = mustRE(`(\s)\n([a-z(])`)
	reNLFollowedByBullet = mustRE(`\n(•)`)
	reQuotationsFirst    = mustRE(`''`)
	reQuotationsSecond   = mustRE("``")
	reHTMLTag            = mustRE(`</?\w+((\s+\w+(\s*=\s*(?:".*?"|'.*?'|[^'">\s]+))?)+\s*|\s*)/?>`)
	reEscapedHTML        = mustRE(`&lt;/?[^gt;]*gt;`)
	rePDFNLMid           = mustRE(`([^\n]\s)\n(\S)`)
	rePDFNLNoSpace       = mustRE(`\n([a-z])`)
	reBrackets           = mustRE(`\[(?:[^\]])*\]`)
	reURLEmail           = mustRE(`@|http|\.com|net|www|//`)

	reAlphaListParens        = mustRE(`(?i)(?:\(([a-z]+)\))|(?:(?:^|\s)([a-z]+)\))`)
	reAlphaListPeriod        = mustRE(`(?:^|\s)([a-z])\.`)
	reAlphaListReplaceParens = mustRE(`(?i)(\()([a-z]+)(\))|(?:^|(\s))([a-z]+)(\))`)
)

var subSymbolRules = []rule{
	{reSubPeriod, "."},
	{reSubArabicComma, "،"},
	{reSubColon, ":"},
	{reSubFullWidthPeriod, "。"},
	{reSubSpecialPeriod, "．"},
	{reSubFullWidthExcl, "！"},
	{reSubExcl, "!"},
	{reSubQuestion, "?"},
	{reSubFullWidthQuestion, "？"},
	{reSubMixedQE, "?!"},
	{reSubMixedQQ, "??"},
	{reSubMixedEQ, "!?"},
	{reSubMixedEE, "!!"},
	{reSubLeftParen, "("},
	{reSubRightParen, ")"},
	{reTmpEnding, ""},
	{reSubNewline, "\n"},
}

var reinsertEllipsisRules = []rule{
	{reEllipsis3, "..."},
	{reEllipsis3spSym, " . . . "},
	{reEllipsis4spSym, ". . . ."},
	{reEllipsis2, ".."},
	{reSubMultiPeriod, "."},
}

var doublePunctRules = []rule{
	{reDoubleQE, subMixedQE},
	{reDoubleEQ, subMixedEQ},
	{reDoubleQQ, subMixedQQ},
	{reDoubleEE, subMixedEE},
}
