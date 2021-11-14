package main

type RememberPayload struct {
	Sentences []string `json:"sentences"`
}

type FindPayload struct {
	Sentence  string  `json:"sentence"`
	Threshold float64 `json:"threshold"`
}
type FindOnePayload struct {
	Sentence string `json:"sentence"`
}
