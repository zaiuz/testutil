package testutil

import "io/ioutil"
import "net/url"
import "net/http"
import "net/http/httptest"
import "testing"
import "regexp"
import z "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"

// ResponseExpectable saves current calling context when chaining exepctation methods. Do
// not use directly. Use one of the Http* method to obtain this and then use Expect method
// to setup test expectations.
type ResponseExpectable struct {
	T        *testing.T
	Response *http.Response
	Error    error
}

// HttpGet() starts a new HTTP GET request and returns an object for setting up
// expectation for the result.
func HttpGet(t *testing.T, url string) *ResponseExpectable {
	response, e := http.Get(url)
	return &ResponseExpectable{t, response, e}
}

// HttpPost() starts a new HTTP POST request with the given data payload and returns an
// object for setting up expectation for the result.
func HttpPost(t *testing.T, url string, data url.Values) *ResponseExpectable {
	response, e := http.PostForm(url, data)
	return &ResponseExpectable{t, response, e}
}

// Expect() reads the response body and tests if the response status code and body content
// matches the supplied values.
func (r *ResponseExpectable) Expect(code int, body string) {
	a.NoError(r.T, r.Error, "error while getting response.")
	a.Equal(r.T, r.Response.StatusCode, code, "invalid status code.")

	if len(body) > 0 {
		raw, e := ioutil.ReadAll(r.Response.Body)
		a.NoError(r.T, e, "error while reading response.")
		a.Equal(r.T, string(raw), body, "wrong response body.")
	}
}

// ExpectPattern() is similar to Expect() but the response body is matched against the
// given regular expression pattern instead.
func (r *ResponseExpectable) ExpectPattern(code int, pattern string) {
	a.NoError(r.T, r.Error, "error while getting response.")
	a.Equal(r.T, r.Response.StatusCode, code, "wrong status code.")

	if len(pattern) > 0 {
		re := regexp.MustCompile(pattern)
		raw, e := ioutil.ReadAll(r.Response.Body)
		a.NoError(r.T, e, "error while reading response.")
		a.True(r.T, re.Match(raw), "response body does not match pattern.")
	}
}

// NewTestRequestPair() creates a new test request/response pair for testing against a
// Context or any roundtripping code. The request is a simple GET / request and the
// response is an instance of httptest.ResponseRecorder.
func NewTestRequestPair() (http.ResponseWriter, *http.Request) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	return response, request
}

// NewTestContext() creates and returns a new test context by using the response and
// request pair from NewTestRequestPair()
func NewTestContext() (*z.Context) {
	response, request := NewTestRequestPair()
	return z.NewContext(response, request)
}
