package requests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

type SResponse struct {
	Resp    *http.Response
	content []byte
}

// Discard 丢弃 Body 中的内容并关闭
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

// Bytes 以 []byte 返回 SResponse 中的内容并关闭
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

// MustBytes 以 []byte 返回 SResponse 中的内容并关闭，如果没有内容返回 nil
func (b *SResponse) MustBytes() (bytes []byte) {
	bytes, _ = b.Bytes()
	return
}

// String 以 string 返回 SResponse 中的内容并关闭
func (b *SResponse) String() (str string, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	str = string(content)
	return
}

// MustString 以 string 返回 SResponse 中的内容并关闭，如果没有内容返回空字符串
func (b *SResponse) MustString() (str string) {
	str, _ = b.String()
	return
}

// Object 以 Object(对象) 返回 SResponse 中的内容并关闭
func (b *SResponse) Object() (obj interface{}, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &obj)
	return
}

// Dict 以 Dict(字典) 返回 SResponse 中的内容并关闭
func (b *SResponse) Dict() (dict Dict, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &dict)
	return
}

// List 以 List(列表) 返回 SResponse 中的内容并关闭
func (b *SResponse) List() (list List, err error) {
	content, err := b.Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &list)
	return
}

// Json 以 Json 对象反序列化 SResponse 中的内容并关闭
func (b *SResponse) Json(objPointer interface{}) (err error) {
	content, err := b.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(content, objPointer)
}

// File 把 SResponse 存到 file 中(不缓存响应内容)
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

// AutoObject 根据 ContentType 自动反序列化 SResponse 中的内容并关闭
func (b *SResponse) AutoObject(objPointer interface{}) (err error) {
	contentType := b.Resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = b.Resp.Header.Get("content-type")
	}
	content, err := b.Bytes()
	if strings.Contains(contentType, "application/xml") {
		if err = xml.Unmarshal(content, objPointer); err != nil {
			return err
		}
		return nil
	} else if strings.Contains(contentType, "application/json") {
		if err = json.Unmarshal(content, objPointer); err != nil {
			return err
		}
		return nil
	} else if strings.Contains(contentType, "multipart/form-data") {
		res, ok := objPointer.(*[]string)
		if !ok {
			return errors.New("decode form-data err")
		}
		flag := "boundary="
		index := strings.Index(contentType, flag)
		if index < 0 {
			return errors.New("decode form-data err")
		}
		index += len(flag)
		boundary := contentType[index:]
		form := multipart.NewReader(bytes.NewReader(content), boundary)
		for {
			part, err := form.NextPart()
			if err != nil {
				break
			}
			bs, err := ioutil.ReadAll(part)
			if err != nil {
				continue
			}
			*res = append(*res, string(bs))
		}
		return nil
	}
	return errors.New("undefined response type")
}
