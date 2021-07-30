package similar

import (
	"github.com/wangbin/jiebago"
)

type Similar struct {
	dict   *SimilarWordDict
	memory *FIFOSentenceList
	seg    *jiebago.Segmenter
}

// NewSimilar instances
func NewSimilar(memorySize int) *Similar {
	seg := &jiebago.Segmenter{}
	seg.LoadDictionary("dict.txt")
	return &Similar{NewSimilarWordDict(), NewSentenceList(memorySize), seg}
}

// Compare sentence in memory, return the distence
func (s *Similar) Compare(sentence string) (rt float64) {
	setVec := SentenceVector{}
	for word := range s.seg.CutAll(sentence) {
		code := s.dict.GetCode(word) // get/create code to dict
		setVec = append(setVec, code)
	}

	rt = s.memory.FindClosestDistance(func(sv SentenceVector) float64 {
		return sv.CosDistance(setVec)
	})

	s.memory.Add(setVec)

	return
}
