package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

// NewMock returns a pointer to http.Client with the mocked Transport.
func NewMock(r ...MockResponse) *http.Client {
	return &http.Client{
		Transport: newRoundTripper(r...),
	}
}

// newRoundTripper initializes a new roundtripper and accepts multiple Response
// structures as variadric arguments.
func newRoundTripper(r ...MockResponse) *RoundTripper {
	return &RoundTripper{
		Responses: r,
	}
}

// RoundTripper is aimed to be used as the Transport property in an http.Client
// in order to mock the responses that it would return in the normal execution.
// If the number of responses that are mocked are not enough, an error with the
// request iteration ID, method and full URL is returned.
type RoundTripper struct {
	Responses []MockResponse

	iteration int32
	mu        sync.RWMutex
}

// MockResponse Wraps the response of the RoundTrip.
type MockResponse struct {
	Response http.Response
	Error    error
}

// Add accepts multiple Response structures as variadric arguments and appends
// those to the current list of Responses.
func (rt *RoundTripper) Add(res ...MockResponse) *RoundTripper {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	rt.Responses = append(rt.Responses, res...)
	return rt
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
//
// RoundTrip should not attempt to interpret the response. In
// particular, RoundTrip must return err == nil if it obtained
// a response, regardless of the response's HTTP status code.
// A non-nil err should be reserved for failure to obtain a
// response. Similarly, RoundTrip should not attempt to
// handle higher-level protocol details such as redirects,
// authentication, or cookies.
//
// RoundTrip should not modify the request, except for
// consuming and closing the Request's Body. RoundTrip may
// read fields of the request in a separate goroutine. Callers
// should not mutate or reuse the request until the Response's
// Body has been closed.
//
// RoundTrip must always close the body, including on errors,
// but depending on the implementation may do so in a separate
// goroutine even after RoundTrip returns. This means that
// callers wanting to reuse the body for subsequent requests
// must arrange to wait for the Close call before doing so.
//
// The Request's URL and Header fields must be initialized.
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.mu.RLock()
	defer func() {
		atomic.AddInt32(&rt.iteration, 1)
		rt.mu.RUnlock()
	}()

	var iteration = atomic.LoadInt32(&rt.iteration)
	if int(iteration) > len(rt.Responses)-1 {
		return nil, fmt.Errorf(
			"failed to obtain response in iteration %d: %s %s",
			iteration+1, req.Method, req.URL,
		)
	}

	// Consume and close the body.
	if req.Body != nil {
		ioutil.ReadAll(req.Body)
		req.Body.Close()
	}

	r := rt.Responses[iteration]
	if r.Error != nil {
		return nil, r.Error
	}

	return &r.Response, nil
}

// NewStringBody creates an io.ReadCloser from a string.
func NewStringBody(b string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(b))
}

// NewByteBody creates an io.ReadCloser from a slice of bytes.
func NewByteBody(b []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(b))
}

// NewStructBody creates an io.ReadCloser from a structure that is attempted
// to be encoded into JSON. In case of failure, it panics.
func NewStructBody(i interface{}) io.ReadCloser {
	var b = new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(i); err != nil {
		panic(fmt.Sprintf("Failed to json.Encode structure %+v", i))
	}
	return ioutil.NopCloser(b)
}
