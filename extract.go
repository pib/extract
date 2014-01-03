// Package extract contains tools for composing extractors to collect
// information from an HTML document in a streaming fashion as it is
// parsed
package extract

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
)

type Extractor interface {
	HandleToken(html.Token)
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

type MultiExtractor []Extractor

func (extractors MultiExtractor) HandleToken(token html.Token) {
	for _, extractor := range extractors {
		extractor.HandleToken(token)
	}
}

type DebugExtractor struct {
	Extractor
}

func (de DebugExtractor) HandleToken(token html.Token) {
	fmt.Printf("Type: %s, DataAtom: \"%s\"(%d), Data: \"%s\", Attr: %v\n", token.Type, token.DataAtom, token.DataAtom, token.Data, token.Attr)
	de.Extractor.HandleToken(token)
}
