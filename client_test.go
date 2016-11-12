package omniboost

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestClientDefaultBaseURL(t *testing.T) {
	c := NewClient(nil)
	if c.BaseURL == nil || c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, defaultBaseURL)
	}
}

func TestClientDefaultUserAgent(t *testing.T) {
	c := NewClient(nil)
	if c.UserAgent != userAgent {
		t.Errorf("NewClick UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}
}

func TestApiToken(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tokentest", func(w http.ResponseWriter, r *http.Request) {
		expected := fmt.Sprintf("Bearer %s", accessToken)
		if got := r.Header.Get("Authorization"); got != expected {
			t.Errorf("Authorization header = %q; want %q", got, expected)
		}
	})

	req, _ := client.NewRequest("GET", "/tokentest", nil)
	_, err := client.Do(req, nil)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}
}

func TestNewRequest(t *testing.T) {
	// c := NewClient(nil)

	// inURL, outURL := "/foo", defaultBaseURL+"foo"
	// inBody, outBody := &DropletCreateRequest{Name: "l"},
	// 	`{"name":"l","region":"","size":"","image":0,`+
	// 		`"ssh_keys":null,"backups":false,"ipv6":false,`+
	// 		`"private_networking":false,"tags":null}`+"\n"
	// req, _ := c.NewRequest("GET", inURL, inBody)

	// // test relative URL was expanded
	// if req.URL.String() != outURL {
	// 	t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	// }

	// // test body was JSON encoded
	// body, _ := ioutil.ReadAll(req.Body)
	// if string(body) != outBody {
	// 	t.Errorf("NewRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	// }

	// // test default user-agent is attached to the request
	// userAgent := req.Header.Get("User-Agent")
	// if c.UserAgent != userAgent {
	// 	t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	// }
}

func TestNewRequest_withUserData(t *testing.T) {
	// c := NewClient(nil)

	// inURL, outURL := "/foo", defaultBaseURL+"foo"
	// inBody, outBody := &DropletCreateRequest{Name: "l", UserData: "u"},
	// 	`{"name":"l","region":"","size":"","image":0,`+
	// 		`"ssh_keys":null,"backups":false,"ipv6":false,`+
	// 		`"private_networking":false,"user_data":"u","tags":null}`+"\n"
	// req, _ := c.NewRequest("GET", inURL, inBody)

	// // test relative URL was expanded
	// if req.URL.String() != outURL {
	// 	t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	// }

	// // test body was JSON encoded
	// body, _ := ioutil.ReadAll(req.Body)
	// if string(body) != outBody {
	// 	t.Errorf("NewRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	// }

	// // test default user-agent is attached to the request
	// userAgent := req.Header.Get("User-Agent")
	// if c.UserAgent != userAgent {
	// 	t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	// }
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_withCustomUserAgent(t *testing.T) {
	ua := "testing"
	c := NewClient(nil)
	c.UserAgent = ua

	req, _ := c.NewRequest("GET", "/foo", nil)

	expected := ua
	if got := req.Header.Get("User-Agent"); got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	_, err := client.Do(req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do()
// function.
func TestDo_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m",
			"errors": [{"resource": "r", "field": "f", "code": "c"}]}`)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Fatalf("Expected error response.")
	}

	expected := &ErrorResponse{
		Response: res,
		Message:  "m",
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}
}

// ensure that we properly handle API errors that do not contain a response
// body
func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	expected := &ErrorResponse{
		Response: res,
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}

func TestDo_completion_callback(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	var completedReq *http.Request
	var completedResp string
	client.OnRequestCompleted(func(req *http.Request, resp *http.Response) {
		completedReq = req
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Errorf("Failed to dump response: %s", err)
		}
		completedResp = string(b)
	})
	_, err := client.Do(req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}
	if !reflect.DeepEqual(req, completedReq) {
		t.Errorf("Completed request = %v, expected %v", completedReq, req)
	}
	expected := `{"A":"a"}`
	if !strings.Contains(completedResp, expected) {
		t.Errorf("expected response to contain %v, Response = %v", expected, completedResp)
	}
}
