package gglmm

import (
	"net/http"
	"testing"
)

type UserID int64

func (id UserID) Authorization() *Authorization {
	return &Authorization{
		Type: "testType",
		ID:   int64(id),
	}
}

func TestAuthorization(t *testing.T) {
	userID := UserID(1)
	authorizationToken, _, err := GenerateAuthorizationToken(userID, JWTExpires, "testSecret")
	if err != nil {
		t.Fatal(err)
	}

	r1, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	r1.Header.Add("Authorization", "Bearer "+authorizationToken)

	authorizationToken = GetAuthorizationToken(r1)
	_, _, err = ParseAuthorizationToken(authorizationToken, "testSecret")
	if err != nil {
		t.Fatal(err)
	}

	r2 := RequestWithAuthorization(r1, userID.Authorization())
	id, err := GetAuthorizationID(r2, "testType")
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatal(id)
	}
}
