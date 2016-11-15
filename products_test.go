package omniboost

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/tim-online/go-omniboost/utils"
)

func TestDroplets_ListDroplets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/products", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"data" : [{
				"id": 1,
				"url" : "http://boulevard.omniboost.io/products/Test123",
				"description" : "long",
				"name" : "Test123",
				"options" : {
					"size" : "",
					"color" : ""
				}
			},{
				"id": 2,
				"url" : "http://boulevard.omniboost.io/products/Test456",
				"description" : "longer",
				"name" : "Test456",
				"options" : {
					"size" : "",
					"color" : ""
				}
			}]
		}`)
	})

	products, _, err := client.Products.List()
	if err != nil {
		t.Errorf("Products.List returned error: %v", err)
	}

	expected := []Product{
		{
			ID: 1,
			Url: utils.NewUrl("http://boulevard.omniboost.io/products/Test123"),
			Description: "long",
			Name: "Test123",
		}, {
			ID: 2,
			Url: utils.NewUrl("http://boulevard.omniboost.io/products/Test456"),
			Description: "longer",
			Name: "Test456",
		},
	}
	if !reflect.DeepEqual(products, expected) {
		t.Errorf("Products.List\n got=%#v\nwant=%#v", products, expected)
	}
}
