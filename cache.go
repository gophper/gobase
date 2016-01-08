package gobase

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	Lock *sync.RWMutex
	Len  uint64
	Item map[interface{}]*Object
}
type Object struct {
	time    time.Time
	timeout time.Duration //if timeout = 0, not gc
	obj     interface{}
}

var (
	CacheTable       *Cache
	ErrTimeOut       error = errors.New("The cache has been timeout.")
	ErrKeyNotFound   error = errors.New("The key was not found.")
	ErrTypeAssertion error = errors.New("Type assertion error.")
)

//新建一个cache
func NewCache() (c *Cache) {
	c = &Cache{
		Lock: new(sync.RWMutex),
		Item: make(map[interface{}]*Object),
	}
	return
}

//cache 垃圾回收
func (c *Cache) gc() {
	for {
		for k, v := range c.Item {

			if v.time.Add(v.timeout).Before(time.Now()) && v.timeout != 0 {
				delete(c.Item, k)
			}

			time.Sleep(time.Microsecond)
		}
		time.Sleep(time.Second)
	}
}

//// 从cache里获取
func (c *Cache) Get(key interface{}) (interface{}, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	if v, ok := c.Item[key]; ok {
		now := time.Now()
		if v.time.Add(v.timeout).Before(now) || v.timeout == 0 {
			v.time = now
			return v.obj, nil
		}
		delete(c.Item, key)
		return nil, ErrTimeOut

	}
	return nil, ErrKeyNotFound

}

// 检查cache是否存在
func (c *Cache) Exists(key interface{}) bool {
	_, ok := c.Item[key]
	return ok
}

// 设置cache
func (c *Cache) Set(key, val interface{}, timeout time.Duration) bool {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if v, ok := c.Item[key]; ok && val == v {
		return false
	}

	o := &Object{
		time:    time.Now(),
		timeout: timeout,
		obj:     val,
	}
	c.Item[key] = o

	c.Len++

	return true
}

// 删除cache
func (c *Cache) Del(key interface{}) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	delete(c.Item, key)
	c.Len--
}

// 清除cache
// 参数是一个函数
//传递外部函数判断cache里数据，决定是否要删除
func (c *Cache) Cleanup(f func(interface{}) bool) {
	for key, value := range c.Item {
		if f(value) {
			c.Del(key)
		}
	}
}
