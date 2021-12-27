package ttlcache

import (
	"errors"
	"sync"
	"time"
)

type ttlCache struct {
	datas map[string]*cacheObject
	mtx   *sync.RWMutex
}
type cacheObject struct {
	data interface{}
	ttl  time.Time
}

type (
	StoreOption struct {
		IsOverride bool
		TTL        *time.Duration
	}
	LoadOption struct {
		Default func() interface{}
		TTL     *time.Duration
	}
	VaridicStoreOption interface{ _StoreOption(*StoreOption) }
	VaridicLoadOption  interface{ _LoadOption(*LoadOption) }
	// Load option
	Default func() interface{}
	// Store, Load option
	TTL time.Duration
	// Store option
	Override bool
)

func FromStoreOptions(vs ...VaridicStoreOption) (res StoreOption) {
	for _, v := range vs {
		v._StoreOption(&res)
	}
	return
}

func FromLoadOptions(vs ...VaridicLoadOption) (res LoadOption) {
	for _, v := range vs {
		v._LoadOption(&res)
	}
	return
}

func (ttl TTL) _StoreOption(opt *StoreOption)      { a := time.Duration(ttl); opt.TTL = &a }
func (ttl TTL) _LoadOption(opt *LoadOption)        { a := time.Duration(ttl); opt.TTL = &a }
func (dft Default) _LoadOption(opt *LoadOption)    { opt.Default = dft }
func (ovd Override) _StoreOption(opt *StoreOption) { opt.IsOverride = bool(ovd) }

var (
	ErrExistKey    = errors.New("exist key")
	ErrNotExistKey = errors.New("not exist key")
)

func TTLCache() *ttlCache {
	return &ttlCache{
		datas: map[string]*cacheObject{},
		mtx:   new(sync.RWMutex),
	}
}
func (ttlc *ttlCache) Store(key string, value interface{}, opts ...VaridicStoreOption) error {
	return ttlc.StoreBy(key, value, FromStoreOptions(opts...))
}
func (ttlc *ttlCache) StoreBy(key string, value interface{}, option StoreOption) error {
	ttlc.mtx.Lock()
	defer ttlc.mtx.Unlock()
	if !option.IsOverride {
		if _, ok := ttlc.datas[key]; ok {
			return ErrExistKey
		}
	}
	node := cacheObject{
		data: value,
	}
	if option.TTL != nil {
		node.ttl = time.Now().Add(*option.TTL)
	}
	ttlc.datas[key] = &node
	return nil
}

func (ttlc *ttlCache) Load(key string, opts ...VaridicLoadOption) (interface{}, error) {
	return ttlc.LoadBy(key, FromLoadOptions(opts...))
}

func (ttlc *ttlCache) LoadBy(key string, option LoadOption) (interface{}, error) {
	ttlc.mtx.RLock()
	defer ttlc.mtx.RUnlock()
	if dat, ok := ttlc.datas[key]; ok {
		if dat.ttl.Unix() != 0 && time.Now().After(dat.ttl) {
			defer func() {
				ttlc.Delete(key)
			}()
			return nil, ErrNotExistKey
		}
		if option.TTL != nil {
			dat.ttl = time.Now().Add(*option.TTL)
		}
		return dat.data, nil
	} else {
		if option.Default != nil {
			return option.Default(), nil
		} else {
			return nil, ErrNotExistKey
		}
	}
}

func (ttlc *ttlCache) Delete(key string) {
	ttlc.mtx.Lock()
	defer ttlc.mtx.Unlock()
	delete(ttlc.datas, key)
}
