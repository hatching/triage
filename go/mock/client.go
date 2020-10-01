// Copyright (C) 2020 Hatching B.V.
// All rights reserved.

package mock_triage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	http "net/http"
)

type ClientMock struct {
	RequestMethod string
	RequestUrl    string
	RequestBody   string
	Response      interface{}
	StatusCode    int
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	var data interface{}
	c.RequestMethod = req.Method
	if req.Method == "POST" {
		err := json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			panic(err)
		}
		js, err := json.Marshal(data)
		c.RequestBody = string(js)
	}
	c.RequestUrl = req.URL.String()

	resp, err := json.Marshal(c.Response)
	if err != nil {
		panic(err)
	}
	var status string
	switch c.StatusCode {
	case 200:
		status = "200 OK"
	case 400:
		status = "400 Bad Request"
	case 500:
		status = "500 Internal Server Error"
	}
	return &http.Response{
		Status:     status,
		StatusCode: c.StatusCode,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       ioutil.NopCloser(bytes.NewBuffer(resp)),
	}, nil
}
