package generator

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

func (sd StructDescription) generateProxyPK2(pk1Field, pk2Field FieldDescriptoin, property Property, srcDir string) bool {
	outputF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("storage_%s_proxy_template.go", strings.ToLower(sd.StructName))), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer outputF.Close()

	s := strings.Replace(ProxyTemplatePK2, "__package__", property.PackageName, -1)
	s = strings.Replace(s, "__Entity__", sd.StructName, -1)
	s = strings.Replace(s, "__PK1Type__", pk1Field.FieldGoType, -1)
	s = strings.Replace(s, "__PK1FieldName__", pk1Field.FieldName, -1)
	s = strings.Replace(s, "__PK2Type__", pk2Field.FieldGoType, -1)
	s = strings.Replace(s, "__PK2FieldName__", pk2Field.FieldName, -1)

	var pk2TypeList string
	if isTypeKeySum(pk2Field.FieldGoType) {
		pk2TypeList = fmt.Sprintf("%sList", pk2Field.FieldGoType)
	} else {
		pk2TypeList = fmt.Sprintf("[]%s", pk2Field.FieldGoType)
	}

	s = strings.Replace(s, "__PK2TypeList__", pk2TypeList, -1)

	formatSrc, _ := format.Source([]byte(s))
	if err == nil {
		outputF.WriteString(string(formatSrc))
	} else {
		outputF.WriteString(s)
	}

	return true
}

const (
	ProxyTemplatePK2 = `// do not edit this file ,this is generated by tools(https://github.com/0studio/go_service_generator)
package __package__

import (
	"github.com/0studio/goutils"
	key "github.com/0studio/storage_key"
	"time"
)

var __importKeyP__Entity__ key.KeyUint64
var __importGoutils__Entity__ goutils.Int32List

type __Entity__Storage interface {
	SetIdListByPK1(__PK1FieldName__ __PK1Type__, __PK2FieldName__List *__PK2TypeList__, now time.Time) (ok bool)
	GetIdListByPK1(__PK1FieldName__ __PK1Type__, now time.Time) (__PK2FieldName__List __PK2TypeList__, ok bool)
	Get(__PK1FieldName__ __PK1Type__, __PK2FieldName__ __PK2Type__, now time.Time) (e __Entity__, ok bool)
	Set(e *__Entity__, now time.Time) (ok bool)
	Add(e *__Entity__, now time.Time) (ok bool)
	MultiGet(__PK1FieldName__ __PK1Type__, __PK2FieldName__List __PK2TypeList__, now time.Time) (eMap __Entity__Map, ok bool)
	MultiUpdate(__PK1FieldName__ __PK1Type__, eMap __Entity__Map, now time.Time) (ok bool)
	MultiAdd(__PK1FieldName__ __PK1Type__, eMap __Entity__Map, now time.Time) (ok bool)
	Delete(__PK1FieldName__ __PK1Type__, __PK2FieldName__ __PK2Type__) (ok bool)
	MultiDelete(__PK1FieldName__ __PK1Type__, __PK2FieldName__List __PK2TypeList__) (ok bool)
}

type __Entity__StorageProxy struct {
	preferedStorage __Entity__Storage
	backupStorage   __Entity__Storage
}

func NewStorageProxy__Entity__(prefered, backup __Entity__Storage) __Entity__Storage {
	return __Entity__StorageProxy{
		preferedStorage: prefered,
		backupStorage:   backup,
	}
}

func (this __Entity__StorageProxy) Get(__PK1FieldName__ __PK1Type__, __PK2FieldName__ __PK2Type__, now time.Time) (e __Entity__, ok bool) {
	e, ok = this.preferedStorage.Get(__PK1FieldName__, __PK2FieldName__, now)
	if ok {
		return
	}
	e, ok = this.backupStorage.Get(__PK1FieldName__, __PK2FieldName__, now)
	if !ok {
		return
	}
	this.preferedStorage.Set(&e, now)
	return
}

func (this __Entity__StorageProxy) Set(e *__Entity__, now time.Time) (ok bool) {
	ok = this.backupStorage.Set(e, now)
	if !ok {
		this.preferedStorage.Set(e, now)
		return ok
	}
	ok = this.preferedStorage.Set(e, now)
	return
}

func (this __Entity__StorageProxy) Add(e *__Entity__, now time.Time) (ok bool) {

	ok = this.backupStorage.Add(e, now)
	if !ok {
		// this.preferedStorage.Add(e, now)
		return ok
	}
	ok = this.preferedStorage.Add(e, now)
	return
}

func (this __Entity__StorageProxy) MultiGet(__PK1FieldName__ __PK1Type__, __PK2FieldName__List __PK2TypeList__, now time.Time) (eMap __Entity__Map, ok bool) {
	eMap, ok = this.preferedStorage.MultiGet(__PK1FieldName__, __PK2FieldName__List, now)
	missedKeyCount := 0
	for _, __PK2FieldName__ := range __PK2FieldName__List {
		if _, find := eMap[__PK2FieldName__]; !find {
			missedKeyCount++
		}
	}
	if missedKeyCount == 0 {
		return
	}

	missedKeys := make(__PK2TypeList__, missedKeyCount)
	i := 0
	for _, __PK2FieldName__ := range __PK2FieldName__List {
		if _, find := eMap[__PK2FieldName__]; !find {
			missedKeys[i] = __PK2FieldName__
			i++
		}
	}

	missedMap, ok := this.backupStorage.MultiGet(__PK1FieldName__, missedKeys, now)
	if !ok {
		return
	}
	this.preferedStorage.MultiUpdate(__PK1FieldName__, missedMap, now)
	for k, v := range missedMap {
		eMap[k] = v
	}
	return
}

func (this __Entity__StorageProxy) MultiUpdate(__PK1FieldName__ __PK1Type__, eMap __Entity__Map, now time.Time) (ok bool) {
	ok = this.backupStorage.MultiUpdate(__PK1FieldName__, eMap, now)
	if !ok {
		this.preferedStorage.MultiUpdate(__PK1FieldName__, eMap, now)
		return
	}
	ok = this.preferedStorage.MultiUpdate(__PK1FieldName__, eMap, now)
	return
}

func (this __Entity__StorageProxy) MultiAdd(__PK1FieldName__ __PK1Type__, eMap __Entity__Map, now time.Time) (ok bool) {
	ok = this.backupStorage.MultiAdd(__PK1FieldName__, eMap, now)
	if !ok {
		// this.preferedStorage.MultiAdd(__PK1FieldName__, eMap, now)
		return
	}
	ok = this.preferedStorage.MultiAdd(__PK1FieldName__, eMap, now)
	return
}
func (this __Entity__StorageProxy) Delete(__PK1FieldName__ __PK1Type__, __PK2FieldName__ __PK2Type__) (ok bool) {
	ok = this.backupStorage.Delete(__PK1FieldName__, __PK2FieldName__)
	if !ok {
		this.preferedStorage.Delete(__PK1FieldName__, __PK2FieldName__)
		return
	}
	ok = this.preferedStorage.Delete(__PK1FieldName__, __PK2FieldName__)
	return
}

func (this __Entity__StorageProxy) MultiDelete(__PK1FieldName__ __PK1Type__, __PK2FieldName__List __PK2TypeList__) (ok bool) {
	ok = this.backupStorage.MultiDelete(__PK1FieldName__, __PK2FieldName__List)
	if !ok {
		this.preferedStorage.MultiDelete(__PK1FieldName__, __PK2FieldName__List)
		return
	}
	ok = this.preferedStorage.MultiDelete(__PK1FieldName__, __PK2FieldName__List)
	return
}

func (this __Entity__StorageProxy) GetIdListByPK1(__PK1FieldName__ __PK1Type__, now time.Time) (idlist __PK2TypeList__, ok bool) {
	idlist, ok = this.preferedStorage.GetIdListByPK1(__PK1FieldName__, now)
	if ok {
		return
	}
	idlist, ok = this.backupStorage.GetIdListByPK1(__PK1FieldName__, now)
	if ok {
		this.preferedStorage.SetIdListByPK1(__PK1FieldName__, &idlist, now)
	}
	return
}

func (this __Entity__StorageProxy) SetIdListByPK1(__PK1FieldName__ __PK1Type__, __PK2FieldName__List *__PK2TypeList__, now time.Time) (ok bool) {
	if len(*__PK2FieldName__List) == 0 {
		return
	}
	ok = this.backupStorage.SetIdListByPK1(__PK1FieldName__, __PK2FieldName__List, now)
	if !ok {
		this.preferedStorage.SetIdListByPK1(__PK1FieldName__, __PK2FieldName__List, now)
		return ok
	}
	ok = this.preferedStorage.SetIdListByPK1(__PK1FieldName__, __PK2FieldName__List, now)
	return
}
`
)
