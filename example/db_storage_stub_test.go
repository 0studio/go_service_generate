package user

// this Test file here for making sure the generated file is working as expected

import (
	"database/sql"
	"fmt"
	"github.com/0studio/databasetemplate"
	"github.com/0studio/goutils"
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

func TestDBUserStorage(t *testing.T) {
	now := time.Now()

	u := User{}
	u.SetId(1)
	u.SetName("hello")
	u.SetIList([]int{1, 2})
	u.SetI2List([]int32{1, 2, 3})
	u.SetI3List([]int8{1, 2, 4})
	u.SetI4List([]int16{1, 2, 5})
	u.SetI5List([]int64{1, 2, 6})
	u.SetI6List([]uint32{1, 2, 7})
	u.SetI7List([]uint8{1, 2, 8})
	u.SetI8List([]uint16{1, 2, 9})
	u.SetI9List([]uint64{1, 2, 10})
	u.SetI10List(goutils.IntList{1, 2, 11})
	u.SetAge(11)
	u.SetSex(true)
	u.SetT(now)
	u.SetT2(now)

	store := NewDBUserStorage(getMockDB(), logger.NewStdoutLogger(), true)

	ok := store.Add(&u, now)
	assert.True(t, ok)

	uRet, ok := store.Get(u.GetId(), u.GetName(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetIList(), u.GetIList())
	assert.Equal(t, uRet.GetI2List(), u.GetI2List())
	assert.Equal(t, uRet.GetI3List(), u.GetI3List())
	assert.Equal(t, uRet.GetI4List(), u.GetI4List())
	assert.Equal(t, uRet.GetI5List(), u.GetI5List())
	assert.Equal(t, uRet.GetI6List(), u.GetI6List())
	assert.Equal(t, uRet.GetI7List(), u.GetI7List())
	assert.Equal(t, uRet.GetI8List(), u.GetI8List())
	assert.Equal(t, uRet.GetI9List(), u.GetI9List())
	assert.Equal(t, uRet.GetI10List(), u.GetI10List())
	assert.Equal(t, uRet.GetAge(), u.GetAge())
	assert.Equal(t, uRet.GetName(), u.GetName())
	assert.Equal(t, uRet.GetSex(), u.GetSex())
	assert.Equal(t, uRet.GetT().Unix(), u.GetT().Unix())
	assert.Equal(t, uRet.GetT2().Unix(), u.GetT2().Unix())

	now = now.Add(time.Second * 10)
	u.SetAge(u.GetAge() + 1)
	u.SetT(now)
	u.SetT2(now)
	u.SetSex(!u.GetSex())
	ok = store.Set(&u, now)
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), u.GetName(), now)
	assert.True(t, ok)
	assert.Equal(t, uRet.GetIList(), u.GetIList())
	assert.Equal(t, uRet.GetI2List(), u.GetI2List())
	assert.Equal(t, uRet.GetI3List(), u.GetI3List())
	assert.Equal(t, uRet.GetI4List(), u.GetI4List())
	assert.Equal(t, uRet.GetI5List(), u.GetI5List())
	assert.Equal(t, uRet.GetI6List(), u.GetI6List())
	assert.Equal(t, uRet.GetI7List(), u.GetI7List())
	assert.Equal(t, uRet.GetI8List(), u.GetI8List())
	assert.Equal(t, uRet.GetI9List(), u.GetI9List())
	assert.Equal(t, uRet.GetI10List(), u.GetI10List())
	assert.Equal(t, uRet.GetAge(), u.GetAge())
	assert.Equal(t, uRet.GetName(), u.GetName())
	assert.Equal(t, uRet.GetSex(), u.GetSex())
	assert.Equal(t, uRet.GetT().Unix(), u.GetT().Unix())
	assert.Equal(t, uRet.GetT2().Unix(), u.GetT2().Unix())

	ok = store.Delete(u.GetId(), u.GetName())
	assert.True(t, ok)

	uRet, ok = store.Get(u.GetId(), u.GetName(), now)
	assert.False(t, ok)
}

func TestDBUserStorageMulti(t *testing.T) {
	now := time.Now()
	var uin int = 1
	u := User{}
	u.SetId(uin)
	u.SetName("n1")
	u.SetAge(111)
	u2 := User{}
	u2.SetId(uin)
	u2.SetName("n2")
	u2.SetAge(222)
	uMap := make(UserMap)
	uMap[u.GetName()] = u
	uMap[u2.GetName()] = u2

	store := NewDBUserStorage(getMockDB(), logger.NewStdoutLogger(), true)

	ok := store.MultiAdd(uin, uMap, now)
	assert.True(t, ok)

	uMapRet, ok := store.MultiGet(uin, []string{u.GetName(), u2.GetName()}, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(uMapRet))

	idList, ok := store.GetIdListByPK1(uin, now)
	assert.True(t, ok)
	assert.Equal(t, 2, len(idList))
	assert.NotEqual(t, 0, idList[0])
	assert.NotEqual(t, 0, idList[1])

	ok = store.MultiDelete(uin, []string{u.GetName(), u2.GetName()})
	assert.True(t, ok)
	// after delete
	uMapRet, ok = store.MultiGet(uin, []string{u.GetName(), u2.GetName()}, now)
	assert.True(t, ok)
	assert.Equal(t, 0, len(uMapRet))

}
