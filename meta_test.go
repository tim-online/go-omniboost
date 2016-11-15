package omniboost

import (
	"encoding/json"
	"testing"
)

var (
	firstPageLinksJSONBlob = []byte(`{
		"meta": {
			"pagination": {
				"total": 51,
				"count": 25,
				"per_page": 25,
				"current_page": 1,
				"total_pages": 3,
				"links": {
					"first": "https://omniboost.io/api/products/?page=1",
					"last": "https://omniboost.io/api/products/?page=3",
					"next": "https://omniboost.io/api/products/?page=2",
					"current": "https://omniboost.io/api/products/?page=1"
				}
			}
		}
	}`)
	otherPageLinksJSONBlob = []byte(`{
		"meta": {
			"pagination": {
				"total": 51,
				"count": 25,
				"per_page": 25,
				"current_page": 2,
				"total_pages": 3,
				"links": {
					"first": "https://omniboost.io/api/products/?page=1",
					"prev": "https://omniboost.io/api/products/?page=1",
					"last": "https://omniboost.io/api/products/?page=3",
					"next": "https://omniboost.io/api/products/?page=3",
					"current": "https://omniboost.io/api/products/?page=2"
				}
			}
		}
	}`)
	lastPageLinksJSONBlob = []byte(`{
		"meta": {
			"pagination": {
				"total": 51,
				"count": 1,
				"per_page": 25,
				"current_page": 3,
				"total_pages": 3,
				"links": {
					"first": "https://omniboost.io/api/products/?page=1",
					"last": "https://omniboost.io/api/products/?page=3",
					"prev": "https://omniboost.io/api/products/?page=2",
					"current": "https://omniboost.io/api/products/?page=3"
				}
			}
		}
	}`)

	missingLinksJSONBlob = []byte(`{ }`)
)

func loadMetaJSON(t *testing.T, j []byte) Meta {
	root := newProductRoot()
	err := json.Unmarshal(j, &root)
	if err != nil {
		t.Fatal(err)
	}

	return *root.Meta
}

func TestLinks_ParseFirst(t *testing.T) {
	meta := loadMetaJSON(t, firstPageLinksJSONBlob)
	pagination := meta.Pagination
	links := meta.Pagination.Links
	if links.Current == nil {
		t.Fatal("links.Current shouln't be empty")
	}

	expectedPage := 1
	if pagination.CurrentPage != expectedPage {
		t.Fatalf("expected current page to be '%d', was '%d'", expectedPage, pagination.CurrentPage)
	}

	if pagination.IsLastPage() {
		t.Fatalf("shouldn't be last page")
	}
}

func TestLinks_ParseMiddle(t *testing.T) {
	meta := loadMetaJSON(t, otherPageLinksJSONBlob)
	pagination := meta.Pagination
	links := meta.Pagination.Links
	if links.Current == nil {
		t.Fatal("links.Current shouln't be empty")
	}

	expectedPage := 2
	if pagination.CurrentPage != expectedPage {
		t.Fatalf("expected current page to be '%d', was '%d'", expectedPage, pagination.CurrentPage)
	}

	if pagination.IsLastPage() {
		t.Fatalf("shouldn't be last page")
	}
}

func TestLinks_ParseLast(t *testing.T) {
	meta := loadMetaJSON(t, lastPageLinksJSONBlob)
	pagination := meta.Pagination
	links := meta.Pagination.Links
	if links.Current == nil {
		t.Fatal("links.Current shouln't be empty")
	}

	expectedPage := 3
	if pagination.CurrentPage != expectedPage {
		t.Fatalf("expected current page to be '%d', was '%d'", expectedPage, pagination.CurrentPage)
	}

	if !pagination.IsLastPage() {
		t.Fatalf("expected last page")
	}
}

func TestLinks_ParseMissing(t *testing.T) {
	meta := loadMetaJSON(t, missingLinksJSONBlob)
	pagination := meta.Pagination
	// links := meta.Pagination.Links
	// if links.Current == nil {
	// 	t.Fatal("links.Current shouln't be empty")
	// }

	expectedPage := 1
	if pagination.CurrentPage != expectedPage {
		t.Fatalf("expected current page to be '%d', was '%d'", expectedPage, pagination.CurrentPage)
	}
}

func TestLinks_ParseURL(t *testing.T) {
	type linkTest struct {
		name, url string
		expected  int
	}

	linkTests := []linkTest{
		{
			name:     "prev",
			url:      "https://omniboost.io/api/products/?page=1",
			expected: 1,
		},
		{
			name:     "last",
			url:      "https://omniboost.io/api/products/?page=5",
			expected: 5,
		},
		{
			name:     "next",
			url:      "https://omniboost.io/api/products/?page=2",
			expected: 2,
		},
	}

	for _, lT := range linkTests {
		p, err := pageForURL(lT.url)
		if err != nil {
			t.Fatal(err)
		}

		if p != lT.expected {
			t.Errorf("expected page for '%s' to be '%d', was '%d'",
				lT.url, lT.expected, p)
		}
	}

}

func TestLinks_ParseEmptyString(t *testing.T) {
	type linkTest struct {
		name, url string
		expected  int
	}

	linkTests := []linkTest{
		{
			name:     "none",
			url:      "http://example.com",
			expected: 0,
		},
		{
			name:     "bad",
			url:      "no url",
			expected: 0,
		},
		{
			name:     "empty",
			url:      "",
			expected: 0,
		},
	}

	for _, lT := range linkTests {
		_, err := pageForURL(lT.url)
		if err == nil {
			t.Fatalf("expected error for test '%s', but received none", lT.name)
		}
	}
}
