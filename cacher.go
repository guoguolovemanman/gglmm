package gglmm

// GGLMMCacheKeyPrefix 缓存前缀
const GGLMMCacheKeyPrefix = "gglmm-cache-"

// Cacher 缓存接口
type Cacher interface {
	Name() string

	SetEx(key string, value interface{}, ex int) error
	Set(key string, value interface{}) error

	Get(key string) (interface{}, error)

	GetInt(key string) (int, error)
	GetInt64(key string) (int64, error)
	GetFloat64(key string) (float64, error)
	GetBytes(key string) ([]byte, error)
	GetString(key string) (string, error)
	GetObj(key string, obj interface{}) error

	Del(key string) error

	Close()
}

var cacher Cacher = nil

// DefaultCacher 默认缓存
func DefaultCacher() Cacher {
	return cacher
}
