package main

import (
	"fmt"
	user "github.com/0studio/go_service_generator/example"
)

func main() {
	var u user.User
	u.SetAge(1)
	u.SetId(1)
	u.SetName("hello")
	fmt.Println(u.GetUpdateSql())

}
