package main

import (
	"fmt"
)

type my_Struct struct {
	userID   int
	DataJSON []byte
}

func testError(isOk bool) error {
	if !isOk {
		return fmt.Errorf("failed")
	}
	return nil
}

func test() {
	flag := true
	result := testError(flag)
	fmt.Printf("result is %s\n", result)
	fmt.Printf("%v is %v", flag, result)

	s := &my_Struct{}
	fmt.Println(s.userID)
}

func main() {
	test()
}
