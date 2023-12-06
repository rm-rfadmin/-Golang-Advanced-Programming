package one

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

//// 生产者
//func Producer(factor int, out chan<- int) {
//	for i := 0; ; i++ {
//		out <- i * factor
//	}
//}
//
//func Consumer(in <-chan int) {
//	for v := range in {
//		fmt.Println(v)
//	}
//}
//
//func TestProducer(t *testing.T) {
//	ch := make(chan int, 64)
//
//	go Producer(3, ch) // 生成3的倍数的序列
//	go Producer(5, ch) // 生产5的倍数序列
//	go Consumer(ch)    // 消费生产的序列
//	time.Sleep(5 * time.Second)
//}

// --------------------------
// 发布订阅
type (
	subscriber chan interface{}         // 订阅者为一个管道
	topicFunc  func(v interface{}) bool // 主题为一个过滤器
)

// 发布者对象
type Publisher struct {
	m           sync.RWMutex             // 读写锁
	buffer      int                      // 订阅队列的缓存大小
	timeout     time.Duration            // 发布超时时间
	subscribers map[subscriber]topicFunc // 订阅者消息
}

// 构建一个发布者对象，可以设置发布超时时间和缓存队列的长度
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 天加一个新的订阅者，订阅全部主题
func (p *Publisher) Subscriber() chan interface{} {
	return p.SubscribeTopic(nil)
}

// 添加一个新的订阅者，订阅过滤器筛选后的主题
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
}

// 退出订阅
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// 发布一个主题
func (p *Publisher) Publish(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// 关闭发布者对象，同事关闭所有的订阅者通道
func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

func TestPubSub(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscriber()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})
	p.Publish("hello, world!")
	p.Publish("hello, golang!")

	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	time.Sleep(3 * time.Second)

}

func worker1(cannel chan bool) {
	for {
		select {
		default:
			fmt.Println("hello")
		case <-cannel:
			// 退出
		}
	}
}

func TestWorker(t *testing.T) {
	cannel := make(chan bool)
	go worker1(cannel)

	time.Sleep(time.Second)
	cannel <- true
}

func worker2(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func TestWorer2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker2(ctx, &wg)

	}
	time.Sleep(time.Second)
	cancel()
	wg.Wait()
}
