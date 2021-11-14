package similar

type FindResult struct {
	Vec        *SentenceVector
	Similarity float64
	dict       *SimilarWordDict
}

func (r *FindResult) ToOriginalSentence() string {
	return VecToSentence(*r.Vec, r.dict)
}

func (r *FindResult) GetVector() *SentenceVector {
	return r.Vec
}

func (r *FindResult) GetSimilarity() float64 {
	return r.Similarity
}
