# rds_cache_go

## 说明
rds_cache_go是基于go-redis的缓存工具包

## 安装
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

// Get
result := caches.Get("key2")

// Delete
result,err := caches.Delete("key2")

// Delete More
result,err = caches.Delete("key2","key3","key4")

// already redis connect client
rdsClient := redis.NewClient(&redis.Options{
    Addr:     fmt.Sprintf("%s:%s", "192.168.0.151", "6379"),
    Password: "",
    DB:       15,
})

// WithRedisClient() 使用已有的redis链接客户端 
// WithOriginDB() 是否使用已经的客户端的数据库 true - 使用已有的客户端链接数据库 | false - 使用WithDB()中的数据库，如果为空默认使用0库
// WithDB() 指定数据库，如果 WithOriginDB() 为false，WithDB()未指定的话默认使用0号库

caches := NewCache(WithRedisClient(rdsClient), WithOriginDB(true), WithDB(11))


```