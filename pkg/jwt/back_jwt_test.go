package jwt_test

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"testing"
	"time"
	"uims/boot"
	"uims/internal/controllers/login_controller"
	"uims/internal/model"
	"uims/internal/service"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

//go test -v internal/controllers/login_controller/login_controller_test.go -test.run TestSetJwt
func TestSetJwt(t *testing.T) {
	user := model.User{}
	_ = service.GetUserService().GetUserInfoByAccount(&user, "uims_super_admin")
	mySigningKey := []byte("UIMS-BACK-JWT")

	type MyCustomClaims struct {
		Account string `json:"account"`
		jwtgo.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		user.Account,
		jwtgo.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "UIMS-BACK",
		},
	}

	fmt.Println(user.Account)
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	data := login_controller.LoginResult{
		User:  user,
		Token: ss,
	}
	fmt.Println(data)
}

//go test -v internal/controllers/login_controller/login_controller_test.go -test.run TestVerifyJwt
func TestVerifyJwt(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoidWltc19zdXBlcl9hZG1pbiIsImV4cCI6MTUwMDAsImlzcyI6IlVJTVMtQkFDSyJ9.6GvbGluJta7NOctjq_11Oq6U0swyKpRPE3ShdwLv_6Y"

	type MyCustomClaims struct {
		Account string `json:"account"`
		jwtgo.StandardClaims
	}

	// sample token is expired.  override time so it parses as valid
	at(time.Unix(0, 0), func() {
		token, err := jwtgo.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwtgo.Token) (interface{}, error) {
			return []byte("UIMS-BACK-JWT"), nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			fmt.Println("ass:", claims.Account)
			fmt.Printf("%v %v", claims.Account, claims.StandardClaims.ExpiresAt)
		} else {
			fmt.Println(err)
		}
	})

}

func at(t time.Time, f func()) {
	jwtgo.TimeFunc = func() time.Time {
		return t
	}
	f()
	jwtgo.TimeFunc = time.Now
}
