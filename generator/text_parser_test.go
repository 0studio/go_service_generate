package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveEmptyLineAndCommentLine(t *testing.T) {
	lines := []string{
		"", " ", "  ", "//hello", " //hello", "hi",
	}
	newLines := removeEmptyLineAndCommentLine(lines)
	assert.Equal(t, 1, len(newLines))
	assert.Equal(t, "hi", newLines[0])
}

func TestRemoveCommentPart(t *testing.T) {
	line := "hello // world"
	newLine, comments := splitCommentPart(line)
	assert.Equal(t, "hello", newLine)
	assert.Equal(t, "world", comments)

	line = "hello"
	newLine, comments = splitCommentPart(line)
	assert.Equal(t, "hello", newLine)
	assert.Equal(t, "", comments)

}

func TestParseStructFileContent(t *testing.T) {
	fileContent :=
		"package myPackageName\n" +
			"import(\"time\")\n" +
			"//hell\n" +
			"type Hello struct{//hi\n " +
			"Name string `mysql:\"pk,default=''\"`\n" +
			"Age string `mysql:\"pk,default=''\"`\n" +
			"\n" +
			"}\n" +
			"type World struct{\n" +
			"Name string \n" +
			"}"

	sdList, prop := ParseStructFileContent(fileContent)
	assert.NotEmpty(t, sdList)
	assert.Equal(t, 2, len(sdList))
	assert.Equal(t, "myPackageName", prop.PackageName)
	assert.Equal(t, "Hello", sdList[0].StructName)
	assert.Equal(t, 2, len(sdList[0].Fields))

	assert.Equal(t, "string", sdList[0].Fields[0].FieldGoType)
	assert.Equal(t, "Name", sdList[0].Fields[0].FieldName)
	assert.Equal(t, "mysql:\"pk,default=''\"", sdList[0].Fields[0].TagString)
	assert.Equal(t, 2, len(sdList[0].Fields[0].MysqlTagFieldList))
	assert.Equal(t, "pk", sdList[0].Fields[0].MysqlTagFieldList[0].TagKey)
	assert.Equal(t, "default", sdList[0].Fields[0].MysqlTagFieldList[1].TagKey)

}
