package generator

import (
	"errors"
	"fmt"
	"strings"
)

var DefaultMysqlTypeMap map[string]string = map[string]string{
	"bool":              "tinyint",
	"int":               "int",
	"int8":              "tinyint",
	"int16":             "smallint",
	"int32":             "int",
	"int64":             "bigint",
	"uint8":             "tinyint unsigned",
	"uint16":            "smallint unsigned",
	"uint32":            "int unsigned",
	"uint64":            "bigint unsigned",
	"float32":           "float",
	"float64":           "double",
	"string":            "varchar(255)",
	"time.Time":         "timestamp",
	"[]int":             "varchar(255)",
	"[]int32":           "varchar(255)",
	"[]int8":            "varchar(255)",
	"[]int16":           "varchar(255)",
	"[]int64":           "varchar(255)",
	"[]uint32":          "varchar(255)",
	"[]uint8":           "varchar(255)",
	"[]uint16":          "varchar(255)",
	"[]uint64":          "varchar(255)",
	"[]string":          "varchar(255)",
	"goutils.Int32List": "varchar(255)",
	"goutils.Int16List": "varchar(255)",
	"goutils.IntList":   "varchar(255)",
	"goutils.Int8List":  "varchar(255)",
	"key.KeyUint64":     "bigint",
	"key.KeyInt":        "int",
	"key.KeyInt32":      "int",
	"key.String":        "varchar(255)",
	"key.KeyString":     "varchar(255)",
}

var DefaultMysqlDefaultValueMap map[string]string = map[string]string{
	"bool":              "0",
	"int":               "0",
	"int8":              "0",
	"int16":             "0",
	"int32":             "0",
	"int64":             "0",
	"uint8":             "0",
	"uint16":            "0",
	"uint32":            "0",
	"uint64":            "0",
	"float32":           "0",
	"float64":           "0",
	"string":            "''",
	"time.Time":         "0",
	"[]int":             "''",
	"[]int32":           "''",
	"[]int8":            "''",
	"[]int16":           "''",
	"[]int64":           "''",
	"[]uint32":          "''",
	"[]uint8":           "''",
	"[]uint16":          "''",
	"[]uint64":          "''",
	"[]string":          "''",
	"goutils.Int32List": "''",
	"goutils.Int16List": "''",
	"goutils.IntList":   "''",
	"goutils.Int8List":  "''",
	"key.KeyUint64":     "0",
	"key.KeyInt":        "0",
	"key.KeyInt32":      "0",
	"key.String":        "''",
	"key.KeyString":     "''",
}

type Property struct {
	PackageName string
}
type FieldDescriptoin struct {
	FieldName            string
	FieldGoType          string
	TagString            string
	MysqlTagFieldList    TagFieldList
	GoTagFieldList       TagFieldList
	ProtoBufTagFieldList TagFieldList
	Comments             string
}

func (fd FieldDescriptoin) IsString() bool {
	if fd.FieldGoType == "string" {
		return true
	}
	if fd.FieldGoType == "key.String" {
		return true
	}
	if fd.FieldGoType == "key.KeyString" {
		return true
	}

	return false

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
		fd.FieldGoType == "uint64" ||
		fd.FieldGoType == "key.KeyUint64" ||
		fd.FieldGoType == "key.KeyInt" ||
		fd.FieldGoType == "key.KeyInt32" {
		return true
	}
	return false
}
func (fd FieldDescriptoin) IsStringList() bool {
	if fd.FieldGoType == "[]string" {
		return true
	}
	return false

}
func (fd FieldDescriptoin) IsIntList() bool {
	if fd.FieldGoType == "[]int" ||
		fd.FieldGoType == "[]int8" ||
		fd.FieldGoType == "[]int16" ||
		fd.FieldGoType == "[]int32" ||
		fd.FieldGoType == "[]int64" ||
		fd.FieldGoType == "[]uint8" ||
		fd.FieldGoType == "[]uint16" ||
		fd.FieldGoType == "[]uint32" ||
		fd.FieldGoType == "[]uint64" ||
		fd.FieldGoType == "goutils.Int32List" ||
		fd.FieldGoType == "goutils.Int16List" ||
		fd.FieldGoType == "goutils.IntList" ||
		fd.FieldGoType == "goutils.Int8List" {

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
	if fd.FieldGoType == "float32" ||
		fd.FieldGoType == "float64" {
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

func (fd FieldDescriptoin) IsTime() bool {
	if fd.FieldGoType == "time.Time" {
		return true
	}
	return false
}
func (fd FieldDescriptoin) IsTimeInt() bool {
	if fd.FieldGoType == "time.Time" && (fd.GetMysqlType() == "int" || fd.GetMysqlType() == "bigint") {
		return true
	}
	return false
}
func (fd FieldDescriptoin) IsEqualableType() bool {
	if fd.IsNumber() || fd.IsBool() || fd.IsInt() || fd.IsString() {
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
func (fd FieldDescriptoin) GetMysqlKey() string {
	key := fd.MysqlTagFieldList.GetValue("key")
	return key
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
func (fd FieldDescriptoin) GetGoDefalutValue() string {
	mysqlDefault := fd.GoTagFieldList.GetValue("default")
	if mysqlDefault != "" {
		return mysqlDefault
	}
	return ""
}
func (fd FieldDescriptoin) GetFieldPosStr() string {

	if fd.IsBool() {
		return "%d"
	}
	if fd.IsInt() {
		return "%d"
	}
	if fd.IsFloat() {
		return "%f"
	}
	if fd.FieldGoType == "time.Time" && fd.GetMysqlType() == "timestamp" {
		return "%s"
	}
	if fd.FieldGoType == "time.Time" && fd.GetMysqlType() == "datetime" {
		return "%s"
	}
	if fd.IsTimeInt() {
		return "%d"
	}
	if fd.FieldGoType == "string" {
		return "'%s'"
	}
	if fd.IsIntList() {
		return "'%s'"
	}
	if fd.IsStringList() {
		return "'%s'"
	}

	// should be here
	fmt.Println("should be here GetFieldPosStr", fd.FieldName, fd.FieldGoType)
	return "%s"

}

func (fd FieldDescriptoin) GetFieldPosValueWithoutPrefix() string {
	s := fd.GetFieldPosValue()
	s = strings.Replace(s, "e.", "", -1)
	return s
}
func (fd FieldDescriptoin) GetFieldPosValue() string {
	if fd.IsBool() {
		return fmt.Sprintf("bool2int(e.%s)", fd.FieldName)
	}

	if fd.IsNumber() {
		return fmt.Sprintf("e.%s", fd.FieldName)
	}
	if fd.FieldGoType == "time.Time" && fd.GetMysqlType() == "timestamp" {
		return fmt.Sprintf("formatTime(e.%s)", fd.FieldName)
	}
	if fd.FieldGoType == "time.Time" && fd.GetMysqlType() == "datetime" {
		return fmt.Sprintf("formatTime(e.%s)", fd.FieldName)
	}
	if fd.IsTimeInt() {
		return fmt.Sprintf("formatTimeUnix(e.%s)", fd.FieldName)
	}
	if fd.FieldGoType == "string" {
		return fmt.Sprintf("e.%s", fd.FieldName)
	}
	if fd.IsIntList() {
		switch fd.FieldGoType {
		case "[]int":
			return fmt.Sprintf("intListJoin(e.%s, `,`)", fd.FieldName)
		case "[]int8":
			return fmt.Sprintf("int8ListJoin(e.%s, `,`)", fd.FieldName)
		case "[]int16":
			return fmt.Sprintf("int16ListJoin(e.%s, `,`)", fd.FieldName)
		case "[]int32":
			return fmt.Sprintf("int32ListJoin(e.%s, `,`)", fd.FieldName)
		case "[]int64":
			return fmt.Sprintf("int64ListJoin(e.%s, `,`)", fd.FieldName)
		case "[]uint8":
			return fmt.Sprintf("uint8ListJoin(e.%s, `,`)", fd.FieldName)
		case "[]uint16":
			return fmt.Sprintf("uint16ListJoin(e.%s, `,`)", fd.FieldName)
		case "[]uint32":
			return fmt.Sprintf("uint32ListJoin(e.%s, `,`)", fd.FieldName)
		case "[]uint64":
			return fmt.Sprintf("uint64ListJoin(e.%s, `,`)", fd.FieldName)
		case "goutils.Int32List":
			return fmt.Sprintf("int32ListJoin(e.%s, `,`)", fd.FieldName)
		case "goutils.Int16List":
			return fmt.Sprintf("int16ListJoin(e.%s, `,`)", fd.FieldName)
		case "goutils.IntList":
			return fmt.Sprintf("intListJoin(e.%s, `,`)", fd.FieldName)
		case "goutils.Int8List":
			return fmt.Sprintf("int8ListJoin(e.%s, `,`)", fd.FieldName)
		default:
			fmt.Println("should be here GetFieldPosValue", fd.FieldGoType, fd.FieldName)
		}
	}
	if fd.IsStringList() {
		return fmt.Sprintf("stringListJoin(e.%s, `,`)", fd.FieldName)
	}

	fmt.Println("should be here GetFieldPosValue", fd.FieldGoType, fd.FieldName)

	return ""
}
func (fd FieldDescriptoin) GetFieldListPosValue() string {
	if fd.IsString() {
		switch fd.FieldGoType {
		case "string":
			return fmt.Sprintf(`sroundJoin(%ss,"'",",")`, fd.FieldName)
		case "key.String":
			return fmt.Sprintf(`sroundJoin2(%ss,"'",",")`, fd.FieldName)
		case "key.KeyString":
			return fmt.Sprintf(`sroundJoin2(%ss,"'",",")`, fd.FieldName)
		}
	}
	if fd.IsInt() {
		switch fd.FieldGoType {
		case "int":
			return fmt.Sprintf("intListJoin(%ss, `,`)", fd.FieldName)
		case "int8":
			return fmt.Sprintf("int8ListJoin(%ss, `,`)", fd.FieldName)
		case "int16":
			return fmt.Sprintf("int16ListJoin(%ss, `,`)", fd.FieldName)
		case "int32":
			return fmt.Sprintf("int32ListJoin(%ss, `,`)", fd.FieldName)
		case "int64":
			return fmt.Sprintf("int64ListJoin(%ss, `,`)", fd.FieldName)
		case "uint8":
			return fmt.Sprintf("uint8ListJoin(%ss, `,`)", fd.FieldName)
		case "uint16":
			return fmt.Sprintf("uint16ListJoin(%ss, `,`)", fd.FieldName)
		case "uint32":
			return fmt.Sprintf("uint32ListJoin(%ss, `,`)", fd.FieldName)
		case "uint64":
			return fmt.Sprintf("uint64ListJoin(%ss, `,`)", fd.FieldName)
		case "key.KeyUint64":
			return fmt.Sprintf("%ss.Join( `,`)", fd.FieldName)
		case "key.KeyInt":
			return fmt.Sprintf("%ss.Join( `,`)", fd.FieldName)
		case "key.KeyInt32":
			return fmt.Sprintf("%ss.Join( `,`)", fd.FieldName)
		default:
			fmt.Println("should be here GetFieldPosValue", fd.FieldGoType, fd.FieldName)
		}
	}

	fmt.Println("should be here GetFieldPosValue", fd.FieldGoType, fd.FieldName)

	return ""
}

type StructDescription struct {
	StructName string
	Fields     []FieldDescriptoin
}

func (sd StructDescription) GetSuggestMapName() string {
	return fmt.Sprintf("%sMap", sd.StructName)
}
func (sd StructDescription) GetSuggestMapKey() string {
	pkFieldList := sd.GetPKFieldList()
	if len(pkFieldList) == 1 { // 一个主键的情况
		return pkFieldList[0].FieldGoType
	} else if len(pkFieldList) == 2 { // 两个主键的情况
		return pkFieldList[1].FieldGoType
	}
	return ""
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
func (sd StructDescription) GetPKFieldList() (pkList []FieldDescriptoin) {
	for _, field := range sd.Fields {
		if field.IsPK() {
			pkList = append(pkList, field)
		}
	}
	return

}
func (sd StructDescription) GetWherePosStr2() (sql string) {
	pkList := sd.GetPKFieldList()
	for idx, field := range pkList {
		sql += fmt.Sprintf("%s=?", field.GetMysqlFieldName())
		if idx != len(pkList)-1 {
			sql += " and "

		}
	}
	return
}
func (sd StructDescription) GetWherePosStr() (sql string) {
	pkList := sd.GetPKFieldList()
	for idx, field := range pkList {
		sql += fmt.Sprintf("%s=%s", field.GetMysqlFieldName(), field.GetFieldPosStr())
		if idx != len(pkList)-1 {
			sql += " and "

		}
	}
	return
}
func (sd StructDescription) GetWherePosValueWithoutThisPrefix() (sql string) {
	sql = sd.GetWherePosValue()
	sql = strings.Replace(sql, "e.", "", -1)
	return
}
func (sd StructDescription) GetWherePosValue() (sql string) {
	pkList := sd.GetPKFieldList()
	for idx, field := range pkList {
		sql += field.GetFieldPosValue()
		if idx != len(pkList)-1 {
			sql += " , "

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
		}
	}
	pkList := sd.GetPK()
	if len(pkList) != 0 {
		sql += ",\nprimary key (" + strings.Join(pkList, ",") + ")"
	}
	if sd.Fields[0].GetMysqlKey() != "" {
		sql += ",\n" + sd.Fields[0].GetMysqlKey() + "\n"

	}

	sql += ");"
	return
}

func (sd StructDescription) GenerateCreateTableFunc() (goCode string) {
	goCode += fmt.Sprintf("func (e %s) GetCreateTableSql() (sql string) {\n", sd.StructName)
	sql, _ := sd.GenerateCreateTableSql()
	sql = strings.Replace(sql, "\n", "", -1)

	goCode += fmt.Sprintf("    sql = \"%s\"\n", sql)
	goCode += "    return\n"
	goCode += "}\n"
	return
}

func (sd StructDescription) JoinMysqlFieldNameList(sep string) (s string) {
	for idx, field := range sd.Fields {
		if idx != len(sd.Fields)-1 {
			s += field.GetMysqlFieldName() + ","
		} else {
			s += field.GetMysqlFieldName()
		}
	}
	return
}
