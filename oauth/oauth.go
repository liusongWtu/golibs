package oauth

import (
	"fmt"
)

var oauthes = make(map[string]OAuth)

type OAuth interface {
	Init(conf map[string]string)
	GetAccesstoken(code string) (map[string]interface{}, error)
	GetUserinfo(accesstoken string, openid string) (map[string]interface{}, error)
	Authorize(code string) (AuthorizeResult, error)
}

func RegisterPlatform(name string, oauth OAuth) {
	if oauth == nil {
		panic("Register simpleoauth instance is nil")
	}
	_, dup := oauthes[name]
	if dup {
		panic("The platform has registered already")
	}
	oauthes[name] = oauth
}

type Manager struct {
	oauth OAuth
}

func NewManager(platformName string, conf map[string]string) (*Manager, error) {
	oauth, ok := oauthes[platformName]
	if !ok {
		return nil, fmt.Errorf("unknown platform %q", platformName)
	}
	oauth.Init(conf)
	return &Manager{oauth}, nil
}

func (m *Manager) Authorize(code string) (AuthorizeResult, error) {
	return m.oauth.Authorize(code)
}

type AuthorizeResult struct {
	Result   bool
	Userinfo map[string]interface{}
}
