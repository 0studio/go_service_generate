package generator

import (
	"errors"
	"fmt"
	"strings"
)

var DefaultMysqlTypeMap map[string]string = map[string]string{
	"bool":      "tinyint",
	"int":       "int",
	"int8":      "tinyint",
	"int16":     "smallint",
	"int32":     "int",
	"int64":     "bigint",
	"uint8":     "tinyint unsigned",
	"uint16":    "smallint unsigned",
	"uint32":    "int unsigned",
	"uint64":    "bigint unsigned",
	"float32":   "float",
	"float64":   "double",
	"string":    "varchar(255)",
	"time.Time": "timestamp",
}
var DefaultProtoBufTypeMap map[string]string = map[string]string{
	"bool":      "bool",
	"int":       "int32",
	"int8":      "int32",
	"int16":     "int32",
	"int32":     "int32",
	"int64":     "int64",
	"uint8":     "uint32",
	"uint16":    "uint32",
	"uint32":    "uint32",
	"uint64":    "uint64",
	"float32":   "double",
	"float64":   "double",
	"string":    "string",
	"time.Time": "int64",
}

var DefaultMysqlDefaultValueMap map[string]string = map[string]string{
	"bool":      "0",
	"int":       "0",
	"int8":      "0",
	"int16":     "0",
	"int32":     "0",
	"int64":     "0",
	"uint8":     "0",
	"uint16":    "0",
	"uint32":    "0",
	"uint64":    "0",
	"float32":   "0",
	"float64":   "0",
	"string":    "''",
	"time.Time": "0",
}

type Property struct {
	PackageName string
}
type FieldDescriptoin struct {
	FieldName            string
	FieldGoType          string
	TagString            string
	MysqlTagFieldList    TagFieldList
	ProtoBufTagFieldList TagFieldList
	Comments             string
}

func (fd FieldDescriptoin) IsInt() bool {
	if fd.FieldGoType == "int" ||
		fd.FieldGoType == "int8" ||
		fd.FieldGoType == "int16" ||
		fd.FieldGoType == "int32" ||
		fd.FieldGoType == "int64" ||
		fd.FieldGoType == "uint8" ||
		fd.FieldGoType == "uint16" ||
		fd.FieldGoType == "uint32" ||
		fd.FieldGoType == "uint64" {
		return true
	}
	return false
}

func (fd FieldDescriptoin) IsBool() bool {
	if fd.FieldGoType == "bool" {
		return true
	}
	return false
}
func (fd FieldDescriptoin) IsFloat() bool {
	if fd.FieldGoType == "flaot32" ||
		fd.FieldGoType == "flaot64" {
		return true
	}
	return false
}
func (fd FieldDescriptoin) IsNumber() bool {
	if fd.IsInt() {
		return true
	}
	if fd.IsFloat() {
		return true
	}
	return false
}
func (fd FieldDescriptoin) IsPK() bool {
	return fd.MysqlTagFieldList.Contains("pk")
}
func (fd FieldDescriptoin) GetMysqlType() string {
	mysqlType := fd.MysqlTagFieldList.GetValue("type")
	if mysqlType != "" {
		return mysqlType
	}
	return DefaultMysqlTypeMap[fd.FieldGoType]
}
func (fd FieldDescriptoin) GetMysqlDefalutValue() string {
	mysqlDefault := fd.MysqlTagFieldList.GetValue("default")
	if mysqlDefault != "" {
		return mysqlDefault
	}
	return DefaultMysqlDefaultValueMap[fd.FieldGoType]

}
func (fd FieldDescriptoin) GetMysqlFieldName() string {
	mysqlFieldName := fd.MysqlTagFieldList.GetValue("name")
	if mysqlFieldName != "" {
		return mysqlFieldName
	}
	return fd.FieldName
}

type StructDescription struct {
	StructName string
	Fields     []FieldDescriptoin
}

func (sd *StructDescription) Reset() {
	sd.StructName = ""
	sd.Fields = nil
}
func (sd StructDescription) GetMysqlTableName() string {
	return sd.StructName
}
func (sd StructDescription) GetPK() (pkList []string) {
	for _, field := range sd.Fields {
		if field.IsPK() {
			pkList = append(pkList, field.GetMysqlFieldName())
		}
	}
	return

}

func (sd StructDescription) GenerateCreateTableSql() (sql string, err error) {
	if len(sd.Fields) == 0 {
		return "", errors.New("no filed found ,generate create table sql error")
	}
	sql += "create table if not exists `" + sd.GetMysqlTableName() + "`(\n"
	for idx, fieldD := range sd.Fields {
		sql += "`" + fieldD.GetMysqlFieldName() + "` " + fieldD.GetMysqlType() + " NOT NULL DEFAULT " + fieldD.GetMysqlDefalutValue()
		if idx != len(sd.Fields)-1 {
			sql += ",\n"
		} else {
			sql += "\n"
		}
	}
	pkList := sd.GetPK()
	if len(pkList) != 0 {
		sql += ",primary key (" + strings.Join(pkList, ",") + ")\n"
	}

	sql += ");"
	return
}

func (sd StructDescription) GenerateCreateTableFunc() (goCode string) {
	goCode += fmt.Sprintf("func (this %s) GetCreateTableSql() (sql string) {\n", sd.StructName)
	sql, _ := sd.GenerateCreateTableSql()
	sql = strings.Replace(sql, "\n", "", -1)

	goCode += fmt.Sprintf("    sql = \"%s\"\n", sql)
	goCode += "    return\n"
	goCode += "}\n"
	return
}

func (sd StructDescription) GenerateInsert() (goCode string) {
	goCode += fmt.Sprintf("func (this %s) GetInsertSql() (sql string) {\n", sd.StructName)
	goCode += fmt.Sprintf("    sql = fmt.Sprintf(\"insert into `%s`(", sd.GetMysqlTableName())
	for idx, field := range sd.Fields {
		if idx != len(sd.Fields)-1 {
			goCode += field.GetMysqlFieldName() + ","
		} else {
			goCode += field.GetMysqlFieldName()
		}
	}
	goCode += ") values ("
	for idx, field := range sd.Fields {
		if field.IsBool() {

			goCode += "%d"
		}
		if field.IsNumber() {
			goCode += "%d"
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "timestamp" {
			goCode += "%s"
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "datetime" {
			goCode += "%s"
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "int" {
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
			goCode += fmt.Sprintf("        bool2int(this.%s)", field.FieldName)
		}

		if field.IsNumber() {
			goCode += fmt.Sprintf("        this.%s", field.FieldName)
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "timestamp" {
			goCode += fmt.Sprintf("        formatTime(this.%s)", field.FieldName)
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "datetime" {
			goCode += fmt.Sprintf("        formatTime(this.%s)", field.FieldName)
		}
		if field.FieldGoType == "time.Time" && field.GetMysqlType() == "int" {
			goCode += fmt.Sprintf("        this.%s.Unix()", field.FieldName)
		}
		if field.FieldGoType == "string" {
			goCode += fmt.Sprintf("        this.%s", field.FieldName)
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
