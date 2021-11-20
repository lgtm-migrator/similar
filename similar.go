package similar

import (
	"math"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wangbin/jiebago"
)

type Similar struct {
	dict   *SimilarWordDict
	memory *SentenceListStore
	seg    *jiebago.Segmenter
}

var SIMILAR_LOGGER = log.WithField("module", "simiar")

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

func VecToSentence(vec SentenceVector, dict *SimilarWordDict) string {
	var parts []string
	for _, wordIndex := range vec {
		w := dict.GetWord(wordIndex)
		if w != nil {
			parts = append(parts, *w)
		}
	}
	return strings.Join(parts, "")
}

// NewSimilar instances
func NewSimilar(memorySize int) *Similar {
	seg := &jiebago.Segmenter{}
	seg.LoadDictionary("dict.txt")
	store, err := NewStore(&StoreConfig{})
	if err != nil {
		SIMILAR_LOGGER.Fatal("store init failed: ", err)
	}
	return &Similar{NewSimilarWordDict(&StoreConfig{}), store, seg}
}

// Remember a sentence text
func (s *Similar) Remember(sentence string) {
	s.memory.Add(SentenceToVec(s.seg, s.dict, sentence))
}

func (s *Similar) FindSimilar(sentence string, threshold float64) (results []*FindResult) {
	if threshold <= 0 {
		return
	}

	setVec := SentenceToVec(s.seg, s.dict, sentence)
	vecs := s.memory.FindAll(func(sv SentenceVector) (bool, float64) {
		similarity := CosDistance(sv, setVec, int(s.dict.GetMaxIndex()))
		if similarity >= threshold {
			return true, similarity
		}
		return false, 0
	})
	for v, sim := range vecs {
		results = append(results, &FindResult{v, sim, s.dict})
	}
	return
}

// FindMostSimilar sentence in memory
func (s *Similar) FindMostSimilar(sentence string) *FindResult {
	setVec := SentenceToVec(s.seg, s.dict, sentence)
	vec, sim := s.memory.FindClosestDistance(func(sv SentenceVector) float64 {
		return CosDistance(sv, setVec, int(s.dict.GetMaxIndex()))
	})
	return &FindResult{vec, sim, s.dict}
}
