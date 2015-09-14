package parse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	//
	FIELD_TYPE_NONE string = ""
	//整数
	FIELD_TYPE_INTEGER string = "n_"
	//整数数组
	FIELD_TYPE_INTEGER_LIST string = "nl_"
	//字符串
	FIELD_TYPE_STRING string = "s_"
	//字符串数组
	FIELD_TYPE_STRING_LIST string = "sl_"
	//时间
	FIELD_TYPE_TIME string = "t_"
	//时间列表
	FIELD_TYPE_TIME_LIST string = "tl_"
)

var (
	//导出类型
	exportTypeMap map[string]bool = map[string]bool{"json": true}
)

type configRow struct {
	//第几行
	line int
	//值列表
	cellValues  []string
	configCells []*configCell
}

type configCell struct {
	index int
	//类型
	fType string
	//名字
	fName string
	name  string
}

func checkExportType(exportType string) bool {
	_, ok := exportTypeMap[exportType]
	return ok
}

//每一行导出
func (cr *configRow) format(buff *bytes.Buffer) error {
	s, err := cr.configCells[0].format(cr.cellValues)
	if err != nil {
		return errors.New(fmt.Sprintf("(cell:%s line:%d) %s", cr.configCells[0].name, cr.line, err.Error()))
	}
	buff.WriteString("{")
	buff.WriteString(s)
	for i := 1; i < len(cr.configCells); i++ {
		s, err = cr.configCells[i].format(cr.cellValues)
		if err != nil {
			return errors.New(fmt.Sprintf("(cell:%s line:%d) %s", cr.configCells[i].name, cr.line, err.Error()))
		}
		buff.WriteString(",")
		buff.WriteString(s)
	}
	buff.WriteString("}\n")
	return nil
}

//每一列导出
func (cc *configCell) format(cellValues []string) (s string, err error) {
	value := cellValues[cc.index]
	switch cc.fType {
	case FIELD_TYPE_NONE:
		s = value
	case FIELD_TYPE_STRING:
		s = "\"" + value + "\""
	case FIELD_TYPE_STRING_LIST:
		s = toStringLists(value)
	case FIELD_TYPE_INTEGER:
		_, err = strconv.Atoi(value)
		if err != nil {
			err = errors.New(fmt.Sprintf("整数类型错误:%s", value))
		} else {
			s = value
		}
	case FIELD_TYPE_INTEGER_LIST:
		s, err = toIntegerLists(value)
	case FIELD_TYPE_TIME:
		s, err = toSecTime(value)
	case FIELD_TYPE_TIME_LIST:
		s, err = toSecTimeLists(value)
	}
	if err == nil {
		s = cc.fName + s
	}
	return
}

func ParseFile(fileName string, outDir string, exportType string, soc string) error {
	defer trap_error()()
	ofn, err := getOutFileName(fileName)
	if err != nil {
		return err
	}
	outFileName := filepath.Join(outDir, ofn)
	f, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}
	//默认excel的第一个sheet为配置
	sheet := file.Sheets[0]
	rows := sheet.Rows
	rowNum := len(rows)
	if rowNum < 3 {
		err = errors.New("validate config excel")
		return err
	}

	firseRowCells := rows[0].Cells
	secondRowCells := rows[1].Cells
	threeRowCells := rows[2].Cells
	//
	titles := make(map[int]string)
	secondConfigCells := make([]*configCell, 0, len(firseRowCells))
	threeConfigCells := make([]*configCell, 0, len(firseRowCells))
	for i := 0; i < len(firseRowCells); i++ {
		if title := strings.TrimSpace(firseRowCells[i].Value); title != "" {
			titles[i] = title
			secondConfigCell, err := make_config_cell(secondRowCells[i].Value, i)
			if err != nil {
				return err
			}
			if secondConfigCell != nil {
				secondConfigCells = append(secondConfigCells, secondConfigCell)
			}
			threeConfigCell, err := make_config_cell(threeRowCells[i].Value, i)
			if err != nil {
				return err
			}
			if threeConfigCell != nil {
				threeConfigCells = append(threeConfigCells, threeConfigCell)
			}
		}
	}
	if exportType == "json" {
		if soc == "server" {
			err = parse2json(f, rows[3:], secondConfigCells, titles, 3)
		} else {
			err = parse2json(f, rows[3:], threeConfigCells, titles, 3)
		}
	}
	return err
}

//导出json
func parse2json(writer io.Writer, rows []*xlsx.Row, configCells []*configCell, titles map[int]string, line int) error {
	if len(configCells) > 0 {
		cr := &configRow{line: line + 1, configCells: configCells}
		bs := make([]byte, 1024)
		buff := bytes.NewBuffer(bs)
		buff.Reset()
		for i := 0; i < len(rows); i++ {
			cellValues, isEmptyRow := make_row_contents(rows[i], titles)
			if isEmptyRow {
				line++
				continue
			}
			cr.cellValues = cellValues
			err := cr.format(buff)
			if err != nil {
				return err
			}
			_, err = buff.WriteTo(writer)
			if err != nil {
				return err
			}
			cr.line++
		}
	}
	return nil
}

func make_row_contents(row *xlsx.Row, titles map[int]string) ([]string, bool) {
	cellValues := make([]string, 0, len(row.Cells))
	b := true
	for i := 0; i < len(row.Cells); i++ {
		if _, ok := titles[i]; ok {
			v := strings.TrimSpace(row.Cells[i].Value)
			cellValues = append(cellValues, v)
			if v != "" {
				b = false
			}
		} else {
			cellValues = append(cellValues, "")
		}
	}
	return cellValues, b
}

func make_config_cell(name string, index int) (*configCell, error) {
	name = strings.TrimSpace(name)
	if name != "|" {
		names := strings.SplitAfterN(name, "_", 2)
		if len(names) == 1 {
			return &configCell{fType: FIELD_TYPE_NONE, fName: "\"" + names[0] + "\":", name: names[0], index: index}, nil
		} else if len(names) == 2 {
			return &configCell{fType: names[0], fName: "\"" + names[1] + "\":", name: names[1], index: index}, nil
		} else {
			err := errors.New(fmt.Sprintf("错误的cell_type_name :%s", name))
			return nil, err
		}
	} else {
		return nil, nil
	}
}

func toStringLists(fValue string) string {
	if fValue == "" {
		return "[]"
	} else {
		flist := strings.Split(fValue, ";")
		l := len(flist)
		if l == 1 {
			return toStringList(fValue)
		} else {
			str := "["
			for i := 0; i < len(flist); i++ {
				s := strings.TrimSpace(flist[i])
				if s != "" {
					s = toStringList(s)
					if i != 0 {
						str = str + "," + s
					} else {
						str = str + s
					}
				}
			}
			str = str + "]"
			return str
		}
	}
}

func toStringList(fValue string) string {
	flist := strings.Split(fValue, ",")
	str := "["
	for i := 0; i < len(flist); i++ {
		if i != 0 {
			str = str + ",\"" + flist[i] + "\""
		} else {
			str = str + "\"" + flist[i] + "\""
		}
	}
	str = str + "]"
	return str
}
func toIntegerLists(fValue string) (string, error) {
	if fValue == "" {
		return "[]", nil
	} else {
		flist := strings.Split(fValue, ";")
		l := len(flist)
		if l == 1 {
			return toIntegerList(fValue)
		} else {
			str := "["
			for i := 0; i < len(flist); i++ {
				s := strings.TrimSpace(flist[i])
				if s != "" {
					s, err := toIntegerList(s)
					if err != nil {
						return "", err
					}
					if i != 0 {
						str = str + "," + s
					} else {
						str = str + s
					}
				}
			}
			str = str + "]"
			return str, nil
		}
	}
}

func toIntegerList(fValue string) (string, error) {
	flist := strings.Split(fValue, ",")
	str := "["
	for i := 0; i < len(flist); i++ {
		s := strings.TrimSpace(flist[i])
		_, err := strconv.Atoi(s)
		if err != nil {
			return "", errors.New(fmt.Sprintf("整数列表类型错误:%s", fValue))
		}
		if i != 0 {
			str = str + "," + s
		} else {
			str = str + s
		}
	}
	str = str + "]"
	return str, nil
}

func toSecTimeLists(fValue string) (string, error) {
	if fValue == "" {
		return "[]", nil
	} else {
		flist := strings.Split(fValue, ";")
		l := len(flist)
		if l == 1 {
			s, err := toSecTimeList(fValue)
			if err != nil {
				return "", errors.New(fmt.Sprintf("时间列表类型错误:%s", fValue))
			}
			return s, err
		} else {
			str := "["
			for i := 0; i < len(flist); i++ {
				s := strings.TrimSpace(flist[i])
				if s != "" {
					s, err := toSecTimeList(s)
					if err != nil {
						return "", errors.New(fmt.Sprintf("时间列表类型错误:%s", fValue))
					}
					if i != 0 {
						str = str + "," + s
					} else {
						str = str + s
					}
				}
			}
			str = str + "]"
			return str, nil
		}
	}
}

func toSecTimeList(fValue string) (string, error) {
	flist := strings.Split(fValue, ",")
	str := "["
	for i := 0; i < len(flist); i++ {
		s := strings.TrimSpace(flist[i])
		st, err := toSecTime(s)
		if err != nil {
			return "", errors.New(fmt.Sprintf("时间列表类型错误:%s", fValue))
		}
		if i != 0 {
			str = str + "," + st
		} else {
			str = str + st
		}
	}
	str = str + "]"
	return str, nil
}

func toSecTime(fValue string) (string, error) {
	flist := strings.Split(fValue, ":")
	if len(flist) != 2 {
		return "", errors.New(fmt.Sprintf("时间类型错误:%s", fValue))
	}
	h, err := strconv.Atoi(flist[0])
	if err != nil {
		return "", errors.New(fmt.Sprintf("时间类型错误:%s", fValue))
	}
	m, err := strconv.Atoi(flist[1])
	if err != nil {
		return "", errors.New(fmt.Sprintf("时间类型错误:%s", fValue))
	}
	return strconv.Itoa(h*60*60 + m*60), nil
}

//根据excel文件名字名字得到输出的文件名字
func getOutFileName(fileName string) (outFileName string, err error) {
	fn := filepath.Base(fileName)
	fns := strings.SplitN(fn, "_", 2)
	if len(fns) == 2 {
		outFileName = fns[1]
		index := strings.LastIndex(outFileName, ".")
		if index != -1 {
			outFileName = outFileName[:index] + ".config"
		}
	} else {
		err = errors.New("配置文件名字格式错误")
	}
	return
}

//func trans2field(name string) (*field, error) {
//	ys := name
//	name = strings.TrimSpace(name)
//	if name == "|" {
//		return nil, nil
//	}
//	index := strings.Index(name, "_")
//	if index == -1 {
//		return &field{name: name, t: FIELD_TYPE_NONE, ys: ys}, nil
//	}
//	t := name[:index+1]
//	name = name[index+1:]
//	if name == "" ||
//		(t != FIELD_TYPE_NONE &&
//			t != FIELD_TYPE_INTEGER &&
//			t != FIELD_TYPE_INTEGER_LIST &&
//			t != FIELD_TYPE_STRING &&
//			t != FIELD_TYPE_STRING_LIST &&
//			t != FIELD_TYPE_TIME &&
//			t != FIELD_TYPE_TIME_LIST) {
//		return nil, errors.New(fmt.Sprintf("field name type error:%v", ys))
//	}
//	return &field{name: "\"" + name + "\":", t: t, ys: ys}, nil
//}
