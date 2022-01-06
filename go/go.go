package g

import (
	"container/list"
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
	"runtime"
	"sync"
)

// one Go per goroutine (goroutine not safe)
type Go struct {
	ChanCb    chan func()
	pendingGo int
}

type LinearGo struct {
	f  func()
	cb func()
}

type LinearContext struct {
	g              *Go
	linearGo       *list.List
	mutexLinearGo  sync.Mutex
	mutexExecution sync.Mutex
}

func New(l int) *Go {
	g := new(Go)
	g.ChanCb = make(chan func(), l)
	return g
}

func (g *Go) Go(f func(), cb func()) {
	g.pendingGo++ //这里go 为啥不会出问题?  ++实在go 外面的,所以这里其实不会被并发的.

	go func() {
		defer func() {
			g.ChanCb <- cb    //go 处理的到的结果,最后在哪里执行了????   只要把结果消息,投递到 想到执行的那个协程  处理的消息队列里 就行. 这个很有意思啊.
			if r := recover(); r != nil {
				if conf.LenStackBuf > 0 {
					buf := make([]byte, conf.LenStackBuf)
					l := runtime.Stack(buf, false)
					log.Error("%v: %s", r, buf[:l])
				} else {
					log.Error("%v", r)
				}
			}
		}()

		f()
	}()
}

func (g *Go) Cb(cb func()) {
	defer func() {
		g.pendingGo--
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	if cb != nil {
		cb()
	}
}

func (g *Go) Close() {
	for g.pendingGo > 0 {  //其实这里也是很巧妙的.
		g.Cb(<-g.ChanCb)
	}
}

func (g *Go) Idle() bool {
	return g.pendingGo == 0
}

func (g *Go) NewLinearContext() *LinearContext {
	c := new(LinearContext)
	c.g = g
	c.linearGo = list.New()
	return c
}

func (c *LinearContext) Go(f func(), cb func()) {
	c.g.pendingGo++

	c.mutexLinearGo.Lock()
	c.linearGo.PushBack(&LinearGo{f: f, cb: cb})
	c.mutexLinearGo.Unlock()

	go func() {
		c.mutexExecution.Lock()
		defer c.mutexExecution.Unlock()

		c.mutexLinearGo.Lock()
		e := c.linearGo.Remove(c.linearGo.Front()).(*LinearGo)
		c.mutexLinearGo.Unlock()

		defer func() {
			c.g.ChanCb <- e.cb
			if r := recover(); r != nil {
				if conf.LenStackBuf > 0 {
					buf := make([]byte, conf.LenStackBuf)
					l := runtime.Stack(buf, false)
					log.Error("%v: %s", r, buf[:l])
				} else {
					log.Error("%v", r)
				}
			}
		}()

		e.f()
	}()
}
