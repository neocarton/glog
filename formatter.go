package glog

import (
	"encoding/json"
	"time"
)

// ToISOTime Convert time to ISO fromat
func ToISOTime(input interface{}, format string) string {
	if input == nil {
		return ""
	}
	t, ok := input.(time.Time)
	if !ok {
		return "!invalid"
	}
	return t.Format(time.RFC3339)
}

// ToJSON Convert to JSON
func ToJSON(input interface{}, format string) string {
	if input == nil {
		return ""
	}
	data, err := json.Marshal(input)
	if err != nil {
		return "!invalid"
	}
	return string(data)
}
