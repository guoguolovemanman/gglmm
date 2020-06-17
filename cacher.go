package gglmm

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

const cacherKeyPrefix = "gglmm:cacher:"

// RegisterCacher --
func RegisterCacher(cacherInstance Cacher) {
	cacherInstance.SetKeyPrefix(cacherKeyPrefix)
	cacher = cacherInstance
}

// DefaultCacher 默认缓存
func DefaultCacher() Cacher {
	return cacher
}
