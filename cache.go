package rds_cache_go

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type config struct {
	Host     string        // redis主机
	Port     string        // redis 端口
	Password string        // redis 密码
	DB       int           // redis库名
	Client   *redis.Client // redis链接客户端 【如果项目中已经有了redis链接可使用该参数】
	OriginDB bool          // 是否强制使用redis客户端【Client】的DB
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

func WithRedisClient(client *redis.Client) Option {
	return func(ca *config) {
		ca.Client = client
	}
}

func WithDB(db int) Option {
	return func(ca *config) {
		ca.DB = db
	}
}

func WithOriginDB(originDB bool) Option {
	return func(ca *config) {
		ca.OriginDB = originDB
	}
}

type Cache interface {
	// Put 放入缓存
	Put(key string, value interface{}, ttl time.Duration) error
	// Exist 判断某个缓存是否存在
	Exist(key string) bool
	// Get 获取某个缓存的值
	Get(key string) (interface{}, error)
	// Delete 删除指定缓存
	Delete(keys ...string) (int64, error)
	// HPut hash put
	HPut(key string, value ...interface{}) error
	// HMPut  hash put 兼容redis v3
	HMPut(key string, value ...interface{}) error
	// HKeyExist 判断hash表中的key是否存在
	HKeyExist(key, field string) bool
	// HGet 获取hash表中指定field的值
	HGet(key, field string) (interface{}, error)
	// HGetAll 获取hash表中的全部值
	HGetAll(key string) (map[string]string, error)
	// HGetKeyAll 获取hash表中的全部key
	HGetKeyAll(key string) ([]string, error)
	// HIncrBy 为哈希表 key 中的指定字段的整数值加上增量 increment 。
	HIncrBy(key, field string, incr int64) (int64, error)
	// HFloatIncrBy 为哈希表 key 中的指定字段的浮点数值加上增量 increment 。
	HFloatIncrBy(key, field string, incr float64) (float64, error)
	// HGetValAll 获取hash表中全部的value值
	HGetValAll(key string) ([]string, error)
	// HDelete 删除hash表中一个或多个字段
	HDelete(key string, fields ...string) (int64, error)
}

// @method NewCache
// @description: 初始化一个cache对象
// @param opts ...Option
// @return Cache
func NewCache(opts ...Option) Cache {

	cha := new(config)

	for _, v := range opts {
		v(cha)
	}

	rClient := new(cache)

	if cha.Client == nil {

		rClient.client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cha.Host, cha.Port),
			Password: cha.Password,
			DB:       cha.DB,
		})
	} else {
		rClient.client = cha.Client
		if cha.OriginDB == false {
			rClient.client.Options().DB = cha.DB
		}
	}

	rClient.ctx = context.Background()

	return rClient

}

type cache struct {
	client *redis.Client
	ctx    context.Context
}

// @method Put
// @description: 放入缓存
// @receiver c
// @param key string key
// @param value interface{} 值
// @param ttl time.Duration 过期时间
// @return error
func (c *cache) Put(key string, value interface{}, ttl time.Duration) error {
	_, err := c.client.Set(c.ctx, key, value, TTL(ttl)).Result()

	if err != nil {
		return err
	}

	return nil
}

// @method Exist
// @description: 判断指定缓存是否存在
// @receiver c
// @param key string 缓存key
// @return bool
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

// @method Get
// @description: 获取指定缓存
// @receiver c
// @param key string 缓存key
// @return interface{}
func (c *cache) Get(key string) (interface{}, error) {
	result, err := c.client.Get(c.ctx, key).Result()

	if err != nil {
		return nil, errors.New("NOT FOUND " + key)
	}

	return result, nil
}

// @method Delete
// @description: 删除指定缓存
// @receiver c
// @param keys ...string
// @return int
func (c *cache) Delete(keys ...string) (int64, error) {
	result, err := c.client.Del(c.ctx, keys...).Result()

	if err != nil {
		return 0, errors.New("DELETE CACHE FAILED " + err.Error())
	}

	return result, nil
}

// @method HPut
// @description: hash put
// @receiver c
// @param key string
// @param value ...interface{}
// @return error
func (c *cache) HPut(key string, value ...interface{}) error {
	_, err := c.client.HSet(c.ctx, key, value...).Result()

	return err
}

// @method HMPut
// @description: hash put 用来兼容redis v3
// @receiver c
// @param key string
// @param value interface{}
// @return error
func (c *cache) HMPut(key string, value ...interface{}) error {
	_, err := c.client.HMSet(c.ctx, key, value).Result()

	return err
}

// @method HKeyExist
// @description: 判断hash表中的key是否存在
// @receiver c
// @param key string
// @return bool
func (c *cache) HKeyExist(key, field string) bool {
	result, err := c.client.HExists(c.ctx, key, field).Result()

	if err != nil {
		return false
	}

	return result

}

// @method HGet
// @description: 获取hash表中指定key，field的值
// @receiver c
// @param key string
// @param field string
// @return interface{}
func (c *cache) HGet(key, field string) (interface{}, error) {
	result, err := c.client.HGet(c.ctx, key, field).Result()

	if err != nil {
		return nil, errors.New("GET HASH CACHE FAILED " + err.Error())
	}

	return result, nil
}

// @method HGetAll
// @description: 获取hash表中的全部数据
// @receiver c
// @param key string
// @return interface{}
func (c *cache) HGetAll(key string) (map[string]string, error) {
	result, err := c.client.HGetAll(c.ctx, key).Result()

	if err != nil {
		return nil, errors.New("GET HASH CACHE ALL FAILED " + err.Error())
	}

	return result, nil
}

// @method HGetKeyAll
// @description: 获取hash表中的全部key
// @receiver c
// @param key string
func (c *cache) HGetKeyAll(key string) ([]string, error) {
	result, err := c.client.HKeys(c.ctx, key).Result()

	if err != nil {
		return nil, errors.New("get hash all key failed :" + err.Error() + " by cache key :" + key)
	}

	return result, nil
}

// @method HIncrBy
// @description: 为哈希表 key 中的指定字段的整数值加上增量 increment
// @receiver c
// @param key string
// @param field string
// @param incr int64
func (c *cache) HIncrBy(key, field string, incr int64) (int64, error) {
	result, err := c.client.HIncrBy(c.ctx, key, field, incr).Result()

	if err != nil {
		return 0, errors.New("hash incr by failed: " + err.Error())
	}

	return result, nil
}

// @method HFloatIncrBy
// @description: 为哈希表 key 中的指定字段的浮点数值加上增量 increment
// @receiver c
// @param key string
// @param field string
// @param incr float64
// @return float64
func (c *cache) HFloatIncrBy(key, field string, incr float64) (float64, error) {
	result, err := c.client.HIncrByFloat(c.ctx, key, field, incr).Result()

	if err != nil {
		return 0, errors.New("hash incr float by failed: " + err.Error())
	}

	return result, nil
}

// @method HGetValAll
// @description: 获取hash表中全部的value 值
// @receiver c
// @param key string
func (c *cache) HGetValAll(key string) ([]string, error) {
	result, err := c.client.HVals(c.ctx, key).Result()

	if err != nil {
		return nil, errors.New("get hash all value failed: " + err.Error())
	}

	return result, nil
}

// @method HDelete
// @description: 删除hash表中一个或多个字段
// @receiver c
// @param key string
// @param fields ...string
func (c *cache) HDelete(key string, fields ...string) (int64, error) {
	result, err := c.client.HDel(c.ctx, key, fields...).Result()

	if err != nil {
		return 0, errors.New("delete hash filed failed: " + err.Error())
	}

	return result, nil
}
