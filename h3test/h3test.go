package h3test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/tsingshaner/go-pkg/log/console"
)

type Request struct {
	Method string      // GET, POST, PUT, PATCH, DELETE
	Path   string      // /api/v1/users
	Body   any         // body
	Query  url.Values  // query
	Header http.Header // header
}

type Engine interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func New(path string) *Request {
	return &Request{
		Method: "GET",
		Path:   path,
		Query:  make(url.Values),
		Header: make(http.Header),
	}
}

func (r *Request) JSON(body any) *Request {
	req := r.SetBody(body)
	if !req.hasHeader("Content-Type", "application/json") {
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
	}

	return req
}

func (r *Request) hasHeader(k string, v string) bool {
	return strings.Contains(strings.Join(r.Header.Values(k), " "), v)
}

// ExtractJSON extracts the JSON response body into a struct pointer.
func ExtractJSON[T any](resp *httptest.ResponseRecorder) *T {
	var target T
	err := json.Unmarshal(resp.Body.Bytes(), &target)
	if err != nil {
		panic(err)
	}
	return &target
}

func (r *Request) POST() *Request {
	req := r.Clone()
	req.Method = "POST"
	return req
}

func (r *Request) GET() *Request {
	req := r.Clone()
	r.Method = "GET"
	return req
}

func (r *Request) PUT() *Request {
	req := r.Clone()
	r.Method = "PUT"
	return req
}

func (r *Request) PATCH() *Request {
	req := r.Clone()
	r.Method = "PATCH"
	return req
}

func (r *Request) DELETE() *Request {
	req := r.Clone()
	req.Method = "DELETE"
	return req
}

func (r *Request) AddHeader(key, value string) *Request {
	req := r.Clone()
	req.Header.Add(key, value)
	return req
}

func (r *Request) AddQuery(key, value string) *Request {
	req := r.Clone()
	req.Query.Add(key, value)
	return req
}

func (r *Request) SetBody(body any) *Request {
	req := r.Clone()
	req.Body = body
	return req
}

func (r *Request) Clone() *Request {
	newReq := &Request{
		Path:   r.Path,
		Method: r.Method,
		Body:   r.Body,
		Query:  cloneMap(r.Query),
		Header: cloneMap(r.Query),
	}

	return newReq
}

func cloneMap(original map[string][]string) map[string][]string {
	clone := make(map[string][]string)
	for key, values := range original {
		newValues := make([]string, len(values))
		copy(newValues, values)
		clone[key] = newValues
	}
	return clone
}

func (r *Request) Send(engine http.Handler) *httptest.ResponseRecorder {
	return jsonRequest(engine, r)
}

func jsonRequest(engine http.Handler, request *Request) *httptest.ResponseRecorder {
	var encodedBody []byte
	if request.hasHeader("Content-Type", "application/json") {
		if json, err := json.Marshal(request.Body); err != nil {
			panic(err)
		} else {
			encodedBody = json
		}
	} else if text, ok := request.Body.(string); ok {
		encodedBody = []byte(text)
	} else if bytes, ok := request.Body.([]byte); ok {
		encodedBody = bytes
	} else if request.Body != nil {
		console.Fatal("h3test this version only support text body or auto encode when Content-Type includes application/json")
	}

	if len(request.Query) != 0 {
		request.Path = request.Path + "?" + request.Query.Encode()
	}

	req := httptest.NewRequest(request.Method, request.Path, bytes.NewReader(encodedBody))
	req.Header = request.Header

	respRecorder := httptest.NewRecorder()
	engine.ServeHTTP(respRecorder, req)
	return respRecorder
}
