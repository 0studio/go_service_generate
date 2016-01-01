package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (sd StructDescription) generateLRUCache1PK(pkField FieldDescriptoin, property Property, srcDir string) bool {
	lruCacheType := sd.getLRUCacheType(pkField)
	if lruCacheType == "" {
		return false
	}
	outputF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("lru_cache_stub.go")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer outputF.Close()

	s := strings.Replace(LRUCacheTemplate, "__package__", property.PackageName, -1)
	s = strings.Replace(s, "__Entity__", sd.StructName, -1)
	s = strings.Replace(s, "__PKType__", pkField.FieldGoType, -1)
	s = strings.Replace(s, "__PKFieldName__", pkField.FieldName, -1)
	if strings.Contains(lruCacheType, "Shard") {
		s = strings.Replace(s, "__NewLRUCacheType__", "New"+lruCacheType+"(shardingCnt,size)", -1)
	} else {
		s = strings.Replace(s, "__NewLRUCacheType__", "New"+lruCacheType+"(size)", -1)
	}
	s = strings.Replace(s, "__LRUCacheType__", lruCacheType, -1)

	var pkTypeList string
	if isTypeKeySum(pkField.FieldGoType) {
		pkTypeList = fmt.Sprintf("%sList", pkField.FieldGoType)
	} else {
		pkTypeList = fmt.Sprintf("[]%s", pkField.FieldGoType)
	}

	s = strings.Replace(s, "__PKTypeList__", pkTypeList, -1)
	outputF.WriteString(s)
	return true
}

const (
	LRUCacheTemplate = `package __package__

import (
	"github.com/0studio/lru"
	key "github.com/0studio/storage_key"
	"time"
)

var __importKeyL key.KeyUint64

type LRULocal__Entity__Storage struct {
	cache *lru.__LRUCacheType__
	size  int64
}

func NewLRULocal__Entity__Storage(shardingCnt int, size int64) (local__Entity__Storage LRULocal__Entity__Storage) {
	local__Entity__Storage = LRULocal__Entity__Storage{
		cache: lru.__NewLRUCacheType__,
		size:  size,
	}
	return
}

func (m LRULocal__Entity__Storage) Get(k __PKType__, now time.Time) (e __Entity__, ok bool) {
	cacheObj, ok := m.cache.Get(k)
	if !ok {
		return
	}
	e = cacheObj.(__Entity__)
	ok = true
	return
}
func (m LRULocal__Entity__Storage) Set(e *__Entity__, now time.Time) bool {
	m.cache.Set(e.__PKFieldName__, *e)
	return true
}
func (m LRULocal__Entity__Storage) Add(e *__Entity__, now time.Time) bool {
	return m.Set(e, now)
}

func (m LRULocal__Entity__Storage) MultiGet(keys __PKTypeList__, now time.Time) (eMap __Entity__Map, ok bool) {
	eMap = make(__Entity__Map)
	var e __Entity__
	for _, k := range keys {
		e, ok = m.Get(k, now)
		if ok {
			eMap[k] = e
		}
	}
	ok = true
	return
}
func (m LRULocal__Entity__Storage) MultiUpdate(eMap __Entity__Map, now time.Time) (ok bool) {
	for _, e := range eMap {
		m.Set(&e, now)
	}
	return true
}

func (m LRULocal__Entity__Storage) Delete(k __PKType__) (ok bool) {
	m.cache.Delete(k)
	return true
}
func (m LRULocal__Entity__Storage) MultiDelete(keys __PKTypeList__) (ok bool) {
	for _, k := range keys {
		m.Delete(k)
	}
	return true
}
func (m LRULocal__Entity__Storage) Len() int {
	return int(m.cache.Size())
}
func (m LRULocal__Entity__Storage) GetAllUin() __PKTypeList__ {
	return m.cache.Keys()
}
func (m LRULocal__Entity__Storage) GetAll() (eMap __Entity__Map) {
	// get all not outdate e
	eMap = make(__Entity__Map)
	for _, item := range m.cache.Items() {
		eMap[item.Key] = item.Value.(__Entity__)
	}
	return
}
`
)