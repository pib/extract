// Package extract contains tools for composing extractors to collect
// information from an HTML document in a streaming fashion as it is
// parsed
package extract

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"net/url"
	"strings"
)

type MetadataField struct {
	Key    string
	Value  string
	Source string
	Weight int
}

type Metadata map[string]MetadataField

var ogUrlFields = map[string]struct{}{
	"url":   {},
	"image": {},
	"audio": {},
	"video": {},
}

func (meta Metadata) Set(commonKey string, field MetadataField) {
	meta[field.Key] = field
	if oldField, exists := meta[commonKey]; !exists || (oldField.Weight < field.Weight) {
		meta[commonKey] = field
	}
}

func (meta Metadata) Dict() map[string]string {
	dict := make(map[string]string)
	for key, val := range meta {
		dict[key] = val.Value
	}
	return dict
}

// A MetadataExtractor gathers metadata from OpenGraph and other meta
// tags, as well as from the title tag and link tags with rel
// attributes. For attributes containing URLs, it will transform them
// from relative to absolute, based on the baseUrl provided.
type MetadataExtractor struct {
	Metadata Metadata
	baseUrl  *url.URL
	inTitle  bool
}

func NewMetadataExtractor(baseUrl *url.URL) *MetadataExtractor {
	return &MetadataExtractor{baseUrl: baseUrl, Metadata: make(map[string]MetadataField)}
}

func (meta *MetadataExtractor) HandleToken(token html.Token) {
	switch token.Type {
	case html.TextToken:
		if _, titleSet := meta.Metadata["title"]; meta.inTitle && !titleSet {
			title := strings.TrimSpace(token.Data)
			meta.Metadata.Set("title",
				MetadataField{Key: "title_tag", Value: title, Source: "title_tag"})
		}
	case html.EndTagToken:
		switch token.DataAtom {
		case atom.Title:
			meta.inTitle = false
		}
	// Self-closing tags may not have "/", so watch start tags as well
	case html.SelfClosingTagToken, html.StartTagToken:
		switch token.DataAtom {
		case atom.Title:
			meta.inTitle = true
		case atom.Meta:
			if prop, _ := Attr(token, "property"); strings.HasPrefix(prop, "og:") {
				key := strings.TrimPrefix(prop, "og:")
				content, _ := Attr(token, "content")
				if _, isUrl := ogUrlFields[key]; isUrl {
					content = meta.absoluteUrl(content)
				}
				meta.Metadata.Set(key, MetadataField{"og_" + key, content, "og", 1})
			}
			if name, exists := Attr(token, "name"); exists {
				content, _ := Attr(token, "content")
				meta.Metadata.Set(name,
					MetadataField{Key: "meta_" + name, Value: content, Source: "meta"})
			}
		case atom.Link:
			if rel, exists := Attr(token, "rel"); exists {
				href, exists := Attr(token, "href")

				if exists {
					href = meta.absoluteUrl(href)
				}

				key := "link_rel_" + rel
				field := MetadataField{Key: key, Value: href, Source: "linkRel"}

				if rel == "canonical" {
					meta.Metadata.Set("url", field)
				} else if rel == "icon" {
					field.Weight = 1
					meta.Metadata.Set("favicon", field)
				} else if rel == "shortcut icon" {
					meta.Metadata.Set("favicon", field)
				} else {
					meta.Metadata[key] = field
				}
			}
		}
	}
}

func (meta *MetadataExtractor) absoluteUrl(href string) string {
	// Absolutize URL
	if meta.baseUrl != nil {
		relUrl, _ := meta.baseUrl.Parse(href)
		if relUrl != nil {
			return relUrl.String()
		}
	}
	return href
}
