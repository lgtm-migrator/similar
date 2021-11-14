package similar

import (
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSimilarWordDict(t *testing.T) {

	assert := assert.New(t)
	dict := NewSimilarWordDict()

	assert.EqualValues(0, dict.GetCode("hello"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(2, dict.GetCode("✅"))

	assert.Nil(dict.Save(nil))

	assert.Nil(os.Remove(similarDictFileDefaultName))

}

func TestLoadAndRestoreDict(t *testing.T) {

	assert := assert.New(t)
	dict := NewSimilarWordDict()

	assert.EqualValues(0, dict.GetCode("hello"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(2, dict.GetCode("✅"))

	tmpFile := path.Join(os.TempDir(), uuid.NewString()+".csv") // tmpFile
	assert.Nil(dict.Save(&tmpFile))

	anotherDict := NewSimilarWordDict()
	assert.Nil(anotherDict.Load(&tmpFile))

	assert.EqualValues(2, anotherDict.GetCode("✅"))
	assert.EqualValues(1, anotherDict.GetCode("word"))
	assert.EqualValues(0, anotherDict.GetCode("hello"))
	assert.EqualValues(1, anotherDict.GetCode("word"))

}
