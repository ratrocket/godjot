package main

import (
	"flag"
	"io"
	"log"
	"os"

	"md0.org/djot/djot_parser"
	"md0.org/djot/html_writer"
)

func main() {
	os.Exit(run())
}

// TODO
// - Files aren't cleaned up when errors are hit.

func run() int {
	var (
		from      = flag.String("from", "", "path to the input djot file (empty or '-' for stdin)")
		to        = flag.String("to", "", "path to the output html file (empty or '-' for stdout)")
		overwrite = flag.Bool("overwrite", false, "overwrite output html file")

		in  io.Reader
		out io.Writer
	)
	flag.Parse()

	if *from == "" || *from == "-" {
		in = os.Stdin
	} else {
		f, err := os.Open(*from)
		if err != nil {
			log.Printf("failed to open input file %v: %v", *from, err)
			return 1
		}
		in = f
		defer f.Close()
	}
	if *to == "" || *to == "-" {
		out = os.Stdout
	} else {
		flags := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		if !*overwrite {
			flags |= os.O_EXCL
		}
		f, err := os.OpenFile(*to, flags, 0640)
		if err != nil {
			log.Printf("failed to open output file %v: %v", *to, err)
			return 1
		}
		defer f.Close()
		out = f
	}
	input, err := io.ReadAll(in)
	if err != nil {
		log.Printf("failed to read input file %v: %v", *from, err)
		return 1
	}
	ast := djot_parser.BuildDjotAst(input)
	context := djot_parser.NewConversionContext("html", djot_parser.DefaultConversionRegistry)
	html := []byte(context.ConvertDjotToHtml(&html_writer.HtmlWriter{}, ast...))
	for len(html) > 0 {
		n, err := out.Write(html)
		if err != nil {
			log.Printf("failed to write output file %v: %v", *to, err)
			return 1
		}
		html = html[n:]
	}

	return 0
}
