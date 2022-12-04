# simcache

* LRU缓存淘汰策略
* 单机并发缓存
* HTTP服务端
* 一致性哈希与分布式节点
* 防止缓存击穿 & 缓存穿透
* protobuf通信


待优化

* sync.mu性能不如sync.map，可将锁机制替换为sync.map
* simcache与httppool虽实现上解耦，但用接口实现关联的方式不太优雅



参考

* https://github.com/golang/groupcache
* bytedance life-algorithm async实现
