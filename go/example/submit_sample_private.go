// Copyright (C) 2020-2022 Hatching B.V.
// All rights reserved.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	triage "github.com/hatching/triage/go"
)

const (
	Token = "<YOUR-APIKEY-HERE>"
)

var password = "password"
var fname = "some-sample-path"

func main() {
	client := triage.NewPrivateClient(Token)
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalln(err)
	}
	sample, err := client.SubmitSampleFile(
		context.Background(),
		fname,
		f,
		false,
		nil,
		&password,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("sample:", sample.ID)
}
