package httpservice

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Name() string
	Errorf(format string, args ...interface{})
}

// TestClient is a wrapper to facilitate easy
// HTTP Gin Handler testing.
type TestClient struct {
	t TestingT
}

// NewTestClient returns a TestClient initialized with
// the given testing.T
func NewTestClient(t TestingT) TestClient {
	return TestClient{t: t}
}

func (t TestClient) newRequest(req *http.Request) *TestRequest {
	gin.SetMode(gin.ReleaseMode)

	var bodyBuffer *bytes.Buffer

	if req.Body != nil {
		bodyBuffer = &bytes.Buffer{}
		req.Body = ioutil.NopCloser(io.TeeReader(req.Body, bodyBuffer))
	}

	rec := httptest.NewRecorder()
	c, e := gin.CreateTestContext(rec)

	c.Request = req

	return &TestRequest{
		Engine:   e,
		Recorder: rec,
		Context:  c,
		t:        t.t,
		body:     bodyBuffer,
	}
}

// Get returns a new TestRequest where the request is a GET.
// The params are added to the request.
func (t TestClient) Get(params map[string]string) *TestRequest {
	p := t.MapToGinParams(params)
	r := httptest.NewRequest("GET", "/", nil)
	tr := t.newRequest(r)
	tr.Context.Params = p

	return tr
}

// GetQuery returns a new TestRequest where the request is a GET.
// The params are added to the request and query items are added
// to the query string of the request. Gin can't 'see' query params
// on a TestRequest with c.Query() unless they are added directly
// to the URL of the request.
func (t TestClient) GetQuery(params, queryItems map[string]string) *TestRequest {
	p := t.MapToGinParams(params)
	r := httptest.NewRequest("GET", "/", nil)
	tr := t.newRequest(r)
	tr.Context.Params = p

	// queryItems
	q := tr.Context.Request.URL.Query()
	for key := range queryItems {
		q.Add(key, queryItems[key])
	}
	tr.Context.Request.URL.RawQuery = q.Encode()

	return tr
}

// PostJSON returns a new TestRequest where the request is a post.
// Takes an interface to marshal into JSON and set as the request body.
func (t TestClient) PostJSON(req interface{}) *TestRequest {
	body, _ := json.Marshal(req)
	r := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	return t.newRequest(r)
}

// PostWithParams returns a new TestRequest where the request is a POST.
// Takes an interface to marshal into JSON & params
// The request body is set and the params are added to the request.
func (t TestClient) PostWithParams(params map[string]string, req interface{}) *TestRequest {
	body, _ := json.Marshal(req)
	r := httptest.NewRequest("POST", "/", bytes.NewReader(body))

	r.Header.Set("Content-Type", "application/json")

	p := t.MapToGinParams(params)

	tr := t.newRequest(r)
	tr.Context.Params = p

	return tr
}

// MapToGinParams maps map[string]string{} to gin.Params
func (t TestClient) MapToGinParams(params map[string]string) []gin.Param {
	l := len(params)
	p := make([]gin.Param, l)

	for key, value := range params {
		p = append(p, gin.Param{Key: key, Value: value})
	}

	return p
}

// TestRequest is a struct to facilitate easy
// HTTP Context testing with Gin Handlers.
type TestRequest struct {
	Recorder *httptest.ResponseRecorder
	Context  *gin.Context
	Engine   *gin.Engine
	body     *bytes.Buffer

	t TestingT
}

// WithT returns a copy of TestRequest with a new T
func (c *TestRequest) WithT(t TestingT) *TestRequest {
	return &TestRequest{
		Recorder: c.Recorder,
		Context:  c.Context,
		Engine:   c.Engine,
		body:     c.body,
		t:        t,
	}
}

// AssertStatus asserts that the recorded status
// is the same as the submitted status.
// The message and args are added to the message when
// the assertion fails.
func (c *TestRequest) AssertStatus(code int, msgAndArgs ...interface{}) bool {
	c.Context.Writer.WriteHeaderNow()

	r := c.Recorder.Result()

	//nolint:errcheck
	defer r.Body.Close()

	return assert.Equal(c.t, code, r.StatusCode, msgAndArgs...)
}

// AssertJSONBody asserts that the recorded response body
// unmarshals to the given interface.
// The message and args are added to the message when
// the assertion fails.
func (c *TestRequest) AssertJSONBody(obj interface{}, msgAndArgs ...interface{}) bool {
	var newInstance reflect.Value
	elem := reflect.TypeOf(obj)

	switch elem.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		theType := elem.Elem()
		newInstance = reflect.New(theType)
	default:
		newInstance = reflect.New(elem)
	}

	out := newInstance.Interface()

	err := json.NewDecoder(c.Recorder.Body).Decode(out)
	if !assert.NoError(c.t, err) {
		return false
	}

	return assert.Equal(c.t, obj, out, msgAndArgs...)
}

// AssertEqualToFile will assert that the HTTP response is equal to the file at the path. Takes an optional filename, otherwise generates a default based upon the test name
func (c *TestRequest) AssertEqualToFile(path ...string) bool {
	f := c.fileName(path...)

	expected, err := ioutil.ReadFile(filepath.Clean(f))
	if !assert.NoError(c.t, err) {
		return false
	}

	r := c.Recorder.Result()

	//nolint:errcheck
	defer r.Body.Close()

	resp, err := httputil.DumpResponse(r, true)
	if !assert.NoError(c.t, err) {
		return false
	}

	return assert.Equal(c.t, string(expected), string(resp))
}

// WriteToFile will write the HTTP response to the file. Takes an optional filename, otherwise generates a default based upon the test name
func (c *TestRequest) WriteToFile(path ...string) bool {
	f := c.fileName(path...)

	r := c.Recorder.Result()

	//nolint:errcheck
	defer r.Body.Close()

	resp, err := httputil.DumpResponse(r, true)
	if !assert.NoError(c.t, err) {
		return false
	}

	err = mkdirAll(f)
	if !assert.NoError(c.t, err) {
		return false
	}

	err = ioutil.WriteFile(f, resp, 0600)

	return assert.NoError(c.t, err)
}

// AssertGoldenTest asserts that both the request and response are equal the the file at `test-fixtures/<test name>.golden`. Will update if `update` is true.
func (c *TestRequest) AssertGoldenTest(update bool, path ...string) bool {
	f := c.fileName(path...)

	output, ok := c.serialize()
	if !ok {
		return ok
	}

	if update {
		err := mkdirAll(f)
		if !assert.NoError(c.t, err) {
			return false
		}

		err = ioutil.WriteFile(f, output, 0600)
		if !assert.NoError(c.t, err) {
			return false
		}
	}

	expected, err := ioutil.ReadFile(filepath.Clean(f))
	if !assert.NoError(c.t, err) {
		return false
	}

	return assert.Equal(c.t, normalize(expected), normalize(output))
}

// Remove return characters from output
func normalize(b []byte) string {
	// Forcing /r/n to /n before /n to /r/n to avoid potential inconsistencies.
	// Ex. /r/n/r/n
	normalizedString := strings.ReplaceAll(string(b), "\r\n", "\n")
	normalizedString = strings.ReplaceAll(normalizedString, "\n", "\r\n")
	return normalizedString
}

func (c *TestRequest) fileName(path ...string) string {
	if len(path) > 0 {
		return filepath.Clean(path[0])
	}

	return filepath.Clean("test-fixtures/" + c.t.Name() + ".golden")
}

// serialize both the request and response to bytes
func (c *TestRequest) serialize() ([]byte, bool) {
	c.Context.Writer.WriteHeaderNow()

	if c.Context.Request.Body != nil {
		c.Context.Request.Body = ioutil.NopCloser(c.body)
	}

	req, err := httputil.DumpRequest(c.Context.Request, true)
	if !assert.NoError(c.t, err) {
		return nil, false
	}

	r := c.Recorder.Result()

	//nolint:errcheck
	defer r.Body.Close()

	resp, err := httputil.DumpResponse(r, true)
	if !assert.NoError(c.t, err) {
		return nil, false
	}

	separator := []byte("\n---\n")

	return append(append(req, separator...), resp...), true
}

// Make all the directories for the files
func mkdirAll(f string) error {
	return os.MkdirAll(path.Dir(f), 0700)
}
