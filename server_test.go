package main

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	const dbFile = "emails.json"
	os.Remove(dbFile)
	defer os.Remove(dbFile)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("Listening error: ", err)
		return
	}

	addr := ln.Addr().String()
	defer ln.Close()
	go serve(ln)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal("Connection error: ", err)
		return
	}

	email := "From: alice@x.com\r\nSubject: hi\r\n\r\nHello\r\n"
	conn.Write([]byte(email))
	conn.(*net.TCPConn).CloseWrite()

	resp, err := io.ReadAll(conn)
	if err != nil {
		t.Fatal("Reading response Error: ", err)
		return
	}

	got := string(resp)
	if !strings.Contains(got, "alice@x.com") {
		t.Errorf("got %q, want it to contain alice@x.com", got)
	}
	if !strings.Contains(got, "hi") {
		t.Errorf("got %q, want it to contain hi", got)
	}

	// handleConn persists a verdict for every parsed email; verify it landed.
	data, err := os.ReadFile(dbFile)
	if err != nil {
		t.Fatalf("reading %s: %v", dbFile, err)
	}

	var record struct {
		Email   Email
		Verdict Verdict
	}
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatalf("unmarshal %s: %v", dbFile, err)
	}

	if record.Verdict.Valid {
		t.Errorf("Verdict.Valid = true, want false (email is missing To/Date headers)")
	}
	for _, want := range []string{"missing To header", "missing Date header"} {
		if !reasonContains(record.Verdict.Reason, want) {
			t.Errorf("Verdict.Reason = %q, want it to contain %q", record.Verdict.Reason, want)
		}
	}
}
