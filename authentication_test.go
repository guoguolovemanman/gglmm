package gglmm

import (
	"net/http"
	"testing"
)

type UserID int64

func (id *UserID) Authorization() *Authorization {
	return &Authorization{
		Type: "test",
		ID:   int64(*id),
	}
}

func TestRequestContext(t *testing.T) {
	userID := UserID(1)
	authorizationToken, jwtClaims, err := GenerateAuthorizationToken(&userID, JWTExpires, "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(authorizationToken, jwtClaims)

	r1, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	r1.Header.Add("Authorization", "Bearer "+authorizationToken)

	authorizationToken = GetAuthorizationToken(r1)
	authorization, jwtClaims, err := ParseAuthorizationToken(authorizationToken, "test")
	t.Log(authorization, jwtClaims)

	r2 := RequestWithAuthorization(r1, userID.Authorization())
	id, err := GetAuthorizationID(r2, "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if id != 1 {
		t.Log(id)
		t.Fail()
	}
}
