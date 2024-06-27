package main

import (
	"flag"
	"fmt"
	"structure_algorithm/huffuman"
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

	switch *action {
	case "zip":
		huffuman.Zip(*srcFileName, *disFileName)
	case "unzip":
		huffuman.Unzip(*srcFileName, *disFileName)
	default:
		flag.PrintDefaults()
	}
}

// 压缩
// go run main.go --action=zip --src=./temp/source.txt --dis=./temp/dis.hfm
// 解压
// go run main.go --action=unzip --src=./temp/dis.hfm --dis=./temp/dis.txt
