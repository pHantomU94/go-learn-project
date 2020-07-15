package lru

import "container/list"

type Cache struct {
	maxBytes int64	//缓存总大小
	nbytes int64	//缓存已使用大小
	ll *list.List	//双向链表，用来维护缓存
	cache map[string]*list.Element	//使用hash索引缓存
	OnEvicted func(key string, value Value)	//记录被移除时的回调函数
}

type entry struct {
	key string	//缓存的键
	value Value	//缓存的值
}

// Value 缓存值接口，需要实现长度方法
type Value interface{
	Len() int
}

// New 构造新的缓存，初始化链表及cache的hash表
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 根据key查询value，将刚被选中的缓存块移动到链表头部，
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return	
}

// RemoveOldest 移除最久未使用的缓存块
func (c *Cache) RemoveOldest() {
	element := c.ll.Back()
	if element != nil {
		// 从链表中移除缓存块
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		// 同时 从map中移除对应的key与缓存块
		delete(c.cache, kv.key)
		// 更新缓存已用大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 删除后的回调
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 添加 / 更新缓存块
func (c *Cache) Add(key string, value Value) {
	// 更新缓存，移动至队首，更新缓存块中的value，同时更新缓存大小
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
	// 添加缓存块，添加缓存块至队首，并在map中添加对应的条目key->ll.element，并更新长度
		element := c.ll.PushFront(&entry{key, value})
		c.cache[key] = element
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 如果更新后的缓存大小超过了最大值，就移除最久未使用的缓存块
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 换粗条目数
func (c *Cache) Len() int {
	return c.ll.Len()
}
