package gglmm

// GGLMMCacheKeyPrefix --
const GGLMMCacheKeyPrefix = "gglmm-cache-"

// Cacher --
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

// DefaultCacher --
func DefaultCacher() Cacher {
	return cacher
}
