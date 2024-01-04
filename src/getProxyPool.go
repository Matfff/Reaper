package src

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetProxyPool() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 创建一个 http 客户端，并设置不跟随重定向
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 构建一个 GET 请求
	req, err := http.NewRequest("GET", HttpProxyPool, nil)
	if err != nil {
		return
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("Error send resp:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("Error read resp.Body:", err)
		return
	}

	// fmt.Println(string(body))
	proxyArray := strings.Split(string(body), "\r\n")
	// fmt.Println(proxyArray)
	fmt.Println("获取IP个数: ", len(proxyArray))

	// 遍历代理IP地址数组
	for _, proxy := range proxyArray {
		// ProxyQueue <- url // 将url加入队列
		ProxyQueue.Push(proxy)
	}
}

func GetProxy() string {
	// 调用代理池
	if ProxyQueue.Len() == 0 {
		GetProxyPool()
	}

	proxyData := ProxyQueue.Pop()
	proxy, _ := proxyData.(string)

	return proxy
}
