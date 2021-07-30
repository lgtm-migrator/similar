package similar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSentenceList(t *testing.T) {
	assert := assert.New(t)

	list := NewSentenceList(5)

	list.Add(SentenceVector{0, 0, 1})
	assert.Equal(1, list.storage[0][2])
	list.Add(SentenceVector{0, 0, 2})
	assert.Equal(2, list.storage[1][2])
	list.Add(SentenceVector{0, 0, 3})
	assert.Equal(3, list.storage[2][2])
	list.Add(SentenceVector{0, 0, 4})
	assert.Equal(4, list.storage[3][2])
	assert.Nil(list.storage[4])
	list.Add(SentenceVector{0, 0, 5})
	assert.Equal(5, list.storage[4][2])

	list.Add(SentenceVector{0, 0, 6})
	assert.Equal(2, list.storage[1][2])
	assert.Equal(3, list.storage[2][2])
	assert.Equal(4, list.storage[3][2])
	assert.Equal(5, list.storage[4][2])
	assert.Equal(6, list.storage[0][2])

	assert.NotNil(list.Find(func(sv SentenceVector) bool {
		return sv[2] == 3
	}))

	assert.Nil(list.Find(func(sv SentenceVector) bool {
		return sv[2] == 7
	}))

	assert.True(list.Exist(SentenceVector{0, 0, 2}))
	assert.False(list.Exist(SentenceVector{0, 0, 1023}))

}

func TestSentenceVector_ToWordFreqVector(t *testing.T) {
	assert := assert.New(t)

	set := SentenceVector{9, 3, 3, 3, 3, 5}
	assert.Equal([]int64{0, 0, 0, 4, 0, 1, 0, 0, 0, 1}, set.ToWordFreqVector(10))

}
