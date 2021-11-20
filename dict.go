package similar

import (
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

var DICT_DEFAULT_FILE_NAME = "similar_dict.csv"
var DICT_HEADERS = []string{"word", "index"}

var DICT_LOGGER = log.WithField("module", "dict")

// SimilarWordDict allocate index number to word
type SimilarWordDict struct {
	dict        map[string]int64
	inverseDict map[int64]string
	dictIdx     int64
	lock        *sync.Mutex
	storeConfig *StoreConfig

	persistence *CSVPersistence
}

func NewSimilarWordDict(storeConfig *StoreConfig) (d *SimilarWordDict) {
	d = &SimilarWordDict{
		map[string]int64{},
		map[int64]string{},
		0,
		&sync.Mutex{},
		storeConfig,
		nil,
	}
	if storeConfig.IsPersistence() {

		persistence, err := NewCSVPersistence(
			storeConfig.GetPersistencePath(DICT_DEFAULT_FILE_NAME),
			DICT_HEADERS...,
		)
		if err != nil {
			DICT_LOGGER.Fatal("dict init failed", err)
		}
		d.persistence = persistence

		DICT_LOGGER.Debug("restore dict data...")

		if err = persistence.Restore(d.recover); err != nil {
			DICT_LOGGER.Fatal("dict restore failed", err)
		}

	}

	return
}

func (d *SimilarWordDict) recover(cells ...string) {

	d.lock.Lock()
	defer d.lock.Unlock()

	word := cells[0]
	index, err := strconv.ParseInt(cells[1], 10, 64)

	if err != nil {
		DICT_LOGGER.Error("restore dict record failed", cells)
		return // skip
	}

	d.dict[word] = index
	d.inverseDict[index] = word

	if index > d.dictIdx {
		d.dictIdx = index
	}

}

// GetCode get the index of word
func (dict *SimilarWordDict) GetCode(word string) int64 {
	dict.lock.Lock()
	defer dict.lock.Unlock()
	if _, exist := dict.dict[word]; !exist {
		dict.dict[word] = dict.dictIdx
		dict.inverseDict[dict.dictIdx] = word
		if dict.storeConfig.IsPersistence() {
			dict.persistence.persistWriter.Write([]string{word, strconv.FormatInt(dict.dictIdx, 10)})
		}
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

func (dict *SimilarWordDict) Close() (err error) {
	if dict.storeConfig.IsPersistence() {
		return dict.persistence.Close()
	}
	return
}
