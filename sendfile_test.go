// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package sendfile

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"testing"
)

func TestSendfile(t *testing.T) {
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

	done := make(chan bool, 1)

	// Start server listening on a socket.
	lis, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		t.Error(err)
	}
	defer lis.Close()
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		b, _ := ioutil.ReadAll(conn)
		if string(b) != contents {
			t.Errorf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		}
		done <- true
	}()
	// Send source file to server.
	conn, err := net.Dial("tcp", lis.Addr().String())
	if err != nil {
		t.Error(err)
	}
	var written int64

	written, err = SendFile(conn, int(0), int64(offset), int64(len(contents)))
	if err == nil {
		t.Error()
	} else if written > 0 {
		t.Error()
	}

	written, err = SendFile(conn, int(srcFile.Fd()), int64(offset), int64(len(contents)))
	if err != nil {
		t.Error(err)
	} else if written != int64(len(contents)) {
		t.Errorf("written %d,len %d", written, len(contents))
	}
	conn.Close() //for returning io.EOF
	<-done
}

func TestSendMmap(t *testing.T) {
	testSendMmap(0, t)
	testSendMmap(10, t)
	pagesize := os.Getpagesize()
	testSendMmap(pagesize-1, t)
	testSendMmap(pagesize, t)
	testSendMmap(pagesize+1, t)
}

func testSendMmap(offset int, t *testing.T) {
	srcName := "srcfile"
	srcFile, err := os.Create(srcName)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(srcName)
	defer srcFile.Close()
	contents := "Hello world"
	if offset > 0 {
		srcFile.Write(make([]byte, offset))
	}
	srcFile.Write([]byte(contents))
	srcFile.Sync()

	done := make(chan bool, 1)

	// Start server listening on a socket.
	lis, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		t.Error(err)
	}
	defer lis.Close()
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		b, _ := ioutil.ReadAll(conn)
		if string(b) != contents {
			t.Errorf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		}
		done <- true
	}()
	// Send source file to server.
	conn, err := net.Dial("tcp", lis.Addr().String())
	if err != nil {
		t.Error(err)
	}
	var written int64

	written, err = sendFile(conn, int(0), int64(offset), int64(len(contents)), maxSendfileSize)
	if err == nil {
		t.Error()
	} else if written > 0 {
		t.Error()
	}
	written, err = sendFile(conn, int(srcFile.Fd()), int64(offset), int64(len(contents)), maxSendfileSize)
	if err != nil {
		t.Error(err)
	} else if written != int64(len(contents)) {
		t.Errorf("written %d,len %d", written, len(contents))
	}
	conn.Close() //for returning io.EOF
	written, err = sendFile(conn, int(srcFile.Fd()), int64(offset), int64(len(contents)), maxSendfileSize)
	if err == nil {
		t.Error()
	} else if written > 0 {
		t.Error()
	}
	<-done
}

func TestSendfileSyscallConn(t *testing.T) {
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

	done := make(chan bool, 1)

	// Start server listening on a socket.
	lis, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		t.Error(err)
	}
	defer lis.Close()
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		b, _ := ioutil.ReadAll(conn)
		if string(b) != contents {
			t.Errorf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		}
		done <- true
	}()
	// Send source file to server.
	conn, err := net.Dial("tcp", lis.Addr().String())
	if err != nil {
		t.Error(err)
	}
	type Conn struct {
		net.Conn
	}
	var written int64
	written, err = SendFile(&Conn{conn}, int(srcFile.Fd()), int64(offset), int64(len(contents)))
	if err != nil {
		t.Error(err)
	} else if written != int64(len(contents)) {
		t.Errorf("written %d,len %d", written, len(contents))
	}
	conn.Close() //for returning io.EOF
	<-done
}

type testConn struct {
	net.Conn
}

func (c *testConn) SyscallConn() (syscall.RawConn, error) {
	return nil, errors.New("")
}

func TestSendfileSyscallConnMore(t *testing.T) {
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

	done := make(chan bool, 1)

	// Start server listening on a socket.
	lis, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		t.Error(err)
	}
	defer lis.Close()
	go func() {
		conn, err := lis.Accept()
		if err != nil {
			t.Error(err)
		}
		defer conn.Close()
		b, _ := ioutil.ReadAll(conn)
		if string(b) != contents {
			t.Errorf("contents not transmitted: got %s (len=%d), want %s\n", string(b), len(b), contents)
		}
		done <- true
	}()
	// Send source file to server.
	conn, err := net.Dial("tcp", lis.Addr().String())
	if err != nil {
		t.Error(err)
	}
	var written int64
	written, err = SendFile(&testConn{conn}, int(srcFile.Fd()), int64(offset), int64(len(contents)))
	if err != nil {
		t.Error(err)
	} else if written != int64(len(contents)) {
		t.Errorf("written %d,len %d", written, len(contents))
	}
	conn.Close() //for returning io.EOF
	<-done
}
