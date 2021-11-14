package main

import "github.com/Soontao/similar"

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

type FindResult struct {
	Sentence   string  `json:"sentence"`
	Similarity float64 `json:"similarity"`
}

func ToFindResult(s *similar.FindResult) *FindResult {
	if s == nil {
		return nil
	}
	rt := &FindResult{}
	rt.Similarity = s.Similarity
	rt.Sentence = s.ToOriginalSentence()
	return rt
}

func ToFindResults(ss []*similar.FindResult) (rs []*FindResult) {
	for _, s := range ss {
		rs = append(rs, ToFindResult(s))
	}
	return
}
