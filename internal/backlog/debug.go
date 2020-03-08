package backlog

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (fn RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewDebugClient() *http.Client {
	logger := log.New(os.Stdout, "Debug: ", 0)

	return NewTestClient(func(req *http.Request) (*http.Response, error) {
		reqBody := []byte{}
		reqBuffer := &bytes.Buffer{}

		if req.Body != nil {
			if _, err := io.Copy(reqBuffer, req.Body); err != nil {
				return nil, err
			}
		}

		reqBody = reqBuffer.Bytes()

		if req.Body != nil {
			req.Body = ioutil.NopCloser(reqBuffer)
		}

		logger.Printf("HTTP Request: %v %v\n", req.Method, req.URL)

		if len(reqBody) > 1024 {
			logger.Printf("HTTP Payload: (%d bytes)\n", len(reqBody))
		} else if len(reqBody) != 0 {
			logger.Printf("HTTP Payload: %s (%d bytes)\n", reqBody, len(reqBody))
		}
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return res, err
		}

		resBuffer := &bytes.Buffer{}

		if _, err = io.Copy(resBuffer, res.Body); err != nil {
			return res, err
		}

		resBody := resBuffer.Bytes()
		res.Body = ioutil.NopCloser(resBuffer)

		if res.Header.Get("Content-Type") == "application/json" {
			logger.Printf("HTTP Response: %v (%d bytes) %s\n", res.StatusCode, len(resBody), resBody)
		} else {
			logger.Printf("HTTP Response: %v (%d bytes)\n", res.StatusCode, len(resBody))
		}

		return res, err
	})
}
