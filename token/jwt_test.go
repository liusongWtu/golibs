package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestJwtHS256Funcs(t *testing.T) {
	tmodel := &Jwt{}
	conf := make(map[string]string)
	conf["hsKey"] = "ABCD"
	tmodel.Init(conf)

	claims := make(jwt.MapClaims)
	claims["username"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 480).Unix()
	//claims["exp"] = time.Now().Unix()+2
	//HS256加解密
	tokenStr, err := tmodel.GetHStoken(claims)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenStr, err)
	claims1, flg := tmodel.ParseHStoken(tokenStr)
	if flg != 1 {
		if flg == 2 {
			t.Fatal("token过期")
		} else if flg == 3 {
			t.Fatal("token解析处理错误")
		} else if flg == 4 {
			t.Fatal("类型转换错误")
		}
	}
	fmt.Println(claims1, flg)
}

func TestJwtES256Funcs(t *testing.T) {
	tmodel := &Jwt{}
	esKeyD, esKeyX, esKeyY := tmodel.GenerateKey()
	conf := make(map[string]string)
	conf["esKeyD"] = esKeyD
	conf["esKeyX"] = esKeyX
	conf["esKeyY"] = esKeyY
	tmodel.Init(conf)

	claims := make(jwt.MapClaims)
	claims["username"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 480).Unix()
	//claims["exp"] = time.Now().Unix()+2
	//ES256加解密
	tokenStr, err := tmodel.GetEStoken(claims)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenStr, err)
	claims1, flg := tmodel.ParseEStoken(tokenStr)
	if flg != 1 {
		if flg == 2 {
			t.Fatal("token过期")
		} else if flg == 3 {
			t.Fatal("token解析处理错误")
		} else if flg == 4 {
			t.Fatal("类型转换错误")
		}
	}
	fmt.Println(claims1, flg)
}
