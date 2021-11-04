// Copyright (C) 2020-2021 Hatching B.V.
// All rights reserved.

package main

import (
	"context"
	"fmt"

	triage "github.com/hatching/triage/go"
)

const (
	Triage = "https://api.tria.ge"
	Token  = "<YOUR-APIKEY-HERE>"
)

func main() {
	client := triage.NewClientWithRootURL(Token, Triage)

	// Submit a sample.
	sample, err := client.SubmitSampleURL(
		context.Background(), "http://www.facebook.com", false, nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("submitted:", sample.ID)
	for msg := range client.SampleEventsByID(context.Background(), sample.ID) {
		fmt.Println(msg.Status)
	}
}
