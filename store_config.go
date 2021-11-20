package similar

import log "github.com/sirupsen/logrus"

type StoreConfig struct {
	InitialSize     int
	MemorySize      int
	MemoryDays      int
	Persistence     bool
	PersistencePath string
}

var STORE_LOGGER = log.WithField("module", "similar_store")

const DEFAULT_INIT_SIZE = 10000
const DEFAULT_MEMORY_SIZE = 1000 * 1000
const DEFAULT_MEMORY_DAYS = 180

func (c *StoreConfig) GetInitialSize() int {
	if c.InitialSize >= 1 {
		return c.InitialSize
	}
	return DEFAULT_INIT_SIZE
}

func (c *StoreConfig) GetMemorySize() int {
	if c.MemorySize > 0 {
		return c.MemorySize
	}
	return DEFAULT_MEMORY_SIZE
}

func (c *StoreConfig) GetMemoryDays() int {
	if c.MemoryDays > 0 {
		return c.MemoryDays
	}
	return DEFAULT_MEMORY_DAYS
}

func (c *StoreConfig) IsPersistence() bool {
	return c.Persistence
}

func (c *StoreConfig) GetPersistencePath(defaultValue string) string {
	if len(c.PersistencePath) > 0 {
		return c.PersistencePath
	}
	return defaultValue
}
