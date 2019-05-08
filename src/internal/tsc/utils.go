// +build windows
// +build amd64

package tsc

import (
	"bytes"
	"io/ioutil"
	"syscall"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// findProcAddr search and catch the proc named procName from handle provided.
// if no proc found, this will return an error
func findProcAddr(handle syscall.Handle, procName string) (proc uintptr, err error) {
	proc, err = syscall.GetProcAddress(handle, procName)
	return
}

func stringToUintptr(content string) uintptr {
	by, err := syscall.BytePtrFromString(content)
	if err != nil {
		panic("syscall: string with NUL passed to StringByteSlice")
	}
	return uintptr(unsafe.Pointer(by))
}

func bytesToUintptr(content []byte) uintptr {
	by := &content[0]

	for i := 0; i < len(content); i++ {
		if content[i] == 0 {
			panic("syscall: string with NUL passed to StringByteSlice")
		}
	}
	return uintptr(unsafe.Pointer(by))
}

func utf8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
