collector--->store

collector 实例里
并且也重写了Collector Interface的规则


collector 实例有:
- `ChunkCollector`{内包含 `Icollector interface`}: chunk批量调用`内包含`的`collector interface`的方法
- `RemoteCollecor`: 远程写入net中
- `CollectorServer`{内包含 `Icollector interface`}:用于从client端接收序列化的,再调用`内包含`的`collector interface`

> RemoteCollecor和CollectorServer就是一对


满足 Collector Interface的实例:
- collector 实例
- 满足deleteStore的实例


``` go
deleteStore interface{
    store interface{
        collector interface
        根据TraceID检索出Tracer     
    }
    delete(..traceID)
}
```

满足 deleteStore 的实例:
- memoryStore
- `RecentStore`{内包含 `IdeleteStore interface`}: 过期清理
- `RingStore`{内包含 `IdeleteStore interface`}: 环形

并且也重写了collector interface的规则


Tracer是一棵树

## 精华:
struct里包含interface并且实现interface,有足够的扩展功能


