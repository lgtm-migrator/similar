package similar

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

// SentenceVector convert sentence to an words index array
type SentenceVector []int64

// ToWordFreqVector, convert word vec to word freq vec
func (sentence *SentenceVector) ToWordFreqVector(dictSize int) []int64 {
	rt := make([]int64, dictSize)
	for _, value := range *sentence {
		if value <= int64(dictSize) {
			rt[value]++
		}
	}
	return rt
}

func (sentence *SentenceVector) ToString() string {
	var content bytes.Buffer
	gob.NewEncoder(&content).Encode(sentence)                 // encode to bytes
	return base64.URLEncoding.EncodeToString(content.Bytes()) // encode to base64 string
}

func NewSentenceVecFromBase64(encoded string) (vec SentenceVector, err error) {
	sevbytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return vec, gob.NewDecoder(bytes.NewReader(sevbytes)).Decode(&vec)
}
