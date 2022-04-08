[![GoDoc](https://pkg.go.dev/badge/github.com/hatching/triage)](https://pkg.go.dev/github.com/hatching/triage@v1.0.0/go#pkg-index)

Hatching Triage
===============

This repository features a command-line client and API for interacting with
[Hatching Triage](https://tria.ge/), an automated malware analysis sandbox.

Our official command-line and API client is implemented in both Golang as
well as Python. For a Java library, please see
[TriageApi by Libra](https://github.com/ThisIsLibra/TriageApi).

## Getting Started

### Installing the Go client

```bash
$ go get github.com/hatching/triage/go/cmd/triage
```

### Installing the Python client

```bash
$ pip install hatching-triage
```

### Go Documentation

[Go documentation can be found here](https://pkg.go.dev/github.com/hatching/triage@v1.0.0/go#pkg-index)

## Setting up authentication

When you have installed the client, you need to set the API key so it can
authenticate itself with Triage. You can do so by using the `authenticate`
subcommand:

```bash
$ triage authenticate <API key>
```

You can find the API key on your [account page](https://tria.ge/account). Do
note that you need to have a Researcher account on your user account for a key
to appear.

## Usage

After installing and authentication, various subcommands are available for
interacting with the Hatching Triage API. Note that any submitted samples will
end up in our public cloud, unless otherwise configured, and therefore will
be accessible by everyone browsing to https://tria.ge/

```bash
$ triage -help
Usage of triage:

  authenticate [token] [flags]

    Stores credentials for Triage.

  submit [url/file] [flags]

    Submit a new sample file or URL.

  select-profile [sample]

    Interactively lets you select profiles for samples that have been submitted
    in interactive mode. If an archive file was submitted, you will also be
    prompted to select the files to analyze from the archive.

  list [flags]

    Show the latest samples that have been submitted.

  file [sample] [task] [file] [flags]

    Download task related files.

  archive [sample] [flags]

    Download all task related files as an archive.

  delete [sample]

    Delete a sample.

  report [sample] [flags]

    Query reports for a (finished) analysis.

  create-profile [flags]

  delete-profile [flags]

  list-profiles [flags]
```
