//excel配置转化json或其他
package main

import (
	"flag"
	"log"
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
	startTime := time.Now()
	err := parse.ParseDir(*inputDir, *outDir, *et, *sc)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(time.Since(startTime))
}
