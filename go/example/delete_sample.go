// Copyright (C) 2020 Hatching B.V.
// All rights reserved.

package main

import (
	"context"
	"fmt"

	"github.com/hatching/triage/go"
)

const (
	Triage = "https://api.tria.ge"
	Token  = "<YOUR-APIKEY-HERE>"
)

func main() {
	client := triage.NewClientWithRootURL(Token, Triage)

	// Submit a URL.
	sample, err := client.SubmitSampleURL(
		context.Background(), "http://www.google.com", false, nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("submitted", sample.ID)
	if err := client.DeleteSample(context.Background(), sample.ID); err != nil {
		panic(err)
	}
	fmt.Println("deleted", sample.ID)
}
