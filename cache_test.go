package rds_cache_go

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var caches = NewCache(WithHost("192.168.0.151"), WithPort("6379"), WithDB(13))

func TestCache(t *testing.T) {

	err := caches.Put("key1", "va", 1*time.Minute)
	if err != nil {
		fmt.Printf("存放失败")
	}

	fmt.Println("放入存放成功")
}

func TestCache_Exist(t *testing.T) {
	exi := caches.Exist("key1")

	fmt.Println(exi)
}

func TestCache_Get(t *testing.T) {
	result := caches.Get("key1")
	fmt.Println(result)
}

func TestCache_Delete(t *testing.T) {
	caches.Put("key1", "value1", 0)
	caches.Put("key2", "value2", 0)

	resunt := caches.Delete("key1", "key2")
	fmt.Println(resunt)
}

func TestCache_RdsClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", "192.168.0.151", "6379"),
		Password: "",
		DB:       15,
	})

	cachess := NewCache(WithRedisClient(client), WithOriginDB(true), WithDB(11))

	cachess.Put("this key", "this aaaaa", 1*time.Minute)
}
