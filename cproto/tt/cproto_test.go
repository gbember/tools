// cproto_test.go
package tt

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gbember/tools/cproto/parse"
)

func TestCProto(t *testing.T) {
	inputDir := "./proto_file"
	outDir := "./proto"
	startTime := time.Now()
	err := parse.CreateProtoPacketFile(outDir)
	if err != nil {
		t.Log(err)
		return
	}
	err = parse.ParseProtoDir(inputDir, outDir)
	if err != nil {
		t.Log(err)
		return
	}
	err = parse.CreateProtoParseFile(outDir)
	if err != nil {
		t.Log(err)
		return
	}
	err = parse.BuildProto(outDir)
	if err != nil {
		errStr := fmt.Sprintf("build proto错误:%v", err)
		t.Log(errors.New(errStr))
		return
	}
	t.Log("解析proto时间:", time.Since(startTime))
}
