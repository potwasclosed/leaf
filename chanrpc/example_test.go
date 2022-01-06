package chanrpc_test

import (
	"fmt"
	"github.com/name5566/leaf/chanrpc"
	"sync"
)

func Example() {
	s := chanrpc.NewServer(10)

	var wg sync.WaitGroup
	wg.Add(1)

	// goroutine 1
	go func() {		//单独开了  一个   协程,  模拟一个 服务
		s.Register("f0", func(args []interface{}) {	// 用这个注册 更多的是 为了 检测 传递的参数吧

		})

		s.Register("f1", func(args []interface{}) interface{} {
			return 1
		})

		s.Register("fn", func(args []interface{}) []interface{} {
			return []interface{}{1, 2, 3}
		})

		s.Register("add", func(args []interface{}) interface{} {
			n1 := args[0].(int)
			n2 := args[1].(int)
			return n1 + n2
		})

		wg.Done()   //这里为啥 应为必须 上面注册好了,才能启动下面的 携程.防止 g2 比 g1 更先执行. 经典啊  经典

		for {
			s.Exec(<-s.ChanCall) 	//不停的从chan里面拿消息,  死循环.
		}
	}()

	wg.Wait()
	wg.Add(1)

	// goroutine 2
	go func() {
		c := s.Open(10)//打开一个client

		// sync
		err := c.Call0("f0")
		if err != nil {
			fmt.Println(err)
		}

		r1, err := c.Call1("f1")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r1)
		}

		rn, err := c.CallN("fn")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(rn[0], rn[1], rn[2])
		}

		ra, err := c.Call1("add", 1, 2)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(ra)
		}

		// asyn
		c.AsynCall("f0", func(err error) {
			if err != nil {
				fmt.Println(err)
			}
		})

		c.AsynCall("f1", func(ret interface{}, err error) {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(ret)
			}
		})

		c.AsynCall("fn", func(ret []interface{}, err error) {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(ret[0], ret[1], ret[2])
			}
		})

		c.AsynCall("add", 1, 2, func(ret interface{}, err error) {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(ret)
			}
		})

		c.Cb(<-c.ChanAsynRet) //因为 异步  所以这么写,等待返回结束.  作者还是很巧妙的啊.
		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)
		c.Cb(<-c.ChanAsynRet)

		// go
		s.Go("f0")

		wg.Done()
	}()

	wg.Wait()

	// Output:
	// 1
	// 1 2 3
	// 3
	// 1
	// 1 2 3
	// 3
}
