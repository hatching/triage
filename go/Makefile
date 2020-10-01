GOPATH ?= $(HOME)/go
ifeq ("$(PREFIX)","")
BIN=$(GOPATH)/bin
SBIN=$(BIN)
else
BIN=$(PREFIX)/bin
SBIN=$(PREFIX)/sbin
endif

all: triage

triage:
	go build -o "$(BIN)/triage" github.com/hatching/triage/go/cmd/triage

test:
	go test -cover