package gglmm

import "testing"

type Test struct {
	Value string `json:"value"`
}

func TestRedisCache(t *testing.T) {

	redisCacher := NewRedisCacher("tcp", "127.0.0.1:6379", 5, 10, 3)
	defer redisCacher.Close()

	err := redisCacher.Set("key", "value")
	if err != nil {
		t.Fatalf(err.Error())
	}

	valueString, err := redisCacher.GetString("key")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if valueString != "value" {
		t.Fatalf("value not match：" + valueString)
	}

	err = redisCacher.Set("key", 1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	valueInt, err := redisCacher.GetInt("key")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if valueInt != 1 {
		t.Fatalf("value not match")
	}

	err = redisCacher.Set("key", 1.0)
	if err != nil {
		t.Fatalf(err.Error())
	}

	valueFloat, err := redisCacher.GetFloat64("key")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if valueFloat != 1.0 {
		t.Fatalf("value not match")
	}

	test := &Test{Value: "value"}

	err = redisCacher.Set("key", test)
	if err != nil {
		t.Fatalf(err.Error())
	}

	objValue := &Test{}
	err = redisCacher.GetObj("key", objValue)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if test.Value != objValue.Value {
		t.Fatalf(objValue.Value)
	}
}
