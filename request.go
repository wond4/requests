package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Headers map[string]string
type Params map[string]string // url query
type Cookies []*http.Cookie   // bytes body

// body types
// only one body can be effective
type Xform map[string]string         // x-www-form-urlencoded
type FormData map[string]interface{} // form data, value can be string or FormDataFile
type FileBody io.Reader              // body from file
type JsonBody interface{}            // json body
type StringBody string               // string body
type BytesBody []byte                // bytes body

// data type
type Dict map[string]interface{}
type List []interface{}

// if File is nil, will read file from Path
type FormDataFile struct {
	Path string
	File io.Reader
}

type SRequest struct {
	Req *http.Request
}

func (r *SRequest) parseArgs(args ...interface{}) (err error) {
	// need set url before call this method
	query := r.Req.URL.Query()
	var bBody BytesBody
	for _, arg := range args {
		switch t := arg.(type) {
		case Headers:
			for k, v := range t {
				r.Req.Header.Set(k, v)
			}
		case Params:
			for k, v := range t {
				query.Add(k, v)
			}
		case Cookies:
			for _, v := range t {
				r.Req.AddCookie(v)
			}
		case Xform:
			r.addXForm(t)
		case FormData:
			err = r.addFormData(t)
			if err != nil {
				return
			}
		case FileBody:
			r.Req.Body = ioutil.NopCloser(t)
		case StringBody:
			bBody = []byte(t)
		case BytesBody:
			bBody = t
		case JsonBody:
			bBody, err = json.Marshal(t)
			if err != nil {
				return
			}
		}
	}
	r.Req.URL.RawQuery = query.Encode()
	if bBody != nil {
		r.Req.Body = ioutil.NopCloser(bytes.NewReader(bBody))
	}
	return
}

func (r *SRequest) addXForm(form Xform) {
	f := &url.Values{}
	for k, v := range form {
		f.Add(k, v)
	}
	r.Req.Body = ioutil.NopCloser(strings.NewReader(f.Encode()))
	r.Req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func (r *SRequest) addFormData(data FormData) (err error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	needClose := true
	defer func() {
		if needClose {
			w.Close()
		}
	}()
	var part io.Writer
	for k, v := range data {
		switch t := v.(type) {
		case string:
			err = w.WriteField(k, t)
			if err != nil {
				return
			}
		case FormDataFile:
			fn := path2filename(t.Path)
			part, err = w.CreateFormFile(k, fn)
			if err != nil {
				return
			}
			if t.File == nil {
				var f *os.File
				f, err = os.Open(t.Path)
				if err != nil {
					return
				}
				_, err = io.Copy(part, f)
				f.Close()
				if err != nil {
					return
				}
			} else {
				_, err = io.Copy(part, t.File)
				if err != nil {
					return
				}
			}
		}
	}
	w.Close()
	needClose = false
	r.Req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
	r.Req.Header.Set("Content-Type", w.FormDataContentType())
	return
}

func path2filename(path string) string {
	index := strings.LastIndex(path, `\`)
	if index == -1 {
		index = strings.LastIndex(path, `/`)
	}
	if index == -1 {
		return path
	}
	return path[index+1:]
}

func NewRequest(method, remoteUrl string, args ...interface{}) (req *SRequest, err error) {
	req = &SRequest{Req: &http.Request{}}
	req.Req.URL, err = url.Parse(remoteUrl)
	if err != nil {
		return
	}
	req.Req.Method = method
	req.Req.Header = map[string][]string{}
	err = req.parseArgs(args...)
	return
}
