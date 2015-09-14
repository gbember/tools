package parse

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ParseDir(dir string, outDir string, exportType string, soc string) error {
	defer trap_error()()
	if !checkExportType(exportType) {
		log.Fatalln(fmt.Sprintf("validate export type:%s", exportType))
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	fiList, err := d.Readdir(99999)
	if err != nil {
		return nil
	}
	errList := make([]error, 0, 10)
	for i := 0; i < len(fiList); i++ {
		if !fiList[i].IsDir() {
			fileName := fiList[i].Name()
			if filepath.Ext(fileName) == ".xlsx" {
				err := ParseFile(filepath.Join(dir, fileName), outDir, exportType, soc)
				if err != nil {
					str := fmt.Sprintf("导出配置文件[%v]错误:%v", fileName, err)
					fmt.Println(str)
					errList = append(errList, errors.New(str))
				}
			}
		}
	}
	if len(errList) == 0 {
		return nil
	}
	errStr := ""
	for i := 0; i < len(errList); i++ {
		s := errList[i].Error()
		errStr = errStr + s + "\n"
	}
	return errors.New(errStr)
}
