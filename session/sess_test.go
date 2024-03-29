package session

import (
	"crypto/aes"
	"encoding/json"
	"testing"
)

func Test_gob(t *testing.T) {
	a := make(map[interface{}]interface{})
	a["username"] = "user001"
	a[12] = 234
	a["user"] = User{"asta", "xie"}
	b, err := EncodeGob(a)
	if err != nil {
		t.Error(err)
	}
	c, err := DecodeGob(b)
	if err != nil {
		t.Error(err)
	}
	if len(c) == 0 {
		t.Error("decodeGob empty")
	}
	if c["username"] != "user001" {
		t.Error("decode string error")
	}
	if c[12] != 234 {
		t.Error("decode int error")
	}
	if c["user"].(User).Username != "asta" {
		t.Error("decode struct error")
	}
}

type User struct {
	Username string
	NickName string
}

func TestGenerate(t *testing.T) {
	str := generateRandomKey(20)
	if len(str) != 20 {
		t.Fatal("generate length is not equal to 20")
	}
}

func TestCookieEncodeDecode(t *testing.T) {
	hashKey := "testhashKey"
	blockkey := generateRandomKey(16)
	block, err := aes.NewCipher(blockkey)
	if err != nil {
		t.Fatal("NewCipher:", err)
	}
	securityName := string(generateRandomKey(20))
	val := make(map[interface{}]interface{})
	val["name"] = "user001"
	val["gender"] = "male"
	str, err := encodeCookie(block, hashKey, securityName, val)
	if err != nil {
		t.Fatal("encodeCookie:", err)
	}
	dst, err := decodeCookie(block, hashKey, securityName, str, 3600)
	if err != nil {
		t.Fatal("decodeCookie", err)
	}
	if dst["name"] != "user001" {
		t.Fatal("dst get map error")
	}
	if dst["gender"] != "male" {
		t.Fatal("dst get map error")
	}
}

func TestParseConfig(t *testing.T) {
	s := `{"cookieName":"gosessionid","gclifetime":3600}`
	cf := new(ManagerConfig)
	cf.EnableSetCookie = true
	err := json.Unmarshal([]byte(s), cf)
	if err != nil {
		t.Fatal("parse json error,", err)
	}
	if cf.CookieName != "gosessionid" {
		t.Fatal("parseconfig get cookiename error")
	}
	if cf.Gclifetime != 3600 {
		t.Fatal("parseconfig get gclifetime error")
	}

	cc := `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"cookiehashkey\"}"}`
	cf2 := new(ManagerConfig)
	cf2.EnableSetCookie = true
	err = json.Unmarshal([]byte(cc), cf2)
	if err != nil {
		t.Fatal("parse json error,", err)
	}
	if cf2.CookieName != "gosessionid" {
		t.Fatal("parseconfig get cookiename error")
	}
	if cf2.Gclifetime != 3600 {
		t.Fatal("parseconfig get gclifetime error")
	}
	if cf2.EnableSetCookie {
		t.Fatal("parseconfig get enableSetCookie error")
	}
	cconfig := new(cookieConfig)
	err = json.Unmarshal([]byte(cf2.ProviderConfig), cconfig)
	if err != nil {
		t.Fatal("parse ProviderConfig err,", err)
	}
	if cconfig.CookieName != "gosessionid" {
		t.Fatal("ProviderConfig get cookieName error")
	}
	if cconfig.SecurityKey != "cookiehashkey" {
		t.Fatal("ProviderConfig get securityKey error")
	}
}
