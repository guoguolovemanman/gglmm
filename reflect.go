package gglmm

import (
	"reflect"
)

// ReflectCache --
func ReflectCache(v reflect.Value) bool {
	if cacheMethod := v.MethodByName("Cache"); cacheMethod.IsValid() {
		results := cacheMethod.Call(nil)
		if results != nil && len(results) == 1 {
			return results[0].Bool()
		}
	}
	return false
}

// ReflectSingleKey --
func ReflectSingleKey(v reflect.Value) string {
	if keyMethod := v.MethodByName("ResponseKey"); keyMethod.IsValid() {
		results := keyMethod.Call(nil)
		if results != nil && len(results) == 1 {
			return results[0].Index(0).String()
		}
	}
	return "record"
}

// ReflectMultiKey --
func ReflectMultiKey(v reflect.Value) string {
	if keyMethod := v.MethodByName("ResponseKey"); keyMethod.IsValid() {
		results := keyMethod.Call(nil)
		if results != nil && len(results) == 1 {
			return results[0].Index(1).String()
		}
	}
	return "records"
}
