# -*- coding:utf-8 -*-
.PHONY: demo
demo:
	go run main.go example/user.txt
	cd example;go build
	go run main.go example2/user.txt
	cd example2;go build
	go run main.go example3/user.txt
	cd example3;go build
	go run main.go example4/user.txt
	cd example4;go build

# 	g++ -o nul  -Wall -Wextra  -pedantic -fsyntax-only -Wno-variadic-macros -std=c99 $(INCLUDES) -S ${CHK_SOURCES}
