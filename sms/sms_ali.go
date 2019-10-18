package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// SmsAliReply短信返回
type SmsAliReply struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

func NewSmsAli() SmsManager {
	return &SmsAli{}
}

type SmsAli struct{}

func (sa *SmsAli) Send(data map[string]string) error {
	//组织参数
	params := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", rand.Int63()),
		"AccessKeyId":      data["access_key"],
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"RegionId":         "cn-hangzhou",
		"PhoneNumbers":     data["phone"], //多个手机号，","相隔
		"SignName":         data["sign_name"],
		"TemplateParam":    data["template_param"],
		"TemplateCode":     data["template_code"],
	}
	//排序key
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//组织字符串
	var sortQueryString string
	for _, v := range keys {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, sa.replace(v), sa.replace(params[v]))
	}
	//组织sign
	stringToSign := fmt.Sprintf("GET&%s&%s", sa.replace("/"), sa.replace(sortQueryString[1:]))
	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", data["access_secret"])))
	mac.Write([]byte(stringToSign))
	sign := sa.replace(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	//发送访问请求
	strUrl := fmt.Sprintf("http://dysmsapi.aliyuncs.com/?Signature=%s%s", sign, sortQueryString)
	resp, err := http.Get(strUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//组织返回数据
	ssr := &SmsAliReply{}
	if err := json.Unmarshal(body, ssr); err != nil {
		return err
	}
	if ssr.Code == "SignatureNonceUsed" {
		return sa.Send(data)
	} else if ssr.Code != "OK" {
		return errors.New(ssr.Code)
	}
	return nil
}

//替换字符串
func (sa *SmsAli) replace(in string) string {
	rep := strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
	return rep.Replace(url.QueryEscape(in))
}

func init() {
	Register("sms_ali", NewSmsAli)
}
