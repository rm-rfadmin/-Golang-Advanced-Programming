package main

/*#include <stdio.h>

static void SayHello(const char* s) {
	puts(s);
}
*/
import "C"

//func main() {
//	println("hello cgo")
//	C.puts(C.CString("hello, world"))
//
//}

func main() {
	C.SayHello(C.CString("hello, world"))
}
