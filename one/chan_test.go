package one

import "testing"

var done1 = make(chan bool)
var msg1 string

func aGroutine() {
	msg1 = "你好，世界"
	close(done1)
}

func TestChan2(t *testing.T) {
	go aGroutine()
	<-done1
	t.Log(msg1)
}

var done3 = make(chan bool)
var msg2 string

func aGoroutine3() {
	msg2 = "hello, world"
	<-done3
}

func TestChan3(t *testing.T) {
	go aGoroutine3()
	done3 <- true
	t.Log(msg2)
}

//var limit = make(chan int,3)
//
//func TestChan4(t *testing.T) {
//	//a := make(func(), 3)
//	work := []int
//	for _, w := range work{
//		go func() {
//			limit <- 1
//			w()
//			<-limit
//		}()
//	}
//	select{}
//}
