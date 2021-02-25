package requests

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// NewClient 使用的选项可使用 | 来设置多个
const (
	OptIdle        = 0         // 不使用选项
	OptCookieJar   = 1 << iota // 使用 cookie
	OptNoSSLVerify             // 不检测 ssl 有效性
)

type SClient struct {
	C *http.Client
}

func Session(option uint) *SClient {
	c := &http.Client{}
	if OptCookieJar&option != 0 {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil
		}
		c.Jar = jar
	}
	if OptNoSSLVerify&option != 0 {
		tp := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		c.Transport = tp
	}
	return &SClient{c}
}

func (c *SClient) Request(method, remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	method = strings.ToUpper(method)
	req, err := NewRequest(method, remoteUrl, args...)
	if err != nil {
		return
	}
	if c.C == nil {
		c.C = &http.Client{}
	}
	return c.Do(req)
}

func (c *SClient) Do(req *SRequest) (resp *SResponse, err error) {
	rawResp, err := c.C.Do(req.Req)
	if err != nil {
		return
	}
	resp = &SResponse{Resp: rawResp}
	return
}

func (c *SClient) Get(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return c.Request("get", remoteUrl, args...)
}

func (c *SClient) Post(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return c.Request("post", remoteUrl, args...)
}

func (c *SClient) Put(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return c.Request("put", remoteUrl, args...)
}

func (c *SClient) Delete(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return c.Request("delete", remoteUrl, args...)
}

var DefaultClient = &SClient{}

func Request(method, remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return DefaultClient.Request(method, remoteUrl, args...)
}

func Get(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return DefaultClient.Get(remoteUrl, args...)
}

func Post(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return DefaultClient.Post(remoteUrl, args...)
}

func Put(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return DefaultClient.Put(remoteUrl, args...)
}

func Delete(remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	return DefaultClient.Delete(remoteUrl, args...)
}
