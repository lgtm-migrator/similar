package similar

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
)

// SentenceListStore cycle override
type SentenceListStore struct {
	storage []*SentenceVector
	idx     int
	lock    *sync.Mutex
}

// NewSentenceList instance
func NewSentenceList(initialSize int) *SentenceListStore {
	if initialSize <= 0 {
		initialSize = 1
	}
	return &SentenceListStore{make([]*SentenceVector, initialSize), 0, &sync.Mutex{}}
}

func (list *SentenceListStore) extendCapacity() {
	capacity := len(list.storage)
	if (list.idx + 1) >= capacity {
		// allocate more resources
		list.storage = append(list.storage, make([]*SentenceVector, (2*capacity)+1)...)
	}
}

func (list *SentenceListStore) Add(vec SentenceVector) {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.internalAdd(vec)
}

func (list *SentenceListStore) internalAdd(vec SentenceVector) {
	list.extendCapacity()
	list.storage[list.idx] = &vec
	list.idx++
}

var sentenceStoreDefaultName = "sentences_vec.csv"
var sentenceStoreDefaultHeaders = []string{"index", "encoded_vec"}

func (list *SentenceListStore) Save(path *string) error {
	list.lock.Lock()
	defer list.lock.Unlock()

	if path == nil {
		path = &sentenceStoreDefaultName
	}

	storageFile, err := os.OpenFile(
		*path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		os.ModePerm,
	)

	if err != nil {
		return nil
	}

	defer storageFile.Close()

	writer := csv.NewWriter(storageFile)

	if err = writer.Write(sentenceStoreDefaultHeaders); err != nil {
		return nil
	}

	for index, vec := range list.storage {
		if vec == nil {
			break
		}
		writer.Write([]string{fmt.Sprintf("%v", index), vec.ToString()})
	}

	writer.Flush()

	return nil
}

func (list *SentenceListStore) Load(path *string) (err error) {
	list.lock.Lock()
	defer list.lock.Unlock()

	if path == nil {
		path = &sentenceStoreDefaultName
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

	if !reflect.DeepEqual(headers, sentenceStoreDefaultHeaders) {
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

		vec, err := NewSentenceVecFromBase64(record[1])

		if err != nil {
			return err
		}

		list.internalAdd(vec)

	}

	return nil
}

func (list *SentenceListStore) Exist(sentence SentenceVector) bool {
	return nil != list.Find(func(sv SentenceVector) bool {
		return reflect.DeepEqual(sentence, sv)
	})
}

func (list *SentenceListStore) Find(predictor func(SentenceVector) bool) SentenceVector {
	for _, vec := range list.storage {
		if vec != nil && predictor(*vec) {
			return *vec
		}
	}
	return nil
}

func (list *SentenceListStore) FindClosestDistance(calculator func(SentenceVector) float64) (rt float64) {
	closestValue := float64(0)
	for _, vec := range list.storage {
		if vec != nil {
			currentValue := calculator(*vec)
			if currentValue < closestValue || closestValue == 0 {
				closestValue = currentValue
				rt = closestValue
			}
		}
	}
	return
}
