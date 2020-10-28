all:
	go build -o triage github.com/hatching/triage/go/cmd/triage

test:
	go test ./go
