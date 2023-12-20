package src

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

func Check(s Resps) {
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
	WriteToBuffer(write)
}
