package tokenizer

import (
	"testing"

	"github.com/ratrocket/godjot/internal/testx"
)

func TestLineTokenizer(t *testing.T) {
	document := []byte("hello\nworld\n!")
	tokenizer := LineTokenizer{Document: document}
	{
		start, end, eof := tokenizer.Scan()
		testx.AssertFalse(t, "", eof)
		testx.AssertEqual(t, "", "hello\n", string(document[start:end]))
	}
	{
		start, end, eof := tokenizer.Scan()
		testx.AssertFalse(t, "", eof)
		testx.AssertEqual(t, "", "world\n", string(document[start:end]))
	}
	{
		start, end, eof := tokenizer.Scan()
		testx.AssertFalse(t, "", eof)
		testx.AssertEqual(t, "", "!", string(document[start:end]))
	}
	{
		_, _, eof := tokenizer.Scan()
		testx.AssertTrue(t, "", eof)
	}
}
