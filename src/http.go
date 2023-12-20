package src

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/twmb/murmur3"
	"golang.org/x/net/html/charset"
)

func rndua() string {
	ua := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_4) AppleWebKit/537.13 (KHTML, like Gecko) Chrome/24.0.1290.1 Safari/537.13",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.0 (KHTML, like Gecko) Chrome/4.0.212.0 Safari/532.0",
		"Mozilla/5.0 (Windows; U; Windows NT 6.0; en-US) AppleWebKit/530.5 (KHTML, like Gecko) Chrome/2.0.172.23 Safari/530.5",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US) AppleWebKit/534.21 (KHTML, like Gecko) Chrome/11.0.682.0 Safari/534.21",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.13 (KHTML, like Gecko) Chrome/24.0.1290.1 Safari/537.13",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US) AppleWebKit/532.2 (KHTML, like Gecko) Chrome/4.0.222.5 Safari/532.2",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US) AppleWebKit/532.1 (KHTML, like Gecko) Chrome/4.0.219.4 Safari/532.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",
		"Mozilla/5.0 (X11; CrOS i686 4319.74.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.57 Safari/537.36",
		"Mozilla/5.0 (X11; U; Linux i686; en-US) AppleWebKit/534.13 (KHTML, like Gecko) Chrome/9.0.597.84 Safari/534.13",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.24 (KHTML, like Gecko) Chrome/19.0.1055.1 Safari/535.24",
		"Mozilla/5.0 (Macintosh; PPC Mac OS X 10_6_7) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.790.0 Safari/535.1",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US) AppleWebKit/533.3 (KHTML, like Gecko) Chrome/5.0.353.0 Safari/533.3",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US) AppleWebKit/532.2 (KHTML, like Gecko) Chrome/4.0.222.4 Safari/532.2",
		"Mozilla/5.0 (Windows NT 6.0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.66 Safari/535.11",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.3 (KHTML, like Gecko) Chrome/6.0.459.0 Safari/534.3",
		"Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US) AppleWebKit/534.3 (KHTML, like Gecko) Chrome/6.0.462.0 Safari/534.3",
		"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36 Mozilla/5.0 (iPad; U; CPU OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B334b Safari/531.21.10",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US) AppleWebKit/530.5 (KHTML, like Gecko) Chrome/2.0.172.43 Safari/530.5",
	}
	n := rand.Intn(20)
	return ua[n]
}

func findJSFiles(httpbody string) []string {
	var jsFiles []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(httpbody))
	if err != nil {
		return jsFiles
	}
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && strings.HasSuffix(src, ".js") {
			jsFiles = append(jsFiles, src)
		}
	})
	return jsFiles
}

func gettitle(httpbody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(httpbody))
	if err != nil {
		return "Not found"
	}
	// title := doc.Find("title").Text()
	title := doc.Find("title").First().Text()
	title = strings.Replace(title, "\n", "", -1)
	title = strings.Trim(title, " ")
	return title
}

func Mmh3Hash32(raw []byte) string {
	var h32 hash.Hash32 = murmur3.New32()
	_, err := h32.Write([]byte(raw))
	if err == nil {
		return fmt.Sprintf("%d", int32(h32.Sum32()))
	} else {
		//log.Println("favicon Mmh3Hash32 error:", err)
		return "0"
	}
}

func StandBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}

func favicon_hash(url string) [2]string {
	result_hash := [2]string{"0", "0"}
	timeout := time.Duration(8 * time.Second)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   timeout,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		//log.Println("favicon client error:", err)
		return result_hash
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//log.Println("favicon file read error: ", err)
		return result_hash
	}
	mmh3icon := Mmh3Hash32(StandBase64(body))
	// 计算 MD5 哈希值
	md5hash := md5.Sum(body)
	md5icon := fmt.Sprintf("%x", md5hash)
	result_hash[0] = mmh3icon
	result_hash[1] = md5icon
	return result_hash
}

func Http(urlstring string) {
	/*
		TLSClientConfig 字段用于配置 HTTPS 连接的 TLS 客户端配置
		InsecureSkipVerify 被设置为 true。这个选项的作用是，它允许客户端跳过对服务端证书的验证。
	*/
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
	req, err := http.NewRequest("GET", urlstring, nil)
	if err != nil {
		return
	}

	parsedURL, err := url.Parse(urlstring)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	host := parsedURL.Host
	cookie := &http.Cookie{
		Name:  "rememberMe",
		Value: "1",
	}
	req.AddCookie(cookie)
	req.Header.Set("Accept", "*/*;q=0.8")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", rndua())
	req.Header.Set("Host", host)

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("Error send resp:", err)
		return
	}
	/*
		resp.Body.Close() 来关闭该响应体
		defer 关键字用于延迟函数的执行，它会在包含它的函数返回之前执行其后的语句。
		在当前函数返回之前（也就是 main() 函数执行完毕之前），无论程序执行过程中是否发生错误，都会执行 resp.Body.Close() 来关闭响应体所占用的资源。
	*/
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("Error read resp.Body:", err)
		return
	}

	/*
		获取 Content-Type 头部信息
		将 Content-Type 转换为小写
	*/
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	contentType = strings.ToLower(contentType)

	// 将响应体转换为指定的编码格式
	httpBodyEncoding, err := charset.NewReader(bytes.NewReader(body), contentType)
	if err != nil {
		// fmt.Println("Error httpBodyEncoding:", err)
		return
	}

	// 读取转换后的内容
	convertedBody, err := io.ReadAll(httpBodyEncoding)
	if err != nil {
		// fmt.Println("Error convertedBody:", err)
		return
	}
	httpBody := string(convertedBody)

	title := gettitle(httpBody)
	httpheader := resp.Header
	var server string
	capital, ok := httpheader["Server"]
	if ok {
		server = capital[0]
	} else {
		Powered, ok := httpheader["X-Powered-By"]
		if ok {
			server = Powered[0]
		} else {
			server = "None"
		}
	}

	jsFiles := findJSFiles(httpBody)
	hashicon := favicon_hash(urlstring + "/favicon.ico")

	s := Resps{urlstring, httpBody, resp.Header, server, resp.StatusCode, len(httpBody), strings.TrimSpace(title), hashicon[0], hashicon[1], jsFiles}

	// 检测
	Check(s)
	// defer wg.Done()

	// js跳转
	reg := `(?s)<script>.*?window\.location=['"](.+?)['"]`
	// 编译正则表达式
	r := regexp.MustCompile(reg)
	// 在 httpBody 中查找匹配项
	match := r.FindStringSubmatch(httpBody)
	// 排除
	// exclude1 := "window.location='/'"
	// exclude2 := "window.location='./m/'"
	// fmt.Println(match)
	if len(match) > 0 {
		if jsJump < 3 {
			jsJump++
			// fmt.Println(urlstring + match[1])
			Http(urlstring + match[1])
		}

	}
}
