// Copyright (C) 2019-2020 Hatching B.V.
// All rights reserved.

package triage

import (
	"context"
	"io"
	"net/http"

	"github.com/hatching/triage/types"
)

func (c *Client) SampleOverviewReport(ctx context.Context, sampleID string) (*types.OverviewReport, error) {
	var resp types.OverviewReport
	path := "/v0/samples/" + sampleID + "/overview.json"
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
