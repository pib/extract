// Package extract contains tools for composing extractors to collect
// information from an HTML document in a streaming fashion as it is
// parsed
package extract

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"io"
)

type Extractor interface {
	HandleToken(html.Token)
}

// An ElementAttributeExtractor extracts the specified attribute from
// each instance of the specified element, and stores their values.
type ElementAttributeExtractor struct {
	Element atom.Atom
	AttrKey string
	Values  *[]string
}

func (e ElementAttributeExtractor) HandleToken(token html.Token) {
	switch token.Type {
	case html.StartTagToken, html.SelfClosingTagToken:
		if token.DataAtom == e.Element {
			if val, exists := Attr(token, e.AttrKey); exists {
				*e.Values = append(*e.Values, val)
			}
		}
	}
}

func Extract(r io.Reader, extractor Extractor) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			switch z.Err() {
			case io.EOF:
				return nil
			default:
				return z.Err()
			}
		default:
			token := z.Token()
			extractor.HandleToken(token)
		}
	}

}

func Attr(token html.Token, key string) (string, bool) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
