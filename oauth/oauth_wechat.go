package oauth

import "errors"

const wechat_getaccesstoken_url = "https://api.weixin.qq.com/sns/oauth2/access_token"
const wechat_getuserinfo_url = "https://api.weixin.qq.com/sns/userinfo"

var wechatOAuth = &WechatOAuth{}

type WechatOAuth struct {
	appkey    string
	appsecret string
}

func (oauth *WechatOAuth) Init(conf map[string]string) {
	oauth.appkey = conf["appkey"]
	oauth.appsecret = conf["appsecret"]
}

func (oauth *WechatOAuth) GetAccesstoken(code string) (map[string]interface{}, error) {
	request := Get(wechat_getaccesstoken_url)
	request.Param("appid", oauth.appkey)
	request.Param("secret", oauth.appsecret)
	request.Param("code", code)
	request.Param("grant_type", "authorization_code")
	var response map[string]interface{}
	err := request.ToJSON(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (oauth *WechatOAuth) GetUserinfo(accesstoken string, openid string) (map[string]interface{}, error) {
	request := Get(wechat_getuserinfo_url)
	request.Param("access_token", accesstoken)
	request.Param("openid", openid)
	var response map[string]interface{}
	err := request.ToJSON(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (oauth *WechatOAuth) Authorize(code string) (AuthorizeResult, error) {
	accesstokenResponse, err := oauth.GetAccesstoken(code)
	if err != nil {
		return AuthorizeResult{false, nil}, err
	}
	if accesstokenResponse == nil {
		return AuthorizeResult{false, nil}, nil
	}
	_, ok := accesstokenResponse["errcode"] //获取accesstoken接口返回错误码
	if ok {
		return AuthorizeResult{false, accesstokenResponse}, errors.New("accesstoken接口返回错误")
	}
	openid := accesstokenResponse["openid"].(string)
	accesstoken := accesstokenResponse["access_token"].(string)
	expire := accesstokenResponse["expires_in"].(float64)
	userInfo, err := oauth.GetUserinfo(accesstoken, openid)
	if err != nil {
		return AuthorizeResult{false, nil}, err
	}
	if userInfo == nil {
		return AuthorizeResult{false, nil}, nil
	}
	_, ok = userInfo["errcode"] //获取用户信息接口返回错误码
	if ok {
		return AuthorizeResult{false, userInfo}, errors.New("用户信息接口返回错误")
	}

	return AuthorizeResult{true, map[string]interface{}{
		"openid":       userInfo["openid"].(string),
		"unionid":      userInfo["unionid"].(string),
		"nickname":     userInfo["nickname"].(string),
		"sex":          userInfo["sex"].(float64),
		"avatar_url":   userInfo["headimgurl"].(string),
		"access_token": accesstoken,
		"expire":       expire,
		"platform":     "wechat",
	}}, nil
}

func init() {
	RegisterPlatform("wechat", wechatOAuth)
}
