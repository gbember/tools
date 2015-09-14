// parse_proto_dir.go
package parse

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	protoIDNameMap  = make(map[uint16]map[uint16]string)
	protoNameMsgMap = make(map[string]*ProtoMsg)
)

func ParseProtoDir(dir string, outdir string) error {
	defer trap_error()()
	file, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		subFiList, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for i := 0; i < len(subFiList); i++ {
			if !subFiList[i].IsDir() && strings.HasSuffix(subFiList[i].Name(), ".proto") {
				protoFile := filepath.Join(dir, subFiList[i].Name())
				err := parseProtoFile(protoFile, outdir)
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}
		}

		err = check_field_type()

		return err
	} else {
		errStr := fmt.Sprintf("filepath not dir:%s", dir)
		return errors.New(errStr)
	}
}

func BuildProto(dir string) error {
	cmd := exec.Command("gofmt", "-w", dir)
	err := cmd.Run()
	if err != nil {
		errStr := fmt.Sprintf("格式化协议输出文件报错:%v", err)
		return errors.New(errStr)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	cmd = exec.Command("go", "build")
	err = cmd.Run()
	os.Chdir(pwd)
	return err
}
