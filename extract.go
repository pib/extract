// Package extract contains tools for composing extractors to collect
// information from an HTML document in a streaming fashion as it is
// parsed
package extract

import (
	"code.google.com/p/go.net/html"
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
