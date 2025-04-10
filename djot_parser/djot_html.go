package djot_parser

import (
	"fmt"
	"strings"

	"md0.org/djot/djot_tokenizer"
	"md0.org/djot/html_writer"
	"md0.org/djot/tokenizer"
)

type (
	ConversionContext struct {
		Format   string
		Registry ConversionRegistry
	}
	ConversionState struct {
		Format string
		Writer *html_writer.HtmlWriter
		Node   TreeNode[DjotNode]
		Parent *TreeNode[DjotNode]
	}
	Conversion         func(state ConversionState, next func(Children))
	ConversionRegistry map[DjotNode]Conversion
	Children           []TreeNode[DjotNode]
)

func (state ConversionState) StandaloneNodeConverter(tag string) *html_writer.HtmlWriter {
	return state.Writer.OpenTag(tag, state.Node.Attributes.Entries()...)
}

func (state ConversionState) InlineNodeConverter(tag string, next func(c Children)) *html_writer.HtmlWriter {
	return state.Writer.InTag(tag, state.Node.Attributes.Entries()...)(func() { next(nil) })
}

func (state ConversionState) BlockNodeConverter(tag string, next func(c Children)) *html_writer.HtmlWriter {
	content := func() {
		state.Writer.WriteString("\n")
		next(nil)
	}
	return state.Writer.InTag(tag, state.Node.Attributes.Entries()...)(content).WriteString("\n")
}

var DefaultSymbolRegistry = map[string]string{}

var htmlReplacer = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
	`>`, "&gt;",
	`–`, `&ndash;`,
	`—`, `&mdash;`,
	`“`, `&ldquo;`,
	`”`, `&rdquo;`,
	`‘`, `&lsquo;`,
	`’`, `&rsquo;`,
	`…`, `&hellip;`,
)

var DefaultConversionRegistry = map[DjotNode]Conversion{
	ThematicBreakNode: func(s ConversionState, n func(c Children)) { s.Writer.OpenTag("hr").WriteString("\n") },
	LineBreakNode:     func(s ConversionState, n func(c Children)) { s.Writer.OpenTag("br").WriteString("\n") },
	TextNode: func(s ConversionState, n func(c Children)) {
		if s.Parent != nil && (s.Parent.Attributes.Get(RawInlineFormatKey) == s.Format || s.Parent.Attributes.Get(RawBlockFormatKey) == s.Format) {
			s.Writer.WriteString(string(s.Node.Text))
		} else {
			s.Writer.WriteString(htmlReplacer.Replace(string(s.Node.Text)))
		}
	},
	SymbolsNode: func(s ConversionState, n func(c Children)) {
		symbol, ok := DefaultSymbolRegistry[string(s.Node.FullText())]
		if ok {
			s.Writer.WriteString(symbol)
		} else {
			s.Writer.WriteString(fmt.Sprintf(":%v:", string(s.Node.FullText())))
		}
	},
	InsertNode:       func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("ins", n) },
	DeleteNode:       func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("del", n) },
	SuperscriptNode:  func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("sup", n) },
	SubscriptNode:    func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("sub", n) },
	HighlightedNode:  func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("mark", n) },
	EmphasisNode:     func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("em", n) },
	StrongNode:       func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("strong", n) },
	ParagraphNode:    func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("p", n).WriteString("\n") },
	ImageNode:        func(s ConversionState, n func(c Children)) { s.StandaloneNodeConverter("img") },
	LinkNode:         func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("a", n) },
	SpanNode:         func(s ConversionState, n func(c Children)) { s.InlineNodeConverter("span", n) },
	DivNode:          func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("div", n) },
	TableCaptionNode: func(s ConversionState, n func(c Children)) { n(nil) },
	TableNode: func(s ConversionState, n func(c Children)) {
		if len(s.Node.Children) > 0 && s.Node.Children[0].Type == TableCaptionNode {
			s.Writer.OpenTag("table", s.Node.Attributes.Entries()...)
			s.Writer.WriteString("\n")
			s.Writer.InTag("caption")(func() { n(s.Node.Children[:1]) })
			s.Writer.WriteString("\n")
			s.Writer.InTag("tbody")(func() { n(s.Node.Children[1:]) })
			s.Writer.CloseTag("table")
		} else {
			s.BlockNodeConverter("table", n)
		}
	},
	TableRowNode: func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("tr", n) },
	TableHeaderNode: func(s ConversionState, n func(c Children)) {
		s.InlineNodeConverter("th", n)
		s.Writer.WriteString("\n")
	},
	TableCellNode: func(s ConversionState, n func(c Children)) {
		s.InlineNodeConverter("td", n)
		s.Writer.WriteString("\n")
	},
	TaskListNode:       func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("ul", n) },
	DefinitionListNode: func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("dl", n) },
	UnorderedListNode:  func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("ul", n) },
	OrderedListNode:    func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("ol", n) },
	ListItemNode: func(s ConversionState, n func(c Children)) {
		class := s.Node.Attributes.Get(djot_tokenizer.DjotAttributeClassKey)
		if class == CheckedTaskItemClass || class == UncheckedTaskItemClass {
			s.Writer.InTag("li")(func() {
				s.Writer.WriteString("\n")
				s.Writer.WriteString("<input disabled=\"\" type=\"checkbox\"")
				if class == CheckedTaskItemClass {
					s.Writer.WriteString(" checked=\"\"")
				}
				s.Writer.WriteString("/>").WriteString("\n")
				n(s.Node.Children)
			}).WriteString("\n")
		} else {
			s.BlockNodeConverter("li", n)
		}
	},
	DefinitionTermNode: func(s ConversionState, n func(c Children)) {
		s.InlineNodeConverter("dt", n)
		s.Writer.WriteString("\n")
	},
	DefinitionItemNode: func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("dd", n) },
	SectionNode:        func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("section", n) },
	QuoteNode:          func(s ConversionState, n func(c Children)) { s.BlockNodeConverter("blockquote", n) },
	DocumentNode:       func(s ConversionState, n func(c Children)) { n(nil) },
	FootnoteDefNode:    func(s ConversionState, n func(c Children)) { n(nil) },
	CodeNode: func(s ConversionState, n func(c Children)) {
		s.Writer.OpenTag("pre").OpenTag("code", s.Node.Attributes.Entries()...)
		n(nil)
		s.Writer.CloseTag("code").CloseTag("pre").WriteString("\n")
	},
	VerbatimNode: func(s ConversionState, n func(c Children)) {
		if _, ok := s.Node.Attributes.TryGet(djot_tokenizer.InlineMathKey); ok {
			attributes := append([]tokenizer.AttributeEntry{{Key: "class", Value: "math inline"}}, s.Node.Attributes.Entries()...)
			s.Writer.InTag("span", attributes...)(func() {
				s.Writer.WriteString("\\(")
				n(nil)
				s.Writer.WriteString("\\)")
			})
		} else if _, ok := s.Node.Attributes.TryGet(djot_tokenizer.DisplayMathKey); ok {
			attributes := append([]tokenizer.AttributeEntry{{Key: "class", Value: "math display"}}, s.Node.Attributes.Entries()...)
			s.Writer.InTag("span", attributes...)(func() {
				s.Writer.WriteString("\\[")
				n(nil)
				s.Writer.WriteString("\\]")
			})
		} else if rawFormat := s.Node.Attributes.Get(RawInlineFormatKey); rawFormat == s.Format {
			n(nil)
		} else {
			s.Writer.InTag("code", s.Node.Attributes.Entries()...)(func() { n(nil) })
		}
	},
	HeadingNode: func(s ConversionState, n func(c Children)) {
		level := len(s.Node.Attributes.Get(HeadingLevelKey))
		s.Writer.InTag(fmt.Sprintf("h%v", level), s.Node.Attributes.Entries()...)(func() { n(nil) }).WriteString("\n")
	},
	RawNode: func(s ConversionState, next func(c Children)) {
		if s.Format == s.Node.Attributes.Get(RawBlockFormatKey) {
			next(nil)
		}
	},
}

func NewConversionContext(format string, converters ...map[DjotNode]Conversion) ConversionContext {
	if len(converters) == 0 {
		converters = []map[DjotNode]Conversion{DefaultConversionRegistry}
	}
	registry := make(map[DjotNode]Conversion)
	for i := 0; i < len(converters); i++ {
		for node, conversion := range converters[i] {
			registry[node] = conversion
		}
	}
	return ConversionContext{
		Format:   format,
		Registry: registry,
	}
}

func (context ConversionContext) ConvertDjotToHtml(
	builder *html_writer.HtmlWriter,
	nodes ...TreeNode[DjotNode],
) string {
	context.convertDjotToHtml(builder, nil, nodes...)
	return builder.String()
}

func (context ConversionContext) convertDjotToHtml(builder *html_writer.HtmlWriter, parent *TreeNode[DjotNode], nodes ...TreeNode[DjotNode]) {
	for _, node := range nodes {
		currentNode := node
		conversion, ok := context.Registry[currentNode.Type]
		if !ok {
			continue
		}
		state := ConversionState{
			Format: context.Format,
			Writer: builder,
			Node:   currentNode,
			Parent: parent,
		}
		conversion(state, func(c Children) {
			if len(c) == 0 {
				context.convertDjotToHtml(builder, &node, currentNode.Children...)
			} else {
				context.convertDjotToHtml(builder, &node, c...)
			}
		})
	}
}
