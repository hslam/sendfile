# sendfile
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/sendfile)](https://pkg.go.dev/github.com/hslam/sendfile)
[![Build Status](https://travis-ci.org/hslam/sendfile.svg?branch=master)](https://travis-ci.org/hslam/sendfile)
[![Go Report Card](https://goreportcard.com/badge/github.com/hslam/sendfile)](https://goreportcard.com/report/github.com/hslam/sendfile)
[![LICENSE](https://img.shields.io/github/license/hslam/sendfile.svg?style=flat-square)](https://github.com/hslam/sendfile/blob/master/LICENSE)

Package sendfile wraps the sendfile system call.

## Get started

### Install
```
go get github.com/hslam/sendfile
```
### Import
```
import "github.com/hslam/sendfile"
```
### Usage
#### Example
```go
package main

import (
	"fmt"
	"github.com/hslam/sendfile"
	"net"
	"os"
)

func main() {
	srcName := "srcfile"
	srcFile, err := os.Create(srcName)
	if err != nil {
		panic(err)
	}
	defer os.Remove(srcName)
	defer srcFile.Close()
	contents := "Hello world"
	srcFile.Write([]byte(contents))
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	done := make(chan bool)
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		buf := make([]byte, len(contents))
		n, _ := conn.Read(buf)
		if string(buf[:n]) != contents {
			fmt.Printf("contents not transmitted: got %s (len=%d), want %s\n", string(buf[:n]), n, contents)
		} else {
			fmt.Println(string(buf[:n]))
		}
		done <- true
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	if _, err = sendfile.SendFile(conn, int(srcFile.Fd()), 0, int64(len(contents))); err != nil {
		panic(err)
	}
	conn.Close()
	<-done
}
```

### Output
```
Hello world
```

### License
This package is licensed under a MIT license (Copyright (c) 2020 Meng Huang)


### Author
sendfile was written by Meng Huang.


