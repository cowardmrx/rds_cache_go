package rds_cache_go

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Option func(ca *config)

func WithHost(host string) Option {
	return func(ca *config) {
		ca.Host = host
	}
}

func WithPort(port string) Option {
	return func(ca *config) {
		ca.Port = port
	}
}

func WithPassword(password string) Option {
	return func(ca *config) {
		ca.Password = password
	}
}

func WithDB(db int) Option {
	return func(ca *config) {
		ca.DB = db
	}
}

type Cache interface {
	// Put 放入缓存
	Put(key string, value interface{}, ttl time.Duration) error
	// Exist 判断某个缓存是否存在
	Exist(key string) bool
	// Get 获取某个缓存的值
	Get(key string) interface{}
	// Delete 删除指定缓存
	Delete(keys ...string) int64
}

//	@method NewCache
//	@description: 初始化一个cache对象
//	@param opts ...Option
//	@return Cache
func NewCache(opts ...Option) Cache {

	cha := new(config)

	for _, v := range opts {
		v(cha)
	}

	rClient := new(cache)

	rClient.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cha.Host, cha.Port),
		Password: cha.Password,
		DB:       cha.DB,
	})

	rClient.ctx = context.Background()

	return rClient

}

type cache struct {
	client *redis.Client
	ctx    context.Context
}

//	@method Put
//	@description: 放入缓存
//	@receiver c
//	@param key string key
//	@param value interface{} 值
//	@param ttl time.Duration 过期时间
//	@return error
func (c *cache) Put(key string, value interface{}, ttl time.Duration) error {
	_, err := c.client.Set(c.ctx, key, value, TTL(ttl)).Result()

	if err != nil {
		return err
	}

	return nil
}

//	@method Exist
//	@description: 判断指定缓存是否存在
//	@receiver c
//	@param key string 缓存key
//	@return bool
func (c *cache) Exist(key string) bool {
	result, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false
	}

	if result <= 0 {
		return false
	}

	return true
}

//	@method Get
//	@description: 获取指定缓存
//	@receiver c
//	@param key string 缓存key
//	@return interface{}
func (c *cache) Get(key string) interface{} {
	result, err := c.client.Get(c.ctx, key).Result()

	if err != nil {
		panic("get cache by " + key + " failed: " + err.Error())
	}

	return result
}

//	@method Delete
//	@description: 删除指定缓存
//	@receiver c
//	@param keys ...string
//	@return int
func (c *cache) Delete(keys ...string) int64 {
	result, err := c.client.Del(c.ctx, keys...).Result()

	if err != nil {
		panic("delete cache failed: " + err.Error())
	}

	return result
}
