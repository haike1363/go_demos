package main

/*
// 连接第三方库
// #cgo CFLAGS: -DPNG_DEBUG=1
// #cgo linux CFLAGS: -DLINUX=1
// #cgo LDFLAGS: -lpng
// 或者
// #cgo pkg-config: png
// #include <png.h>
#include <stdio.h>
int cAdd(int a, int b) {
    return a + b;
}

void hello() {
    printf("hello, Cgo! -- From C world.\n");
}
*/
import "C"
import "fmt"

func main() {
    C.hello()
    var added int = C.cAdd(1, 2)
    fmt.Println(added)
}
