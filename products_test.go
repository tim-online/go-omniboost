package omniboost

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDroplets_ListDroplets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/products", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"products": [{"id":1},{"id":2}]}`)
	})

	products, _, err := client.Products.List()
	if err != nil {
		t.Errorf("Products.List returned error: %v", err)
	}

	expected := []Product{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Products.List\n got=%#v\nwant=%#v", products, expected)
	}
}
