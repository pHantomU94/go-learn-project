package geecache

import pb "geecache/geecachepb"

// PeerPicker 远程节点选择器接口
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 某个节点的缓存获取API接口
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}