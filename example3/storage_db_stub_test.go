package user

// this Test file here for making sure the generated file is working as expected

import (
	"database/sql"
	"fmt"
	"github.com/0studio/databasetemplate"
	"github.com/0studio/logger"
	"github.com/0studio/storage_key"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

// GRANT ALL PRIVILEGES ON *.* TO 'th_dev'@'127.0.0.1'     IDENTIFIED BY 'th_devpass' WITH GRANT OPTION;

func getMockDB() (dt databasetemplate.DatabaseTemplate) {
	var ok bool
	db1, ok := databasetemplate.NewDBInstance(
		databasetemplate.DBConfig{
			Host: "127.0.0.1",
			User: "th_dev",
			Pass: "th_devpass",
			Name: "test",
		}, true)
	if !ok {
		fmt.Println("initmock_databasetemplate_fail")
	}
	dt = databasetemplate.NewDatabaseTemplateSharding([]*sql.DB{db1})

	return
}

func TestDBUser3Storage(t *testing.T) {
	now := time.Now()

	u := User3{}
	u.SetId(1)
	u.SetKey32(13)
	u.SetIkey(2)
	u.SetAge(11)
	u.SetName("hello")

	store := NewDBUser3Storage(getMockDB(), logger.NewStdoutLogger(), true)

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

func TestDBUser3StorageMulti(t *testing.T) {
	now := time.Now()
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

	store := NewDBUser3Storage(getMockDB(), logger.NewStdoutLogger(), true)

	ok := store.MultiAdd(uin, uMap, now)
	assert.True(t, ok)

	uMapRet, ok := store.MultiGet(uin, key.KeyInt32List{u.GetKey32(), u2.GetKey32()}, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(uMapRet))

	idList, ok := store.GetIdListByPK1(uin, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(idList))
	assert.NotEqual(t, 0, idList[0])
	assert.NotEqual(t, 0, idList[1])

	ok = store.MultiDelete(uin, key.KeyInt32List{u.GetKey32(), u2.GetKey32()})
	assert.True(t, ok)
	// after delete
	uMapRet, ok = store.MultiGet(uin, key.KeyInt32List{u.GetKey32(), u2.GetKey32()}, now)
	assert.True(t, ok)
	assert.Equal(t, 0, len(uMapRet))

}
