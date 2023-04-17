GOPATH ?= $(HOME)/go

ifeq ("$(PREFIX)","")
	BIN=$(GOPATH)/bin
	SBIN=$(BIN)
else
	BIN=$(PREFIX)/bin
	SBIN=$(PREFIX)/sbin
endif

all:
	go build -o "$(BIN)/triage" "$(CURDIR)/cmd/triage"

test:
	go test -cover