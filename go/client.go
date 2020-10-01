// Copyright (C) 2019-2020 Hatching B.V.
// All rights reserved.

package triage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
)

const maxLimit = 200
const clientVersion = "alpha"

// A Client can be used to make requests to the Triage API.
type Client struct {
	rootURL    string
	token      string
	httpClient HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewClient creates a new client with the specified access token.
//
// You can find this token on your user account page on https://tria.ge if you
// have been granted API permissions.
//
// No attempt to validate the token is made.
//
// The client will make requests on behalf of the user that owns the token.
func NewClient(token string) *Client {
	return NewClientWithRootURL(token, "https://api.tria.ge")
}

// NewClientWithRootURL creates a client with a non-standard API root.
func NewClientWithRootURL(token, rootURL string) *Client {
	return &Client{
		rootURL:    rootURL,
		token:      token,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, method, c.rootURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("triage: %w", err)
	}
	r.Header.Set("Authorization", "Bearer "+c.token)
	r.Header.Set("User-Agent", fmt.Sprintf(
		"Triage Go Client/%s Go/%s",
		clientVersion,
		runtime.Version(),
	))
	return r, nil
}

// jsonRequestJSON crafts a request with the specified value as JSON body and
// parses the response as JSON.
func (c *Client) jsonRequestJSON(ctx context.Context, method, path string, reqVal, respVal interface{}) error {
	var body io.Reader
	if reqVal != nil {
		b, err := json.Marshal(reqVal)
		if err != nil {
			return fmt.Errorf("triage: %w", err)
		}
		body = bytes.NewReader(b)
	}

	r, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return fmt.Errorf("triage: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	return c.requestJSON(r, respVal)
}

// requestJSON performs the specified request and parses the response as JSON
// into the specified receiver.
func (c *Client) requestJSON(r *http.Request, respVal interface{}) error {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return fmt.Errorf("triage: %w", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode > 299 {
		return decodeAPIError(r, resp)
	}

	if respVal != nil {
		if err := json.NewDecoder(resp.Body).Decode(respVal); err != nil {
			return fmt.Errorf("triage: %w", err)
		}
	}
	return nil
}

func (c *Client) requestRawFile(ctx context.Context, method, path string) (io.ReadCloser, error) {
	r, err := c.newRequest(ctx, method, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("triage: %w", err)
	}
	if resp.StatusCode > 299 {
		return nil, decodeAPIError(r, resp)
	}
	return httpResponseReader{resp}, nil
}

// Error represents APi errors in a generic form.
//
// The Status and Kind fields can be used to programmaticaly determine the type
// of error to appropriately handle it. The Message field contains a human
// readable cause.
type Error struct {
	Status  int
	Kind    string `json:"error"`
	Message string `json:"message,omitempty"`
}

func decodeAPIError(r *http.Request, resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("triage: %s %s: status=%d", r.Method, r.URL, resp.StatusCode)
	}
	var apiErr Error
	if err := json.Unmarshal(b, &apiErr); err != nil {
		return fmt.Errorf("triage: %s %s: status=%d. body=%s", r.Method, r.URL, resp.StatusCode, b)
	}
	apiErr.Status = resp.StatusCode
	return &apiErr
}

func (err Error) Error() string {
	return fmt.Sprintf("triage: %d %s: %v", err.Status, err.Kind, err.Message)
}

// ErrorKind returns the kind of the specified error.
//
// If the error is not a Triage error, an empty string is returned.
func ErrorKind(err error) string {
	var apiErr Error
	if errors.As(err, &apiErr) {
		return apiErr.Kind
	}
	return ""
}

type httpResponseReader struct {
	*http.Response
}

func (resp httpResponseReader) Read(b []byte) (int, error) {
	return resp.Body.Read(b)
}

func (resp httpResponseReader) Close() error {
	io.Copy(ioutil.Discard, resp.Body)
	return resp.Body.Close()
}
