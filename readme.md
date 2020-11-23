# requests for Golang
一个 Golang 版本的 requests 网络请求库(类似于 python 版本)，方便苦于 Golang 繁琐的 http 请求的 Gopher 使用

## 构思
- 关键字参数  
Golang 没有关键字参数，通过自定义类型 + 参数类型断言模拟。
如地址栏参数在 python 版本中使用 params = {"abc":456} 传入，
在 Golang 版本中通过 Params{"abc":456} 传入。

## 使用
- 请求体使用
可以传入 Headers、Params、Cookies、Xform、FormData、FileBody、JsonBody、StringBody、BytesBody 等多种参数


- 响应体使用  
为响应体添加了 String、Dict、List、Object、Json、Discard、Bytes、File 等方法，用于不同的解析方式。
File 方法用于保存文件或流式接收，不会缓存响应内容（不限响应体大小）。
其他方法在解析同时会保存响应内容，可以多次使用不同解析方式处理（建议用于响应体不大的情况）。

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

