package generator

import (
	"strconv"
	"strings"
)

type TagField struct {
	TagKey   string
	TagValue string
	IsKV     bool // if true TagFile 形如 foo=bar,if false ,only TagKey 有用
}
type TagFieldList []TagField

func (l TagFieldList) Contains(field string) bool {
	for _, tagField := range l {
		if tagField.TagKey == field {
			return true
		}
	}
	return false
}
func (l TagFieldList) GetValue(tagKey string) string {
	for _, tagField := range l {
		if tagField.TagKey == tagKey && tagField.IsKV {
			return tagField.TagValue
		}
	}
	return ""
}

// 比如`mysql:"type=int,pk"` 分别对应 TagField{TagKey:"type",TagValue:"int",IsKV:true} {TagKey:"pk"}

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
// func parseTag(tag string) (tagFields []TagField) {
// 	for _, tagFieldStr := range strings.Split(tag, ",") {
// 		i := strings.Index(tagFieldStr, "=")
// 		if i != -1 {
// 			tagField := TagField{
// 				IsKV:     true,
// 				TagKey:   tagFieldStr[:i],
// 				TagValue: tagFieldStr[i+1:],
// 			}
// 			tagFields = append(tagFields, tagField)
// 		} else {
// 			tagField := TagField{
// 				IsKV:   false,
// 				TagKey: tagFieldStr,
// 			}
// 			tagFields = append(tagFields, tagField)

// 		}

// 	}
// 	return
// }
func parseTag(tag string) (tagFields []TagField) {
	// When modifying this code, also update the validateStructTag code
	// in golang.org/x/tools/cmd/vet/structtag.go.

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		quoteChar := ""
		for i < len(tag) && tag[i] > ' ' && (tag[i] != ',' || (tag[i] == ',' && quoteChar != "")) && tag[i] != 0x7f {
			if tag[i] == '"' || tag[i] == '\'' {
				if i == 0 || (i != 0 && tag[i-1] != '\\') {
					if quoteChar == "" {
						quoteChar = tag[i : i+1]
					} else if quoteChar == tag[i:i+1] {
						quoteChar = ""
					}
				}
			}
			i++
		}
		if i == 0 || (i < len(tag) && tag[i] != ',') {
			break
		}

		j := strings.Index(tag[:i], "=")
		if j != -1 {
			unQuoted, err := strconv.Unquote(tag[j+1 : i])
			if err != nil {
				unQuoted = tag[j+1 : i]
				if unQuoted[0] == '\'' && unQuoted[len(unQuoted)-1] == '\'' {
					unQuoted = unQuoted[1 : len(unQuoted)-1]
				}

			}

			tagField := TagField{
				IsKV:     true,
				TagKey:   tag[:j],
				TagValue: unQuoted,
			}
			tagFields = append(tagFields, tagField)
		} else {
			tagField := TagField{
				IsKV:   false,
				TagKey: tag[:i],
			}
			tagFields = append(tagFields, tagField)
		}

		if i < len(tag) {
			tag = tag[i+1:]
		} else {
			break
		}

	}
	return
}
