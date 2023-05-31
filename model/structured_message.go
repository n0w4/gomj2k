package model

import "encoding/json"

type StructuredMessage struct {
	Topic   string          `json:"topic,omitempty"`
	Key     json.RawMessage `json:"key,omitempty"`
	Payload json.RawMessage `json:"payload"`
	Headers []HeaderContent `json:"headers,omitempty"`
}

type HeaderContent struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}
