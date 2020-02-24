package gglmm

import (
	"reflect"
	"testing"
)

type Model1 struct {
}

type Model2 struct {
}

func (test Model2) Preloads() []string {
	return []string{"test"}
}

func (test Model2) Cache() bool {
	return true
}

func (test Model2) ResponseKey() [2]string {
	return [...]string{"test", "tests"}
}

func TestModel1(t *testing.T) {
	var model1 Model1
	rType := reflect.TypeOf(model1)
	rValue := reflect.ValueOf(model1)

	model := ReflectNew(rType)
	if _, ok := model.(*Model1); !ok {
		t.Log("ReflectNew")
		t.Fail()
	}

	modleSlice := ReflectNewSliceOfPtrTo(rType)
	if _, ok := modleSlice.(*[]*Model1); !ok {
		t.Log("ReflectNewSliceOfPtrTo")
		t.Fail()
	}

	cache := reflectCache(rValue)
	if cache {
		t.Log("reflectCache")
		t.Fail()
	}

	key := reflectSingleKey(rValue)
	if key != "model" {
		t.Log("reflectSingleKey")
		t.Fail()
	}

	keys := reflectMultiKey(rValue)
	if keys != "models" {
		t.Log("reflectMultiKey")
		t.Fail()
	}
}

func TestModel2(t *testing.T) {
	var model2 Model2
	rValue := reflect.ValueOf(model2)

	cache := reflectCache(rValue)
	if !cache {
		t.Log("reflectCache")
		t.Fail()
	}

	key := reflectSingleKey(rValue)
	if key != "test" {
		t.Log("reflectSingleKey")
		t.Fail()
	}

	keys := reflectMultiKey(rValue)
	if keys != "tests" {
		t.Log("reflectMultiKey")
		t.Fail()
	}
}
