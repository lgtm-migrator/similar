package similar

import (
	"math"

	"github.com/wangbin/jiebago"
)

type Similar struct {
	dict   *SimilarWordDict
	memory *SentenceListStore
	seg    *jiebago.Segmenter
}

func CosDistance(sv SentenceVector, anothersv SentenceVector, dictSize int) float64 {
	wordFreq1 := sv.ToWordFreqVector(dictSize)
	wordFreq2 := anothersv.ToWordFreqVector(dictSize)
	topSum := int64(0)
	for idx, value := range wordFreq1 {
		topSum += (value * wordFreq2[idx])
	}
	bottomLeftSum := int64(0)
	for _, value := range wordFreq1 {
		bottomLeftSum += int64(math.Pow(float64(value), 2))
	}
	bottomRightSum := int64(0)
	for _, value := range wordFreq2 {
		bottomRightSum += int64(math.Pow(float64(value), 2))
	}
	return (float64(topSum) / (math.Sqrt(float64(bottomLeftSum) * float64(bottomRightSum))))
}

func SentenceToVec(seg *jiebago.Segmenter, dict *SimilarWordDict, sentence string) (vec SentenceVector) {
	for word := range seg.CutAll(sentence) {
		code := dict.GetCode(word)
		vec = append(vec, code)
	}
	return
}

// NewSimilar instances
func NewSimilar(memorySize int) *Similar {
	seg := &jiebago.Segmenter{}
	seg.LoadDictionary("dict.txt")
	return &Similar{NewSimilarWordDict(), NewSentenceList(memorySize), seg}
}

// Compare sentence in memory, return the distence
func (s *Similar) Compare(sentence string) (rt float64) {
	setVec := SentenceToVec(s.seg, s.dict, sentence)

	rt = s.memory.FindClosestDistance(func(sv SentenceVector) float64 {
		return CosDistance(sv, setVec, int(s.dict.GetMaxIndex()))
	})

	s.memory.Add(setVec)

	return
}
