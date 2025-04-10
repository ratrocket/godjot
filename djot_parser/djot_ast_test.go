package djot_parser

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"md0.org/djot/djot_tokenizer"
	"md0.org/djot/html_writer"
	"md0.org/djot/internal/testx"
	"md0.org/djot/tokenizer"
)

func printDjot(text string) string {
	document := []byte(text)
	ast := BuildDjotAst(document)
	fmt.Printf("ast: %v\n", ast)
	return NewConversionContext("html", DefaultConversionRegistry).ConvertDjotToHtml(&html_writer.HtmlWriter{}, ast...)
}

const examplesDir = "examples"

func TestDownloadExample(t *testing.T) {
	normalize := func(line string) string {
		line = strings.Trim(line, "\r\n\t")
		line = strings.TrimPrefix(line, "<pre><code>")
		line = strings.TrimSuffix(line, "</code></pre>")
		return line
	}

	response, err := http.Get("https://raw.githubusercontent.com/jgm/djot/main/doc/syntax.html")
	testx.AssertNil(t, "", err)
	docBytes, err := io.ReadAll(response.Body)
	testx.AssertNil(t, "", err)
	var (
		djotStartToken = []byte(`<div class="djot">`)
		htmlStartToken = []byte(`<div class="html">`)
		endToken       = []byte(`</div>`)
	)
	example := 0
	for {
		djotStart := bytes.Index(docBytes, djotStartToken)
		if djotStart == -1 {
			break
		}
		djotEnd := djotStart + bytes.Index(docBytes[djotStart:], endToken)
		djotExample := html.UnescapeString(normalize(string(docBytes[djotStart+len(djotStartToken) : djotEnd])))
		docBytes = docBytes[djotEnd+len(endToken):]

		htmlStart := bytes.Index(docBytes, htmlStartToken)
		testx.AssertNotEqual(t, "", htmlStart, -1)
		htmlEnd := htmlStart + bytes.Index(docBytes[htmlStart:], endToken)
		htmlExample := html.UnescapeString(normalize(string(docBytes[htmlStart+len(htmlStartToken) : htmlEnd])))
		docBytes = docBytes[htmlEnd+len(endToken):]

		// Ignore 64th example because it's not self-contained and requires additional definition of table
		if example != 64 {
			testx.AssertNil(t, "", os.WriteFile(path.Join(examplesDir, fmt.Sprintf("%02d.html", example)), []byte(htmlExample), 0660))
			testx.AssertNil(t, "", os.WriteFile(path.Join(examplesDir, fmt.Sprintf("%02d.djot", example)), []byte(djotExample), 0660))
		}
		example++
	}
}

func TestStartSymbol(t *testing.T) {
	dir, err := os.ReadDir(examplesDir)
	testx.AssertNil(t, "", err)
	for _, entry := range dir {
		name := entry.Name()
		example, ok := strings.CutSuffix(name, ".html")
		if !ok {
			continue
		}
		djotExample, err := os.ReadFile(path.Join(examplesDir, fmt.Sprintf("%v.djot", example)))
		testx.AssertNil(t, "", err)
		_ = BuildDjotAst(djotExample)
	}
	symbols := make([]byte, 0)
	for s := range djot_tokenizer.StartSymbols {
		if !tokenizer.SpaceNewLineByteMask.Has(s) {
			symbols = append(symbols, s)
		}
	}
	sort.Slice(symbols, func(i, j int) bool { return symbols[i] < symbols[j] })
	t.Logf("%#v", string(symbols))
}

func TestDjotDocExample(t *testing.T) {
	dir, err := os.ReadDir(examplesDir)
	testx.AssertNil(t, "", err)
	for _, entry := range dir {
		name := entry.Name()
		example, ok := strings.CutSuffix(name, ".html")
		if !ok {
			continue
		}
		htmlExample, err := os.ReadFile(path.Join(examplesDir, fmt.Sprintf("%v.html", example)))
		testx.AssertNil(t, "", err)
		djotExample, err := os.ReadFile(path.Join(examplesDir, fmt.Sprintf("%v.djot", example)))
		testx.AssertNil(t, "", err)
		t.Run(example+":"+string(djotExample), func(t *testing.T) {
			result := printDjot(string(djotExample))
			testx.AssertEqual(
				t,
				fmt.Sprintf("invalid html (%v != %v), djot tokens: %v",
					string(htmlExample),
					result,
					djot_tokenizer.BuildDjotTokens(djotExample)),
				string(htmlExample),
				result,
			)
		})
	}
}

func TestManualExamples(t *testing.T) {
	t.Run("link in text", func(t *testing.T) {
		result := printDjot("link http://localhost:3000/debug/pprof/profile?seconds=10 -o profile.pprof")
		testx.AssertEqual(t, "", "<p>link http://localhost:3000/debug/pprof/profile?seconds=10 -o profile.pprof</p>\n", result)
	})
	t.Run("block attributes", func(t *testing.T) {
		result := printDjot(`{key="value"}
# Header`)
		testx.AssertEqual(t, "", `<section id="Header">
<h1 key="value">Header</h1>
</section>
`, result)
	})
	t.Run("inline attributes", func(t *testing.T) {
		result := printDjot(`![img](link){key="value"}`)
		t.Log(result)
	})
}
