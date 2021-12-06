package rds_cache_go

import (
	"fmt"
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
