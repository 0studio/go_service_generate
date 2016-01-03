package user

// this Test file here for making sure the generated file is working as expected

import (
	"github.com/0studio/logger"
	key "github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProxyUser2Storage(t *testing.T) {
	now := time.Now()

	u := User2{}
	u.SetId(1)
	u.SetName("hello")
	u.SetAge(11)

	mc := NewMCUser2Storage(memcache.NewMockClient(), 10, "user2")
	p1 := NewStorageProxy(NewLRUCacheUser2Storage(1, 10), mc)
	store := NewStorageProxy(p1, NewDBUser2Storage(getMockDB(), logger.NewStdoutLogger(), true))

	ok := store.Add(&u, now)
	assert.True(t, ok)

	uRet, ok := store.Get(u.GetId(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetAge(), u.GetAge())
	assert.Equal(t, uRet.GetName(), u.GetName())

	u.SetAge(12)
	u.SetT(now)
	u.SetT2(now)
	u.SetSex(true)
	ok = store.Set(&u, now)
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetAge(), u.GetAge())
	assert.Equal(t, uRet.GetName(), u.GetName())
	assert.Equal(t, uRet.GetSex(), u.GetSex())
	assert.Equal(t, uRet.GetT().Unix(), u.GetT().Unix())
	assert.Equal(t, uRet.GetT2().Unix(), u.GetT2().Unix())

	ok = store.Delete(u.GetId())
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), now)
	assert.False(t, ok)
}

func TestProxyCacheUser2StorageMulti(t *testing.T) {
	now := time.Now()
	var uin key.KeyUint64 = 1
	var uin2 key.KeyUint64 = 2
	u := User2{}
	u.SetId(uin)
	u.SetName("n1")
	u.SetAge(111)
	u2 := User2{}
	u2.SetId(uin2)
	u2.SetName("n2")
	u2.SetAge(222)
	uMap := make(User2Map)
	uMap[u.GetId()] = u
	uMap[u2.GetId()] = u2

	mc := NewMCUser2Storage(memcache.NewMockClient(), 10, "user2")
	p1 := NewStorageProxy(NewLRUCacheUser2Storage(1, 10), mc)
	store := NewStorageProxy(p1, NewDBUser2Storage(getMockDB(), logger.NewStdoutLogger(), true))

	ok := store.MultiAdd(uMap, now)
	assert.True(t, ok)

	uMapRet, ok := store.MultiGet([]key.KeyUint64{uin, uin2}, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(uMapRet))

	ok = store.MultiDelete([]key.KeyUint64{uin, uin2})
	assert.True(t, ok)
	// after delete
	uMapRet, ok = store.MultiGet([]key.KeyUint64{uin, uin2}, now)
	assert.True(t, ok)
	assert.Equal(t, 0, len(uMapRet))

}
