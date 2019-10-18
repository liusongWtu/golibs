package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/big"
)

const (
	Ok                = 1 //正确
	Expired           = 2 //token过期
	HandleError       = 3 //token解析处理错误
	TypeExchangeError = 4 //类型转换错误
)

type Jwt struct {
	//ES256 keys---调用Jwt.GenerateKey()生成
	esKeyD string
	esKeyX string
	esKeyY string
	//HS256 signed key
	hsKey string
}

func (jw *Jwt) Init(data map[string]string) {
	jw.esKeyD = data["esKeyD"]
	jw.esKeyX = data["esKeyX"]
	jw.esKeyY = data["esKeyY"]
	jw.hsKey = data["hsKey"]
}

//生成esKeyD,esKeyX,esKeyY
func (jw *Jwt) GenerateKey() (string, string, string) {
	randKey := rand.Reader
	var err error
	prk, err := ecdsa.GenerateKey(elliptic.P256(), randKey)
	if err != nil {
		fmt.Println("generate key error", err)
	}

	puk := prk.PublicKey
	fmt.Printf("esKeyD=%X\n", prk.D)
	fmt.Printf("esKeyX=%X\n", prk.X)
	fmt.Printf("esKeyY=%X\n", prk.Y)
	fmt.Println("prk", prk, " \npbk", puk)

	return fmt.Sprintf("%X", prk.D), fmt.Sprintf("%X", prk.X), fmt.Sprintf("%X", prk.Y)
}

//获取签名算法为ES256的token
func (jw *Jwt) GetEStoken(claims jwt.MapClaims) (string, error) {
	keyD := new(big.Int)
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyD.SetString(jw.esKeyD, 16)
	keyX.SetString(jw.esKeyX, 16)
	keyY.SetString(jw.esKeyY, 16)

	tokn := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	privateKey := ecdsa.PrivateKey{D: keyD, PublicKey: publicKey}
	tokenStr, err := tokn.SignedString(&privateKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

//获取签名算法为HS256的token
func (jw *Jwt) GetHStoken(claims jwt.MapClaims) (string, error) {
	tokn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密算法是HS256时，这里的SignedString必须是[]byte（）类型
	tokenStr, err := tokn.SignedString([]byte(jw.hsKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

//解析签名算法为ES256的token
func (jw *Jwt) ParseEStoken(tokenES string) (jwt.MapClaims, int) {
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyX.SetString(jw.esKeyX, 16)
	keyY.SetString(jw.esKeyY, 16)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	tokn, err := jwt.Parse(tokenES, func(token *jwt.Token) (interface{}, error) {
		return &publicKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, Expired
			} else {
				return nil, HandleError
			}
		} else {
			return nil, HandleError
		}
		return nil, Ok
	}

	claims, ok := tokn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, TypeExchangeError
	}

	return claims, Ok
}

//解析签名算法为HS256的token
func (jw *Jwt) ParseHStoken(tokenString string) (jwt.MapClaims, int) {
	tokn, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jw.hsKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, Expired
			} else {
				return nil, HandleError
			}
		} else {
			return nil, HandleError
		}
		return nil, Ok
	}

	claims, ok := tokn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, TypeExchangeError
	}

	return claims, Ok
}
