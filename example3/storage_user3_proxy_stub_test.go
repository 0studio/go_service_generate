package user

// this Test file here for making sure the generated file is working as expected

import (
	"github.com/0studio/logger"
	"github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestProxyUser3Storage(t *testing.T) {
	now := time.Now()
	rand.Seed(now.UnixNano())

	var uin key.KeyUint64 = key.KeyUint64(rand.Int())
	u := User3{}
	u.SetId(uin)
	u.SetKey32(13)
	u.SetIkey(2)
	u.SetAge(11)
	u.SetName("hello")

	local := NewLRUCacheUser3Storage(1, 10)
	mc := NewMCUser3Storage(memcache.NewMockClient(), 10, "user3")
	p := NewStorageProxyUser3(local, mc)
	store := NewStorageProxyUser3(p, NewDBUser3Storage(getMockDB(), logger.NewStdoutLogger(), true))

	ok := store.Add(&u, now)
	assert.True(t, ok)

	uRet, ok := store.Get(u.GetId(), u.GetKey32(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetKey32(), u.GetKey32())
	assert.Equal(t, uRet.GetAge(), u.GetAge())
	assert.Equal(t, uRet.GetName(), u.GetName())
	assert.Equal(t, uRet.GetIkey(), u.GetIkey())

	u.SetAge(12)
	u.SetIkey(13)
	u.SetName("world")
	u.SetT(now)
	u.SetT2(now)
	u.SetSex(true)
	ok = store.Set(&u, now)
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), u.GetKey32(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetAge(), u.GetAge())
	assert.Equal(t, uRet.GetName(), u.GetName())
	assert.Equal(t, uRet.GetIkey(), u.GetIkey())
	assert.Equal(t, uRet.GetSex(), u.GetSex())
	assert.Equal(t, uRet.GetT().Unix(), u.GetT().Unix())
	assert.Equal(t, uRet.GetT2().Unix(), u.GetT2().Unix())

	ok = store.Delete(u.GetId(), u.GetKey32())
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), u.GetKey32(), now)
	assert.False(t, ok)
}

func TestProxyUser3StorageMulti(t *testing.T) {
	now := time.Now()
	rand.Seed(now.UnixNano())
	var uin key.KeyUint64 = key.KeyUint64(rand.Int())
	u := User3{}
	u.SetId(uin)
	u.SetKey32(111)
	u.SetAge(111)
	u2 := User3{}
	u2.SetId(uin)
	u2.SetKey32(222)
	u2.SetAge(222)
	uMap := make(User3Map)
	uMap[u.GetKey32()] = u
	uMap[u2.GetKey32()] = u2

	local := NewLRUCacheUser3Storage(1, 10)
	mc := NewMCUser3Storage(memcache.NewMockClient(), 10, "user3")
	p := NewStorageProxyUser3(local, mc)
	store := NewStorageProxyUser3(p, NewDBUser3Storage(getMockDB(), logger.NewStdoutLogger(), true))

	ok := store.MultiAdd(uin, uMap, now)
	assert.True(t, ok)

	uMapRet, ok := store.MultiGet(uin, key.KeyInt32List{u.GetKey32(), u2.GetKey32()}, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(uMapRet))

	ok = store.MultiDelete(uin, key.KeyInt32List{u.GetKey32(), u2.GetKey32()})
	assert.True(t, ok)
	// after delete
	uMapRet, ok = store.MultiGet(uin, key.KeyInt32List{u.GetKey32(), u2.GetKey32()}, now)
	assert.True(t, ok)
	assert.Equal(t, 0, len(uMapRet))

}
