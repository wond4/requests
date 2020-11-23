package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type SResponse struct {
	Resp    *http.Response
	content []byte
}

// 丢弃 Body 中的内容并关闭
func (b *SResponse) Discard() {
	if b.Resp.Body == nil {
		return
	}
	if b.content != nil {
		b.content = nil
		return
	}
	defer func() {
		_ = b.Resp.Body.Close()
	}()
	_, _ = io.Copy(ioutil.Discard, b.Resp.Body)
}

// 以 []byte 返回 SResponse 中的内容并关闭
func (b *SResponse) Bytes() (bytes []byte, err error) {
	if b.Resp.Body == nil {
		return nil, fmt.Errorf("body is null")
	}
	if b.content != nil {
		return b.content, nil
	}
	defer func() {
		err = b.Resp.Body.Close()
	}()
	b.content, err = ioutil.ReadAll(b.Resp.Body)
	if err != nil {
		b.content = nil
		return
	}
	return b.content, nil
}

// 以 string 返回 SResponse 中的内容并关闭
func (b *SResponse) String() (str string, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	str = string(content)
	return
}

// 以 Object(对象) 返回 SResponse 中的内容并关闭
func (b *SResponse) Object() (obj interface{}, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &obj)
	return
}

// 以 Dict(字典) 返回 SResponse 中的内容并关闭
func (b *SResponse) Dict() (dict Dict, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &dict)
	return
}

// 以 List(列表) 返回 SResponse 中的内容并关闭
func (b *SResponse) List() (list List, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &list)
	return
}

// 以 Json 对象序列化 SResponse 中的内容并关闭
func (b *SResponse) Json(objPointer interface{}) (err error) {
	content, err := b.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(content, objPointer)
}

// 把 SResponse 存到 file 中(不缓存响应内容)
func (b *SResponse) File(f io.Writer) (n int64, err error) {
	if b.content != nil {
		nn, err := f.Write(b.content)
		return int64(nn), err
	}
	if b.Resp.Body == nil {
		return 0, fmt.Errorf("body is null")
	}
	defer func() {
		err = b.Resp.Body.Close()
	}()
	return io.Copy(f, b.Resp.Body)
}
