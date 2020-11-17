// Copyright (c) 2020 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

// +build !darwin,!linux,!windows,!dragonfly,!freebsd,!netbsd,!openbsd

package sendfile

// SendFile wraps the sendfile system call.
func SendFile(conn net.Conn, src int, pos, remain int64) (written int64, err error) {
	return sendFile(conn, src, pos, remain)
}
