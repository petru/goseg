package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"goseg"
)

func main() {
	lang := flag.String("lang", "en", "ISO 639-1 language code (default en)")
	docType := flag.String("doc-type", "", "document type (e.g. pdf)")
	noClean := flag.Bool("no-clean", false, "disable text cleaning/preprocessing")
	cleanOnly := flag.Bool("clean-only", false, "print cleaned text without segmenting")
	workers := flag.Int("workers", 0, "goroutines for segmentation (0 = GOMAXPROCS, 1 = sequential)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: goseg [options] [file]\n\n")
		fmt.Fprintf(os.Stderr, "Rule-based sentence boundary detection. Reads stdin if no file is given.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	text, err := readInput(flag.Arg(0))
	if err != nil {
		fatal(err)
	}

	opts := []goseg.Option{goseg.Language(*lang), goseg.CleanText(!*noClean), goseg.Workers(*workers)}
	if *docType != "" {
		opts = append(opts, goseg.DocType(*docType))
	}

	if *cleanOnly {
		out := goseg.Clean(text, opts...)
		if !strings.HasSuffix(out, "\n") {
			out += "\n"
		}
		if _, err := io.WriteString(os.Stdout, out); err != nil {
			fatal(err)
		}
		return
	}

	w := bufio.NewWriter(os.Stdout)
	for _, s := range goseg.Segment(text, opts...) {
		if _, err := fmt.Fprintln(w, s); err != nil {
			fatal(err)
		}
	}
	if err := w.Flush(); err != nil {
		fatal(err)
	}
}

func readInput(path string) (string, error) {
	if path == "" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return string(b), nil
	}
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %q: %w", path, err)
	}
	b, err := io.ReadAll(f)
	cerr := f.Close()
	if err != nil {
		return "", fmt.Errorf("read %q: %w", path, err)
	}
	if cerr != nil {
		return "", fmt.Errorf("close %q: %w", path, cerr)
	}
	return string(b), nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "goseg: %v\n", err)
	os.Exit(1)
}
