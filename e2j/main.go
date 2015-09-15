//excel配置转化json或其他
package main

import (
	"flag"
	"fmt"
	"time"
	"util/tools/e2j/parse"
)

var (
	inputDir *string = flag.String("inputDir", "./", "input dir")
	outDir   *string = flag.String("outDir", "./", "out dir")
	et       *string = flag.String("et", "json", "export type")
	sc       *string = flag.String("sc", "server", "server or client")
)

func main() {
	flag.Parse()
	defer wait_exit()
	startTime := time.Now()
	err := parse.ParseDir(*inputDir, *outDir, *et, *sc)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("耗时:", time.Since(startTime))
}

func wait_exit() {
	fmt.Print("按任意键关闭....")
	var k int
	fmt.Scan(&k)
}
