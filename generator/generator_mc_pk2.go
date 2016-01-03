package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (sd StructDescription) generateMCPK2(pkField, pk2Field FieldDescriptoin, property Property, srcDir string) bool {
	outputF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("storage_%s_mc_template.go", strings.ToLower(sd.StructName))), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer outputF.Close()

	s := strings.Replace(MCTemplatePK2, "__package__", property.PackageName, -1)
	s = strings.Replace(s, "__Entity__", sd.StructName, -1)
	s = strings.Replace(s, "__PK1Type__", pkField.FieldGoType, -1)
	s = strings.Replace(s, "__PK2Type__", pk2Field.FieldGoType, -1)
	s = strings.Replace(s, "__PK1FieldName__", pkField.FieldName, -1)
	s = strings.Replace(s, "__PK2FieldName__", pk2Field.FieldName, -1)
	var pk2TypeList string
	if isTypeKeySum(pk2Field.FieldGoType) {
		pk2TypeList = fmt.Sprintf("%sList", pk2Field.FieldGoType)
	} else {
		pk2TypeList = fmt.Sprintf("[]%s", pk2Field.FieldGoType)
	}

	s = strings.Replace(s, "__PK2TypeList__", pk2TypeList, -1)

	outputF.WriteString(s)
	outputF.WriteString(sd.generateMCPK2getMCKey(pkField, pk2Field))
	outputF.WriteString(sd.generateMCPK2getRawKey(pkField, pk2Field))

	return true
}
func (sd StructDescription) generateMCPK2getMCKey(pkField, pk2Field FieldDescriptoin) string {
	var pk2String string
	if pk2Field.FieldGoType == "string" {
		pk2String = pk2Field.FieldName
	} else if pk2Field.FieldGoType == "int" {
		pk2String = fmt.Sprintf("int2str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int8" {
		pk2String = fmt.Sprintf("int82str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int16" {
		pk2String = fmt.Sprintf("int162str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int32" {
		pk2String = fmt.Sprintf("int322str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int64" {
		pk2String = fmt.Sprintf("int642str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint8" {
		pk2String = fmt.Sprintf("uint82str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint16" {
		pk2String = fmt.Sprintf("uint162str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint32" {
		pk2String = fmt.Sprintf("uint322str(%s)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint64" {
		pk2String = fmt.Sprintf("uint642str(%s)", pk2Field.FieldName)
	} else {
		pk2String = fmt.Sprintf("%s.String()", pk2Field.FieldName)
	}
	var pk1String string
	if pkField.FieldGoType == "string" {
		pk1String = pkField.FieldName
	} else if pkField.FieldGoType == "int" {
		pk1String = fmt.Sprintf("int2str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "int8" {
		pk1String = fmt.Sprintf("int82str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "int16" {
		pk1String = fmt.Sprintf("int162str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "int32" {
		pk1String = fmt.Sprintf("int322str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "int64" {
		pk1String = fmt.Sprintf("int642str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint8" {
		pk1String = fmt.Sprintf("uint82str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint16" {
		pk1String = fmt.Sprintf("uint162str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint32" {
		pk1String = fmt.Sprintf("uint322str(%s)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint64" {
		pk1String = fmt.Sprintf("uint642str(%s)", pkField.FieldName)

	} else {
		pk1String = fmt.Sprintf("%s.String()", pkField.FieldName)
	}

	s := fmt.Sprintf(`
func (m MC%sStorage) getMCKey(%s %s,%s %s) string {
	return fmt.Sprintf("%%s^%%s_%%s", %s,%s, m.prefix)
}
`, sd.StructName, pkField.FieldName, pkField.FieldGoType, pk2Field.FieldName, pk2Field.FieldGoType,
		pk2String, pk1String)
	return s
}
func (sd StructDescription) generateMCPK2getRawKey(pkField, pk2Field FieldDescriptoin) string {
	var pk2String string
	if pk2Field.FieldGoType == "string" {
		pk2String = fmt.Sprintf("%s = pk2Str", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int" {
		pk2String = fmt.Sprintf("%s = str2int(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int8" {
		pk2String = fmt.Sprintf("%s = str2int8(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int16" {
		pk2String = fmt.Sprintf("%s = str2int16(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int32" {
		pk2String = fmt.Sprintf("%s = str2int32(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "int64" {
		pk2String = fmt.Sprintf("%s = str2int64(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint8" {
		pk2String = fmt.Sprintf("%s = str2uint8(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint16" {
		pk2String = fmt.Sprintf("%s = str2uint16(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint32" {
		pk2String = fmt.Sprintf("%s = str2uint32(pk2Str)", pk2Field.FieldName)
	} else if pk2Field.FieldGoType == "uint64" {
		pk2String = fmt.Sprintf("%s = str2uint64(pk2Str)", pk2Field.FieldName)
	} else {
		pk2String = fmt.Sprintf("%s.FromString(pk2Str)", pk2Field.FieldName)
	}
	var pk1String string
	if pkField.FieldGoType == "string" {
		pk1String = fmt.Sprintf("%s = pk1Str", pkField.FieldName)
	} else if pkField.FieldGoType == "int" {
		pk1String = fmt.Sprintf("%s = str2int(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "int8" {
		pk1String = fmt.Sprintf("%s = str2int8(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "int16" {
		pk1String = fmt.Sprintf("%s = str2int16(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "int32" {
		pk1String = fmt.Sprintf("%s = str2int32(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "int64" {
		pk1String = fmt.Sprintf("%s = str2int64(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint8" {
		pk1String = fmt.Sprintf("%s = str2uint8(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint16" {
		pk1String = fmt.Sprintf("%s = str2uint16(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint32" {
		pk1String = fmt.Sprintf("%s = str2uint32(pk1Str)", pkField.FieldName)
	} else if pkField.FieldGoType == "uint64" {
		pk1String = fmt.Sprintf("%s = str2uint64(pk1Str)", pkField.FieldName)
	} else {
		pk1String = fmt.Sprintf("%s.FromString(pk1Str)", pkField.FieldName)
	}

	s := fmt.Sprintf(`
func (m MC%sStorage) getRawKey(prefixKey string) (%s %s,%s %s) {
	var pk1Str string
	var pk2Str string
	char1Idx := strings.Index(prefixKey, "^")
	char2Idx := strings.Index(prefixKey, "_")
    if char1Idx != -1 && char2Idx != -1 && char2Idx > char1Idx {
		pk2Str = prefixKey[:char1Idx]
		pk1Str = prefixKey[char1Idx+1 : char2Idx]
		%s
		%s
    }
	return
}
`, sd.StructName, pkField.FieldName, pkField.FieldGoType, pk2Field.FieldName, pk2Field.FieldGoType,
		pk2String, pk1String)
	return s
}

const (
	MCTemplatePK2 = `// do not edit this file ,this is generated by tools(https://github.com/0studio/go_service_generator)
package __package__

import (
	"fmt"
	"github.com/0studio/goutils"
	key "github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
    "strings"
	"time"
)

var __importKeyMC key.KeyUint64
var __importGoutilsMC goutils.Int32List

type MC__Entity__Storage struct {
	expireSeconds uint32
	prefix        string
	client        memcache.Client
}

func NewMC__Entity__Storage(client memcache.Client, expireSeconds uint32, prefix string) MC__Entity__Storage {
	return MC__Entity__Storage{expireSeconds: expireSeconds, prefix: prefix, client: client}
}

func (m MC__Entity__Storage) Get(__PK1FieldName__ __PK1Type__, __PK2FieldName__ __PK2Type__, now time.Time) (e __Entity__, ok bool) {
	item := m.client.Get(m.getMCKey(__PK1FieldName__, __PK2FieldName__))
	if item.Error() != nil || item.Status() != memcache.StatusNoError {
		return
	}
	byteData := item.Value()
	if byteData == nil {
		return
	}

	e.UnSerial(byteData)
	if !(e.__PK1FieldName__ == __PK1FieldName__ ){
		ok = false
		return
	}

	ok = true
	return
}
func (m MC__Entity__Storage) Set(e *__Entity__, now time.Time) (ok bool) {
	item := memcache.Item{Key: m.getMCKey(e.__PK1FieldName__, e.__PK2FieldName__), Value: e.Serial(), Expiration: m.expireSeconds}
	res := m.client.Set(&item)
	return res.Error() == nil

}
func (m MC__Entity__Storage) Add(e *__Entity__, now time.Time) (ok bool) {
	return m.Set(e, now)
}

func (m MC__Entity__Storage) MultiGet(__PK1FieldName__ __PK1Type__, keys __PK2TypeList__, now time.Time) (eMap __Entity__Map, ok bool) {
	prefixKeys := make([]string, len(keys))
	for idx, __PK2FieldName__ := range keys {
		prefixKeys[idx] = m.getMCKey(__PK1FieldName__, __PK2FieldName__)
	}
	itemMap := m.client.GetMulti(prefixKeys)

	eMap = make(__Entity__Map)
	var e __Entity__
	for k, item := range itemMap {
		if len(item.Value()) == 0 {
			continue
		}
		if e.UnSerial(item.Value()) {
			if e.__PK1FieldName__== __PK1FieldName__ {
				_, __PK2FieldName__ := m.getRawKey(k)
				eMap[__PK2FieldName__] = e
			}
		}
	}
	ok = true
	return
}
func (m MC__Entity__Storage) MultiUpdate(__PK1FieldName__ __PK1Type__, eMap __Entity__Map, now time.Time) (ok bool) {
	items := make([]memcache.Item, 0, len(eMap))
	itemsPointer := make([]*memcache.Item, 0, len(eMap))
	var idx int = 0
	for _, e := range eMap {
		items = append(items, memcache.Item{Key: m.getMCKey(__PK1FieldName__, e.__PK2FieldName__), Value: e.Serial(), Expiration: m.expireSeconds})
		itemsPointer = append(itemsPointer, &(items[idx]))
		idx++
	}
	responses := m.client.SetMulti(itemsPointer)
	for _, res := range responses {
		if res.Error() != nil {
			return false
		}
	}
	return true
}
func (m MC__Entity__Storage) MultiAdd(__PK1FieldName__ __PK1Type__, eMap __Entity__Map, now time.Time) (ok bool) {
	return m.MultiUpdate(__PK1FieldName__, eMap, now)
}
func (m MC__Entity__Storage) Delete(__PK1FieldName__ __PK1Type__, __PK2FieldName__ __PK2Type__) (ok bool) {
	res := m.client.Delete(m.getMCKey(__PK1FieldName__, __PK2FieldName__))
	return res.Error() == nil || res.Status() == memcache.StatusKeyNotFound
}
func (m MC__Entity__Storage) MultiDelete(__PK1FieldName__ __PK1Type__, keys __PK2TypeList__) (ok bool) {
	Prefixkeys := make([]string, len(keys), len(keys))
	for idx, __PK2FieldName__ := range keys {
		Prefixkeys[idx] = m.getMCKey(__PK1FieldName__, __PK2FieldName__)
	}
	responses := m.client.DeleteMulti(Prefixkeys)
	for _, res := range responses {
		if res.Error() != nil && res.Status() != memcache.StatusKeyNotFound {
			return false
		}
	}
	return true
}
func (m MC__Entity__Storage) GetIdListByPK1(key __PK1Type__, now time.Time) (list __PK2TypeList__, ok bool) {
	return
}
func (m MC__Entity__Storage) SetIdListByPK1(key __PK1Type__, idlist *__PK2TypeList__, now time.Time) (ok bool) {
	return true
}
`
)
