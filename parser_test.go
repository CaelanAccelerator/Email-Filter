package main

import (
	"strings"
	"testing"
)

func TestTrimCR(t *testing.T) {
	cases := []struct {
		name string //name of the test case
		in   string // input
		want string // expect
	}{
		{name: "trailing CR removed", in: "hello\r", want: "hello"},
		{name: "case empty", in: "", want: ""},
		{name: "case no \r", in: "hello", want: "hello"},
		{name: "mutiple \r", in: "hello \r\r", want: "hello "},
	}

	for _, c := range cases {
		got := trimCR(c.in)
		if got != c.want {
			t.Errorf("%s: trimCR(%q) = %q, want %q", c.name, c.in, got, c.want)
		}
	}
}

func TestSplitHeader(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantKey string
		wantVal string
		wantOK  bool
	}{
		{name: "one \" : \"", in: "Subject: hello", wantKey: "subject", wantVal: "hello", wantOK: true},
		{name: "mutiple \" : \"", in: "Subject: :hello", wantKey: "subject", wantVal: ":hello", wantOK: true},
		{name: "no space between Key and Value", in: "From:hello", wantKey: "from", wantVal: "hello", wantOK: true},
		{name: "garbage line", in: "garbage line", wantKey: "", wantVal: "", wantOK: false},
		//{name: "empty line", in: "", wantKey: "", wantVal: "", wantOK: false},
		{name: "empty val", in: "Subject:", wantKey: "subject", wantVal: "", wantOK: true},
	}

	for _, c := range cases {
		key, val, ok := splitHeader(c.in)
		if key != c.wantKey {
			t.Errorf("%s: the key of splitHeader(%q) = %q, want %q", c.name, c.in, key, c.wantKey)
		}
		if val != c.wantVal {
			t.Errorf("%s: the val of splitHeader(%q) = %q, want %q", c.name, c.in, val, c.wantVal)
		}
		if ok != c.wantOK {
			t.Errorf("%s: the ok of splitHeader(%q) = %t, want %t", c.name, c.in, ok, c.wantOK)
		}
	}
}
func TestParseEmailStream(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		wantFrom string
		wantBody string
	}{
		{
			name:     "complete email",
			in:       "From: alice\r\nSubject: hi\r\n\r\nHello\r\n",
			wantFrom: "alice",
			wantBody: "Hello",
		},
		{
			name:     "headers only, no body",
			in:       "From: alice\r\n\r\n",
			wantFrom: "alice",
			wantBody: "",
		},
		{
			name:     "no blank line, all headers",
			in:       "From: alice\r\nTo: bob\r\n",
			wantFrom: "alice",
			wantBody: "",
		},
		{
			name:     "leading blank line, body only",
			in:       "\r\njust body\r\n",
			wantFrom: "",
			wantBody: "just body",
		},
		{
			name:     "body with blank line",
			in:       "\r\nline1\n\nline2",
			wantFrom: "",
			wantBody: "line1\n\nline2",
		},
	}

	for _, c := range cases {
		got := parseEmailStream(strings.NewReader(c.in))
		if got.Headers["from"] != c.wantFrom {
			t.Errorf("%s: from = %q, want %q", c.name, got.Headers["from"], c.wantFrom)
		}
		if got.Body != c.wantBody {
			t.Errorf("%s: body = %q, want %q", c.name, got.Body, c.wantBody)
		}
	}
}
