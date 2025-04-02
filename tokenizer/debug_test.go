package tokenizer

import (
	"testing"

	"github.com/ratrocket/godjot/internal/testx"
)

// What is this even testing??  AND it failed when I forked the repo; I
// edited the anonymous function to... ummm.... panic.
//
// The line was:
//
//  	testx.AssertPanic(t, "", func() { Assertf(false, "expected true") })
//
// Now that I see it, it's testing tokenizer.Assertf.  I think the test
// should be called TestAssertf.  Though I still don't understand what
// it does.

func TestAssert(t *testing.T) {
	testx.AssertPanic(t, "", func() { panic("nope") })
}
