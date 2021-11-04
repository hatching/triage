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

	sample, err := client.SubmitSampleURL(
		context.Background(), "http://www.google.nl", true, nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("sample:", sample.ID)

	err = client.SetSampleProfileAutomatically(
		context.Background(), sample.ID, []string{},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Profile set for", sample.ID)
}
