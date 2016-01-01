package main

import (
	"fmt"
	"github.com/0studio/go_service_generator/generator"
	"os"
	"path/filepath"
	"strings"
)

// give abc.txt return abc
// give a/b/abc.txt return abc
func getFileName(fileName string) string {
	fileName = filepath.Base(fileName)
	idx := strings.Index(fileName, ".")
	if idx != -1 {
		return fileName[:idx]
	}
	return fileName
}

// go run /main.go example/example_1.go
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("please give a go struct defintion file as params like this : %s\n go_struct.go", os.Args[0])
		return
	}
	goStructFile := os.Args[1]
	srcDir := filepath.Dir(goStructFile)
	// if !strings.HasSuffix(goStructFile, ".go") {
	// 	fmt.Printf("the first param must be a go source file ,and some struct are defined there\n")
	// 	return
	// }
	structDescriptionList, property := generator.ParseStructFile(goStructFile)
	if len(structDescriptionList) == 0 {
		fmt.Println("no struct found in ", goStructFile)
		return
	}

	generator.GenerateUtils(property, srcDir)
	structDescriptionList[0].GenerateEntity(property, srcDir)
	structDescriptionList[0].GenerateDBStorage(property, srcDir)
	structDescriptionList[0].GenerateLRUCache(property, srcDir)

	sqlF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("%s_create_table.sql", getFileName(goStructFile))), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sqlF.Close()
	for _, sd := range structDescriptionList {
		sql, err := sd.GenerateCreateTableSql()
		if err != nil {
			fmt.Println(err)
			continue
		}

		sqlF.WriteString(sql)
		sqlF.WriteString("\n")
	}
}
