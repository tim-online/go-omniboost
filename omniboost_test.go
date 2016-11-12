package omniboost

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"golang.org/x/oauth2"
)

var (
	mux         *http.ServeMux
	client      *Client
	accessToken string
	server      *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// build oauth http client
	accessToken = "mysecrettoken"
	token := &oauth2.Token{AccessToken: accessToken}
	ts := oauth2.StaticTokenSource(token)
	oauthClient := oauth2.NewClient(oauth2.NoContext, ts)

	client = NewClient(oauthClient)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	expected := url.Values{}
	for k, v := range values {
		expected.Add(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatalf("parseForm(): %v", err)
	}

	if !reflect.DeepEqual(expected, r.Form) {
		t.Errorf("Request parameters = %v, expected %v", r.Form, expected)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

// json := []byte(`{
// 	"name": "Kahn test 1",
// 	"url": "http://www.tim-online.nl/kahn.html",
// 	"description": "Kaaaaaaaaaaaahn",
// 	"variant": [
// 	{
// 		"sku": "kahn-test-1",
// 		"price": 40,
// 		"images": []
// 	}
// 	]
// }`)

// omni := NewClient()
// createRequest := &ProductCreateRequest{}
// json.Unmarshal(json, createRequest)
// omni.Products.Create(createRequest)

// Kijk:

// import "golang.org/x/oauth2"

// pat := "mytoken"
// type TokenSource struct {
//     AccessToken string
// }

// func (t *TokenSource) Token() (*oauth2.Token, error) {
//     token := &oauth2.Token{
//         AccessToken: t.AccessToken,
//     }
//     return token, nil
// }

// tokenSource := &TokenSource{
//     AccessToken: pat,
// }
// oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
// client := godo.NewClient(oauthClient)
