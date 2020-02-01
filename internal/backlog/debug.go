package backlog

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type DebugClient struct {
	logger     *log.Logger
	httpClient *http.Client
}

func (d *DebugClient) Do(req *http.Request) (*http.Response, error) {
	reqBody := []byte{}
	reqBuffer := &bytes.Buffer{}

	if req.Body == nil {
		goto LOGGER
	}
	if _, err := io.Copy(reqBuffer, req.Body); err != nil {
		return nil, err
	}

	reqBody = reqBuffer.Bytes()
	req.Body = ioutil.NopCloser(reqBuffer)

LOGGER:

	d.logger.Printf("HTTP Request: %v %v\n", req.Method, req.URL)

	if len(reqBody) > 0 {
		d.logger.Printf("HTTP Payload: (%d bytes) %s\n", len(reqBody), reqBody)
	}

	res, err := d.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	resBuffer := &bytes.Buffer{}

	if _, err = io.Copy(resBuffer, res.Body); err != nil {
		return nil, err
	}

	resBody := resBuffer.Bytes()
	res.Body = ioutil.NopCloser(resBuffer)

	d.logger.Printf("HTTP Response: %v (%d bytes) %s\n", res.StatusCode, len(resBody), resBody)

	return res, nil
}

func NewDebugClient() *DebugClient {
	return &DebugClient{
		logger:     log.New(os.Stdout, "debug: ", 0),
		httpClient: &http.Client{},
	}
}
