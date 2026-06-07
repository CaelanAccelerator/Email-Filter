package main

import (
	"strconv"
	"strings"
)

type Verdict struct {
	Valid  bool
	Spam   bool
	Reason []string
}

const SPAM_COUNT = 5

var spamWords = []string{"win", "free", "prize", "money", "offer"}

func classify(email Email) Verdict {
	v := Verdict{Valid: true, Spam: false}
	//rule 1: header checking
	validateHeaders(email, &v)

	//rule 2: content checking
	if len(email.Body) == 0 {
		v.Valid = false
		v.Reason = append(v.Reason, "empty body")
	}

	//rule 3: spam word count checking
	spamCheck(email, &v)

	return v
}

// helper func, validate headers
func validateHeaders(email Email, v *Verdict) {
	if email.Headers["from"] == "" {
		v.Valid = false
		v.Reason = append(v.Reason, "missing From header")
	}
	if email.Headers["to"] == "" {
		v.Valid = false
		v.Reason = append(v.Reason, "missing To header")
	}
	if email.Headers["subject"] == "" {
		v.Valid = false
		v.Reason = append(v.Reason, "missing Subject header")
	}
	if email.Headers["date"] == "" {
		v.Valid = false
		v.Reason = append(v.Reason, "missing Date header")
	}
}

// helper func, check spam using word count
func spamCheck(email Email, v *Verdict) {
	totalCount := 0
	body := strings.ToLower(email.Body)
	for _, word := range spamWords {
		count := strings.Count(body, word)
		if count == 0 {
			continue
		}
		totalCount += count
		v.Reason = append(v.Reason, "Spam word "+word+" appears "+strconv.Itoa(count)+" times")
	}
	if totalCount > SPAM_COUNT {
		v.Spam = true
		v.Reason = append(v.Reason, "Spam words appear "+strconv.Itoa(totalCount)+" times, probably spam")
	}
}
