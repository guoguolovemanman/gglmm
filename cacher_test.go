package gglmm

import (
	"log"
	"reflect"
	"testing"

	redis "github.com/weihongguo/gglmm-redis"
)

type CacheModel struct {
	Model
}

func (model CacheModel) Cache() bool {
	return true
}

func TestCahcer(t *testing.T) {
	redisCacher := redis.NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 30)
	defer redisCacher.Close()
	RegisterCacher(redisCacher)

	model := CacheModel{Model: Model{ID: 1}}
	if err := CacherSetByIDRequest(reflect.ValueOf(model), reflect.TypeOf(model), model, IDRequest{ID: 1}); err != nil {
		log.Fatal(err)
	}

	if err := CacherGetByIDRequest(reflect.ValueOf(model), reflect.TypeOf(model), &model, IDRequest{ID: 1}); err != nil {
		log.Fatal(err)
	}

	if err := CacherDelPattern(reflect.ValueOf(model), reflect.TypeOf(model), 1); err != nil {
		log.Fatal(err)
	}

	if err := CacherGetByIDRequest(reflect.ValueOf(model), reflect.TypeOf(model), &model, IDRequest{ID: 1}); err == nil {
		log.Fatal("不应该有数据啊！")
	}
}
