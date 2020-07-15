package singleflight

import "sync"


// call 用来存储缓存访问请求
type call struct {
	wg sync.WaitGroup	// 锁避免重复请求
	val interface{}	
	err error
}

type Group struct {
	mu sync.Mutex	// 锁避免同时修改请求映射
	kcmap map[string]*call	// 请求映射表 key -> call
}

// 防止缓存击穿，处理瞬时并发的请求
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	// group加锁保护kcmap不被同时写入
	g.mu.Lock()
	// 延迟初始化请求映射表，减少内存占用
	if g.kcmap == nil {
		g.kcmap = make(map[string]*call)
	}

	// 如果请求已经存在，解锁对kcmap的保护，等待请求的同步锁解锁，直接返回请求数据
	if c, ok := g.kcmap[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// 如果请求不存在，新建请求，加入同步锁，并将请求写入请求映射表，而后解锁对kcmap的保护
	c := new(call)
	c.wg.Add(1)
	g.kcmap[key] = c
	g.mu.Unlock() 
	
	// 等待请求响应结束
	c.val, c.err = fn()
	// 释放同步锁
	c.wg.Done()

	// 加锁操作kcmap，将请求记录删除
	g.mu.Lock()
	delete(g.kcmap, key)
	g.mu.Unlock()

	return c.val, c.err
}