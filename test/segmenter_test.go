package goseg_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"goseg"
)

type goldenCase struct {
	File     string      `json:"file"`
	Name     string      `json:"name"`
	Kind     string      `json:"kind"`
	Text     *string     `json:"text"`
	Language *string     `json:"language"`
	DocType  *string     `json:"doc_type"`
	Clean    *bool       `json:"clean"`
	Expected interface{} `json:"expected"`
}

func TestGoldenCases(t *testing.T) {
	data, err := os.ReadFile("testdata/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []goldenCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}

	fails := 0
	for i, c := range cases {
		if c.Text == nil {
			got := goseg.Segment("")
			if len(got) != 0 {
				t.Errorf("%s / %s: nil text: got %#v", c.File, c.Name, got)
				fails++
			}
			continue
		}
		lang := "en"
		if c.Language != nil && *c.Language != "" {
			lang = *c.Language
		}
		opts := []goseg.Option{goseg.Language(lang)}
		if c.Clean != nil {
			opts = append(opts, goseg.CleanText(*c.Clean))
		}
		if c.DocType != nil {
			opts = append(opts, goseg.DocType(*c.DocType))
		}

		if c.Kind == "clean" {
			got := goseg.Clean(*c.Text, opts...)
			want, ok := c.Expected.(string)
			if !ok {
				t.Errorf("#%d %s / %s: expected string, got %T", i, c.File, c.Name, c.Expected)
				fails++
				continue
			}
			if got != want {
				t.Errorf("#%d %s / %s\n got: %q\nwant: %q", i, c.File, c.Name, got, want)
				fails++
			}
			continue
		}

		got := goseg.Segment(*c.Text, opts...)
		want := expectedStrings(t, c.Expected)
		if !equalStrings(got, want) {
			t.Errorf("#%d %s / %s\n got: %#v\nwant: %#v", i, c.File, c.Name, got, want)
			fails++
			if fails > 40 {
				t.Fatalf("too many failures, stopping")
			}
		}
	}
}

func expectedStrings(t *testing.T, v interface{}) []string {
	t.Helper()
	switch x := v.(type) {
	case []interface{}:
		out := make([]string, 0, len(x))
		for _, item := range x {
			s, ok := item.(string)
			if !ok {
				t.Fatalf("expected string in golden list, got %T", item)
			}
			out = append(out, s)
		}
		return out
	case []string:
		return x
	default:
		t.Fatalf("expected list of strings, got %T", v)
		return nil
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSegment(t *testing.T) {
	tests := []struct {
		name string
		text string
		opts []goseg.Option
		want []string
	}{
		{
			name: "empty",
			text: "",
			want: []string{},
		},
		{
			name: "newline only",
			text: "\n",
			want: []string{},
		},
		{
			name: "english abbreviations",
			text: "Hello world. My name is Mr. Smith. I work for the U.S. Government and I live in the U.S. I live in New York.",
			want: []string{
				"Hello world.",
				"My name is Mr. Smith.",
				"I work for the U.S. Government and I live in the U.S.",
				"I live in New York.",
			},
		},
		{
			name: "html stripped to empty",
			text: "<b></b>",
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goseg.Segment(tt.text, tt.opts...)
			if !equalStrings(got, tt.want) {
				t.Fatalf("got %#v want %#v", got, tt.want)
			}
		})
	}
}

func TestDoesNotMutateInput(t *testing.T) {
	in := "It was a cold \nnight in the city."
	orig := in
	_ = goseg.Segment(in)
	if in != orig {
		t.Fatalf("input changed: %q", in)
	}
}

func TestParallelMatchesSequential(t *testing.T) {
	var b strings.Builder
	for i := 0; i < 40; i++ {
		b.WriteString("Hello world. My name is Mr. Smith. I live in the U.S. Government district.\n\n")
		b.WriteString("She turned to him, \"This is great.\" She held the book out.\n\n")
	}
	text := b.String()
	seq := goseg.Segment(text, goseg.Language("en"), goseg.Workers(1))
	par := goseg.Segment(text, goseg.Language("en"), goseg.Workers(8))
	if !equalStrings(seq, par) {
		t.Fatalf("parallel != sequential\n seq=%d par=%d", len(seq), len(par))
	}
	if len(seq) == 0 {
		t.Fatal("expected sentences")
	}
}

func TestParallelSingleParagraph(t *testing.T) {
	text := "Hello world. My name is Mr. Smith. I work for the U.S. Government."
	seq := goseg.Segment(text, goseg.Workers(1))
	par := goseg.Segment(text, goseg.Workers(8))
	if !equalStrings(seq, par) {
		t.Fatalf("got %#v want %#v", par, seq)
	}
}
