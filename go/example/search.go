// Copyright (C) 2020 Hatching B.V.
// All rights reserved.

package main

import (
	"context"
	"fmt"

	"github.com/hatching/triage/go"
)

const (
	Triage = "https://api.tria.ge/v0"
	Token  = "<YOUR-APIKEY-HERE>"
)

func main() {
	client := triage.NewClientWithRootURL(Token, Triage)

	for sample := range client.Search(context.Background(), "NOT family:emotet", 20) {
		fmt.Println("found ", sample.ID)
	}
}
