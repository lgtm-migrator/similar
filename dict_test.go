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
	dict := NewSimilarWordDict(&StoreConfig{Persistence: true})

	assert.EqualValues(0, dict.GetCode("hello"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(2, dict.GetCode("✅"))

	assert.Nil(dict.Close())

	_, err := os.Stat(DICT_DEFAULT_FILE_NAME)

	assert.Nil(err)

	assert.Nil(os.Remove(DICT_DEFAULT_FILE_NAME))

}

func TestLoadAndRestoreDict(t *testing.T) {

	assert := assert.New(t)
	tmpFile := path.Join(os.TempDir(), uuid.NewString()+".csv") // tmpFile
	dict := NewSimilarWordDict(&StoreConfig{Persistence: true, PersistencePath: tmpFile})

	assert.EqualValues(0, dict.GetCode("hello"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(1, dict.GetCode("word"))
	assert.EqualValues(2, dict.GetCode("✅"))

	assert.Nil(dict.Close())
	defer os.Remove(tmpFile)

	anotherDict := NewSimilarWordDict(&StoreConfig{Persistence: true, PersistencePath: tmpFile})

	assert.EqualValues(2, anotherDict.GetCode("✅"))
	assert.EqualValues(1, anotherDict.GetCode("word"))
	assert.EqualValues(0, anotherDict.GetCode("hello"))
	assert.EqualValues(1, anotherDict.GetCode("word"))

}
