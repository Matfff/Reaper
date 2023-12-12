package src

import (
	"bufio"
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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/twmb/murmur3"
	"golang.org/x/net/html/charset"
)

type resps struct {
	url        string
	body       string
	header     map[string][]string
	server     string
	statuscode int
	length     int
	title      string
	mmh3icon   string
	md5icon    string
	jsfiles    []string
}

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

func writeToBuffer(data string) {
	BufferMutex.Lock()
	defer BufferMutex.Unlock()

	Buffer <- data // 写入数据到缓冲区

	// 当缓冲区大小达到一定条件时触发写入
	// fmt.Println(len(Buffer))
	if len(Buffer) >= 1024 {
		// WriteTrigger <- struct{}{}
		WriteFile()
	}
}

func WriteFile() {
	// 打开文件，使用 os.O_APPEND 模式打开，表示在文件末尾追加内容
	file, err := os.OpenFile(Output, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		// fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for len(Buffer) > 0 {
		data := <-Buffer
		_, err := writer.WriteString(data + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	fmt.Println("正在刷新缓冲区到文件...")
	writer.Flush() // 写入到文件
}

func removeDuplicates(arr []string) []string {
	encountered := map[string]bool{} // 创建一个 map 用于记录元素是否已经出现过
	result := []string{}             // 保存去重后的结果

	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true      // 将元素标记为已经出现过
			result = append(result, v) // 将不重复的元素加入到结果切片中
		}
	}
	return result
}

func Check(s resps) {
	// fileMutex.Lock()
	// defer fileMutex.Unlock()

	var headers []string
	var fps []string
	// 遍历 map 类型的数据
	for key, values := range s.header {
		header := key + ": " + values[0]
		headers = append(headers, header)
	}
	// fmt.Println(headers)
	str_headers := strings.Join(headers, ",")
	str_js := strings.Join(s.jsfiles, ",")

	// 输出解析后的数据示例
	for _, fingerprint := range Fps {
		haveFingerprint := true

		if len(fingerprint.Headers) == 0 && len(fingerprint.Body) == 0 && len(fingerprint.Icon) == 0 && len(fingerprint.JS) == 0 && len(fingerprint.Title) == 0 {
			continue
		}

		// headers判断
		if len(fingerprint.Headers) > 0 {
			var escapedHeader string
			for _, header := range fingerprint.Headers {
				// 判断header是否存在于str_headers中
				if fingerprint.Regexp == "true" {
					escapedHeader = header
				} else {
					escapedHeader = regexp.QuoteMeta(header)
				}
				match, _ := regexp.MatchString(escapedHeader, str_headers)
				if !match {
					haveFingerprint = false
					continue
				}
			}
		}

		// body判断
		if len(fingerprint.Body) > 0 {
			var escapedBody string
			for _, body := range fingerprint.Body {
				if fingerprint.Regexp == "true" {
					escapedBody = body
				} else {
					escapedBody = regexp.QuoteMeta(body)
				}
				match, _ := regexp.MatchString(escapedBody, s.body)
				if !match {
					haveFingerprint = false
					continue
				}
			}
		}

		// icon判断
		if len(fingerprint.Icon) > 0 {
			for _, icon := range fingerprint.Icon {
				if !strings.Contains(s.mmh3icon, icon) && !strings.Contains(s.md5icon, icon) {
					haveFingerprint = false
					continue
				}
			}
		}

		// js判断
		if len(fingerprint.JS) > 0 {
			var escapedJs string
			for _, js := range fingerprint.JS {
				if fingerprint.Regexp == "true" {
					escapedJs = js
				} else {
					escapedJs = regexp.QuoteMeta(js)
				}
				match, _ := regexp.MatchString(escapedJs, str_js)
				if !match {
					haveFingerprint = false
					continue
				}
			}
		}

		// title判断
		if len(fingerprint.Title) > 0 {
			var escapedtitle string
			for _, title := range fingerprint.Title {
				if fingerprint.Regexp == "true" {
					escapedtitle = title
				} else {
					escapedtitle = regexp.QuoteMeta(title)
				}
				match, _ := regexp.MatchString(escapedtitle, s.title)
				if !match {
					haveFingerprint = false
					continue
				}
			}
		}

		if haveFingerprint {
			fps = append(fps, fingerprint.Fp)
			// fmt.Println(fingerprint)
			// break
		}

		// fmt.Println("Fp:", fingerprint.Fp)
		// fmt.Println("Headers:", fingerprint.Headers)
		// fmt.Println("Body:", fingerprint.Body)
		// fmt.Println("Icon:", fingerprint.Icon)
		// fmt.Println("JS:", fingerprint.JS)
	}

	remove_duplicates_fps := removeDuplicates(fps)
	fp := strings.Join(remove_duplicates_fps, "|")
	write := fmt.Sprintf("%s,%s,%s,%s,%s,%s", s.url, fp, s.server, strconv.Itoa(s.statuscode), strconv.Itoa(s.length), s.title)
	fmt.Println(write)
	writeToBuffer(write)
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

	s := resps{urlstring, httpBody, resp.Header, server, resp.StatusCode, len(httpBody), strings.TrimSpace(title), hashicon[0], hashicon[1], jsFiles}

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
