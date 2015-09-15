//导出自定义protobuf
package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"time"
	"util/tools/cproto/parse"
)

var (
	//*.proto协议文件夹
	inputDir *string = flag.String("inputDir", "./", "input dir")
	//解析后生成的*_proto.go的文件夹
	outDir *string = flag.String("outDir", "./", "out dir")
)

func main() {
	flag.Parse()
	defer wait_exit()
	dir, err := filepath.Abs(*outDir)
	if err != nil {
		fmt.Println("outDir 文件夹错误:", err)
		return
	}
	startTime := time.Now()
	err = parse.ParseProtoDir(*inputDir, dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = parse.CreateProtoPacketFile(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = parse.CreateProtoParseFile(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = parse.BuildProto(*outDir)
	if err != nil {
		errStr := fmt.Sprintf("build proto错误:%v", err)
		fmt.Println(errors.New(errStr))
		return
	}
	fmt.Println("解析proto时间:", time.Since(startTime))
}

func wait_exit() {
	fmt.Print("按任意键关闭....")
	var k int
	fmt.Scan(&k)
}
