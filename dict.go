package similar

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"sync"
)

var similarDictFileDefaultName = "similar_dict.csv"
var similarDictHeader = []string{"word", "index"}

// SimilarWordDict allocate index number to word
type SimilarWordDict struct {
	dict        map[string]int64
	inverseDict map[int64]string
	dictIdx     int64
	lock        *sync.Mutex
}

func NewSimilarWordDict() *SimilarWordDict {
	return &SimilarWordDict{
		map[string]int64{},
		map[int64]string{},
		0,
		&sync.Mutex{},
	}
}

// GetCode get the index of word
func (dict *SimilarWordDict) GetCode(word string) int64 {
	dict.lock.Lock()
	defer dict.lock.Unlock()
	if _, exist := dict.dict[word]; !exist {
		dict.dict[word] = dict.dictIdx
		dict.inverseDict[dict.dictIdx] = word
		dict.dictIdx++
	}
	return dict.dict[word]
}

// GetWord get the word of index
func (dict *SimilarWordDict) GetWord(index int64) *string {
	if value, exist := dict.inverseDict[index]; exist {
		return &value
	}
	return nil
}

func (dict *SimilarWordDict) GetMaxIndex() int64 {
	return dict.dictIdx
}

// Load dict from file
func (dict *SimilarWordDict) Load(path *string) (err error) {
	dict.lock.Lock()
	defer dict.lock.Unlock()

	if path == nil {
		path = &similarDictFileDefaultName
	}

	storageFile, err := os.OpenFile(
		*path,
		os.O_RDONLY,
		os.ModePerm,
	)

	if err != nil {
		return
	}

	defer storageFile.Close()

	reader := csv.NewReader(storageFile)
	headers, err := reader.Read()

	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("no content in file: %v", storageFile.Name())
		}
		return
	}

	if !reflect.DeepEqual(headers, similarDictHeader) {
		err = fmt.Errorf("the header of target file is not correct: %v, required: %v", headers, similarDictHeader)
		return
	}

	for {
		record, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				// file end
				break
			} else {
				return err
			}
		}

		word := record[0]
		index, err := strconv.ParseInt(record[1], 10, 64)

		if err != nil {
			return err
		}

		dict.dict[word] = index

		if dict.dictIdx < index {
			dict.dictIdx = index
		}

	}

	return nil
}

// Save dict to storage file
func (dict *SimilarWordDict) Save(path *string) (err error) {
	dict.lock.Lock()
	defer dict.lock.Unlock()

	if path == nil {
		path = &similarDictFileDefaultName
	}

	storageFile, err := os.OpenFile(
		*path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		os.ModePerm,
	)

	if err != nil {
		return
	}

	defer storageFile.Close()

	writer := csv.NewWriter(storageFile)

	if err = writer.Write(similarDictHeader); err != nil {
		return
	}

	for word, index := range dict.dict {
		writer.Write([]string{word, strconv.FormatInt(index, 10)})
	}

	writer.Flush()

	return nil
}
