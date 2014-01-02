package extract

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
)

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

func Attr(token html.Token, key string) (string, bool) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
