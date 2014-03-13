package extract

import (
	"strings"
	"testing"
)

func TestCategoryExtractor(t *testing.T) {
	doc := `
<html>
<head>
<title>Hello there.</title>
<meta property="article:section" content="Category 1" />
<meta property="article:section" content="Category 2" />
<meta property="article:tag" content="Tag 1" />
<meta property="article:tag" content="Tag 2" />
</head>
</html>
`
	catExt := NewCategoryExtractor()
	err := Extract(strings.NewReader(doc), catExt)
	if err != nil {
		t.Error(err)
	}
	expectedCategories := []string{"category 1", "category 2"}
	expectedTags := []string{"tag 1", "tag 2"}
	for _, expected := range expectedCategories {
		if _, gotCategory := catExt.Categories[expected]; !gotCategory {
			t.Error("Expected category", expected, "not in", catExt.Categories)
		}
	}
	for _, expected := range expectedTags {
		if _, gotTag := catExt.Tags[expected]; !gotTag {
			t.Error("Expected tag", expected, "not in", catExt.Tags)
		}
	}
}
