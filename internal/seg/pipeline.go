// Package seg implements the pragmatic_segmenter rule pipeline.
// Importers should use package goseg.
package seg

import "runtime"

// Segment splits text into sentences for the given ISO 639-1 language code.
// Unknown codes use the common English-like rules. If clean is true, the
// preprocessing pass runs first. docType selects document-specific cleaning
// (for example "pdf").
//
// workers is the maximum number of goroutines. 0 means GOMAXPROCS when the
// input is large enough to split at blank lines; 1 forces a single pass.
func Segment(text, language, docType string, clean bool, workers int) []string {
	if text == "" {
		return []string{}
	}
	lang := newLanguage(language)
	n := workers
	if n < 0 {
		n = 1
	}
	if n == 0 {
		n = runtime.GOMAXPROCS(0)
	}
	chunks := []string{text}
	if n > 1 && (workers > 1 || len(text) >= parallelMinBytes) {
		chunks = splitChunks(text, n)
	}
	if len(chunks) == 1 {
		return segmentChunk(chunks[0], lang, docType, clean)
	}
	return mapChunks(chunks, n, func(c string) []string {
		return segmentChunk(c, lang, docType, clean)
	})
}

func segmentChunk(text string, lang *language, docType string, clean bool) []string {
	if clean {
		text = cleanText(text, docType, lang)
		if text == "" {
			return []string{}
		}
	}
	return process(text, lang)
}

// Clean runs only the cleaning/preprocessing pass.
func Clean(text, language, docType string) string {
	if text == "" {
		return ""
	}
	return cleanText(text, docType, newLanguage(language))
}
