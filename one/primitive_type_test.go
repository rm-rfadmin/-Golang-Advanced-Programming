package one

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
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

func TestSliceOpertion(t *testing.T) {
	var a []int
	// 尾部追加
	a = append(a, 1)
	a = append(a, 1, 2, 3)           // 手动解包
	a = append(a, []int{1, 2, 3}...) // 追加一个切片，切片需要解包

	var b = []int{1, 2, 3}
	// 头部追加
	// 头部追加一般都会导致内存重新分配，而且会导致已有的元素全部复制一次
	// 因此头部添加性能比尾部添加性能差很多
	b = append([]int{0}, b...)
	b = append([]int{-3, -2, -1}, b...)

	//var c []int
	// append 支持链式操作
	//c = append(c[:i], append([]int{x}, a[i:]...)...)  // 在第i个位置插入x
	//c = append(a[:i], append([]int{1,2,3}, a[i:]...)...)  // 在第i个位置插入切片

	// 使用copy() 和 append组合，可以避免使用临时切片
	//a = append(a, 0)  // 切片扩展一个空间
	//copy(a[i+1:], a[i:])  // a[i:] 向后偏移
	//a[i] = x

	// 使用copy和append组合，实现中间插入多个元素
	//a = append(a, x...)       // 为x切片扩展足够的空间
	//copy(a[i+len(x):], a[i:]) // a[i]向后移动Len（x）个位置
	//copy(a[i:], x)            // 复制新切片

	// 删除切片元素
	// 尾部删除效率最快
	a = []int{1, 2, 3}
	a = a[:len(a)-1] // 删除尾部一个元素
	//a = a[:len(a)-N]  // 删除尾部N个元素

	// 删除头部元素
	//a = []int{1,2,3}
	//a = a[1:] // 删除开头第一个元素
	//a = a[N:]  // 删除开头N个元素

	// 使用append 原地完成
	//a = []int{1,2,3}
	//a = append(a[:0], a[1:]...)  // 删除开头第一个元素
	//a = append(a[:0], a[N:]...)  // 删除开头N个元素

	// 使用copy删除开头的元素
	a = []int{1, 2, 3}
	//t.Log(copy(a, a[1:]))
	a = a[:copy(a, a[1:])] // 删除开头1个元素
	//a = a[:copy(a, a[N:]...)]  // 删除开头N个元素

	// 删除中间元素
	a = []int{1, 2, 3}
	//a = append(a[:i], a[i+1:]...)  // 删除第i个元素
	//a = append(a[:1], a[1+1:]...)
	//a = append(a[:i], a[i+N]...)  // 删除第N个元素
	//
	//// 使用copy删除中间的元素
	//a = a[:i+copy(a[i:], a[i+1:])]  // 删除中间1个元素
	//a = a[:i+copy(a[i:], a[i+N:])]  // 删除中间N个元素

	// GC
	a = []int{1, 2, 3}
	a = a[:len(a)-1] // 被删除的最后一个元素依然被引用，可能导致垃圾回收操作呗阻碍
	t.Log(a)

	// 保险的做法
	//f := []int{1,2,3}
	//f[len(f)-1] = nil // 垃圾回收最后一个元素内存
	//f = f[:len(a)-1]
}

func TrimSpace(s []byte) []byte {
	b := s[:0]
	for _, x := range s {
		if x != ' ' {
			b = append(b, x)
		}
	}
	return b
}

func Filter(s []byte, fn func(x byte) bool) []byte {
	b := s[:0]
	for _, x := range s {
		if !fn(x) {
			b = append(b, x)
		}
	}
	return b
}

func FindPhoneNumber(filename string) []byte {
	// 返回的byte指向了保存整个文件的数组
	// 由于切片引用了整个原始数组，导致垃圾回收不能及时释放底层数组的空间
	b, _ := ioutil.ReadFile(filename)
	return regexp.MustCompile("[0-9]+").Find(b)

}

func FindPhoneNumberV1(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	b = regexp.MustCompile("[0-9]+").Find(b)
	// 将原来的切片复制到新的切片上面，有传值的代价，但是换取了切断了原始数据的依赖
	return append([]byte{}, b...)
}

func Swap(a, b int) (int, int) {
	return b, a
}

func Sum(a int, more ...int) int {
	for _, v := range more {
		a += v
	}
	return a
}

func Print(a ...interface{}) {
	fmt.Println(a...)
}

func TestFunc(t *testing.T) {
	var a = []interface{}{123, "abc"}

	Print(a...)
	Print(a)

	inc := Inc()
	t.Log(inc)

	// 因为defe语句引用的都是同一个i变量，所以最后都为3
	for i := 0; i < 3; i++ {
		defer func() { t.Log(i) }()
	}

	// 重新赋值，生成独有变量规避
	for i := 0; i < 3; i++ {
		i := i
		defer func() { t.Log(i) }()
	}

	for i := 0; i < 3; i++ {
		// 通过函数传入i
		// defer 语句马上对调用参数求值
		defer func(i int) { t.Log(i) }(i)
	}

}

func Find(m map[int]int, key int) (value int, ok bool) {
	value, ok = m[key]
	return
}

func Inc() (v int) {
	defer func() { v++ }()
	return 42
}

func twice(x []int) {
	for i := range x {
		x[i] *= 2
	}
}

type IntSliceHeader struct {
	Data []int
	Len  int
	Cap  int
}

func twiceV1(x IntSliceHeader) {
	for i := 0; i < x.Len; i++ {
		x.Data[i] *= 2
	}
}

type Point struct {
	X, Y float64
}

type ClorePoint struct {
	Point
	Color color.RGBA
}

type Cache struct {
	m map[string]string
	sync.Mutex
}

func (p *Cache) Lookup(key string) string {
	p.Lock()
	defer p.Unlock()

	return p.m[key]
}

func TestFuncV1(t *testing.T) {
	s1 := []int{1, 2, 3}
	twice(s1)
	t.Log(s1)

	var cp ClorePoint
	cp.X = 1
	cp.Point.Y = 2
	t.Log(cp.X, cp.Y)

}

type UpperWriter struct {
	io.Writer
}

func (p *UpperWriter) Write(data []byte) (n int, err error) {
	return p.Writer.Write(bytes.ToUpper(data))
}

type UpperString string

func (s UpperString) String() string {
	return strings.ToUpper(string(s))
}

//type fmt.Stringer interface {
//
//}

type TB struct {
	testing.TB
}

func (p *TB) Fatal(args ...interface{}) {
	fmt.Println("TB.Fatal disabled!")
}

func TestInterface(t *testing.T) {
	//fmt.Fprintln(&UpperWriter{os.Stdout}, "hello, world")

	//var (
	//	a io.ReadCloser = (*os.File)(f)
	//)

	var tb testing.TB = new(TB)
	tb.Fatal("hello.playground!")
}

var total struct {
	sync.Mutex
	value int
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= 100; i++ {
		total.Lock()
		total.value += i
		total.Unlock()
	}
}

func TestGoRoutine(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()

	t.Log(total.value)
}

var totalV1 uint64

func workerV1(wg *sync.WaitGroup) {
	defer wg.Done()

	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&totalV1, i)
	}

}

func TestGoRouotine(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go workerV1(&wg)
	go workerV1(&wg)

	wg.Wait()

	fmt.Println(totalV1)
}

type singleton struct{}

var (
	instance    *singleton
	initialized uint32
	mu          sync.Mutex
)

func Instance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}

	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		defer atomic.StoreUint32(&initialized, 1)
		instance = &singleton{}
	}
	return instance
}

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}

	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

var (
	instanceV1 *singleton
	once       sync.Once
)

func InstanceV1() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})

	return instance
}

//var config atomic.Value  // 保存当前配置信息
//
//config.Store(loadConfig())
//
//go func() {
//	for {
//		time.Sleep(time.Second)
//		config.Store(loadConfig())
//	}
//}()
//
//for i := 0; i < 10; i++ {
//	go func() {
//		for r := range requests() {
//			c := config.Load()
//			// ...
//		}
//}
//}

var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func TestSetup(t *testing.T) {
	go setup()
	for !done {
	}
	print(a)
}

func TestChan(t *testing.T) {
	done1 := make(chan int)

	go func() {
		println("你好，世界")
		done1 <- 1
	}()

	<-done1
}

var done2 = make(chan bool)

var msg string

func aGoruntine() {
	msg = "你好，世界"
	done2 <- true
}

func TestChan1(t *testing.T) {
	go aGoruntine()
	<-done2

}
