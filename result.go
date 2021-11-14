package similar

type FindResult struct {
	vec        *SentenceVector
	similarity float64
	dict       *SimilarWordDict
}

func (r *FindResult) ToOriginalSentence() string {
	return VecToSentence(*r.vec, r.dict)
}

func (r *FindResult) GetVector() *SentenceVector {
	return r.vec
}

func (r *FindResult) GetSimilarity() float64 {
	return r.similarity
}
