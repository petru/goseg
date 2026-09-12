# goseg

Rule-based sentence boundary detection for Go. A port of the Ruby gem
[pragmatic_segmenter](https://github.com/diasks2/pragmatic_segmenter),
designed for speed on large corpora and for texts whose format and domain
are unknown.

It does not use machine learning and does not require training data. The
algorithm is conservative: when a boundary is ambiguous it prefers to keep
the text as one sentence. Parentheticals and quotations inside a sentence
stay with that sentence.

**Dependencies:** none (Go standard library only).

**License:** MIT. Copyright Kevin S. Dias (pragmatic_segmenter) and
Petru Madar (goseg).

## Build

Requires [Go](https://go.dev/dl/) 1.22 or later. There are no extra
dependencies.

```
git clone https://github.com/petru/goseg.git
cd goseg
go test ./...
go build -o goseg ./cmd/goseg
```

That writes a `goseg` binary in the current directory. To install it on
your `PATH`:

```
go install ./cmd/goseg
```

From a clone, import the library as `goseg` (see below).

## Library

```go
package main

import (
    "fmt"
    "goseg"
)

func main() {
    text := "Hello world. My name is Mr. Smith. I work for the U.S. Government and I live in the U.S. I live in New York."
    for _, s := range goseg.Segment(text) {
        fmt.Println(s)
    }
    // Hello world.
    // My name is Mr. Smith.
    // I work for the U.S. Government and I live in the U.S.
    // I live in New York.
}
```

Options:

```go
goseg.Segment(text, goseg.Language("de"))
goseg.Segment(text, goseg.Language("en"), goseg.DocType("pdf"))
goseg.Segment(text, goseg.CleanText(false))
goseg.Segment(text, goseg.Workers(1)) // sequential
goseg.Segment(text, goseg.Workers(8)) // at most 8 goroutines

cleaned := goseg.Clean(text, goseg.DocType("pdf"))
```

`Language` takes an ISO 639-1 code. Unknown codes use the common
English-like rules. Cleaning (errant newlines, HTML, table-of-contents
dots, and similar) is on by default.

## Parallelization

The abbreviation and punctuation rules on a single string have to run
in order. Parallelism is over **independent pieces of the file**, not
inside one paragraph.

On a large input (32 KB or more), goseg splits at blank lines, packs
those paragraphs into about as many chunks as there are CPUs, and
segments each chunk on its own goroutine. Results are concatenated in
document order. A file with no blank lines stays sequential so a
sentence is never cut in half.

| Option | Behavior |
|--------|----------|
| `Workers(0)` (default) | `GOMAXPROCS` goroutines when the text is large enough to split |
| `Workers(1)` | always one pass |
| `Workers(n)` for `n > 1` | at most `n` goroutines, even on smaller multi-paragraph text |

On a ~550 KB Project Gutenberg book this is about **6×** faster than a
single core (same sentences).

## Command line

After `go build -o goseg ./cmd/goseg` (or `go run ./cmd/goseg`):

```
./goseg -lang en file.txt
cat file.txt | ./goseg -lang ja
./goseg -doc-type pdf -clean-only file.txt
./goseg -workers 1 file.txt    # sequential
./goseg -workers 0 file.txt    # all CPUs (default)
```

Prints one sentence per line. Flags: `-lang`, `-doc-type`, `-no-clean`,
`-clean-only`, `-workers`.

## Languages

| Code | Language   | Code | Language  |
|------|------------|------|-----------|
| `am` | Amharic    | `hy` | Armenian  |
| `ar` | Arabic     | `it` | Italian   |
| `bg` | Bulgarian  | `ja` | Japanese  |
| `da` | Danish     | `kk` | Kazakh    |
| `de` | German     | `my` | Burmese   |
| `el` | Greek      | `nl` | Dutch     |
| `en` | English    | `pl` | Polish    |
| `es` | Spanish    | `ru` | Russian   |
| `fa` | Persian    | `ur` | Urdu      |
| `fr` | French     | `zh` | Chinese   |
| `hi` | Hindi      |      |           |

## Compatibility

Behavior is intended to match pragmatic_segmenter 0.3.24. The test suite
under `test/` replays 583 cases from that gem's RSpec suite.

Go's `regexp` engine (RE2) does not support lookaround, so those Ruby
patterns are implemented with equivalent prefix/suffix checks rather than
copied as-is. Placeholder characters used while hiding non-boundary
punctuation are the same as the gem's.

```
go test ./...
go test ./test -bench . -benchmem
```
