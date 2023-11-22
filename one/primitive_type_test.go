package one

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"reflect"
	"testing"
	"unicode/utf8"
	"unsafe"
)

func TestList(t *testing.T) {
	var a [3]int                    // 定义了一个长度为3的数组，元素为0值
	var b = [...]int{1, 2, 3}       // 定义了一个长度为3的数组，元素为1,2,3
	var c = [...]int{2: 3, 1: 2}    // 定义了一个长度为3的数组，元素为0,2
	var d = [...]int{1, 2, 4: 5, 6} // 定义了一个长度为6的数组 1， 2,0,0,5,6

	t.Log(a)
	t.Log(b)
	t.Log(c)
	t.Log(d)

	var e = [...]int{1, 2, 3}
	var f = &e
	t.Log(e[0], e[1])
	t.Log(f[0], f[1])
	t.Log("===================")
	for i, v := range f {
		t.Log(i, v)
	}

	t.Log("============================")
	for i := range e {
		t.Logf("e[%d]: %d\n", i, a[i])
	}

	for i, v := range f {
		t.Logf("f[%d]: %d\n", i, v)
	}

	for i := 0; i < len(f); i++ {
		t.Logf("ff[%d]: %d\n", i, f[i])

	}

	var times [5][0]int
	t.Log(times)
}

func TestStr(t *testing.T) {
	// 字符串数组
	var s1 = [2]string{"hello", "wrold"}
	var s2 = [...]string{"hello", "wrold"}
	var s3 = [...]string{1: "wrod", 2: "hello"}
	t.Log(s1, s2, s3)

	// 结构体数组
	var line1 [2]image.Point
	var line2 = [...]image.Point{image.Point{X: 0, Y: 0}, image.Point{X: 1, Y: 1}}
	var line3 = [...]image.Point{{0, 0}, {1, 1}}
	t.Log(line1, line2, line3)

	// 函数数组
	var decoder1 [2]func(io.Reader) (image.Image, error)
	var decoder2 = [...]func(reader io.Reader) (image.Image, error){
		png.Decode,
		jpeg.Decode,
	}
	t.Log(decoder1, decoder2)

	// 接口数组
	var unknown1 [2]interface{}
	var unknown2 = [...]interface{}{123, "你号"}
	//var unknown2 = []interface{}{123, "你号"}
	t.Log(unknown1, unknown2)

	// 通道数组
	var chanList = [2]chan int{}
	t.Log(chanList)

	// 定义空数组
	// 长度为0的数组在内存中不占用空间，可以用于强调某种特优类型的操作时使用，避免分配格外的内存空间，如通道同步等
	var d [0]int
	var e = [0]int{}
	var f = [...]int{}
	t.Log(d, e, f)

	c1 := make(chan [0]int)
	go func() {
		t.Log("c1")
		c1 <- [0]int{}
	}()
	<-c1

	c2 := make(chan struct{})
	go func() {
		t.Log("C2")
		c2 <- struct{}{} // struct{} 部分是类型，{}表示结构体的值
	}()
	<-c2

	t.Logf("b: %T\n", d)
	t.Logf("b: %#v\n", d)
}

func TestStrOpertion(t *testing.T) {
	s := "hello, world"
	hello := s[:5]
	world := s[7:]
	s1 := "hello, world"[:5]
	s2 := "hello, world"[7:]
	t.Log(hello)
	t.Log(world)
	t.Log(s1)
	t.Log(s2)
	t.Log("len(s):", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)
	t.Log("len(s1):", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len)
	t.Log("len(s2):", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)

	t.Logf("%#v\n", []byte("hello, 世界"))
	// []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
	//--- PASS: TestStrOpertion (0.00s)

	t.Log("\xe4\xb8\x96\xe7\x95\x8c")
	t.Log("\xe4\xb0\x96\xe7\x95\x8cabc")
	for i, c := range "\xe4\x00\x00\xe7\x95\x8cabc" { // 故意损坏编码
		t.Log(i, c)
	}

	for i, c := range []byte("世界abc") {
		t.Log(i, c)
	}

	t.Log("==========================================")
	t.Logf("%#v\n", []rune("世界"))
	t.Logf("%#v\n", string([]rune{19990, 30028})) // ????
}

func forOnString(s string, forBody func(i int, r rune)) {
	for i := 0; len(s) > 0; {
		r, size := utf8.DecodeRuneInString(s)
		forBody(i, r)
		s = s[size:]
		i += size
	}
}

func str2Bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		p[i] = c
	}
	return p
}

//func bytes2str(s []byte) (p string) {
//	data := make([]byte, len(s))
//	for i, c := range s {
//		data[i] = c
//	}

//hdr := (reflect.StringHeader)(unsafe.Pointer(&p))
//hdr.Data = uintptr(unsafe.Pointer(&data[0]))
//hdr.Len = len(s)

//}

func str2runes(s []byte) []rune {
	var p []int32
	for len(s) > 0 {
		r, size := utf8.DecodeRune(s)
		p = append(p, int32(r))
		s = s[size:]
	}
	return []rune(p)
}

func TestStr2Rune(t *testing.T) {
	//b := "\xe4\xb8\x96\xe7\x95\x8c"
	b := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
	r := str2runes(b)
	t.Log(r)

	s := runes2string(r)
	t.Log(s)
}

func runes2string(s []int32) string {
	var p []byte
	buf := make([]byte, 3)
	for _, r := range s {
		n := utf8.EncodeRune(buf, r)
		p = append(p, buf[:n]...)
	}
	return string(p)
}

func TestSlice(t *testing.T) {
	//reflect.SliceHeader{}
	//a []int  // nil 切片，和nil相等，一般用来表示一个不存在的切片
	//b := []int{} // 空切片，和nil不相等，一般用来表示一个空集合
	//c := []int{1, 2, 3} // 有3个元素的切片， len和cap都为3
	//d := c[:2]          // 有两个元素的切片，len2，cap3
	//e := c[0:2:cap(c)]  // 有两个元素的切片，len为2，cap为3
	////t.Log(e)
	////t.Log(d)
	//f := c[:0]  // 有0个元素的切片，len为0，cap3
	//g := make([]int, 3)  // 有3个元素的切片，len和cap都为3
	//h := make([]int, 2, 3) // 有2个元素的切片，len2，cap3
	//i := make([]int, 0, 3)  // 有0个元素的切片， len0 cap3

	// 遍历切片的方式
	a := make([]int, 3)
	for i := range a {
		t.Logf("a[%d]: %d\n", i, a[i])
	}

	for i, v := range a {
		t.Logf("a[%d]: %d\n", i, v)
	}

	for i := 0; i < len(a); i++ {
		t.Logf("a[%d]: %d\n", i, a[i])
	}
}
