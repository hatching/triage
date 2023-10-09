// Copyright (C) 2022 Hatching B.V.
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
var timeout = 150
var network = "internet"

func main() {
	client := triage.NewClient(Token)
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
		&timeout,
		&network,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("sample:", sample.ID)
}
