package requests

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
)

// NewClient 使用的选项可使用 | 来设置多个
const (
	OptIdle        = 0         // 不使用选项
	OptCookieJar   = 1 << iota // 使用 cookie
	OptNoSSLVerify             // 不检测 ssl 有效性
)

type SClient struct {
	C     *http.Client
	cOnce sync.Once

	middlewareReq  *reqMiddleware
	mReqOnce       sync.Once
	middlewareResp *respMiddleware
	mRespOnce      sync.Once
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
	return &SClient{C: c}
}

func (c *SClient) Request(method, remoteUrl string, args ...interface{}) (resp *SResponse, err error) {
	method = strings.ToUpper(method)
	req, err := NewRequest(method, remoteUrl, args...)
	if err != nil {
		return
	}
	if c.C == nil {
		c.cOnce.Do(func() {
			c.C = &http.Client{}
		})
	}
	return c.Do(req)
}

func (c *SClient) Do(req *SRequest) (resp *SResponse, err error) {
	if c.middlewareReq != nil {
		for _, f := range c.middlewareReq.middleware {
			f(req)
		}
	}
	rawResp, err := c.C.Do(req.Req)
	if err != nil {
		return
	}
	resp = &SResponse{Resp: rawResp}
	if c.middlewareResp != nil {
		for _, f := range c.middlewareResp.middleware {
			f(resp)
		}
	}
	return
}

func (c *SClient) AddReqMiddleware(f func(req *SRequest)) (id int) {
	if c.middlewareReq == nil {
		c.mReqOnce.Do(func() {
			c.middlewareReq = new(reqMiddleware)
		})
	}
	return c.middlewareReq.Add(f)
}

func (c *SClient) RemoveReqMiddleware(id int) {
	if c.middlewareReq != nil {
		c.middlewareReq.Remove(id)
	}
}
func (c *SClient) ClearReqMiddleware() {
	if c.middlewareReq != nil {
		c.middlewareReq.Clear()
	}
}

func (c *SClient) AddRespMiddleware(f func(req *SResponse)) (id int) {
	if c.middlewareResp == nil {
		c.mRespOnce.Do(func() {
			c.middlewareResp = new(respMiddleware)
		})
	}
	return c.middlewareResp.Add(f)
}

func (c *SClient) RemoveRespMiddleware(id int) {
	if c.middlewareResp != nil {
		c.middlewareResp.Remove(id)
	}
}
func (c *SClient) ClearRespMiddleware() {
	if c.middlewareResp != nil {
		c.middlewareResp.Clear()
	}
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
