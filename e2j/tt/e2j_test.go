package test

import (
	"testing"
	"time"

	"github.com/gbember/tools/e2j/parse"
)

//func TestParseFile(t *testing.T) {
//	startTime := time.Now()
//	inputDir := "./excel"
//	outDir := "./config"
//	err := parse.ParseFile(filepath.Join(inputDir, "测试_data_test.xlsx"), outDir, "json", "server")
//	if err != nil {
//		t.Log(err)
//		return
//	}
//	t.Log("测试时间:", time.Since(startTime))
//}

func TestParseDir(t *testing.T) {
	start := time.Now()
	inputDir := "./excel"
	outDir := "./config"
	err := parse.ParseDir(inputDir, outDir, "json", "server")
	n := time.Now().Sub(start).Nanoseconds()
	t.Log(n, err)
}
