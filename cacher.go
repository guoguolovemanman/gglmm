package gglmm

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// Cacher 缓存接口
type Cacher interface {
	SetExpires(expires int)
	Expires() int

	SetKeyPrefix(keyPrefix string)
	KeyPrefix() string

	SetEx(key string, value interface{}, ex int) error
	Set(key string, value interface{}) error

	Get(key string) (interface{}, error)

	GetInt(key string) (int, error)
	GetInt64(key string) (int64, error)
	GetFloat64(key string) (float64, error)
	GetBytes(key string) ([]byte, error)
	GetString(key string) (string, error)
	GetObj(key string, obj interface{}) error

	Del(key string) (int, error)
	DelPattern(pattern string) (int, error)

	Close()
}

var cacher Cacher = nil

const cacherKeyPrefix = "gglmm:cache:"

// RegisterCacher 注册缓存
func RegisterCacher(cacherInstance Cacher) {
	cacherInstance.SetKeyPrefix(cacherKeyPrefix)
	cacher = cacherInstance
}

// DefaultCacher 默认缓存
func DefaultCacher() Cacher {
	return cacher
}

// CacherGetByIDRequest --
func CacherGetByIDRequest(v reflect.Value, t reflect.Type, model interface{}, idRequest IDRequest) error {
	if cacher != nil {
		if SupportCache(v) {
			cacheKey := t.Name() + ":" + strconv.FormatInt(idRequest.ID, 10)
			if len(idRequest.Preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(idRequest.Preloads, "-")
			}
			if err := cacher.GetObj(cacheKey, model); err != nil {
				return err
			}
			return nil
		}
		return errors.New("模型不支持缓存")
	}
	return errors.New("请注册Cacher")
}

// CacherSetByIDRequest --
func CacherSetByIDRequest(v reflect.Value, t reflect.Type, model interface{}, idRequest IDRequest) error {
	if cacher != nil {
		if SupportCache(v) {
			cacheKey := t.Name() + ":" + strconv.FormatInt(idRequest.ID, 10)
			if len(idRequest.Preloads) > 0 {
				cacheKey = cacheKey + ":" + strings.Join(idRequest.Preloads, "-")
			}
			cacher.Set(cacheKey, model)
			return nil
		}
		return errors.New("模型不支持缓存")
	}
	return errors.New("请注册Cacher")
}

// CacherDelPattern --
func CacherDelPattern(v reflect.Value, t reflect.Type, id int64) error {
	if cacher != nil {
		if SupportCache(v) {
			cacheKey := t.Name() + ":" + strconv.FormatInt(id, 10)
			cacher.DelPattern(cacheKey)
			return nil
		}
		return errors.New("模型不支持缓存")
	}
	return errors.New("请注册Cacher")
}
