package src

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reaper/src/queue"
	"strings"
	"sync"
)

var (
	Buffer      = make(chan string, 2048) // 全局缓冲区
	BufferMutex sync.Mutex                // 用于保护缓冲区的互斥锁
	IsFlush     bool                      // 手动执行flush
	Fps         []Fingerprint             // 指纹数组
	Thread      int                       // 最大线并发数
	jsJump      int                       // js最大跳转次数
	List        string                    // 扫描的目标URL/主机列表的文件路径（一行一个）
	Output      string                    // 输出结果（csv格式）
	Wg          sync.WaitGroup
	UrlQueue    *queue.Queue
)

func read_json() []Fingerprint {
	// 读取 JSON 文件
	filePath := "reaper.json"

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("reaper.json 文件不存在")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening reaper.json: ", err)
		return nil
	}
	defer file.Close()

	var data Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	fps := data.Fingerprints
	return fps
}

func read_url() {
	// filePath := "url.txt"
	UrlQueue = queue.NewQueue()
	// maxThreads := Thread
	// semaphore := make(chan struct{}, maxThreads)

	// 打开文件
	file, err := os.Open(List)
	if err != nil {
		fmt.Println("Error read file: ", err)
		return
	}
	defer file.Close()

	// 创建 Scanner 对象
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容并去除换行符
	for scanner.Scan() {
		url := scanner.Text()
		url = strings.TrimSpace(url)       // 去除每行两端的空白字符
		url = strings.TrimSuffix(url, "/") // 去掉 "/"

		// 判断是否以 "http" 开头，如果不是则添加 "http"
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}

		// urlQueue <- url // 将url加入队列
		UrlQueue.Push(url)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanner: ", err)
	}

	// Wg.Wait()
}

func sendUrl() {
	defer Wg.Done()
	for UrlQueue.Len() != 0 {
		dataface := UrlQueue.Pop()
		url, _ := dataface.(string)
		// fmt.Println(url)
		Http(url)
	}
}

func scan() {
	// fmt.Printf("thread: %d", Thread)
	for i := 0; i <= Thread; i++ {
		Wg.Add(1)
		go sendUrl()
	}
	Wg.Wait()
}

func Begin() {
	flag.IntVar(&Thread, "t", 100, "并发数")
	flag.StringVar(&List, "l", "", "扫描的目标URL/主机列表的文件路径(一行一个)")
	flag.StringVar(&Output, "o", "result.csv", "输出结果(csv格式)")
	flag.Parse()

	// 检查必填参数是否已设置
	if List == "" {
		fmt.Println("必填参数 'List' 未设置")
		flag.PrintDefaults()
		os.Exit(1)
	}

	read_url()
	Fps = read_json()
	// 创建或覆盖文件
	file, err := os.Create(Output)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 写入字符串到文件
	content := "url,fingerprint,server,status,length,title\n"
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	scan()

	defer func() {
		WriteFile()
	}()
}
