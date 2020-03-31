package gglmm

import (
	"reflect"
	"testing"
)

type Model1 struct {
}

type Model2 struct {
}

func (test Model2) Cache() bool {
	return true
}

func (test Model2) ResponseKey() [2]string {
	return [...]string{"test", "tests"}
}

var model1 = Model1{}
var model2 = Model2{}

func TestReflectCache(t *testing.T) {
	tests := []struct {
		model  interface{}
		result bool
	}{
		{model: model1, result: false},
		{model: model2, result: true},
	}
	for _, test := range tests {
		result := ReflectCache(reflect.ValueOf(test.model))
		if result != test.result {
			t.Fatalf("fail: %v <=> %v\n", result, test.result)
		}
	}
}

func TestReflectSingleKey(t *testing.T) {
	tests := []struct {
		model  interface{}
		result string
	}{
		{model: model1, result: "record"},
		{model: model2, result: "test"},
	}
	for _, test := range tests {
		result := ReflectSingleKey(reflect.ValueOf(test.model))
		if result != test.result {
			t.Fatalf("fail: %v <=> %v\n", result, test.result)
		}
	}
}

func TestReflectMultiKey(t *testing.T) {
	tests := []struct {
		model  interface{}
		result string
	}{
		{model: model1, result: "records"},
		{model: model2, result: "tests"},
	}
	for _, test := range tests {
		result := ReflectMultiKey(reflect.ValueOf(test.model))
		if result != test.result {
			t.Fatalf("fail: %v <=> %v\n", result, test.result)
		}
	}
}
