// Package extract contains tools for composing extractors to collect
// information from an HTML document in a streaming fashion as it is
// parsed
package extract

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"strings"
)

type MetadataField struct {
	Key    string
	Value  string
	Source string
	Weight int
}

type Metadata map[string]MetadataField

func (meta Metadata) Set(commonKey string, field MetadataField) {
	meta[field.Key] = field
	if oldField, exists := meta[commonKey]; !exists || (oldField.Weight < field.Weight) {
		meta[commonKey] = field
	}
}

// A MetadataExtractor gathers metadata from OpenGraph and other meta
// tags, as well as from the title tag and link tags with rel
// attributes.
type MetadataExtractor struct {
	Metadata Metadata
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
			meta.Metadata.Set("title",
				MetadataField{Key: "title_tag", Value: token.Data, Source: "title_tag"})
		}
	case html.EndTagToken:
		switch token.DataAtom {
		case atom.Title:
			meta.inTitle = false
		}
	case html.SelfClosingTagToken:
		switch token.DataAtom {
		case atom.Meta:
			if prop, _ := Attr(token, "property"); strings.HasPrefix(prop, "og:") {
				key := strings.TrimPrefix(prop, "og:")
				content, _ := Attr(token, "content")
				meta.Metadata.Set(key, MetadataField{"og_" + key, content, "og", 1})
			}
			if name, exists := Attr(token, "name"); exists {
				content, _ := Attr(token, "content")
				meta.Metadata.Set(name,
					MetadataField{Key: "meta_" + name, Value: content, Source: "meta"})
			}
		case atom.Link:
			if rel, exists := Attr(token, "rel"); exists {
				href, _ := Attr(token, "href")
				key := "link_rel_" + rel
				field := MetadataField{Key: key, Value: href, Source: "linkRel"}

				if rel == "canonical" {
					meta.Metadata.Set("url", field)
				} else {
					meta.Metadata[key] = field
				}
			}
		}
	}
}
