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

	for sample := range client.SamplesForUser(context.Background(), 20) {
		client.DeleteSample(context.Background(), sample.ID)
		fmt.Println("deleted:", sample.ID)
	}
}
