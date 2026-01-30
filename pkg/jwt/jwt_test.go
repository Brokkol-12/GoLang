package jwt_test

import (
	"golang/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "a21@a.ru"
	jwtService := jwt.NewJWT("M9sVthtJq9WFpadl5UtSNAm7EKTPACrmpUXZyOuDkVq")
	token, err := jwtService.Create(jwt.JWTData{
		Email: "a21@a.ru",
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is Invalid")
	}
	if data.Email != email {
		t.Fatalf("Email %s not equal %s", data.Email, email)
	}
}
