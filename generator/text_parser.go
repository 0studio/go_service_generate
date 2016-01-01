package generator

import (
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
)

func ParseStructFile(fileName string) (sdList []StructDescription, prop Property) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	return ParseStructFileContent(string(data))

}
func ParseStructFileContent(content string) (sdList []StructDescription, prop Property) {
	var structFirstLineRegex *regexp.Regexp
	structFirstLineRegex, _ = regexp.Compile("^type ([A-Za-z0-9_]+) struct[ \t]?{")
	var isStartingParseStruct bool
	var sd StructDescription
	var lineComments string

	content = strings.Replace(content, "\r", "", -1) //
	lines := strings.Split(content, "\n")
	lines = removeEmptyLineAndCommentLine(lines)
	for _, line := range lines {
		if prop.PackageName == "" {
			idx := strings.Index(line, "package ")
			if idx == 0 {
				prop.PackageName = strings.TrimSpace(line[len("package ")+idx:])
			}
			continue
		}
		line, lineComments = splitCommentPart(line)
		if !isStartingParseStruct {
			matched := structFirstLineRegex.FindStringSubmatch(line)
			if len(matched) >= 2 {
				sd.Reset()
				sd.StructName = matched[1]
				isStartingParseStruct = true
			}
			continue
		}
		if line == "}" { // 解析完一个struct
			isStartingParseStruct = false
			if sd.StructName != "" && len(sd.Fields) != 0 {
				sdList = append(sdList, sd)
				sd.Reset()
			}

			continue
		}
		if isStartingParseStruct { // 在解析field 中
			fd := FieldDescriptoin{}
			tagStartIdx := strings.Index(line, "`")
			tagEndIdx := strings.LastIndex(line, "`")
			if tagStartIdx != -1 && tagEndIdx != -1 && tagEndIdx != tagStartIdx {
				fd.TagString = line[tagStartIdx+1 : tagEndIdx]
				fd.MysqlTagFieldList = parseTag(reflect.StructTag(fd.TagString).Get("mysql"))
				fd.ProtoBufTagFieldList = parseTag(reflect.StructTag(fd.TagString).Get("pb"))
				fd.GoTagFieldList = parseTag(reflect.StructTag(fd.TagString).Get("go"))

				line = line[:tagStartIdx]
			}
			tokens := strings.Fields(line)
			if len(tokens) < 2 {
				continue
			}
			fd.FieldName = tokens[0]
			fd.FieldGoType = tokens[1]
			fd.Comments = lineComments
			sd.Fields = append(sd.Fields, fd)
		}

	}

	return
}
func splitCommentPart(line string) (content, comments string) {
	idx := strings.Index(line, "//")
	if idx != -1 {
		return strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+2:])
	}
	return line, ""
}
func removeEmptyLineAndCommentLine(lines []string) (newLines []string) {
	newLines = make([]string, 0, len(lines))
	for _, line := range lines {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" {
			continue
		}
		if len(trimLine) > 2 && trimLine[0] == '/' && trimLine[1] == '/' { // if line starts with // ,it is comments line
			continue
		}
		newLines = append(newLines, trimLine)
	}
	return
}
