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
	d.logger.Printf("HTTP Request: %v %v\n", req.Method, req.URL)

	res, err := d.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	if _, err = io.Copy(buf, res.Body); err != nil {
		return nil, err
	}

	res.Body = ioutil.NopCloser(buf)

	d.logger.Printf("HTTP Response: %v %s\n", res.StatusCode, buf.Bytes())

	return res, nil
}

func NewDebugClient() *DebugClient {
	return &DebugClient{
		logger:     log.New(os.Stdout, "debug: ", 0),
		httpClient: &http.Client{},
	}
}
