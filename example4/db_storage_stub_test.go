package user

// this Test file here for making sure the generated file is working as expected

import (
	"database/sql"
	"fmt"
	"github.com/0studio/databasetemplate"
	"github.com/0studio/logger"
	"github.com/stretchr/testify/assert"
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

func TestDBUser4Storage(t *testing.T) {
	now := time.Now()

	u := User4{}
	u.SetId("key1")
	u.SetName("hello")

	store := NewDBUser4Storage(getMockDB(), logger.NewStdoutLogger(), true)

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

func TestDBUser4StorageMulti(t *testing.T) {
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

	store := NewDBUser4Storage(getMockDB(), logger.NewStdoutLogger(), true)

	ok := store.MultiAdd(uin, uMap, now)
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
