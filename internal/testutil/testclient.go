package testutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type TestClient struct {
	data              []byte
	enableHTTPRequest bool
}

func NewTestClient(data []byte, enableHTTPRequest bool) *TestClient {
	return &TestClient{
		data:              data,
		enableHTTPRequest: enableHTTPRequest,
	}
}

func (t *TestClient) doFakeRequest(req *http.Request) (*http.Response, error) {
	log.Printf("Fake Request: %v %v\n", req.Method, req.URL)

	res := &http.Response{}

	res.StatusCode = http.StatusOK
	res.Body = ioutil.NopCloser(bytes.NewBuffer(t.data))

	log.Printf("Fake Response: %v %s\n", res.StatusCode, t.data)

	return res, nil
}

func (t *TestClient) doHTTPRequest(req *http.Request) (*http.Response, error) {
	log.Printf("HTTP Request: %v %v\n", req.Method, req.URL)

	hc := &http.Client{}

	res, err := hc.Do(req)

	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	if _, err = io.Copy(buf, res.Body); err != nil {
		return nil, err
	}

	res.Body = ioutil.NopCloser(buf)

	log.Printf("HTTP Response: %v %s\n", res.StatusCode, buf.Bytes())

	return res, nil
}

func (t *TestClient) Do(req *http.Request) (*http.Response, error) {
	if t.enableHTTPRequest {
		return t.doHTTPRequest(req)
	}

	return t.doFakeRequest(req)
}
