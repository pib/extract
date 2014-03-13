package extract

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"strings"
)

type Classification map[string]struct{}

var metaCategoryFields = map[string]struct{}{
	"article:section": {},
}

var metaTagFields = map[string]struct{}{
	"article:tag": {},
}

func (c Classification) Set(key string) {
	if normKey := strings.TrimSpace(strings.ToLower(key)); normKey != "" {
		c[normKey] = struct{}{}
	}
}

// A CategoryExtractor gathers a list of categories and tags from
// OpenGraph and other meta tags.
type CategoryExtractor struct {
	Categories Classification
	Tags       Classification
}

func NewCategoryExtractor() *CategoryExtractor {
	return &CategoryExtractor{
		Categories: make(Classification),
		Tags:       make(Classification),
	}
}

func (c *CategoryExtractor) HandleToken(token html.Token) {
	switch token.Type {
	// Self-closing tags may not have "/", so watch start tags as well
	case html.SelfClosingTagToken, html.StartTagToken:
		switch token.DataAtom {
		case atom.Meta:
			content, _ := Attr(token, "content")
			if prop, exists := Attr(token, "property"); exists {
				if _, isCategory := metaCategoryFields[prop]; isCategory {
					c.Categories.Set(content)
				} else if _, isTag := metaTagFields[prop]; isTag {
					c.Tags.Set(content)
				}
			}
		}
	}
}

func (c CategoryExtractor) Dict() map[string][]string {
	categories := make([]string, len(c.Categories))
	tags := make([]string, len(c.Tags))

	for key, _ := range c.Categories {
		categories = append(categories, key)
	}
	for key, _ := range c.Tags {
		tags = append(tags, key)
	}
	return map[string][]string{
		"categories": categories,
		"tags":       tags,
	}
}
