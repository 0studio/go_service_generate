package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (sd StructDescription) generateProxyPK1(pkField FieldDescriptoin, property Property, srcDir string) bool {
	outputF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("storage_proxy_stub.go")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer outputF.Close()

	s := strings.Replace(ProxyTemplatePK1, "__package__", property.PackageName, -1)
	s = strings.Replace(s, "__Entity__", sd.StructName, -1)
	s = strings.Replace(s, "__PKType__", pkField.FieldGoType, -1)
	s = strings.Replace(s, "__PK1FieldName__", pkField.FieldName, -1)

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
	ProxyTemplatePK1 = `// do not edit this file ,this is generated by tools(https://github.com/0studio/go_service_generator)
package __package__

import (
	"github.com/0studio/goutils"
	key "github.com/0studio/storage_key"
	"time"
)

var __importKeyP key.KeyUint64
var __importGoutilsP goutils.Int32List

type __Entity__Storage interface {
	Get(__PK1FieldName__ __PKType__, now time.Time) (e __Entity__, ok bool)
	Set(e *__Entity__, now time.Time) (ok bool)
	Add(e *__Entity__, now time.Time) bool
	MultiGet(__PK1FieldName__List __PKTypeList__, now time.Time) (eMap __Entity__Map, ok bool)
	MultiUpdate(eMap __Entity__Map, now time.Time) (ok bool)
	Delete(__PK1FieldName__ __PKType__) (ok bool)
	MultiDelete(__PK1FieldName__List __PKTypeList__) (ok bool)
}

type __Entity__StorageProxy struct {
	preferedStorage __Entity__Storage
	backupStorage   __Entity__Storage
}

func NewStorageProxy(prefered, backup __Entity__Storage) __Entity__Storage  {
	return __Entity__StorageProxy{
		preferedStorage: prefered,
		backupStorage:   backup,
	}
}

func (this __Entity__StorageProxy) Get(__PK1FieldName__ __PKType__, now time.Time) (e __Entity__, ok bool) {
	e, ok = this.preferedStorage.Get(__PK1FieldName__, now)
	if ok {
		return
	}
	e, ok = this.backupStorage.Get(__PK1FieldName__, now)
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

func (this __Entity__StorageProxy) MultiGet(__PK1FieldName__List __PKTypeList__, now time.Time) (eMap __Entity__Map, ok bool) {
	eMap, ok = this.preferedStorage.MultiGet(__PK1FieldName__List, now)
	missedKeyCount := 0
	for _, __PK1FieldName__ := range __PK1FieldName__List {
		if _, find := eMap[__PK1FieldName__]; !find {
			missedKeyCount++
		}
	}
	if missedKeyCount == 0 {
		return
	}

	missedKeys := make(__PKTypeList__, missedKeyCount)
	i := 0
	for _, __PK1FieldName__ := range __PK1FieldName__List {
		if _, find := eMap[__PK1FieldName__]; !find {
			missedKeys[i] = __PK1FieldName__
			i++
		}
	}

	missedMap, ok := this.backupStorage.MultiGet(missedKeys, now)
	if !ok {
		return
	}
	this.preferedStorage.MultiUpdate(missedMap, now)
	for k, v := range missedMap {
		eMap[k] = v
	}
	return
}

func (this __Entity__StorageProxy) MultiUpdate(eMap __Entity__Map, now time.Time) (ok bool) {
	ok = this.backupStorage.MultiUpdate(eMap, now)
	if !ok {
		this.preferedStorage.MultiUpdate(eMap, now)
		return
	}
	ok = this.preferedStorage.MultiUpdate(eMap, now)
	return
}
func (this __Entity__StorageProxy) Delete(__PK1FieldName__ __PKType__) (ok bool) {
	ok = this.backupStorage.Delete(__PK1FieldName__)
	if !ok {
		this.preferedStorage.Delete(__PK1FieldName__)
		return
	}
	ok = this.preferedStorage.Delete(__PK1FieldName__)
	return
}

func (this __Entity__StorageProxy) MultiDelete(__PK1FieldName__List __PKTypeList__) (ok bool) {
	ok = this.backupStorage.MultiDelete(__PK1FieldName__List)
	if !ok {
		this.preferedStorage.MultiDelete(__PK1FieldName__List)
		return
	}
	ok = this.preferedStorage.MultiDelete(__PK1FieldName__List)
	return

}
`
)
