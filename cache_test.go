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

func TestCache_HPut(t *testing.T) {

	//value1 := map[string]interface{}{
	//	"data": "shabi",
	//}
	//
	//jsonValue, _ := json.Marshal(value1)

	if err := caches.HPut("h_key", "hvalue_1", 1, "key_2", 3.14); err != nil {
		t.Logf("hash put failed: %v", err)
		return
	}

	t.Log("hash put success")
}

func TestCache_HMPut(t *testing.T) {

	if err := caches.HMPut("h_m_put", "hm_key_1", "hm_value_1", "hm_key_2", "hm_value_2"); err != nil {
		t.Logf("hm put failed: %v", err.Error())
		return
	}

	t.Log("hm put success")
}

func TestCache_HKeyExist(t *testing.T) {
	if ok := caches.HKeyExist("h_m_put", "hm_key_1"); !ok {
		t.Log("key is not exist")
		return
	}

	t.Log("key is exist")
}

func TestCache_HGet(t *testing.T) {
	result := caches.HGet("h_key", "hvalue_1")

	t.Logf("result is :%v", result)
}

func TestCache_HGetAll(t *testing.T) {
	result := caches.HGetAll("h_keys")

	t.Logf("result is : %v", result)
}

func TestCache_HGetKeyAll(t *testing.T) {
	result := caches.HGetKeyAll("h_key")
	t.Logf("all key is : %v", result)
}

func TestCache_HIncrBy(t *testing.T) {
	result := caches.HIncrBy("h_key", "hvalue_1", 3)

	t.Logf("result is : %v", result)
}

func TestCache_HFloatIncrBy(t *testing.T) {
	result := caches.HFloatIncrBy("h_key", "key_2", 3.14)

	t.Logf("result is : %v", result)
}

func TestCache_HGetValAll(t *testing.T) {
	result := caches.HGetValAll("h_key")

	t.Logf("result is : %v", result)
}

func TestCache_HDelete(t *testing.T) {
	result := caches.HDelete("h_key", "key_2")

	t.Logf("result is : %v", result)
}
