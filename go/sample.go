// Copyright (C) 2019-2021 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type SampleStatus string

const (
	// A sample has been submitted and is queued for static analysis or the
	// static analys is in progress.
	SampleStatusPending SampleStatus = "pending"
	// The static analysis report is ready. The sample will remain in this
	// state until a profile is selected.
	SampleStatusStaticAnalysis SampleStatus = "static_analysis"
	// All parameters for sandbox analysis have been selected. The sample is
	// scheduled for running on the sandbox.
	SampleStatusSheduled SampleStatus = "scheduled"
	// The sample is being ran by the sandbox.
	SampleStatusRunning SampleStatus = "running"
	// The sandbox has finished running the sample and the resulting metrics
	// are being processed into reports.
	SampleStatusProcessing SampleStatus = "processing"
	// The sample has reports that can be retrieved. This state is terminal.
	SampleStatusReported SampleStatus = "reported"
	// Analysis of the sample has failed. Any other state may transition into
	// this state. This state is terminal.
	SampleStatusFailed SampleStatus = "failed"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusScheduled  TaskStatus = "scheduled"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusReported   TaskStatus = "reported"
	TaskStatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID     string     `json:"id"`
	Status TaskStatus `json:"status"`
}

type Sample struct {
	ID          string       `json:"id"`
	Private     bool         `json:"private"`
	Status      SampleStatus `json:"status"`
	Kind        string       `json:"kind"`
	Filename    string       `json:"filename"`
	URL         string       `json:"url"`
	Tasks       []Task       `json:"tasks"`
	SubmittedAt time.Time    `json:"submitted"`
	CompletedAt *time.Time   `json:"completed"`
}

type ProfileSelection struct {
	Profile string `json:"profile"`
	Pick    string `json:"pick"`
}

type SampleEvent struct {
	Sample
	Error error
}

func (c *Client) SubmitSampleFile(ctx context.Context, filename string, file io.Reader, interactive bool, profiles []ProfileSelection, password *string) (*Sample, error) {
	var pw string
	if password != nil && *password != "" {
		pw = *password
	}
	request, err := json.Marshal(struct {
		Kind        string             `json:"kind"`
		Interactive bool               `json:"interactive"`
		Profiles    []ProfileSelection `json:"profiles"`
		Password    string             `json:"password,omitempty"`
	}{
		Kind:        "file",
		Interactive: interactive,
		Profiles:    profiles,
		Password:    pw,
	})
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	mw := multipart.NewWriter(w)
	go func() {
		defer func() {
			mw.Close()
			w.Close()
		}()

		jsonField, _ := mw.CreateFormField("_json")
		jsonField.Write(request)

		fileField, _ := mw.CreateFormFile("file", filename)
		io.Copy(fileField, file)
	}()

	req, err := c.newRequest(ctx, http.MethodPost, "/v0/samples", r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	var ret Sample
	err = c.requestJSON(req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (c *Client) SubmitSampleURL(ctx context.Context, url string, interactive bool, profiles []ProfileSelection) (*Sample, error) {
	req := struct {
		Kind        string             `json:"kind"`
		URL         string             `json:"url"`
		Interactive bool               `json:"interactive"`
		Profiles    []ProfileSelection `json:"profiles"`
	}{
		Kind:        "url",
		URL:         url,
		Interactive: interactive,
		Profiles:    profiles,
	}
	var resp Sample
	err := c.jsonRequestJSON(ctx, http.MethodPost, "/v0/samples", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SetSampleProfile(ctx context.Context, sampleID string, profiles []ProfileSelection) error {
	req := struct {
		Auto     bool               `json:"auto"`
		Profiles []ProfileSelection `json:"profiles"`
	}{
		Auto:     false,
		Profiles: profiles,
	}
	return c.jsonRequestJSON(
		ctx, http.MethodPost, "/v0/samples/"+sampleID+"/profile", req, nil,
	)
}

func (c *Client) SetSampleProfileAutomatically(ctx context.Context, sampleID string, pick []string) error {
	req := struct {
		Auto bool     `json:"auto"`
		Pick []string `json:"pick"`
	}{
		Auto: true,
		Pick: pick,
	}
	return c.jsonRequestJSON(
		ctx, http.MethodPost, "/v0/samples/"+sampleID+"/profile", req, nil,
	)
}

type SampleResp struct {
	Data []Sample `json:"data"`
	Next *string  `json:"next"`
}

func (c *Client) samples(ctx context.Context, subset, search *string, max int) <-chan Sample {
	samples := make(chan Sample)
	go func() {
		defer close(samples)
		var response SampleResp
		var counter int
		for {
			limit := max
			if limit > maxLimit {
				limit = maxLimit
			}
			var u string
			if subset != nil {
				u = fmt.Sprint(
					"/v0/samples?subset=", *subset, "&limit=", limit,
				)
			} else {
				u = fmt.Sprint(
					"/v0/search?query=", url.QueryEscape(*search),
					"&limit=", limit,
				)
			}
			if response.Next != nil {
				u += fmt.Sprint("&offset=", *response.Next)
			}
			err := c.jsonRequestJSON(ctx, http.MethodGet, u, nil, &response)
			if err != nil {
				break
			}
			if len(response.Data) == 0 {
				break
			}
			for _, sample := range response.Data {
				samples <- sample
				counter++
				if counter >= max {
					return
				}
			}
			if response.Next == nil {
				break
			}
		}
	}()
	return samples
}

func (c *Client) PublicSamples(ctx context.Context, max int) <-chan Sample {
	subset := "public"
	return c.samples(ctx, &subset, nil, max)
}

func (c *Client) SamplesForUser(ctx context.Context, max int) <-chan Sample {
	subset := "owned"
	return c.samples(ctx, &subset, nil, max)
}

func (c *Client) Search(ctx context.Context, query string, max int) <-chan Sample {
	return c.samples(ctx, nil, &query, max)
}

func (c *Client) SampleByID(ctx context.Context, sampleID string) (*Sample, error) {
	var resp Sample
	err := c.jsonRequestJSON(
		ctx, http.MethodGet, "/v0/samples/"+sampleID, nil, &resp,
	)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteSample(ctx context.Context, sampleID string) error {
	return c.jsonRequestJSON(
		ctx, http.MethodDelete, "/v0/samples/"+sampleID, nil, nil,
	)
}

func (c *Client) SampleEventsByID(ctx context.Context, sampleID string) <-chan SampleEvent {
	out := make(chan SampleEvent, 1)
	go func() {
		defer close(out)

		r, err := c.newRequest(
			ctx, http.MethodGet, "/v0/samples/"+sampleID+"/events", nil,
		)
		if err != nil {
			out <- SampleEvent{Error: err}
			return
		}
		resp, err := c.httpClient.Do(r)
		if err != nil {
			out <- SampleEvent{Error: err}
			return
		}

		if resp.StatusCode > 299 {
			out <- SampleEvent{Error: decodeAPIError(r, resp)}
			return
		}

		dec := json.NewDecoder(resp.Body)
		for {
			var sample Sample
			if err := dec.Decode(&sample); err == io.EOF {
				return
			} else if err != nil {
				out <- SampleEvent{Error: err}
				return
			}
			out <- SampleEvent{Sample: sample}
		}
	}()
	return out
}

func (c *Client) SampleEvents(ctx context.Context) <-chan SampleEvent {
	out := make(chan SampleEvent, 128)
	go func() {
		defer close(out)

		r, err := c.newRequest(ctx, http.MethodGet, "/v0/samples/events", nil)
		if err != nil {
			out <- SampleEvent{Error: err}
			return
		}
		resp, err := c.httpClient.Do(r)
		if err != nil {
			out <- SampleEvent{Error: err}
			return
		}

		if resp.StatusCode > 299 {
			out <- SampleEvent{Error: decodeAPIError(r, resp)}
			return
		}

		dec := json.NewDecoder(resp.Body)
		for {
			var sample Sample
			if err := dec.Decode(&sample); err == io.EOF {
				return
			} else if err != nil {
				out <- SampleEvent{Error: err}
				return
			}
			out <- SampleEvent{Sample: sample}
		}
	}()
	return out
}
