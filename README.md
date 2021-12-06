#rds_cache_go

##说明
rds_cache_go是基于go-redis的缓存工具包

##安装
go get -u github.com/cowardmrx/rds_cache_go

## 使用
```go
var caches = NewCache(WithHost("192.168.0.151"), WithPort("6379"), WithDB(13))

// PUT 
err := caches.Put("key1", "va", 1*time.Minute)

// PUT json
amap := map[string]interface{}{
	"data"  :   "data1"
}

amapJson,_ := json.Marshal(amap)

err := caches.Put("key2",amapJson,1 * time.Minute)

// Exist 
exist := caches.Exist("key2")

// Delete
result,err := caches.Delete("key2")

// Delete More
result,err = caches.Delete("key2","key3","key4")
```