# requests for Golang
一个 Golang 版本的 requests 网络请求库(类似于 python 版本)，方便苦于 Golang 繁琐的 http 请求的 Gopher 使用

## 构思
- 关键字参数  
Golang 没有关键字参数，通过自定义类型 + 参数类型断言模拟。
如地址栏参数在 python 版本中使用 params = {"abc":456} 传入，
在 Golang 版本中通过 Params{"abc":456} 传入。

## 使用
### 类型和常用方法
- 通用类型
  - Dict: map[string]interface{} 的别名
  - List: []interface{} 的别名
  
- 请求时可以使用的类型
  - Headers: 请求头字典，map[string]string 的别名
  - Params: 地址栏参数字典，map[string]string 的别名
  - Cookies: Cookie 数组，[]*http.Cookie 的别名，可以使用 Map2Cookies 快速生成
  - Xform: x-www-form-urlencoded 格式的 body 字典，map[string]string 的别名
  - FormData: form-data 格式的 body 字典，map[string]interface{} 的别名，interface{} 可以选择 string（一般场景）或 FormDataFile（上传文件）
  - FileBody: 文件流形式的 body，io.Reader 的别名
  - JsonBody: 任意对象类型的 body，会被转换成 json 字符串，interface{} 的别名
  - StringBody: 字符串类型的 body，string 的别名
  - BytesBody: 二进制数组类型的 body，[]byte 的别名
> PS：body 只能传入一个，否则行为无法保证

- 响应中可以使用类型和方法
 - String() 会返回一个响应 body 的 string 副本
 - Bytes() 会返回一个响应 body 的 []byte 副本
 - Object() 会返回一个响应 body 的 interface{} 副本（json 解码）
 - Dict() 会返回一个响应 body 的 Dict 副本（json 解码）
 - List() 会返回一个响应 body 的 List 副本（json 解码）
 - Json(ptr) 会尝试把 body 解码到一个传入的指针中
 - File(io.Writer) 会尝试把 body 保存到一个传入的文件流中（整个过程为流式，不需考虑内存占用）
 - Discard() 会丢弃整个响应
> PS: 除 File 、 Discard 之外的方法会缓存响应的 body 内容，可以重复调用
> PS: 大部分方法会返回一个 error 值用来标识是否调用成功

- Find 方法  
使用多级参数查找一个对象中的元素，可用于读取配置文件或网络返回的 JSON 对象
  - 函数定义： func Find(obj interface{}, args ...interface{}) (interface{}, bool)  
    - args 中的元素可以是 string 或 int, string 代表查询的是字典， int 代表查询的是数组
    - 如 JSON 对象 {"abc":[123,456,789]} 中的 789 可以使用参数 ["abc",2] 查出
    - 再如 JSON 对象 ["abc",[123,456,789,{"def":111}]] 中的 111 可以使用参数 [1,3,"def"] 查出
    - 再如 JSON 对象 ["abc",[123,456,789,{"def":111}]] 中的 {"def":111} 可以使用参数 [1,3] 查出
    - 返回值分别为查找结果和是否找到
  > 包装了常用的 int int64 uint string float32 之类的方法 

- 使用例子  
获取 `https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=7` 并解析
```
package main

import (
	"fmt"

	"gitee.com/binshub/requests"
)

func main() {
	fmt.Println("start")
	resp, _ := requests.Get("https://cn.bing.com/HPImageArchive.aspx",
		requests.Params{"format": "js", "idx": "0", "n": "7"})
	if resp.Resp.StatusCode != 200 {
		fmt.Printf("response code err. expect %v, get %v", 200, resp.Resp.StatusCode)
	}
	res, _ := resp.Object()
	url, ok := requests.FindString(res, "images", 1, "url")
	if !ok {
		fmt.Println("data parse err.")
	}
	fmt.Println(url)
}
```

