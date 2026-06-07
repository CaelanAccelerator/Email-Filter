package main

import (
	"encoding/json"
	"os"
)

func persistVerdict(email Email, verdict Verdict) error {
	record := struct {
		Email   Email
		Verdict Verdict
	}{email, verdict}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	f, err := os.OpenFile("emails.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write(data)

	return err
}
