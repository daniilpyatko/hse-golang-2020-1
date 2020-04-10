package main

/*
// Всё что в комментари над import "C" является кодом на C code и будет скомпилирован при помощи GCC.
// У вас должен быть установлен GCC

#include <stdio.h>

int multiply(int a, int b) {
	int c = a * b;
	printf("%d\n", c);
	return c;
}
*/
import "C" //это псевдо-пакет, он реализуется компилятором
import "fmt"

func main() {
	a := 2
	b := 3
	// для того чтобы вызвать СИшный крод надо добавить префикс "С."
	// там же туда надо передать СИшные переменные
	res := C.multiply(C.int(a), C.int(b))
	fmt.Printf("Multiply in C (%T): %d * %d = %d\n", res, a, b, int(res))
	fmt.Printf("с-var internals %T = %+v\n", res, res)
}
