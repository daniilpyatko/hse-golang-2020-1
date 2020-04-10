/*
	https://github.com/mailru/easyjson/blob/master/jlexer/bytestostr.go

*/

package main

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

// bytesToStr сорздаёт строку, указывающую на слайс байт, чтобы избежать копирования.
//
// Warning: the string returned by the function should be used with care, as the whole input data
// chunk may be either blocked from being freed by GC because of a single string or the buffer.Data
// may be garbage-collected even when the string exists.
func bytesToStr(data []byte) string {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	fmt.Printf("type: %T, value: %+v\n", h, h)
	fmt.Printf("type: %T, value: %+v\n", h.Data, h.Data)
	shdr := reflect.StringHeader{Data: h.Data, Len: h.Len}
	fmt.Printf("type: %T, value: %+v\n", shdr, shdr)
	return *(*string)(unsafe.Pointer(&shdr))
}

func resData() string {

	data := []byte(`some test`)
	str := bytesToStr(data)
	//str := string(data)
	fmt.Printf("type: %T, value: %+v\n", str, str)

	fmt.Printf("String before modify: %s\n", str)
	data[1] = 65
	fmt.Printf("String after modify: %s\n", str)

	return str
}

func main() {
	a := resData()
	time.Sleep(10 * time.Second)
	fmt.Printf("GetData: %s", a)
}
