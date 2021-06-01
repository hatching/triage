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

func (c *Client) SampleOverviewReport(ctx context.Context, sampleID string) (*types.OverviewReport, error) {
	var resp types.OverviewReport
	path := "/v1/samples/" + sampleID + "/overview.json"
	if err := c.jsonRequestJSON(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SampleStaticReport(ctx context.Context, sampleID string) (*types.StaticReport, error) {
	var resp types.StaticReport
	path := "/v0/samples/" + sampleID + "/reports/static"
	if err := c.jsonRequestJSON(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SampleTaskKernelReport(ctx context.Context, sampleID, taskID string) ([]json.RawMessage, error) {
	var resp []json.RawMessage
	overview, err := c.SampleOverviewReport(ctx, sampleID)
	if err != nil {
		return resp, err
	}
	var task *types.TaskSummary
	for _, t := range overview.Tasks {
		if t.Name == taskID {
			task = &t
			break
		}
	}
	if task == nil {
		return resp, fmt.Errorf("Task does not exist")
	}
	var file io.ReadCloser
	if strings.Contains(task.Platform, "windows") {
		if file, err = c.requestRawFile(ctx, "GET", fmt.Sprintf("/v0/samples/%s/%s/logs/onemon.json", sampleID, taskID)); err != nil {
			return resp, err
		}
	} else if strings.Contains(task.Platform, "linux") {
		if file, err = c.requestRawFile(ctx, "GET", fmt.Sprintf("/v0/samples/%s/%s/logs/stahp.json", sampleID, taskID)); err != nil {
			return resp, err
		}
	} else {
		return resp, fmt.Errorf("Platform not supported")
	}
	d := json.NewDecoder(file)
	for {
		var event json.RawMessage
		if err := d.Decode(&event); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		resp = append(resp, event)
	}
	return resp, nil
}

func (c *Client) SampleTaskReport(ctx context.Context, sampleID, taskID string) (*types.TriageReport, error) {
	var resp types.TriageReport
	path := "/v0/samples/" + sampleID + "/" + taskID + "/report_triage.json"
	if err := c.jsonRequestJSON(ctx, http.MethodGet, path, nil, &resp); err != nil {
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
