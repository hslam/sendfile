// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build darwin linux dragonfly freebsd netbsd openbsd

package sendfile

import (
	"net"
	"syscall"
)

// SendFile wraps the sendfile system call.
func SendFile(conn net.Conn, src int, pos, remain int64) (written int64, err error) {
	var dst int
	if syscallConn, ok := conn.(syscall.Conn); ok {
		raw, err := syscallConn.SyscallConn()
		if err != nil {
			return sendFile(conn, src, pos, remain, maxSendfileSize)
		}
		raw.Control(func(fd uintptr) {
			dst = int(fd)
		})
	} else {
		return sendFile(conn, src, pos, remain, maxSendfileSize)
	}
	for remain > 0 {
		n := maxSendfileSize
		if int(remain) < maxSendfileSize {
			n = int(remain)
		}
		position := pos
		n, errno := syscall.Sendfile(dst, src, &position, n)
		if n > 0 {
			pos += int64(n)
			written += int64(n)
			remain -= int64(n)
		} else if n == 0 && errno == nil {
			break
		}
		if errno == syscall.EAGAIN {
			continue
		}
		if errno != nil {
			err = errno
			break
		}
	}
	return written, err
}
