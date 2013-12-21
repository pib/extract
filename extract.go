// Package extract contains tools for composing extractors to collect
// information from an HTML document in a streaming fashion as it is
// parsed
package extract

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"io"
	"strings"
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

type MetadataField struct {
	Key    string
	Val    string
	Source string
}

type MetadataExtractor struct {
	Metadata map[string]MetadataField
	inTitle  bool
}

func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{Metadata: make(map[string]MetadataField)}
}

func (meta *MetadataExtractor) HandleToken(token html.Token) {
	switch token.Type {
	case html.StartTagToken:
		switch token.DataAtom {
		case atom.Title:
			meta.inTitle = true
		}
	case html.TextToken:
		if _, titleSet := meta.Metadata["title"]; meta.inTitle && !titleSet {
			meta.Metadata["title"] = MetadataField{Key: "title", Val: token.Data, Source: "titleTag"}
		}
	case html.EndTagToken:
		switch token.DataAtom {
		case atom.Title:
			meta.inTitle = false
		}
	case html.SelfClosingTagToken:
		if token.DataAtom == atom.Meta {
			if prop, _ := Attr(token, "property"); strings.HasPrefix(prop, "og:") {
				key := strings.TrimPrefix(prop, "og:")
				content, _ := Attr(token, "content")
				meta.Metadata[key] = MetadataField{Key: key, Val: content, Source: "og"}
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
