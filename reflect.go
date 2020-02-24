package gglmm

import (
	"reflect"
)

// ReflectNew 新建模型
func ReflectNew(t reflect.Type) interface{} {
	value := reflect.New(t)
	return value.Interface()
}

// ReflectNewSliceOfPtrTo 新建模型列表
func ReflectNewSliceOfPtrTo(t reflect.Type) interface{} {
	value := reflect.New(reflect.SliceOf(reflect.PtrTo(t)))
	return value.Interface()
}

func reflectCache(v reflect.Value) bool {
	if cacheMethod := v.MethodByName("Cache"); cacheMethod.IsValid() {
		results := cacheMethod.Call(nil)
		if results != nil && len(results) == 1 {
			return results[0].Bool()
		}
	}
	return false
}

func reflectSingleKey(v reflect.Value) string {
	if keyMethod := v.MethodByName("ResponseKey"); keyMethod.IsValid() {
		results := keyMethod.Call(nil)
		if results != nil && len(results) == 1 {
			return results[0].Index(0).String()
		}
	}
	return "record"
}

func reflectMultiKey(v reflect.Value) string {
	if keyMethod := v.MethodByName("ResponseKey"); keyMethod.IsValid() {
		results := keyMethod.Call(nil)
		if results != nil && len(results) == 1 {
			return results[0].Index(1).String()
		}
	}
	return "records"
}
