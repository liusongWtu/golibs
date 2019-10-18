package oauth

import (
	"encoding/json"
	"errors"
	"strings"
)

const qq_getaccesstoken_url = "https://graph.qq.com/oauth2.0/token"
const qq_getuserinfo_url = "https://graph.qq.com/user/get_user_info"
const qq_openid_url = "https://graph.qq.com/oauth2.0/me"

var qqOAuth = &QQOAuth{}

type QQOAuth struct {
	appkey       string
	appsecret    string
	redirect_url string
}

func (oauth *QQOAuth) Init(conf map[string]string) {
	oauth.appkey = conf["appkey"]
	oauth.appsecret = conf["appsecret"]
	oauth.redirect_url = conf["redirect_url"]
}

func (oauth *QQOAuth) GetAccesstoken(code string) (map[string]interface{}, error) {
	request := Get(qq_getaccesstoken_url)
	request.Param("grant_type", "authorization_code")
	request.Param("client_id", oauth.appkey)
	request.Param("client_secret", oauth.appsecret)
	request.Param("code", code)
	request.Param("redirect_uri", oauth.redirect_url)

	response, err := request.String()
	if err != nil {
		return nil, err
	}
	if strings.Contains(response, "callback") {
		return nil, nil
	}
	temp := strings.Split(response, "&")[0]
	accesstoken := strings.Split(temp, "=")[1]
	return map[string]interface{}{"access_token": accesstoken}, nil
}

func (oauth *QQOAuth) GetUserinfo(accesstoken string, openid string) (map[string]interface{}, error) {
	request := Get(qq_getuserinfo_url)
	request.Param("access_token", accesstoken)
	request.Param("oauth_consumer_key", oauth.appkey)
	request.Param("openid", openid)
	var response map[string]interface{}
	err := request.ToJSON(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (oauth *QQOAuth) Authorize(code string) (AuthorizeResult, error) {
	accesstokenResponse, err := oauth.GetAccesstoken(code)
	if err != nil {
		return AuthorizeResult{false, nil}, err
	}
	if accesstokenResponse == nil {
		return AuthorizeResult{false, nil}, nil
	}
	accesstoken := accesstokenResponse["access_token"].(string)
	expire := accesstokenResponse["expires_in"].(float64)
	openidResponse := oauth.GetOpenid(accesstoken)
	if _, ok := openidResponse["error"]; ok { //获取openid接口返回错误
		return AuthorizeResult{false, openidResponse}, errors.New("openid接口返回错误")
	}
	openid := openidResponse["openid"].(string)
	unionid := openidResponse["unionid"].(string)

	userInfo, err := oauth.GetUserinfo(accesstoken, openid)
	if err != nil {
		return AuthorizeResult{false, nil}, err
	}
	if userInfo == nil {
		return AuthorizeResult{false, nil}, nil
	}
	var sex int
	gender, ok := userInfo["gender"]
	if !ok {
		sex = 1
	}
	if gender.(string) == "女" {
		sex = 2
	} else {
		sex = 1
	}
	return AuthorizeResult{true, map[string]interface{}{
		"openid":       openid,
		"unionid":      unionid,
		"nickname":     userInfo["nickname"].(string),
		"sex":          sex,
		"avatar_url":   userInfo["figureurl_qq_1"].(string), // QQ头像 40x40尺寸
		"access_token": accesstoken,
		"expire":       expire,
		"platform":     "qq",
	}}, nil
}

func (oauth *QQOAuth) GetOpenid(accesstoken string) map[string]interface{} {
	request := Get(qq_openid_url)
	request.Param("access_token", accesstoken)
	request.Param("unionid", "1")
	responseStr, _ := request.String()
	var response map[string]interface{}
	json.Unmarshal([]byte(responseStr[10:len(responseStr)-3]), &response)
	return response
}

func init() {
	RegisterPlatform("qq", qqOAuth)
}
