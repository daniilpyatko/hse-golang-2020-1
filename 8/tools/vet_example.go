package main

import (
	"fmt"
	"time"
)

func test_error(is_ok bool) error {
	if !is_ok {
		return fmt.Errorf("failed")
	}
	return nil
}

func main() {

	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("%d\n", i)
		}()
	}

	flag := true
	result := test_error(flag)
	fmt.Printf("result is %d\n", result)
	fmt.Printf("%v is %d\n", flag, result)
	time.Sleep(10 * time.Second)
}
