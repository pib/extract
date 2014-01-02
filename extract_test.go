package extract

import (
	"reflect"
	"strings"
	"testing"
)

type ExpectedField struct {
	keys []string
	meta MetadataField
}

func e(keys []string, key, val, source string) ExpectedField {
	return ExpectedField{keys, MetadataField{key, val, source, 0}}
}

type c []string

func TestMetadataExtractor(t *testing.T) {
	doc := `
<html>
<head>
<title>Hello there.</title>
<link rel="canonical" href="http://www.example.com/" />
<meta property="og:site_name" content="Hello Site" />
<meta property="og:title" content="Also hello." />
<meta property="og:description" content="This page doesn't have a whole lot going on." />
<meta name="description" content="Different description" />
<meta property="og:type" content="article" />
<meta property="og:url" content="http://example.com/" />
<meta property="og:image" content="http://placehold.it/80x60" />
</head>
</html>
`
	meta := NewMetadataExtractor()
	err := Extract(strings.NewReader(doc), meta)
	if err != nil {
		t.Error(err)
	}
	expectedFields := []ExpectedField{
		e(c{"link_rel_canonical"}, "link_rel_canonical", "http://www.example.com/", "linkRel"),
		e(c{"site_name", "og_site_name"}, "og_site_name", "Hello Site", "og"),
		e(c{"title", "og_title"}, "og_title", "Also hello.", "og"),
		e(c{"title_tag"}, "title_tag", "Hello there.", "title_tag"),
		e(c{"description"}, "og_description", "This page doesn't have a whole lot going on.", "og"),
		e(c{"meta_description"}, "meta_description", "Different description", "meta"),
		e(c{"type"}, "og_type", "article", "og"),
		e(c{"url"}, "og_url", "http://example.com/", "og"),
		e(c{"image"}, "og_image", "http://placehold.it/80x60", "og"),
	}
	for _, expected := range expectedFields {
		for _, key := range expected.keys {
			actual := meta.Metadata[key]
			expected.meta.Weight = actual.Weight // no expected weight
			if !reflect.DeepEqual(actual, expected.meta) {
				t.Error("For", key, "expected", expected.meta, "got", actual)
			}
		}
	}

}
