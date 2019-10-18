package oauth

import "errors"

const weibo_getaccesstoken_url = "https://api.weibo.com/oauth2/access_token"
const weibo_getuserinfo_url = "https://api.weibo.com/2/users/show.json"

var weiboOAuth = &WeiboOAuth{}

type WeiboOAuth struct {
	appkey       string
	appsecret    string
	redirect_url string
}

func (oauth *WeiboOAuth) Init(conf map[string]string) {
	oauth.appkey = conf["appkey"]
	oauth.appsecret = conf["appsecret"]
	oauth.redirect_url = conf["redirect_url"]
}

func (oauth *WeiboOAuth) GetAccesstoken(code string) (map[string]interface{}, error) {
	request := Post(weibo_getaccesstoken_url)
	request.Param("client_id", oauth.appkey)
	request.Param("client_secret", oauth.appsecret)
	request.Param("grant_type", "authorization_code")
	request.Param("code", code)
	request.Param("redirect_uri", oauth.redirect_url)
	var response map[string]interface{}
	err := request.ToJSON(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (oauth *WeiboOAuth) GetUserinfo(accesstoken string, openid string) (map[string]interface{}, error) {
	request := Get(weibo_getuserinfo_url)
	request.Param("access_token", accesstoken)
	request.Param("uid", openid)
	var response map[string]interface{}
	err := request.ToJSON(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (oauth *WeiboOAuth) Authorize(code string) (AuthorizeResult, error) {
	accesstokenResponse, err := oauth.GetAccesstoken(code)
	if err != nil {
		return AuthorizeResult{false, nil}, err
	}
	if accesstokenResponse == nil {
		return AuthorizeResult{false, nil}, nil
	}
	_, ok := accesstokenResponse["error_code"] //获取accesstoken接口返回错误码
	if ok {
		return AuthorizeResult{false, accesstokenResponse}, errors.New("accesstoken接口返回错误码")
	}
	openid := accesstokenResponse["uid"].(string)
	accesstoken := accesstokenResponse["access_token"].(string)
	expire := accesstokenResponse["expires_in"].(float64)
	userInfo, err := oauth.GetUserinfo(accesstoken, openid)
	if err != nil {
		return AuthorizeResult{false, nil}, err
	}
	if userInfo == nil {
		return AuthorizeResult{false, nil}, nil
	}
	_, ok = userInfo["error_code"] //获取用户信息接口返回错误码
	if ok {
		return AuthorizeResult{false, userInfo}, errors.New("用户信息接口返回错误")
	}
	var sex int
	if userInfo["gender"].(string) == "m" {
		sex = 1
	} else if userInfo["gender"].(string) == "f" {
		sex = 2
	} else if userInfo["gender"].(string) == "n" {
		sex = 0
	}
	return AuthorizeResult{true, map[string]interface{}{
		"openid":       openid,
		"unionid":      "",
		"nickname":     userInfo["screen_name"].(string),
		"sex":          sex,
		"avatar_url":   userInfo["profile_image_url"].(string),
		"access_token": accesstoken,
		"expire":       expire,
		"platform":     "weibo",
	}}, nil
}

func init() {
	RegisterPlatform("weibo", weiboOAuth)
}
