package gglmm

import (
	"net/http"
	"testing"
)

type UserID int64

func (id UserID) AuthInfo() *AuthInfo {
	return &AuthInfo{
		Type: "testType",
		ID:   int64(id),
	}
}

func TestAuthorization(t *testing.T) {
	userID := UserID(1)
	authJWT, _, err := GenerateAuthJWT(userID, JWTExpires, "testSecret")
	if err != nil {
		t.Fatal(err)
	}

	r1, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	r1.Header.Add("Authorization", "Bearer "+authJWT)

	authJWT = GetAuthJWT(r1)
	_, _, err = ParseAuthJWT(authJWT, "testSecret")
	if err != nil {
		t.Fatal(err)
	}

	r2 := RequestWithAuthInfo(r1, userID.AuthInfo())
	id, err := GetAuthID(r2, "testType")
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatal(id)
	}
}
