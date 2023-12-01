package cache

import "sync"

type geecache struct {
	c     map[string][]byte//保存kv值
	mutex sync.RWMutex//读写锁
}

func (c *geecache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[k] = v

	return nil
}

func (c *geecache) Get(k string) ([]byte, error) {
	c.mutex.RLock()//只读锁
	defer c.mutex.RUnlock()//
	return c.c[k], nil
}

func (c *geecache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.c, k)

	return nil
}
//申请内存缓存
func newgeecache() *geecache {
	return &geecache{make(map[string][]byte), sync.RWMutex{}}
}
//创建cache并返回cache接口
func New() Cache {
	c := newgeecache()
	return c
}
