package geecache

import (
	"fmt"
	"geecache/singleflight"
	"log"
	"sync"
)

// Getter 数据获取接口（如访问数据库）
type Getter interface {
	Get(Key string) ([]byte, error)
}

// GetterFunc 实现Getter接口的函数
type GetterFunc func(key string) ([]byte, error)

// Get 为GetterFunc实现Getter接口，在接口函数中调用自己，免去使用struct实现接口时需要实例化的过程
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group is a cache namesapce and associated data loaded spread over
type Group struct {
	name string	// group的名称
	getter Getter // 本地数据获取方法
	mainCache cache // 本地缓存
	peers PeerPicker // 远程节点选择器
	loader *singleflight.Group // 并发请求组
}

var (
	mu sync.RWMutex					 // 读写锁
	groups = make(map[string]*Group) // 缓存组，使用Group的名字区分
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or nil if there's no such group
func GetGroup(name string) *Group {
	mu.RLock()
	g :=groups[name]
	mu.RUnlock()
	return g
}

// Get 通过本地缓存获取key对应的值，如果本地缓存没有调用load方法
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

// getLocally 本地数据获取方法，调用本地Gettter来获取数据，并添加到本地缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache 添加本地缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeers 指定Group的PeerPicker
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// load 加载数据，首先判断是否在远程缓存节点，如果有返回值，没有则选用本地加载策略（如： 访问数据库）
func (g *Group) load(key string) (value ByteView, err error) {

	view, err := g.loader.Do(key, func()(interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, err
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err != nil{
		log.Println("[GeeCache] Failed to get from peer", err)
		return
	}
	return view.(ByteView), err
}

// getFromPeer 根绝key选择节点位置，获取指定PeerGetter，通过API访问数据
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, err
}