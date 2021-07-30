package similar

type SimilarWordDict struct {
	dict    map[string]int64
	dictIdx int64
}

func NewSimilarWordDict() *SimilarWordDict {
	return &SimilarWordDict{
		map[string]int64{},
		0,
	}
}

// GetCode of word
func (dict *SimilarWordDict) GetCode(word string) int64 {
	if _, exist := dict.dict[word]; !exist {
		dict.dict[word] = dict.dictIdx
		dict.dictIdx++
	}
	return dict.dict[word]
}

func (dict *SimilarWordDict) GetMaxIndex() int64 {
	return dict.dictIdx
}
