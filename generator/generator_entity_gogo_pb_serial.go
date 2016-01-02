package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func (sd StructDescription) GenerateEntitySerialUnSerial(property Property, srcDir string) {
	outputF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("entity_serial_stub.go")), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer outputF.Close()
	outputF.WriteString(fmt.Sprintf(`// do not edit this file ,this is generated by tools(https://github.com/0studio/go_service_generator)
package %s

import (
    "time"
    "github.com/0studio/goutils"
    key "github.com/0studio/storage_key"
	"github.com/gogo/protobuf/proto"
)
var ___importTimeS time.Time
var ___importGoutilsS goutils.Int32List
var ___importKeyS key.KeyUint64

`, property.PackageName))
	outputF.WriteString(sd.generateEntitySerial(property, srcDir))
	outputF.WriteString(sd.generateEntityUnSerial(property, srcDir))

}
func (sd StructDescription) generateEntitySerial(property Property, srcDir string) (goCode string) {
	var fieldStr string
	for _, fd := range sd.Fields {
		if fd.IsEqualableType() {
			fieldStr += fmt.Sprintf("    pb.%s = %s(e.%s)\n", Camelize(fd.FieldName), fd.getPBInfo().goType, fd.FieldName)
		} else if fd.IsTime() {
			fieldStr += fmt.Sprintf("    pb.%s = %s(formatTimeUnix(e.%s))\n", Camelize(fd.FieldName), fd.getPBInfo().goType, fd.FieldName)
		} else if fd.IsIntList() {
			switch fd.FieldGoType {
			case "[]int":
				fieldStr += fmt.Sprintf("    pb.%s = intList2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)
			case "[]int8":
				fieldStr += fmt.Sprintf("    pb.%s = int8List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)
			case "[]int16":
				fieldStr += fmt.Sprintf("    pb.%s = int16List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)

			case "[]int32":
				fieldStr += fmt.Sprintf("    pb.%s = int32List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)

			case "[]int64":
				fieldStr += fmt.Sprintf("    pb.%s = int64List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)

			case "[]uint8":
				fieldStr += fmt.Sprintf("    pb.%s = uint8List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)

			case "[]uint16":
				fieldStr += fmt.Sprintf("    pb.%s = uint16List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)

			case "[]uint32":
				fieldStr += fmt.Sprintf("    pb.%s = uint32List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)

			case "[]uint64":
				fieldStr += fmt.Sprintf("    pb.%s = uint64List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)

			case "goutils.Int32List":
				fieldStr += fmt.Sprintf("    pb.%s = int32List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)
			case "goutils.Int16List":
				fieldStr += fmt.Sprintf("    pb.%s = int16List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)
			case "goutils.IntList":
				fieldStr += fmt.Sprintf("    pb.%s = intList2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)
			case "goutils.Int8List":
				fieldStr += fmt.Sprintf("    pb.%s = int8List2%sList(e.%s)\n",
					Camelize(fd.FieldName), fd.getPBInfo().dataType, fd.FieldName)
			default:
				fmt.Println("should be here GetFieldPosValue", fd.FieldGoType, fd.FieldName)
			}
		} else {
			fmt.Printf(" Serial() %s.%s is not handled \n", sd.StructName, fd.FieldName)
		}

	}
	return fmt.Sprintf(
		`
func (e %s) Serial() (data []byte) {
	var pb %sPB
%s
    data, _ = pb.Marshal()
    return
}
`, sd.StructName, sd.StructName, fieldStr)

}
func (sd StructDescription) generateEntityUnSerial(property Property, srcDir string) (goCode string) {
	var fieldStr string
	for _, fd := range sd.Fields {
		if fd.IsEqualableType() {
			fieldStr += fmt.Sprintf("    e.%s = %s(pb.Get%s())\n", fd.FieldName, fd.FieldGoType, Camelize(fd.FieldName))
		} else if fd.IsTime() {
			fieldStr += fmt.Sprintf("    e.%s = newTime(pb.Get%s())\n", fd.FieldName, Camelize(fd.FieldName))
		} else if fd.IsIntList() {
			switch fd.FieldGoType {
			case "[]int":
				fieldStr += fmt.Sprintf("    e.%s = %sList2intList(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))
			case "[]int8":
				fieldStr += fmt.Sprintf("    e.%s = %sList2int8List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))
			case "[]int16":
				fieldStr += fmt.Sprintf("    e.%s = %sList2int16List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))

			case "[]int32":
				fieldStr += fmt.Sprintf("    e.%s = %sList2int32List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))

			case "[]int64":
				fieldStr += fmt.Sprintf("    e.%s = %sList2int64List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))

			case "[]uint8":
				fieldStr += fmt.Sprintf("    e.%s = %sList2uint8List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))

			case "[]uint16":
				fieldStr += fmt.Sprintf("    e.%s = %sList2uint16List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))

			case "[]uint32":
				fieldStr += fmt.Sprintf("    e.%s = %sList2uint32List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))

			case "[]uint64":
				fieldStr += fmt.Sprintf("    e.%s = %sList2uint64List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))

			case "goutils.Int32List":
				fieldStr += fmt.Sprintf("    e.%s = %sList2int32List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))
			case "goutils.Int16List":
				fieldStr += fmt.Sprintf("    e.%s = %sList2int16List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))
			case "goutils.IntList":
				fieldStr += fmt.Sprintf("    e.%s = %sList2intList(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))
			case "goutils.Int8List":
				fieldStr += fmt.Sprintf("    e.%s = %sList2int8List(pb.Get%s())\n",
					fd.FieldName, fd.getPBInfo().dataType, Camelize(fd.FieldName))
			default:
				fmt.Println("should be here GetFieldPosValue", fd.FieldGoType, fd.FieldName)
			}
		} else {
			fmt.Printf(" Serial() %s.%s is not handled \n", sd.StructName, fd.FieldName)
		}

	}
	return fmt.Sprintf(
		`
func (e *%s) UnSerial(data []byte) bool {
	var pb %sPB
	err := proto.Unmarshal(data, &pb)
	if err != nil {
		return false
	}

%s
    return true
}

`, sd.StructName, sd.StructName, fieldStr)

}