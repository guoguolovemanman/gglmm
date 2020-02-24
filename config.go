package gglmm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

// ConfigAPI --
type ConfigAPI struct {
	Address string
}

// Check --
func (config ConfigAPI) Check() bool {
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	log.Println("config api check valid")
	return true
}

// ConfigDB --
type ConfigDB struct {
	Dialect         string
	Address         string
	User            string
	Password        string
	Database        string
	MaxOpen         int
	MaxIdel         int
	ConnMaxLifetime int
}

// Check --
func (config ConfigDB) Check() bool {
	if config.Dialect == "" {
		return false
	}
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	if config.User == "" || config.Password == "" {
		return false
	}
	if config.Database == "" {
		return false
	}
	if config.MaxOpen <= 0 || config.MaxIdel <= 0 || config.ConnMaxLifetime <= 0 {
		return false
	}
	log.Println("config db check valid")
	return true
}

// ConfigRedis --
type ConfigRedis struct {
	Network     string
	Address     string
	MaxActive   int
	MaxIdel     int
	IdelTimeout int
}

// Check --
func (config ConfigRedis) Check() bool {
	if config.Network == "" {
		return false
	}
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	if config.MaxActive <= 0 || config.MaxIdel <= 0 || config.IdelTimeout <= 0 {
		return false
	}
	log.Println("config redis check valid")
	return true
}

// ConfigRPC --
type ConfigRPC struct {
	Network string
	Address string
}

// Check --
func (config ConfigRPC) Check() bool {
	if config.Network == "" {
		return false
	}
	if config.Address == "" {
		return false
	}
	return true
}

type configChecker interface {
	Check() bool
}

// ParseConfigFile --
func ParseConfigFile(file string, config configChecker) bool {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		log.Fatal(err)
	}
	return config.Check()
}
