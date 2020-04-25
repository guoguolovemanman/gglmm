package gglmm

import "testing"

func TestBcrypt(t *testing.T) {
	test := "test"

	hashed, err := GeneratePassword(test)
	if err != nil {
		t.Fatal(err)
	}

	err = ComparePassword(hashed, test)
	if err != nil {
		t.Fatal(err)
	}
}
