package extract

import (
	"reflect"
	"strings"
	"testing"
)

func TestMetadataExtractor(t *testing.T) {
	doc := `
<html>
<head>
<title>Hello there.</title>
<meta property="og:title" content="Also hello." />
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
	expectedFields := []MetadataField{
		{Key: "title", Val: "Also hello.", Source: "og"},
	}
	for _, expected := range expectedFields {
		actual := meta.Metadata[expected.Key]
		if !reflect.DeepEqual(actual, expected) {
			t.Error("Expected", expected, "got", actual)
		}
	}

}
