package extract

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"strings"
)

type TextExtractor struct {
	bytes.Buffer
	inScript bool
	inStyle  bool
}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

func (t *TextExtractor) HandleToken(token html.Token) {
	switch token.Type {
	case html.StartTagToken:
		switch token.DataAtom {
		// Upon hitting a opening body tag, clear the content seen so
		// far. This works around pages that don't have a body tag for
		// some reason by just grabbing all the textual content.
		case atom.Body:
			t.Reset()
		case atom.Script:
			t.inScript = true
		case atom.Style:
			t.inStyle = true
		}
	case html.EndTagToken:
		switch token.DataAtom {
		case atom.Head:
			t.Reset()
		case atom.Script:
			t.inScript = false
		case atom.Style:
			t.inStyle = false
		}
	case html.TextToken:
		if t.inScript || t.inStyle {
			return
		}
		if t.Len() > 0 {
			t.WriteString(" ")
		}
		words := strings.Fields(token.Data)
		if len(words) > 0 {
			t.WriteString(words[0])
			for _, word := range words[1:] {
				t.WriteString(" ")
				t.WriteString(word)
			}
		}
	}
}
