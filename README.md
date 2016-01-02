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

give a file like example/user.txt
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
go run main.go  example/user.txt
```

# how to run test in example{1,2,3,4} 
1. create a mysql user
```
 GRANT ALL PRIVILEGES ON *.* TO 'th_dev'@'127.0.0.1'     IDENTIFIED BY 'th_devpass' WITH GRANT OPTION;
```
2. run make 

