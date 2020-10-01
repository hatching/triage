// Copyright (C) 2019-2020 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
	TaskStatusSheduled   TaskStatus = "scheduled"
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
	SubmittedAt time.Time    `json:"completed"`
	CompletedAt *time.Time   `json:"submitted"`
}

type ProfileSelection struct {
	Profile string `json:"profile"`
	Pick    string `json:"pick"`
}

type SampleEvent struct {
	Sample
	Error error
}

func (c *Client) SubmitSampleFile(ctx context.Context, filename string, file io.Reader, interactive bool, profiles []ProfileSelection) (*Sample, error) {
	jsonReq, err := json.Marshal(struct {
		Kind        string             `json:"kind"`
		Interactive bool               `json:"interactive"`
		Profiles    []ProfileSelection `json:"profiles"`
	}{
		Kind:        "file",
		Interactive: interactive,
		Profiles:    profiles,
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
		_, _ = jsonField.Write(jsonReq)
		fileField, _ := mw.CreateFormFile("file", filename)
		_, _ = io.Copy(fileField, file)
	}()

	req, err := c.newRequest(ctx, http.MethodPost, "/v0/samples", r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	var respSample Sample
	if err := c.requestJSON(req, &respSample); err != nil {
		return nil, err
	}
	return &respSample, nil
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
	if err := c.jsonRequestJSON(ctx, http.MethodPost, "/v0/samples", req, &resp); err != nil {
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
	return c.jsonRequestJSON(ctx, http.MethodPost, "/v0/samples/"+sampleID+"/profile", req, nil)
}

func (c *Client) SetSampleProfileAutomatically(ctx context.Context, sampleID string, pick []string) error {
	req := struct {
		Auto bool     `json:"auto"`
		Pick []string `json:"pick"`
	}{
		Auto: true,
		Pick: pick,
	}
	return c.jsonRequestJSON(ctx, http.MethodPost, "/v0/samples/"+sampleID+"/profile", req, nil)
}

type SampleResp struct {
	Data []Sample `json:"data"`
	Next string   `json:"next"`
}

func (c *Client) samples(ctx context.Context, subset string, max int) <-chan Sample {
	if max > maxLimit {
		max = maxLimit
	}

	samples := make(chan Sample)
	go func() {
		var response SampleResp
		var counter int
		for {
			url := fmt.Sprintf("/v0/samples?subset=%s&limit=%v", subset, max)
			if response.Next != "" {
				url = fmt.Sprintf("%s&offset=%s", url, response.Next)
			}
			if err := c.jsonRequestJSON(ctx, http.MethodGet, url, nil, &response); err != nil {
				panic(err)
			}
			if len(response.Data) == 0 {
				break
			}
			for _, sample := range response.Data {
				samples <- sample
				counter++
				if counter >= max {
					close(samples)
					return
				}
			}
		}
		close(samples)
	}()
	return samples
}

func (c *Client) PublicSamples(ctx context.Context, max int) <-chan Sample {
	return c.samples(ctx, "public", max)
}

func (c *Client) SamplesForUser(ctx context.Context, max int) <-chan Sample {
	return c.samples(ctx, "owned", max)
}

func (c *Client) SampleByID(ctx context.Context, sampleID string) (*Sample, error) {
	var resp Sample
	if err := c.jsonRequestJSON(ctx, http.MethodGet, "/v0/samples/"+sampleID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) DeleteSample(ctx context.Context, sampleID string) error {
	return c.jsonRequestJSON(ctx, http.MethodDelete, "/v0/samples/"+sampleID, nil, nil)
}

func (c *Client) SampleEventsByID(ctx context.Context, sampleID string) <-chan SampleEvent {
	out := make(chan SampleEvent, 1)
	go func() {
		defer close(out)

		r, err := c.newRequest(ctx, http.MethodGet, "/v0/samples/"+sampleID+"/events", nil)
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
