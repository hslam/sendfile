# sendfile
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hslam/sendfile)](https://pkg.go.dev/github.com/hslam/sendfile)
[![Build Status](https://github.com/hslam/sendfile/workflows/build/badge.svg)](https://github.com/hslam/sendfile/actions)
[![codecov](https://codecov.io/gh/hslam/sendfile/branch/master/graph/badge.svg)](https://codecov.io/gh/hslam/sendfile)
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
		conn, _ := lis.Accept()
		defer conn.Close()
		buf := make([]byte, len(contents))
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))
		close(done)
	}()
	conn, _ := net.Dial("tcp", "127.0.0.1:9999")
	if _, err = sendfile.SendFile(conn, int(srcFile.Fd()), 0, int64(len(contents))); err != nil {
		fmt.Println(err)
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


