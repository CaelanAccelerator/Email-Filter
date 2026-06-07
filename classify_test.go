package main

import (
	"strings"
	"testing"
)

func TestClassify(t *testing.T) {
	cases := []struct {
		name       string //name of the test case
		email      Email
		wantValid  bool
		wantSpam   bool
		wantReason []string // substrings we expect to find among got.Reason
	}{
		{
			name: "complete email, no spam",
			email: Email{
				Headers: map[string]string{"from": "alice", "to": "bob", "subject": "hi", "date": "2026-06-07"},
				Body:    "just saying hello",
			},
			wantValid:  true,
			wantSpam:   false,
			wantReason: nil,
		},
		{
			name: "missing From header",
			email: Email{
				Headers: map[string]string{"to": "bob", "subject": "hi", "date": "2026-06-07"},
				Body:    "hello",
			},
			wantValid:  false,
			wantSpam:   false,
			wantReason: []string{"missing From header"},
		},
		{
			name: "missing multiple headers",
			email: Email{
				Headers: map[string]string{"from": "alice"},
				Body:    "hello",
			},
			wantValid:  false,
			wantSpam:   false,
			wantReason: []string{"missing To header", "missing Subject header", "missing Date header"},
		},
		{
			name: "empty body",
			email: Email{
				Headers: map[string]string{"from": "alice", "to": "bob", "subject": "hi", "date": "2026-06-07"},
				Body:    "",
			},
			wantValid:  false,
			wantSpam:   false,
			wantReason: []string{"empty body"},
		},
		{
			name: "spam words push count over threshold",
			email: Email{
				Headers: map[string]string{"from": "alice", "to": "bob", "subject": "hi", "date": "2026-06-07"},
				Body:    "win win win free free free prize prize prize",
			},
			wantValid:  true,
			wantSpam:   true,
			wantReason: []string{"Spam word win appears 3 times", "probably spam"},
		},
	}

	for _, c := range cases {
		got := classify(c.email)
		if got.Valid != c.wantValid {
			t.Errorf("%s: Valid = %t, want %t", c.name, got.Valid, c.wantValid)
		}
		if got.Spam != c.wantSpam {
			t.Errorf("%s: Spam = %t, want %t", c.name, got.Spam, c.wantSpam)
		}
		for _, want := range c.wantReason {
			if !reasonContains(got.Reason, want) {
				t.Errorf("%s: Reason = %q, want it to contain %q", c.name, got.Reason, want)
			}
		}
	}
}

func reasonContains(reasons []string, substr string) bool {
	for _, r := range reasons {
		if strings.Contains(r, substr) {
			return true
		}
	}
	return false
}
