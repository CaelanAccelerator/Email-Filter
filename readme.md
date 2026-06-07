# go-email-filter

A small email parser and inbound spam filter in Go.
Built to get hands-on with Go and learn how email parsing works at the protocol level.

## What it does
Accepts an email over TCP → parses headers/body with a state machine →
runs a simple spam check → persists each verdict as JSON.

## Run
```bash
go run .
# in another terminal:
nc -N localhost 8080 < testdata/spam.eml
```

## Tests
```bash
go test ./...
```

> Note: the spam check is a naive keyword heuristic — kept intentionally simple.