package main

import (
	"flag"
	"fmt"
)

var action = flag.String("action", "", "zip: 压缩 unzip: 解压")
var srcFileName = flag.String("src", "", "输入文件")
var disFileName = flag.String("dis", "", "输出文件")

func main() {
	flag.Parse()
	fmt.Println(*srcFileName, *action, *disFileName)
	if len(*srcFileName) == 0 || len(*action) == 0 || len(*disFileName) == 0 {
		flag.PrintDefaults()
		return
	}
	// str := "2233366666688888888999999999"
	// bts := []byte(str)
	// huffuman.Encode(bts)

	// 读取文件

	// 压缩文件
}
