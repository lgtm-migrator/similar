package similar

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSentenceList(t *testing.T) {
	assert := assert.New(t)

	list := NewSentenceList(5)
	assert.EqualValues(5, len(list.storage))
	list.Add(SentenceVector{0, 0, 1})
	assert.EqualValues(1, (*list.storage[0])[2])
	list.Add(SentenceVector{0, 0, 2})
	assert.EqualValues(2, (*list.storage[1])[2])
	list.Add(SentenceVector{0, 0, 3})
	assert.EqualValues(3, (*list.storage[2])[2])
	list.Add(SentenceVector{0, 0, 4})
	assert.EqualValues(4, (*list.storage[3])[2])
	assert.Nil(list.storage[4])
	assert.EqualValues(5, len(list.storage))
	list.Add(SentenceVector{0, 0, 5})
	assert.EqualValues(16, len(list.storage))
	assert.EqualValues(5, (*list.storage[4])[2])

	list.Add(SentenceVector{0, 0, 6})
	assert.EqualValues(6, (*list.storage[5])[2])

	assert.NotNil(list.Find(func(sv SentenceVector) bool {
		return sv[2] == 3
	}))

	assert.Nil(list.Find(func(sv SentenceVector) bool {
		return sv[2] == 7
	}))

	assert.True(list.Exist(SentenceVector{0, 0, 2}))
	assert.False(list.Exist(SentenceVector{0, 0, 1023}))

}

func TestStoreSaveLoad(t *testing.T) {
	assert := assert.New(t)
	list := NewSentenceList(5)
	for i := 0; i < 1000; i++ {
		list.Add(SentenceVector{0, 0, int64(i)})
	}
	assert.Nil(list.Save(nil))

	anotherList := NewSentenceList(1000)
	assert.Nil(anotherList.Load(nil))
	assert.EqualValues(&SentenceVector{0, 0, 999}, anotherList.storage[999])

	assert.Nil(os.Remove(sentenceStoreDefaultName))
}

func TestSentenceVector_ToWordFreqVector(t *testing.T) {
	assert := assert.New(t)

	set := SentenceVector{9, 3, 3, 3, 3, 5}
	assert.EqualValues([]int64{0, 0, 0, 4, 0, 1, 0, 0, 0, 1}, set.ToWordFreqVector(10))

}
