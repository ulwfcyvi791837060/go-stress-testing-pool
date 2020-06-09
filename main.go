/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 13:44
 */

package main

import (
	"flag"
	"fmt"
	"go-stress-testing-pool/model"
	"go-stress-testing-pool/server"
	"runtime"
	"strings"
)

type array []string

func (a *array) String() string {
	return fmt.Sprint(*a)
}

func (a *array) Set(s string) error {
	*a = append(*a, s)

	return nil
}


/**

./go_stress_testing_linux -c 62500 -n 1  -u ws://192.168.0.74:443/acc
`62500*16 = 100W `正好可以达到我们的要求

# 查看用法
go run main.go

# 使用请求百度页面
go run main.go -c 1 -n 100 -u https://www.baidu.com/

# 使用debug模式请求百度页面
go run main.go -c 1 -n 1 -d true -u https://www.baidu.com/

# 使用 curl文件(文件在curl目录下) 的方式请求
go run main.go -c 1 -n 1 -p curl/baidu.curl.txt

# 压测webSocket连接
go run main.go -c 10 -n 10 -u ws://127.0.0.1:8089/acc


示例: go run main.go -c 3000 -n 1 -u https://www.baidu.com/

 go run main.go -c 3000 -n 1 -p curl/test.chrome.curl.txt


压测地址或curl路径必填
当前请求参数: -c 1 -n 1 -d false -u
Usage of C:\Users\Administrator\AppData\Local\Temp\___go_build_main_go.exe:
  -H value
    	自定义头信息传递给服务器 示例:-header 'Content-Type: application/json'
  -c uint
    	并发数 (default 1)
  -d string
    	调试模式 (default "false")
  -data string
    	HTTP POST方式传送数据
  -n uint
    	请求总数 (default 1)
  -p string
    	curl文件路径
  -u string
    	压测地址
  -v string
    	验证方法 http 支持:statusCode、json webSocket支持:json
 */

// go 实现的压测工具
//
// 编译可执行文件
//go:generate go build main.go
func main() {

	runtime.GOMAXPROCS(1)

	var (
		concurrency uint64 // 并发数
		cycleNumber uint64 // 请求总数(单个并发)
		debugStr    string // 是否是debug
		requestUrl  string // 压测的url 目前支持，http/https ws/wss
		path        string // curl文件路径 http接口压测，自定义参数设置
		verify      string // verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
		headers     array  // 自定义头信息传递给服务器
		body        string // HTTP POST方式传送数据
	)

	flag.Uint64Var(&concurrency, "c", 1, "并发数")
	flag.Uint64Var(&cycleNumber, "n", 1, "请求总数")
	flag.StringVar(&debugStr, "d", "false", "调试模式")
	flag.StringVar(&requestUrl, "u", "", "压测地址")
	flag.StringVar(&path, "p", "", "curl文件路径")
	flag.StringVar(&verify, "v", "", "验证方法 http 支持:statusCode、json webSocket支持:json")
	flag.Var(&headers, "H", "自定义头信息传递给服务器 示例:-header 'Content-Type: application/json'")
	flag.StringVar(&body, "data", "", "HTTP POST方式传送数据")

	// 解析参数
	flag.Parse()


	// 可注释 go run main.go -c 3000 -n 1 -p curl/test.chrome.curl.txt -v json
	//go run main.go -c 3000 -n 1 -p curl/test.chrome.curl.txt -v json
	//go run main.go -c 3000 -n 1 -p curl/aws.heartBeat.chrome.curl.txt -v json

	concurrency = 1000
	cycleNumber = 1
	//path =  "curl/local.heartBeat.chrome.curl.txt"
	//path =  "curl/local.heartBeat.chrome.curl.txt"
	path =  "curl/aws.heartBeat.chrome.curl.txt"
	verify = "json"
	//debugStr = "true"


	if concurrency == 0 || cycleNumber == 0 || (requestUrl == "" && path == "") {
		fmt.Printf("示例: go run main.go -c 1 -n 1 -u https://www.baidu.com/ \n")
		fmt.Printf("压测地址或curl路径必填 \n")
		fmt.Printf("当前请求参数: -c %d -n %d -d %v -u %s \n", concurrency, cycleNumber, debugStr, requestUrl)

		flag.Usage()

		return
	}

	debug := strings.ToLower(debugStr) == "true"
	request, err := model.NewRequest(requestUrl, verify, 0, debug, path, headers, body)
	if err != nil {
		fmt.Printf("参数不合法 %v \n", err)

		return
	}

	fmt.Printf("\n 开始启动  并发数:%d  循环次数:%d \n", concurrency, cycleNumber)
	request.Print()

	// 开始处理
	//server.Dispose(concurrency, cycleNumber, request)
	server.DisposePool(concurrency, cycleNumber, request)

	return
}
