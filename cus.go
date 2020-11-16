package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("http://www.bilibili.com")
	if err != nil {
		fmt.Println("get err=", err)
		return
	}

	var content string
	buf := make([]byte, 4*1024)
	for {
		n, _ := resp.Body.Read(buf)
		defer resp.Body.Close()
		if n == 0 {
			fmt.Println("读取完成！")
			break
		}
		content += string(buf[:n])

	}
	fmt.Println(content)

}
