package src

import (
	"bufio"
	"fmt"
	"os"
)

func WriteToBuffer(data string) {
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
