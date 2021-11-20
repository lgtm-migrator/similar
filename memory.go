package similar

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"sync"
)

// SentenceListStore cycle override
type SentenceListStore struct {
	storage       []*SentenceVector
	idx           int
	lock          *sync.Mutex
	config        *StoreConfig
	persistWriter *csv.Writer
	persistFile   *os.File
	persistence   *CSVPersistence
}

// NewStore instance
func NewStore(c *StoreConfig) (s *SentenceListStore, err error) {

	s = &SentenceListStore{
		make([]*SentenceVector, c.GetInitialSize()),
		0,
		&sync.Mutex{},
		c,
		nil,
		nil,
		nil,
	}

	if c.IsPersistence() {
		STORE_LOGGER.Debug("init store with persistence")
		s.persistence, err = NewCSVPersistence(
			c.GetPersistencePath(sentenceStoreDefaultName),
			sentenceStoreDefaultHeaders...,
		)
		if err != nil {
			return nil, err
		}
		err = s.persistence.Restore(func(cells ...string) {
			vec, err := NewSentenceVecFromBase64(cells[1])
			if err != nil {
				STORE_LOGGER.Error("load data item failed", cells, err)
			}
			s.memoryAdd(vec)
		})
		if err != nil {
			return
		}
	}

	return
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
	idx := list.memoryAdd(vec)
	if list.config.IsPersistence() {
		list.persistence.Append(fmt.Sprintf("%v", idx), vec.ToString())
	}
}

func (list *SentenceListStore) memoryAdd(vec SentenceVector) (idx int) {
	list.extendCapacity()
	list.storage[list.idx] = &vec
	idx = list.idx
	list.idx++
	return
}

var sentenceStoreDefaultName = "sentences_vec.csv"
var sentenceStoreDefaultHeaders = []string{"index", "encoded_vec"}

// Close storage, if persistence required, flush all data
func (list *SentenceListStore) Close() (err error) {
	STORE_LOGGER.Debug("closing...")
	if list.config.IsPersistence() {
		if list.persistence != nil {
			STORE_LOGGER.Debug("flush persistence data")
			err = list.persistence.Close()
			if err != nil {
				return
			}
		}
	}
	return
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

func (list *SentenceListStore) FindAll(calculator func(SentenceVector) (bool, float64)) map[*SentenceVector]float64 {
	rs := map[*SentenceVector]float64{}
	for _, vec := range list.storage {
		if vec != nil {
			if put, similarity := calculator(*vec); put {
				rs[vec] = similarity
			}
		}
	}
	return rs
}

func (list *SentenceListStore) FindClosestDistance(calculator func(SentenceVector) float64) (
	sVec *SentenceVector,
	similarity float64,
) {
	closestValue := float64(0)
	for _, vec := range list.storage {
		if vec != nil {
			currentValue := calculator(*vec)
			if currentValue < closestValue || closestValue == 0 {
				closestValue = currentValue
				sVec = vec
				similarity = closestValue
			}
		}
	}
	return
}
