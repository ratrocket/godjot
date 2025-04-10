# godjot

[Djot](https://github.com/jgm/djot) markup language parser implemented
in Go language.

## Note on fork

This is my fork of
[sivukhin/godjot](https://github.com/sivukhin/godjot).  I removed
testify, ripped out all the CI stuff, took out something about "idea"
(an IDE?), altered the Makefile to my liking, etc, etc.  See commit
0016b84a41c5b3591f9d56e792f318285ed363e7 for my initial changes.

I'm evaluating if I can use this library as the basis for my own
djot'ing.  I've used [jotdown](https://github.com/hellux/jotdown) (also
pandoc) as a command line tool to simply render html from djot, but I
want to take a page from the book of
[jonashietala.se](https://www.jonashietala.se/blog/2024/02/02/blogging_in_djot_instead_of_markdown/)
and do some pre/alternative processing/transforming on my djot, the
first "victim" being generating a table of contents (see their
[mod.rs](https://github.com/treeman/jonashietala/blob/master/src/markup/djot/mod.rs)
and
[table_of_content.rs](https://github.com/treeman/jonashietala/blob/master/src/markup/djot/table_of_content.rs)
for the basic outline/idea).

The *problem* with just copying that stuff is... I don't work in rust, I
work in go.  So... jotdown seems like the closest / most approachable
library to what I want, but I'm not switching to rust.  So I'm going to
evaluate if *this* library (rather, my fork of it) is up to the task.
If not I'll change tack.

This fork will probably migrate over to sourcehut soon-ish and the
module will become md0.org/djot (I can't stand the naming convention of
`gothing`, or worse, `go-thing`).

I'll (try to, haha) update this README as things change.  *For now* what
follows is from the original repo (except for `s/sivukhin/ratrocket`).

## Installation

You can install **godjot** as a standalone binary:
```shell
$> go install github.com/ratrocket/godjot@latest
$> echo '*Hello*, _world_' | godjot
<p><strong>Hello</strong>, <em>world</em></p>
```

## Usage

**godjot** provides API to parse AST from djot string

``` go
var djot []byte
ast := djot_parser.BuildDjotAst(djot)
```

AST is loosely typed and described with following simple struct:

```go
type TreeNode[T ~int] struct {
    Type       T                     // one of DjotNode options
    Attributes tokenizer.Attributes  // string attributes of node
    Children   []TreeNode[T]         // list of child
    Text       []byte                // not nil only for TextNode
}
```

You can transform AST to HTML with predefined set of rules:

```go
content := djot_parser.NewConversionContext(
    "html",
    djot_parser.DefaultConversionRegistry,
    map[djot_parser.DjotNode]djot_parser.Conversion{
        /*
            You can overwrite default conversion rules with custom map
            djot_parser.ImageNode: func(state djot_parser.ConversionState, next func(c djot_parser.Children)) {
                state.Writer.
                    OpenTag("figure").
                    OpenTag("img", state.Node.Attributes.Entries()...).
                    OpenTag("figcaption").
                    WriteString(state.Node.Attributes.Get(djot_parser.ImgAltKey)).
                    CloseTag("figcaption").
                    CloseTag("figure")
            }
        */
    }
).ConvertDjotToHtml(&html_writer.HtmlWriter{}, ast...)
```

This implementation passes all examples provided in the
[spec](https://htmlpreview.github.io/?https://github.com/jgm/djot/blob/master/doc/syntax.html)
but can diverge from original javascript implementation in some cases.
