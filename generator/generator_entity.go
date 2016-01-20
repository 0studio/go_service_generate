package generator

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

func (sd StructDescription) GenerateEntity(property Property, srcDir string) {
	outputF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("entity_%s_template.go", strings.ToLower(sd.StructName))), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer outputF.Close()
	s := "package "
	s += property.PackageName

	s +=
		fmt.Sprintf(
			` // do not edit this file ,this is generated by tools(https://github.com/0studio/go_service_generator)

import (
    "bytes"
    "fmt"
    "time"
    "github.com/0studio/bit"
    "github.com/0studio/goutils"
    key "github.com/0studio/storage_key"
)
var ___importTime%s time.Time
var ___importBit%s bit.BitInt32
var ___importGoutils%s goutils.Int32List
var ___importKey%s key.KeyUint64
var ___importBytes%s bytes.Buffer

`, sd.StructName, sd.StructName, sd.StructName, sd.StructName, sd.StructName)

	// for _, sd := range structDescriptionList {
	// 生成结构体
	s += fmt.Sprintf("type %s struct{\n", sd.StructName)
	if len(sd.Fields) < 8 {
		s += "    flag            bit.BitInt8\n"
	} else if len(sd.Fields) < 16 {
		s += "    flag            bit.BitInt16\n"
	} else if len(sd.Fields) < 32 {
		s += "    flag            bit.BitInt32\n"
	} else {
		s += "    flag            bit.BitInt64\n"
	}

	for _, field := range sd.Fields {
		s += fmt.Sprintf("    %s%s%s", field.FieldName, makeString(" ", 16-len(field.FieldName)), field.FieldGoType)
		if field.TagString != "" {
			s += fmt.Sprintf(" `%s`", field.TagString)
		}
		if field.Comments != "" {
			s += fmt.Sprintf(" // %s", field.Comments)
		}

		s += "\n"
	}
	s += "}\n\n"
	s += fmt.Sprintf("type %sList []%s\n\n", sd.StructName, sd.StructName)

	s += sd.generateEntityMap()
	s += sd.generateNewEntity()

	// 生成setter getter
	var flagPos int = 0
	for _, field := range sd.Fields {
		s += fmt.Sprintf("func(e *%s) Set%s(value %s) {\n", sd.StructName, Camelize(field.FieldName), field.FieldGoType)
		if field.IsEqualableType() {
			s += fmt.Sprintf("    if e.%s != value {\n", field.FieldName)
			s += fmt.Sprintf("        e.%s = value\n", field.FieldName)
			if !field.IsPK() {
				s += fmt.Sprintf("        e.flag.SetFlag(%d)\n", flagPos)
			}
			s += "    }\n"
		} else {
			s += fmt.Sprintf("    e.%s = value\n", field.FieldName)
			if !field.IsPK() {
				s += fmt.Sprintf("    e.flag.SetFlag(%d)\n", flagPos)
			}
		}

		s += "}\n"
		s += fmt.Sprintf("func(e %s) Get%s()(value %s) {\n", sd.StructName, Camelize(field.FieldName), field.FieldGoType)
		s += fmt.Sprintf("    return e.%s\n", field.FieldName)
		s += "}\n"
		if !field.IsPK() {
			s += fmt.Sprintf("func(e %s) Is%sModified() bool {\n", sd.StructName, Camelize(field.FieldName))
			s += fmt.Sprintf("    return e.flag.IsPosTrue(%d)\n", flagPos)
			s += "}\n"
			flagPos++
		}
	}
	s += fmt.Sprintf("func(e %s) IsFlagDirty() bool {\n", sd.StructName)
	s += fmt.Sprintf("    return e.flag!=0\n")
	s += "}\n"
	s += fmt.Sprintf("func(e *%s) ClearFlag() {\n", sd.StructName)
	s += fmt.Sprintf("    e.flag=0\n")
	s += "}\n"

	s += fmt.Sprintf("func(e %s) IsFlagNew() bool {\n", sd.StructName)
	s += fmt.Sprintf("    return e.flag.IsPosTrue(%d)\n", len(sd.Fields))
	s += "}\n"
	s += fmt.Sprintf("func(e *%s) ClearFlagNew() {\n", sd.StructName)
	s += fmt.Sprintf("    e.flag.UnSetFlag(%d)\n", len(sd.Fields))
	s += "}\n"

	s += sd.GenerateInsert()
	s += sd.GenerateInsertWithArgs()

	s += "\n"
	s += sd.GenerateCreateTableFunc()
	s += "\n"
	s += sd.GenerateUpdate()
	s += "\n"
	s += sd.GenerateUpdateWithArgs()
	s += "\n"

	s += sd.GenerateInsertForMap()
	s += sd.GenerateInsertForMapWithArgs()

	// }

	formatSrc, err := format.Source([]byte(s))
	if err != nil {
		outputF.WriteString(s)
	} else {
		outputF.WriteString(string(formatSrc))
	}

}
func (sd StructDescription) generateNewEntity() (s string) {
	var fields string
	for _, fd := range sd.Fields {
		if fd.IsPK() {
			fields += "        " + fd.FieldName + ":" + fd.FieldName + ",\n"
		} else {
			if fd.GetGoDefalutValue() != "" {
				fields += "        " + fd.FieldName + ":" + fd.GetGoDefalutValue() + ",\n"
			}
		}
	}
	s = fmt.Sprintf(
		`func New%s(%s) (e %s) {
    e = %s{
%s
    }
    e.flag.SetFlag(%d)
    return
}
`, sd.StructName, sd.GetPKVarDeclear(), sd.StructName,
		sd.StructName,
		fields, len(sd.Fields))
	return
}
func (sd StructDescription) generateEntityMap() (s string) {
	pkList := sd.GetPKFieldList()
	if len(pkList) == 0 || len(pkList) > 2 {
		return ""
	}

	return fmt.Sprintf("type %s map[%s]%s\n\n", sd.GetSuggestMapName(), sd.GetSuggestMapKey(), sd.StructName)
}

func makeString(s string, n int) (r string) {
	for i := 0; i < n; i++ {
		r += s
	}
	return

}
func (sd StructDescription) GenerateInsertWithArgs() (goCode string) {
	goCode += fmt.Sprintf("func (e %s) GetInsertSqlWithArgs() (sql string,args []interface{}) {\n", sd.StructName)
	goCode += fmt.Sprintf("    sql = \"insert into `%s`(", sd.GetMysqlTableName())
	goCode += sd.JoinMysqlFieldNameList(",")
	goCode += ") values ("
	for idx, _ := range sd.Fields {
		goCode += "?"
		if idx != len(sd.Fields)-1 {
			goCode += ","
		}
	}
	goCode += ");\"\n"
	goCode += fmt.Sprintf("    args = []interface{}{\n")
	for idx, field := range sd.GetMysqlFieldList() {
		if field.IsBool() {
			goCode += fmt.Sprintf("        bool2int(e.%s)", field.FieldName)
		}
		if field.IsIntList() {
			switch field.FieldGoType {
			case "[]int":
				goCode += fmt.Sprintf("        intListJoin(e.%s, `,`)", field.FieldName)
			case "[]int8":
				goCode += fmt.Sprintf("        int8ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int16":
				goCode += fmt.Sprintf("        int16ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int32":
				goCode += fmt.Sprintf("        int32ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int64":
				goCode += fmt.Sprintf("        int64ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint8":
				goCode += fmt.Sprintf("        uint8ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint16":
				goCode += fmt.Sprintf("        uint16ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint32":
				goCode += fmt.Sprintf("        uint32ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint64":
				goCode += fmt.Sprintf("        uint64ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int32List":
				goCode += fmt.Sprintf("        int32ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int16List":
				goCode += fmt.Sprintf("        int16ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.IntList":
				goCode += fmt.Sprintf("        intListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int8List":
				goCode += fmt.Sprintf("        int8ListJoin(e.%s, `,`)", field.FieldName)
			default:
				fmt.Println("should be here generateDBMapRowSpicialTypeTrans", field.FieldGoType)
			}
		}
		if field.IsStringList() {
			goCode += fmt.Sprintf("        stringListJoin(e.%s, `,`)", field.FieldName)
		}

		if field.IsNumber() {
			goCode += fmt.Sprintf("        e.%s", field.FieldName)
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "timestamp" {
			goCode += fmt.Sprintf("        e.%s", field.FieldName)
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "datetime" {
			goCode += fmt.Sprintf("        e.%s", field.FieldName)
		}
		if field.IsTimeInt() {
			goCode += fmt.Sprintf("        formatTimeUnix(e.%s)", field.FieldName)
		}
		if field.FieldGoType == "string" {
			goCode += fmt.Sprintf("        e.%s", field.FieldName)
		}

		if idx != len(sd.Fields)-1 {
			goCode += ",\n"
		}
	}

	goCode += "}\n"

	goCode += "    return\n"
	goCode += "}\n"
	return
}
func (sd StructDescription) GenerateInsert() (goCode string) {
	goCode += fmt.Sprintf("func (e %s) GetInsertSql() (sql string) {\n", sd.StructName)
	goCode += fmt.Sprintf("    sql = fmt.Sprintf(\"insert into `%s`(", sd.GetMysqlTableName())
	goCode += sd.JoinMysqlFieldNameList(",")
	goCode += ") values ("
	for idx, field := range sd.Fields {
		if field.IsBool() {

			goCode += "%d"
		}
		if field.IsInt() {
			goCode += "%d"
		}
		if field.IsFloat() {
			goCode += "%f"
		}
		if field.IsIntList() {
			goCode += "'%s'"
		}
		if field.IsStringList() {
			goCode += "'%s'"
		}

		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "timestamp" {
			goCode += "%s"
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "datetime" {
			goCode += "%s"
		}
		if field.IsTimeInt() {
			goCode += "%d"
		}
		if field.FieldGoType == "string" {
			goCode += "'%s'"
		}

		if idx != len(sd.Fields)-1 {
			goCode += ","
		}
	}
	goCode += ");\",\n"
	for idx, field := range sd.Fields {
		if field.IsBool() {
			goCode += fmt.Sprintf("        bool2int(e.%s)", field.FieldName)
		}
		if field.IsIntList() {
			switch field.FieldGoType {
			case "[]int":
				goCode += fmt.Sprintf("        intListJoin(e.%s, `,`)", field.FieldName)
			case "[]int8":
				goCode += fmt.Sprintf("        int8ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int16":
				goCode += fmt.Sprintf("        int16ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int32":
				goCode += fmt.Sprintf("        int32ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int64":
				goCode += fmt.Sprintf("        int64ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint8":
				goCode += fmt.Sprintf("        uint8ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint16":
				goCode += fmt.Sprintf("        uint16ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint32":
				goCode += fmt.Sprintf("        uint32ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint64":
				goCode += fmt.Sprintf("        uint64ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int32List":
				goCode += fmt.Sprintf("        int32ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int16List":
				goCode += fmt.Sprintf("        int16ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.IntList":
				goCode += fmt.Sprintf("        intListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int8List":
				goCode += fmt.Sprintf("        int8ListJoin(e.%s, `,`)", field.FieldName)
			default:
				fmt.Println("should be here generateDBMapRowSpicialTypeTrans", field.FieldGoType)
			}
		}
		if field.IsStringList() {
			goCode += fmt.Sprintf("        stringListJoin(e.%s, `,`)", field.FieldName)
		}

		if field.IsNumber() {
			goCode += fmt.Sprintf("        e.%s", field.FieldName)
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "timestamp" {
			goCode += fmt.Sprintf("        formatTime(e.%s)", field.FieldName)
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "datetime" {
			goCode += fmt.Sprintf("        formatTime(e.%s)", field.FieldName)
		}
		if field.IsTimeInt() {
			goCode += fmt.Sprintf("        formatTimeUnix(e.%s)", field.FieldName)
		}
		if field.FieldGoType == "string" {
			goCode += fmt.Sprintf("        e.%s", field.FieldName)
		}

		if idx != len(sd.Fields)-1 {
			goCode += ",\n"
		}
	}

	goCode += ")\n"

	goCode += "    return\n"
	goCode += "}\n"
	return
}
func (sd StructDescription) GenerateInsertForMap() (goCode string) {

	pkList := sd.GetPKFieldList()
	if len(pkList) == 0 || len(pkList) > 2 {
		return ""
	}

	var values string
	var valuesPos string
	for idx, field := range sd.Fields {
		if idx == 0 {
			valuesPos += "\"("
		}

		if field.IsBool() {

			valuesPos += "%d"
		}
		if field.IsInt() {
			valuesPos += "%d"
		}
		if field.IsFloat() {
			valuesPos += "%f"
		}
		if field.IsIntList() {
			valuesPos += "'%s'"
		}
		if field.IsStringList() {
			valuesPos += "'%s'"
		}

		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "timestamp" {
			valuesPos += "%s"
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "datetime" {
			valuesPos += "%s"
		}
		if field.IsTimeInt() {
			valuesPos += "%d"
		}
		if field.FieldGoType == "string" {
			valuesPos += "'%s'"
		}

		if idx != len(sd.Fields)-1 {
			valuesPos += ","
		}
	}
	valuesPos += ")\",\n"
	for idx, field := range sd.Fields {
		if field.IsBool() {
			values += fmt.Sprintf("        bool2int(e.%s)", field.FieldName)
		} else if field.IsIntList() {
			switch field.FieldGoType {
			case "[]int":
				values += fmt.Sprintf("        intListJoin(e.%s, `,`)", field.FieldName)
			case "[]int8":
				values += fmt.Sprintf("        int8ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int16":
				values += fmt.Sprintf("        int16ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int32":
				values += fmt.Sprintf("        int32ListJoin(e.%s, `,`)", field.FieldName)
			case "[]int64":
				values += fmt.Sprintf("        int64ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint8":
				values += fmt.Sprintf("        uint8ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint16":
				values += fmt.Sprintf("        uint16ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint32":
				values += fmt.Sprintf("        uint32ListJoin(e.%s, `,`)", field.FieldName)
			case "[]uint64":
				values += fmt.Sprintf("        uint64ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int32List":
				values += fmt.Sprintf("        int32ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int16List":
				values += fmt.Sprintf("        int16ListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.IntList":
				values += fmt.Sprintf("        intListJoin(e.%s, `,`)", field.FieldName)
			case "goutils.Int8List":
				values += fmt.Sprintf("        int8ListJoin(e.%s, `,`)", field.FieldName)
			default:
				fmt.Println("should be here generateDBMapRowSpicialTypeTrans", field.FieldGoType)
			}
		} else if field.IsStringList() {
			values += fmt.Sprintf("        stringListJoin(e.%s, `,`)", field.FieldName)
		} else if field.IsNumber() {
			values += fmt.Sprintf("        e.%s", field.FieldName)
		} else if field.FieldGoType == "time.Time" && field.GetMysqlType() == "timestamp" {
			values += fmt.Sprintf("        formatTime(e.%s)", field.FieldName)
		} else if field.FieldGoType == "time.Time" && field.GetMysqlType() == "datetime" {
			values += fmt.Sprintf("        formatTime(e.%s)", field.FieldName)
		} else if field.IsTimeInt() {
			values += fmt.Sprintf("        formatTimeUnix(e.%s)", field.FieldName)
		} else if field.FieldGoType == "string" {
			values += fmt.Sprintf("        e.%s", field.FieldName)
		}

		if idx != len(sd.Fields)-1 {
			values += ",\n"
		}
	}

	// values += ")\n"

	goCode += fmt.Sprintf(
		`
func (eMap %sMap) GetInsertSql() (sql string) {
    if len(eMap) == 0 {
        return ""
    }
    var strBuffer bytes.Buffer
    strBuffer.WriteString("insert into %s(%s)values")

    var idx int
	for _, e := range eMap {
        strBuffer.WriteString(fmt.Sprintf(%s %s))
        if idx != len(eMap) - 1 {
             strBuffer.WriteString(",")
        }
        idx++
	}
    return strBuffer.String()
}
`, sd.StructName, sd.GetMysqlTableName(), sd.JoinMysqlFieldNameList(","), valuesPos, values)

	return
}
func (sd StructDescription) GenerateInsertForMapWithArgs() (goCode string) {

	pkList := sd.GetPKFieldList()
	if len(pkList) == 0 || len(pkList) > 2 {
		return ""
	}

	var valuesPos string
	for idx, _ := range sd.GetMysqlFieldList() {
		if idx == 0 {
			valuesPos += "\"("
		}
		valuesPos += "?"

		if idx != len(sd.Fields)-1 {
			valuesPos += ","
		}
	}
	valuesPos += ")\",\n"

	// values += ")\n"

	goCode += fmt.Sprintf(
		`
func (eMap %sMap) GetInsertSqlWithArgs() (sql string,args []interface{}) {
    if len(eMap) == 0 {
        return "", nil
    }
    var strBuffer bytes.Buffer
    strBuffer.WriteString("insert into %s(%s)values")

    var idx int=0
	for _, _ = range eMap {
        strBuffer.WriteString(fmt.Sprintf(%s ))
        if idx != len(eMap) - 1 {
             strBuffer.WriteString(",")
        }
        idx++
	}
    var eArgs []interface{}
    idx = 0
	for _, e := range eMap {
        _,eArgs = e.GetInsertSqlWithArgs()
        if idx ==0 {
            args = make([]interface{},0,len(eMap)*len(eArgs))
         }
        args = append(args, eArgs...)
        idx++
	}

    return strBuffer.String(),args
}
`, sd.StructName, sd.GetMysqlTableName(), sd.JoinMysqlFieldNameList(","), valuesPos)

	return
}
func (sd StructDescription) GenerateUpdate() (goCode string) {
	if len(sd.GetPKFieldList()) == 0 {
		return
	}

	goCode += fmt.Sprintf("func (e %s) GetUpdateSql() (sql string) {\n", sd.StructName)
	goCode += fmt.Sprintf("    if !e.IsFlagDirty(){\n")
	goCode += fmt.Sprintf("        return\n")
	goCode += fmt.Sprintf("    }\n\n")
	goCode += fmt.Sprintf("    if e.IsFlagNew(){\n")
	goCode += fmt.Sprintf("        return\n")
	goCode += fmt.Sprintf("    }\n\n")

	goCode += fmt.Sprintf("    var isFirstField bool = true\n")
	goCode += fmt.Sprintf("    var updateBuffer bytes.Buffer\n")

	for _, field := range sd.Fields {
		if field.IsPK() {
			continue
		}

		goCode += fmt.Sprintf("    if e.Is%sModified(){\n", Camelize(field.FieldName))
		goCode += fmt.Sprintf("        if !isFirstField{\n")
		goCode += fmt.Sprintf("            updateBuffer.WriteString(`,`)\n")
		goCode += fmt.Sprintf("        }\n")
		goCode += fmt.Sprintf("        isFirstField=false\n")
		goCode += fmt.Sprintf("        updateBuffer.WriteString(fmt.Sprintf(`%s=%s`,%s))\n",
			field.GetMysqlFieldName(), field.GetFieldPosStr(), field.GetFieldPosValue())

		goCode += "    }\n"

		// if idx != len(sd.Fields)-1 {
		// 	goCode += ","
		// }
	}

	goCode += fmt.Sprintf("    sql=fmt.Sprintf(`update %s set %%s where %s`, updateBuffer.String(),%s)\n", sd.StructName, sd.GetWherePosStr(), sd.GetWherePosValue())
	goCode += "    return\n"
	goCode += "}\n"
	return
}
func (sd StructDescription) GenerateUpdateWithArgs() (goCode string) {
	pkList := sd.GetPKFieldList()
	if len(pkList) == 0 {
		return
	}

	goCode += fmt.Sprintf("func (e %s) GetUpdateSqlArgs() (sql string,args []interface{}) {\n", sd.StructName)
	goCode += fmt.Sprintf("    if !e.IsFlagDirty(){\n")
	goCode += fmt.Sprintf("        return\n")
	goCode += fmt.Sprintf("    }\n\n")
	goCode += fmt.Sprintf("    if e.IsFlagNew(){\n")
	goCode += fmt.Sprintf("        return\n")
	goCode += fmt.Sprintf("    }\n\n")

	goCode += fmt.Sprintf("    var isFirstField bool = true\n")
	goCode += fmt.Sprintf("    var updateBuffer bytes.Buffer\n")
	goCode += fmt.Sprintf(fmt.Sprintf("    args = make([]interface{},0,e.flag.GetTrueLen()+%d)\n\n", len(pkList)))

	for _, field := range sd.Fields {
		if field.IsPK() {
			continue
		}

		goCode += fmt.Sprintf("    if e.Is%sModified(){\n", Camelize(field.FieldName))
		goCode += fmt.Sprintf("        if !isFirstField{\n")
		goCode += fmt.Sprintf("            updateBuffer.WriteString(`,`)\n")
		goCode += fmt.Sprintf("        }\n")
		goCode += fmt.Sprintf("        isFirstField=false\n")
		goCode += fmt.Sprintf("        updateBuffer.WriteString(`%s=?`)\n", field.GetMysqlFieldName())
		goCode += fmt.Sprintf("        args = append(args,%s)\n", field.GetFieldPosValue4Args())

		goCode += "    }\n"

		// if idx != len(sd.Fields)-1 {
		// 	goCode += ","
		// }
	}

	for _, field := range pkList {
		goCode += fmt.Sprintf("        args = append(args,%s)\n", field.GetFieldPosValue4Args())
	}

	goCode += fmt.Sprintf("    sql=fmt.Sprintf(`update %s set %%s where %s`, updateBuffer.String())\n", sd.StructName, sd.GetWherePosStr2())
	goCode += "    return\n"
	goCode += "}\n"
	return
}
func Camelize(s string) (ret string) {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}
func LowerCaseFirstChar(s string) (ret string) {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[0:1]) + s[1:]
}
