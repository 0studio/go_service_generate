package generator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateCreateTableSql(t *testing.T) {
	idFd := FieldDescriptoin{FieldName: "id", FieldGoType: "int", TagString: `mysql:"pk,default=1,type=bigint"`,
		MysqlTagFieldList: TagFieldList{
			TagField{TagKey: "pk"},
			TagField{TagKey: "default", TagValue: "1", IsKV: true},
			TagField{TagKey: "type", TagValue: "bigint", IsKV: true},
		},
	}
	nameFd := FieldDescriptoin{FieldName: "name", FieldGoType: "string", TagString: `mysql:"default='hello',type=varchar(100)"`,
		MysqlTagFieldList: TagFieldList{
			TagField{TagKey: "default", TagValue: "''", IsKV: true},
			TagField{TagKey: "type", TagValue: "varchar(100)", IsKV: true},
		},
	}
	sd := StructDescription{StructName: "test", Fields: []FieldDescriptoin{idFd, nameFd}}
	createTableSql, err := sd.GenerateCreateTableSql()
	assert.NoError(t, err)
	fmt.Println(createTableSql)
	expectSql := "create table if not exists `test`(\n" +
		"`id` bigint NOT NULL DEFAULT 1,\n" +
		"`name` varchar(100) NOT NULL DEFAULT ''\n" +
		",primary key (id)\n" +
		");"
	assert.Equal(t, expectSql, createTableSql)

}
