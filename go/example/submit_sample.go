// Copyright (C) 2020-2021 Hatching B.V.
// All rights reserved.

package main

import (
	"bytes"
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

	sample, err := client.SubmitSampleFile(
		context.Background(), "testfile", bytes.NewBufferString("hello\n"),
		false, nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("sample:", sample.ID)
}
