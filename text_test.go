package extract

import (
	"strings"
	"testing"
)

type TestPair struct {
	in  string
	out string
}

func TestPlaintextTextExtractor(t *testing.T) {
	tests := []TestPair{
		{"Not really HTML, even.", "Not really HTML, even."},
		{"Ignore <!-- blah blah --> the <!-- pointless --> comments", "Ignore the comments"},
		{"<html><body>This one has a body, at least.</body></html>", "This one has a body, at least."},
		{`<html><head><title>Title isn't part of the text.</title></head>
          <body>Woop, getting a bit tricker.⚃</body></html>`,
			"Woop, getting a bit tricker.⚃"},
		{`<html><head><title>Title isn't part of the text.</title>
            </head>Somebody messed this file up severely.`,
			"Somebody messed this file up severely."},
		{`<html><head><title>Title isn't part of the text.</title></head>
          <body><h1>Heyo.</h1>Implicit<p>whitespace!</p></body>`,
			"Heyo. Implicit whitespace!"},
		{`<h1>Whoops</h1> <div>space</div>
`,
			"Whoops space"},
		{`<html><head><title>Title isn't part of the text.</title></head>
          <body><b>Heyo.</b>No Implicit<span>whitespace!</span></body>`,
			"Heyo.No Implicitwhitespace!"},
		{"Compressed   \t\n  whitespace", "Compressed whitespace"},
		{`<html><head><title>Title isn't part of the text.</title></head>
          <body>
           <script>
             var foo = "ignore this!";
           </script>
           <h1>That script tag.</h1>
           It should be ignored!</body>`,
			"That script tag. It should be ignored!"},
		{`<?xml version="1.0" encoding="UTF-8"?>
<html>
          <body>
           <style>
             body {
               background: chartreuse;
             }
           </style>
           <h1>That style tag.</h1>
           It should be ignored!</body>`,
			"That style tag. It should be ignored!"},
		{"<ul><li><a>one</a></li><li><a>two</a></li></ul>", "one two"},
		{"<strong> Blah<i> blah</i></strong>, blah? Blah blah. <a>Blar.</a>",
			"Blah blah, blah? Blah blah. Blar."},
		{`<p>a <em>b</em> <img/>  c d</p>`, "a b c d"},
		{`<p>a <em>b</em> <img alt=";)"/>  c d</p>`, "a b ;) c d"},
		{`<img alt="a b" /><br/>c d`, "a b c d"},
		{`Hello<div id="disqus_thread">This should be ignored</div>There`, "Hello There"},
		{`Cookie Cookie <em>Cookie</em><div id="disqus_thread">DAMNIT</div>`, "Cookie Cookie Cookie"},
	}

	for _, test := range tests {
		text := NewTextExtractor()
		err := Extract(strings.NewReader(test.in), text)
		if err != nil {
			t.Error(err)
		}
		if s := text.String(); s != test.out {
			t.Errorf("Expected \"%s\", got \"%s\"", test.out, s)
		}
	}
}
