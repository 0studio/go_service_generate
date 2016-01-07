package user

// this Test file here for making sure the generated file is working as expected

import (
	"github.com/0studio/logger"
	"github.com/dropbox/godropbox/memcache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProxyUser4Storage(t *testing.T) {
	now := time.Now()

	u := User4{}
	u.SetId("key1")
	u.SetName("hello")

	local := NewLRUCacheUser4Storage(1, 10)
	mc := NewMCUser4Storage(memcache.NewMockClient(), 10, "user4")
	p := NewStorageProxyUser4(local, mc)
	store := NewStorageProxyUser4(p, NewDBUser4Storage(getMockDB(), logger.NewStdoutLogger(), true))

	ok := store.Add(&u, now)
	assert.True(t, ok)

	uRet, ok := store.Get(u.GetId(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetName(), u.GetName())

	ok = store.Set(&u, now)
	assert.True(t, ok)

	u.SetName("wooooo")
	ok = store.Set(&u, now)
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetName(), u.GetName())

	ok = store.Delete(u.GetId())
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), now)
	assert.False(t, ok)
}

func TestProxyUser4StorageMulti(t *testing.T) {
	now := time.Now()
	var uin string = "u1"
	var uin2 string = "u2"
	u := User4{}
	u.SetId(uin)
	u.SetName("n1")
	u2 := User4{}
	u2.SetId(uin2)
	u2.SetName("n2")
	uMap := make(User4Map)
	uMap[u.GetId()] = u
	uMap[u2.GetId()] = u2

	local := NewLRUCacheUser4Storage(1, 10)
	mc := NewMCUser4Storage(memcache.NewMockClient(), 10, "user4")
	p := NewStorageProxyUser4(local, mc)
	store := NewStorageProxyUser4(p, NewDBUser4Storage(getMockDB(), logger.NewStdoutLogger(), true))

	ok := store.MultiAdd(uMap, now)
	assert.True(t, ok)

	uMapRet, ok := store.MultiGet([]string{uin, uin2}, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(uMapRet))

	ok = store.MultiDelete([]string{uin, uin2})
	assert.True(t, ok)
	// after delete
	uMapRet, ok = store.MultiGet([]string{uin, uin2}, now)
	assert.True(t, ok)
	assert.Equal(t, 0, len(uMapRet))

}
