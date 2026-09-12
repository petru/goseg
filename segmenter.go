// Package goseg is a rule-based sentence boundary detector ported from the
// Ruby pragmatic_segmenter gem. It works out of the box across many languages
// and does not require training data.
package goseg

import "goseg/internal/seg"

// An Option configures Segment or Clean.
type Option func(*config)

type config struct {
	language string
	docType  string
	clean    bool
	workers  int
}

func defaultConfig() config {
	return config{
		language: "en",
		clean:    true,
	}
}

// Language sets the ISO 639-1 language code (default "en"). Unknown codes
// fall back to the common English-like rules.
func Language(code string) Option {
	return func(c *config) {
		if code == "" {
			c.language = "en"
			return
		}
		c.language = code
	}
}

// DocType enables document-type-specific cleaning. The Ruby gem documents
// "pdf"; "docx" is also accepted by some of its tests.
func DocType(t string) Option {
	return func(c *config) { c.docType = t }
}

// CleanText enables or disables the cleaning/preprocessing pass (default true).
func CleanText(clean bool) Option {
	return func(c *config) { c.clean = clean }
}

// Workers sets how many goroutines Segment may use. 0 (the default) uses
// GOMAXPROCS when the text is large enough to split on blank lines. 1
// forces a single sequential pass.
func Workers(n int) Option {
	return func(c *config) { c.workers = n }
}

// Segment splits text into sentences.
func Segment(text string, opts ...Option) []string {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return seg.Segment(text, cfg.language, cfg.docType, cfg.clean, cfg.workers)
}

// Clean runs only the cleaning/preprocessing pass and returns the resulting
// text. Empty input returns an empty string.
func Clean(text string, opts ...Option) string {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return seg.Clean(text, cfg.language, cfg.docType)
}
