# sendfile
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
	"io/ioutil"
	"net"
	"os"
)

func main() {
	// Set up source data file.
	srcName := "srcfile"
	srcFile, err := os.Create(srcName)
	if err != nil {
		panic(err)
	}
	defer os.Remove(srcName)
	defer srcFile.Close()
	contents := "Hello world"
	offset := 10
	if offset > 0 {
		srcFile.Write(make([]byte, offset))
	}
	srcFile.Write([]byte(contents))
	srcFile.Sync()

	done := make(chan bool)

	// Start server listening on a socket.
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		b, err := ioutil.ReadAll(conn)
		if string(b) != contents {
			fmt.Printf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		} else {
			fmt.Println(string(b))
		}
		done <- true
	}()

	// Send source file to server.
	conn, err := net.Dial("tcp", lis.Addr().String())
	if err != nil {
		panic(err)
	}
	_, err = sendfile.SendFile(conn, int(srcFile.Fd()), int64(offset), int64(len(contents)))
	if err != nil {
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


