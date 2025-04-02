package djot_parser

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ratrocket/godjot/djot_tokenizer"
)

func seedFuzz(f *testing.F) {
	dir, err := os.ReadDir(examplesDir)
	if err != nil {
		f.Fatal(err)
	}
	for _, entry := range dir {
		name := entry.Name()
		example, ok := strings.CutSuffix(name, ".html")
		if !ok {
			continue
		}
		djotExample, err := os.ReadFile(path.Join(examplesDir, fmt.Sprintf("%v.djot", example)))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(djotExample))
	}
}

func FuzzDjotE2E(f *testing.F) {
	seedFuzz(f)
	f.Fuzz(func(t *testing.T, input string) { _ = printDjot(input) })
}

func FuzzDjotTokenizer(f *testing.F) {
	seedFuzz(f)
	f.Fuzz(func(t *testing.T, input string) { _ = djot_tokenizer.BuildDjotTokens([]byte(input)) })
}
