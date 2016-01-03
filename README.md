# go_service_generate
 this tool can generate service for db/memcache/lrucache
 you just need defined a golang struct
 and use this tool ,
 it will generate 3 struct implements  an interface like 
```
  // this is generated interface in example4/
 type User4Storage interface {
	Get(id string, now time.Time) (e User4, ok bool)
	Set(e *User4, now time.Time) (ok bool)
	Add(e *User4, now time.Time) bool
	MultiGet(idList []string, now time.Time) (eMap User4Map, ok bool)
	MultiUpdate(eMap User4Map, now time.Time) (ok bool)
	Delete(id string) (ok bool)
	MultiDelete(idList []string) (ok bool)
}
```

how to use

give a file like example/user4.txt

```
// this an example in example4/user.txt
package user

import (
	"time"
)

type User4 struct {
	id   string    `mysql:"pk,defalut='',type=varchar(100)"` // id
	name string `mysql:"defalut='hello',name=helloName,type=varchar(10)"`

}
```

then run 

```
go run main.go  example4/user.txt
```
then it will generate some file in example4/
```
   build_proto.sh
   entity_serial_stub.go
   entity_stub.go
   example4.test
   serial.pb.go
   serial.proto
   storage_db_stub.go
   storage_db_stub_test.go
   storage_lru_cache_stub.go
   storage_mc_stub.go
   storage_proxy_stub.go
   user_create_table.sql
   utils_stub.go
```
before that 
you need make sure you have install protobuf(on mac:brew install protobuf)
and run these commands on your shell
```
	go get github.com/0studio/goutils
	go get github.com/0studio/bit
	go get github.com/0studio/storage_key
	go get github.com/0studio/databasetemplate
	go get github.com/0studio/lru
	go get github.com/dropbox/godropbox/memcache
	go get github.com/0studio/logger
	go get -u github.com/gogo/protobuf/proto
	go install github.com/gogo/protobuf/proto
	go get -u github.com/gogo/protobuf/protoc-gen-gogo
	go install github.com/gogo/protobuf/protoc-gen-gogo
	go get -u github.com/gogo/protobuf/gogoproto
	go install  github.com/gogo/protobuf/gogoproto
    
```

# how to run test in example{1,2,3,4} 
1. create a mysql user
```
 GRANT ALL PRIVILEGES ON *.* TO 'th_dev'@'127.0.0.1'     IDENTIFIED BY 'th_devpass' WITH GRANT OPTION;
```
2. run make 

