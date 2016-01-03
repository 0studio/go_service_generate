# -*- coding:utf-8 -*-
.PHONY: demo
demo:
	go run main.go example/user.txt
	cd example;go build;go test
	go run main.go example2/user.txt
	cd example2;go build;go test
	go run main.go example3/user.txt
	cd example3;go build;go test
	go run main.go example4/user.txt
	cd example4;go build; go test

get-deps:
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


