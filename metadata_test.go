package extract

import (
	"net/url"
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
<link rel="canonical" href="http://www.example.com/posts/2" />
<link rel="next" href="./3" />
<link rel="all" href="all" />
<link rel="all_short" href="all/short" />
<link rel="home" href="../.." />
<link rel="shortcut icon" href="http://example.com/favicon.ico" />
<link rel="icon" href="http://example.com/favicon.png" />
<meta property="og:site_name" content="Hello Site" />
<meta property="og:title" content="Also hello." />
<meta property="og:description" content="This page doesn't have a whole lot going on." />
<meta name="description" content="Different description" />
<meta property="og:type" content="article" />
<meta property="og:url" content="http://example.com/" />
<meta property="og:image" content="/80x60.png" />
</head>
</html>
`
	url, _ := url.Parse("http://example.com/posts/2")
	meta := NewMetadataExtractor(url)
	err := Extract(strings.NewReader(doc), meta)
	if err != nil {
		t.Error(err)
	}
	expectedFields := []ExpectedField{
		e(c{"link_rel_canonical"}, "link_rel_canonical", "http://www.example.com/posts/2", "linkRel"),
		e(c{"link_rel_next"}, "link_rel_next", "http://example.com/posts/3", "linkRel"),
		e(c{"link_rel_all"}, "link_rel_all", "http://example.com/posts/all", "linkRel"),
		e(c{"link_rel_all_short"}, "link_rel_all_short", "http://example.com/posts/all/short", "linkRel"),
		e(c{"link_rel_home"}, "link_rel_home", "http://example.com/", "linkRel"),
		e(c{"site_name", "og_site_name"}, "og_site_name", "Hello Site", "og"),
		e(c{"title", "og_title"}, "og_title", "Also hello.", "og"),
		e(c{"title_tag"}, "title_tag", "Hello there.", "title_tag"),
		e(c{"description"}, "og_description", "This page doesn't have a whole lot going on.", "og"),
		e(c{"meta_description"}, "meta_description", "Different description", "meta"),
		e(c{"type"}, "og_type", "article", "og"),
		e(c{"url"}, "og_url", "http://example.com/", "og"),
		e(c{"image"}, "og_image", "http://example.com/80x60.png", "og"),
		e(c{"favicon", "link_rel_icon"}, "link_rel_icon", "http://example.com/favicon.png", "linkRel"),
		e(c{"link_rel_shortcut icon"}, "link_rel_shortcut icon", "http://example.com/favicon.ico", "linkRel"),
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
