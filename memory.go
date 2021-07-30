package similar

import (
	"math"
	"reflect"
)

type SentenceVector []int64

// ToWordFreqVector
func (sentence *SentenceVector) ToWordFreqVector(dictSize int) []int64 {
	rt := make([]int64, dictSize)
	for _, value := range *sentence {
		if value <= int64(dictSize) {
			rt[value]++
		}
	}
	return rt
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

// FIFOSentenceList cycle override
type FIFOSentenceList struct {
	size    int
	storage []SentenceVector
	idx     int
}

// NewSentenceList instance
func NewSentenceList(size int) *FIFOSentenceList {
	return &FIFOSentenceList{size, make([]SentenceVector, size), 0}
}

func (list *FIFOSentenceList) Add(sentence SentenceVector) {
	if list.idx >= list.size {
		list.idx = list.idx % list.size
	}
	list.storage[list.idx] = sentence
	list.idx++
}

func (list *FIFOSentenceList) Exist(sentence SentenceVector) bool {
	return nil != list.Find(func(sv SentenceVector) bool {
		return reflect.DeepEqual(sentence, sv)
	})
}

func (list *FIFOSentenceList) Find(predictor func(SentenceVector) bool) SentenceVector {
	for _, vec := range list.storage {
		if vec != nil && predictor(vec) {
			return vec
		}
	}
	return nil
}

func (list *FIFOSentenceList) FindClosestDistance(calculator func(SentenceVector) float64) (rt float64) {
	closestValue := float64(0)
	for _, vec := range list.storage {
		if vec != nil {
			currentValue := calculator(vec)
			if currentValue < closestValue || closestValue == 0 {
				closestValue = currentValue
				rt = closestValue
			}
		}
	}
	return
}
