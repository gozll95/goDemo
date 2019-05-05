url:http://colobu.com/2017/07/11/dive-into-sync-Map/

github: http://colobu.com/2017/07/11/dive-into-sync-Map/


思想:

可以说，上面的解决方案相当简洁，并且利用读写锁而不是Mutex可以进一步减少读写的时候因为锁带来的性能。

但是，它在一些场景下也有问题，如果熟悉Java的同学，可以对比一下java的ConcurrentHashMap的实现，在map的数据非常大的情况下，一把锁会导致大并发的客户端共争一把锁，Java的解决方案是shard, 内部使用多个锁，每个区间共享一把锁，这样减少了数据共享一把锁带来的性能影响，orcaman提供了这个思路的一个实现： concurrent-map，他也询问了Go相关的开发人员是否在Go中也实现这种方案，由于实现的复杂性，答案是Yes, we considered it.,但是除非有特别的性能提升和应用场景，否则没有进一步的开发消息。