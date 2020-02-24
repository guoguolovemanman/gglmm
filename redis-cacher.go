package gglmm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisCahcerName --
const RedisCahcerName = "redis"

var defaultRedisCacheExpires = 24 * 60 * 60

var (
	// ErrConn --
	ErrConn = errors.New("连接错误")
)

// RedisCacher 缓存
type RedisCacher struct {
	redisPool *redis.Pool
	keyPrefix string
}

// RegisterRedisCacher 注册全局RedisCache
func RegisterRedisCacher(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration) {
	cacher = NewRedisCacher(network, address, maxActive, maxIdle, idleTimeout)
}

// RegisterRedisCacherConfig --
func RegisterRedisCacherConfig(config ConfigRedis) {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("RedisConfig check invalid")
	}
	idelTimeout, err := time.ParseDuration(fmt.Sprintf("%ds", config.IdelTimeout))
	if err != nil {
		log.Fatal(err)
	}
	cacher = NewRedisCacher(
		config.Network,
		config.Address,
		config.MaxActive,
		config.MaxIdel,
		idelTimeout,
	)
}

// NewRedisCacher --
func NewRedisCacher(network string, address string, maxActive int, maxIdle int, idleTimeout time.Duration) *RedisCacher {
	return &RedisCacher{
		redisPool: &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: idleTimeout * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				return redis.Dial(network, address)
			},
		},
		keyPrefix: GGLMMCacheKeyPrefix,
	}
}

// SetDefaultRedisCacheExpires 注册全局DefaultRedisCacheExpires
func SetDefaultRedisCacheExpires(expires int) {
	defaultRedisCacheExpires = expires
}

// CloseRedisCacher --
func CloseRedisCacher() {
	if cacher != nil && cacher.Name() == RedisCahcerName {
		cacher.Close()
	}
}

// SetKeyPrefix --
func (cacher *RedisCacher) SetKeyPrefix(keyPrefix string) {
	cacher.keyPrefix = keyPrefix
}

// Name --
func (cacher *RedisCacher) Name() string {
	return RedisCahcerName
}

// SetEx --
func (cacher *RedisCacher) SetEx(key string, value interface{}, ex int) error {
	redisConn := cacher.redisPool.Get()
	if redisConn == nil {
		return ErrConn
	}
	defer redisConn.Close()

	var err error
	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map, reflect.Ptr:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		_, err = redisConn.Do("set", cacher.keyPrefix+key, jsonValue, "ex", ex)
	default:
		_, err = redisConn.Do("set", cacher.keyPrefix+key, value, "ex", ex)
	}
	return err
}

// Set --
func (cacher *RedisCacher) Set(key string, value interface{}) error {
	return cacher.SetEx(key, value, defaultRedisCacheExpires)
}

// Get --
func (cacher *RedisCacher) Get(key string) (interface{}, error) {
	redisConn := cacher.redisPool.Get()
	if redisConn == nil {
		return nil, ErrConn
	}
	defer redisConn.Close()

	return redisConn.Do("get", cacher.keyPrefix+key)
}

// Del --
func (cacher *RedisCacher) Del(key string) error {
	redisConn := cacher.redisPool.Get()
	if redisConn == nil {
		return ErrConn
	}
	defer redisConn.Close()

	_, err := redisConn.Do("del", cacher.keyPrefix+key)
	return err
}

// Close --
func (cacher *RedisCacher) Close() {
	cacher.redisPool.Close()
}

// GetInt --
func (cacher *RedisCacher) GetInt(key string) (int, error) {
	return redis.Int(cacher.Get(key))
}

// GetInt64 --
func (cacher *RedisCacher) GetInt64(key string) (int64, error) {
	return redis.Int64(cacher.Get(key))
}

// GetFloat64 --
func (cacher *RedisCacher) GetFloat64(key string) (float64, error) {
	return redis.Float64(cacher.Get(key))
}

// GetBytes --
func (cacher *RedisCacher) GetBytes(key string) ([]byte, error) {
	return redis.Bytes(cacher.Get(key))
}

// GetString --
func (cacher *RedisCacher) GetString(key string) (string, error) {
	return redis.String(cacher.Get(key))
}

// GetObj --
func (cacher *RedisCacher) GetObj(key string, obj interface{}) error {
	value, err := cacher.GetBytes(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, obj)
}
