// Copyright (C) 2019-2021 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hatching/triage/types"
)

func (c *Client) SamplePath(ctx context.Context, sampleID, path string, resp interface{}) error {
	path = "/v0/samples/" + sampleID + path
	return c.jsonRequestJSON(ctx, http.MethodGet, path, nil, &resp)
}

func (c *Client) SampleSample(ctx context.Context, sampleID string) (io.ReadCloser, error) {
	path := "/v0/samples/" + sampleID + "/sample"
	return c.requestRawFile(ctx, http.MethodGet, path)
}

func (c *Client) SampleOverviewReport(ctx context.Context, sampleID string) (*types.OverviewReport, error) {
	var resp types.OverviewReport
	path := "/v1/samples/" + sampleID + "/overview.json"
	err := c.jsonRequestJSON(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SampleStaticReport(ctx context.Context, sampleID string) (*types.StaticReport, error) {
	var resp types.StaticReport
	path := "/v0/samples/" + sampleID + "/reports/static"
	err := c.jsonRequestJSON(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SampleTaskKernelReport(ctx context.Context, sampleID, taskID string) (<-chan json.RawMessage, error) {
	overview, err := c.SampleOverviewReport(ctx, sampleID)
	if err != nil {
		return nil, err
	}

	var task *types.TaskSummary
	for _, t := range overview.Tasks {
		if t.Name == taskID {
			task = &t
			break
		}
	}
	if task == nil {
		return nil, fmt.Errorf("Task does not exist")
	}

	var proto string
	switch {
	case strings.Contains(task.OS, "windows"),
		strings.Contains(task.Platform, "windows"):
		proto = "onemon.json"
	case strings.Contains(task.OS, "linux"),
		strings.Contains(task.OS, "ubuntu"),
		strings.Contains(task.Platform, "linux"),
		strings.Contains(task.Platform, "ubuntu"):
		proto = "stahp.json"
	case strings.Contains(task.OS, "macos"),
		strings.Contains(task.Platform, "macos"):
		proto = "bigmac.json"
	case strings.Contains(task.OS, "android"),
		strings.Contains(task.Platform, "android"):
		proto = "droidy.json"
	default:
		return nil, fmt.Errorf("Platform not supported")
	}

	file, err := c.requestRawFile(ctx, "GET", fmt.Sprint(
		"/v0/samples/", sampleID, "/", taskID, "/logs/", proto,
	))
	if err != nil {
		return nil, err
	}

	ret := make(chan json.RawMessage)
	go func() {
		defer close(ret)
		d := json.NewDecoder(file)
		for {
			var event json.RawMessage
			err := d.Decode(&event)
			if err != nil {
				break
			}
			ret <- event
		}
	}()
	return ret, nil
}

func (c *Client) SampleTaskReport(ctx context.Context, sampleID, taskID string) (*types.TriageReport, error) {
	var resp types.TriageReport
	path := "/v0/samples/" + sampleID + "/" + taskID + "/report_triage.json"
	err := c.jsonRequestJSON(ctx, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SampleTaskFile(ctx context.Context, sampleID, taskID, filename string) (io.ReadCloser, error) {
	path := "/v0/samples/" + sampleID + "/" + taskID + "/" + filename
	return c.requestRawFile(ctx, http.MethodGet, path)
}

func (c *Client) SampleArchiveTAR(ctx context.Context, sampleID string) (io.ReadCloser, error) {
	path := "/v0/samples/" + sampleID + "/archive"
	return c.requestRawFile(ctx, http.MethodGet, path)
}

func (c *Client) SampleArchiveZIP(ctx context.Context, sampleID string) (io.ReadCloser, error) {
	path := "/v0/samples/" + sampleID + "/archive.zip"
	return c.requestRawFile(ctx, http.MethodGet, path)
}

func (c *Client) SampleTaskPCAP(ctx context.Context, sampleID, taskID string) (io.ReadCloser, error) {
	return c.SampleTaskFile(ctx, sampleID, taskID, "dump.pcap")
}

func (c *Client) SampleTaskPCAPNG(ctx context.Context, sampleID, taskID string) (io.ReadCloser, error) {
	return c.SampleTaskFile(ctx, sampleID, taskID, "dump.pcapng")
}

func (c *Client) SampleURLScanScreenshot(ctx context.Context, sampleID string) (io.ReadCloser, error) {
	return c.SampleTaskFile(ctx, sampleID, "urlscan1", "screenshot.png")
}
