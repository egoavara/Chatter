package ttlcache

import (
	"sync"
	"time"
)

type ttlCache struct {
	datas map[string]ttl
	mtx   *sync.RWMutex
}
type ttl struct {
	data interface{}
	ttl  time.Time
}
type LoadResult int32

const (
	Ok = 0
)
const (
	NotExist LoadResult = 1
	Timeout  LoadResult = 2
)

type StoreResult int32

const (
	Exist StoreResult = 1
)

func TTLCache() *ttlCache {
	return &ttlCache{
		datas: map[string]ttl{},
		mtx:   new(sync.RWMutex),
	}
}
func (ttlc *ttlCache) Store(key string, value interface{}, ttl time.Duration) StoreResult        {}
func (ttlc *ttlCache) StoreIfNot(key string, value interface{}, ttl time.Duration) StoreResult   {}
func (ttlc *ttlCache) StoreIfExist(key string, value interface{}, ttl time.Duration) StoreResult {}
func (ttlc *ttlCache) Load(key string) (interface{}, LoadResult)                                 {}
