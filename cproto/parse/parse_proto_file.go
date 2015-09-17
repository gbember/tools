// parse_proto_file.go
package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	PROTO_PARSE_OUTPUT_FILE = "proto_parse.go"
)

var (
	protoFieldType1 = map[string]string{"required": "", "optional": "[]"}
	protoFieldType2 = map[string]string{"Bool": "bool", "Int8": "int8", "Uint8": "uint8",
		"Int16": "int16", "Uint16": "uint16", "Int32": "int32", "String": "string",
		"Uint32": "uint32", "Int64": "int64", "Uint64": "uint64"}

	parttenMHead     = regexp.MustCompile(`^[\t\s]*([a-zA-Z0-9_]*)[\t\s]*([a-zA-Z][a-zA-Z_]*)[\t\s]*\[[\t\s]*id=([1-9][0-9]*)[\t\s]*][\t\s]*{[\t\s\r\n]*$`)
	parttenComment   = regexp.MustCompile(`^[\t\s]*\/\/(.*)\r\n$`)
	parttenField     = regexp.MustCompile(`^[\t\s]*([a-zA-Z0-9_]*)[\t\s]*([a-zA-Z0-9_]*)[\t\s]*([a-zA-Z0-9_]*)[\t\s]*=[\t\s0-9]*;(.*)\r\n$`)
	parttenMEnd      = regexp.MustCompile(`^[\t\s]*}[\t\s]*$`)
	parttenBlankLine = regexp.MustCompile(`^[\t\s\r\n]*$`)
)

type ProtoMsg struct {
	Id       uint16
	Name     string
	Comments string
	Fields   []ProtoMsgField
}

type ProtoMsgField struct {
	Type1    string
	Type2    string
	Name     string
	Comments string
	FilePos  string
}

//解析一个协议
func (protoMsg *ProtoMsg) parse(buf *bufio.Reader, fileName string, lineNum int) (int, error) {
	protoMsg.Fields = []ProtoMsgField{}
	msgStr := ""
	state := 0
	field := ProtoMsgField{}
	for {
		line, err := buf.ReadString('\n')

		if err != nil {
			return lineNum, err
		}
		lineNum++
		msgStr += line
		if state == 0 {
			if parttenComment.MatchString(line) {
				ms := parttenComment.FindStringSubmatch(line)
				protoMsg.Comments = protoMsg.Comments + "//" + ms[1] + "\n"
			} else if parttenMHead.MatchString(line) {
				ms := parttenMHead.FindStringSubmatch(line)
				protoMsg.Name = strings.Title(ms[2])
				//protoMsg.Name = firstToUpper(ms[2])
				id1, err := strconv.Atoi(ms[3])
				if err != nil {
					errStr := fmt.Sprintf("parse proto error:(%s:%d) fileid only uint16:%s",
						fileName, lineNum, line)
					return lineNum, errors.New(errStr)
				}
				id := uint16(id1)
				if int(id) != id1 {
					errStr := fmt.Sprintf("parse proto error:(%s:%d) fileid only uint16:%s",
						fileName, lineNum, line)
					return lineNum, errors.New(errStr)
				}
				protoMsg.Id = id
				state = 1
			} else if parttenBlankLine.MatchString(line) {
			} else {
				errStr := fmt.Sprintf("parse proto error:(%s:%d) %s:格式错误",
					fileName, lineNum, line)
				return lineNum, errors.New(errStr)
			}
		} else {
			if parttenMEnd.MatchString(line) {
				return lineNum, nil
			} else if parttenField.MatchString(line) {
				ms := parttenField.FindStringSubmatch(line)
				field.Type1 = ms[1]
				field.Type2 = firstToUpper(ms[2])
				field.Name = firstToUpper(ms[3])
				filePos := fmt.Sprintf("%s:%d", fileName, lineNum)
				field.FilePos = filePos
				if ms[4] != "" {
					field.Comments = field.Comments + "//" + ms[4] + "\n"
				}
				protoMsg.Fields = append(protoMsg.Fields, field)
				field = ProtoMsgField{}
			} else if parttenComment.MatchString(line) {
				ms := parttenComment.FindStringSubmatch(line)
				field.Comments = field.Comments + "//" + ms[1] + "\n"
			} else if parttenBlankLine.MatchString(line) {
			} else {
				errStr := fmt.Sprintf("parse proto error:(%s:%d) %s:格式错误",
					fileName, lineNum, line)
				return lineNum, errors.New(errStr)
			}
		}
	}
}

//协议输出
func (protoMsg *ProtoMsg) write(writer *bufio.Writer) error {
	if protoMsg.Comments != "" {
		_, err := writer.WriteString(protoMsg.Comments)
		if err != nil {
			return err
		}
	}
	_, err := writer.WriteString("type " + protoMsg.Name + " struct {\n")
	if err != nil {
		return err
	}
	for i := 0; i < len(protoMsg.Fields); i++ {
		fp := &protoMsg.Fields[i]
		ftype := protoFieldType1[fp.Type1]
		if t, ok := protoFieldType2[fp.Type2]; ok {
			ftype += t
		} else {
			ftype += fp.Type2
		}
		if fp.Comments != "" {
			_, err = writer.WriteString(fp.Comments)
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString(fp.Name + " " + ftype + "\n")
		if err != nil {
			return err
		}
	}
	_, err = writer.WriteString("}\n")
	if err != nil {
		return err
	}
	//生成读写消息方法
	_, err = writer.WriteString("func (r *" + protoMsg.Name + ")Read(p *Packet) error {\n")
	if err != nil {
		return err
	}
	for i := 0; i < len(protoMsg.Fields); i++ {
		iStr := strconv.Itoa(i)
		fp := &protoMsg.Fields[i]
		if protoFieldType1[fp.Type1] == "[]" {
			if _, ok := protoFieldType2[fp.Type2]; ok {
				writer.WriteString("value" + iStr + ",err := p.readArray" + fp.Type2 + "()\nif err != nil { return err }\nr." + fp.Name + "=value" + iStr + "\n")
			} else {
				writer.WriteString("len" + iStr + ",err := p.readUint16()\nif err != nil {return err}\nr." + fp.Name + "=make([]" + fp.Type2 + ",0,int(len" + iStr + "))\n")
				writer.WriteString("for i:=uint16(0);i<len" + iStr + ";i++ {\nv := " + fp.Type2 + "{}\n")
				writer.WriteString("err = v.Read(p)\nif err != nil {return err}\nr." + fp.Name + " = append(r." + fp.Name + ",v)\n}\n")
			}
		} else {
			if _, ok := protoFieldType2[fp.Type2]; ok {
				writer.WriteString("value" + iStr + ",err := p.read" + fp.Type2 + "()\nif err != nil { return err }\nr." + fp.Name + "=value" + iStr + "\n")
			} else {
				writer.WriteString("r." + fp.Name + "= " + fp.Type2 + "{}\n")
				writer.WriteString("err = r." + fp.Name + ".Read(p)\nif err != nil {return err}\n")
			}
		}
	}
	_, err = writer.WriteString("return nil\n}\n")

	writer.WriteString("func (r *" + protoMsg.Name + ")WriteMsgID(p *Packet){\n")
	writer.WriteString("p.writeUint16(" + strings.ToUpper(protoMsg.Name) + ")\n}\n")

	_, err = writer.WriteString("func (r *" + protoMsg.Name + ")Write(p *Packet) {\n")
	if err != nil {
		return err
	}

	for i := 0; i < len(protoMsg.Fields); i++ {
		fp := &protoMsg.Fields[i]
		if protoFieldType1[fp.Type1] == "[]" {
			if _, ok := protoFieldType2[fp.Type2]; ok {
				_, err = writer.WriteString("p.writeArray" + fp.Type2 + "(r." + fp.Name + ")\n")
				if err != nil {
					return err
				}
			} else {
				writer.WriteString("p.writeUint16(uint16(len(r." + fp.Name + ")))\n")
				writer.WriteString("for i:=0;i<len(r." + fp.Name + ");i++{\n")
				writer.WriteString("r." + fp.Name + "[i].Write(p)\n")
				_, err = writer.WriteString("}\n")
				if err != nil {
					return err
				}
			}
		} else {
			if _, ok := protoFieldType2[fp.Type2]; ok {
				_, err = writer.WriteString("p.write" + fp.Type2 + "(r." + fp.Name + ")\n")
				if err != nil {
					return err
				}
			} else {
				_, err = writer.WriteString("r." + fp.Name + ".Write(p)\n")
				if err != nil {
					return err
				}
			}
		}
	}
	_, err = writer.WriteString("}\n")
	return nil
}

//解析单个proto文件
func parseProtoFile(protoFile string, outdir string) error {
	defer trap_error()()
	file, err := os.Open(protoFile)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	lineNum := 0

	pIDNameMap := make(map[uint16]string)
	pNameMsgMap := make(map[string]*ProtoMsg)

	for {
		protoMsg := &ProtoMsg{}
		lineNum1, err := protoMsg.parse(buf, filepath.Base(file.Name()), lineNum)
		lineNum = lineNum1
		if (err == io.EOF || err == nil) && protoMsg.Name != "" {
			err = add_proto_msg(protoMsg)
			if err != nil {
				return err
			}
			pIDNameMap[protoMsg.Id] = protoMsg.Name
			pNameMsgMap[protoMsg.Name] = protoMsg
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	err = createProtoFile(outdir, protoFile, pIDNameMap, pNameMsgMap)
	return err
}

func add_proto_msg(protoMsg *ProtoMsg) error {
	if m, ok := protoIDNameMap[protoMsg.Id/100]; ok {
		if _, ok := m[protoMsg.Id]; ok {
			errStr := fmt.Sprintf("duplicate message id:%d", protoMsg.Id)
			return errors.New(errStr)
		} else {
			m[protoMsg.Id] = protoMsg.Name
		}
	} else {
		protoIDNameMap[protoMsg.Id/100] = map[uint16]string{protoMsg.Id: protoMsg.Name}
	}
	if _, ok := protoNameMsgMap[protoMsg.Name]; ok {
		errStr := fmt.Sprintf("duplicate message name:%s", protoMsg.Name)
		return errors.New(errStr)
	}
	protoNameMsgMap[protoMsg.Name] = protoMsg
	return nil
}

//创建proto_parse.go文件
func CreateProtoParseFile(dir string) error {
	if len(protoIDNameMap) > 0 {
		fileName := filepath.Join(dir, PROTO_PARSE_OUTPUT_FILE)
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		_, err = writer.WriteString("package " + filepath.Base(dir) + "\n")
		if err != nil {
			return err
		}
		writer.WriteString("import (\"errors\"\n\"fmt\")\n")

		writer.WriteString("type Messager interface{\n")
		writer.WriteString("Read(p *Packet) error \n")
		writer.WriteString("WriteMsgID(p *Packet)\n")
		writer.WriteString("Write(p *Packet)\n")
		writer.WriteString("}\n")

		writer.WriteString("type PMessage struct{\n")
		writer.WriteString("ID uint16\n")
		writer.WriteString("Msg Messager\n")
		writer.WriteString("}\n")

		writer.WriteString("func EncodeProto(msg Messager)[]byte{\n")
		writer.WriteString("p := NewWriter()\n")
		writer.WriteString("msg.WriteMsgID(p)\n")
		writer.WriteString("msg.Write(p)\n")
		writer.WriteString("return p.Data()\n")
		writer.WriteString("}\n")

		writer.WriteString("func DecodeProto(bin []byte) (msgID uint16, msg Messager, err error) {\n" +
			"p := NewReader(bin)\n" +
			"msgID, err = p.readUint16()\n" +
			"if err != nil {return}\n" +
			"mid := msgID / 100\n" +
			"switch mid {\n")

		for mid, m := range protoIDNameMap {
			writer.WriteString("case " + strconv.Itoa(int(mid)) + ":\n")
			writer.WriteString("switch msgID{\n")
			for _, name := range m {
				_, err = writer.WriteString(fmt.Sprintf("case %s:\nv := &%s{}\nerr = v.Read(p)\nif err==nil{msg = v}\n",
					strings.ToUpper(name), name))
				if err != nil {
					return err
				}
			}
			_, err = writer.WriteString("default:\nerr = errors.New(fmt.Sprintf(\"error: invalid msgID %d\", msgID))\n}\n")
		}

		_, err = writer.WriteString("default:\nerr = errors.New(fmt.Sprintf(\"error: invalid msgID %d\", msgID))\n}\nreturn\n}")
		if err != nil {
			return err
		}
		writer.Flush()
	}
	return nil
}

//解析输出proto文件
func createProtoFile(dir string, protoFile string, protoIDNameMap map[uint16]string, protoNameMsgMap map[string]*ProtoMsg) error {
	len := len(protoIDNameMap)
	if len > 0 {
		fileName := strings.Replace(filepath.Base(protoFile), ".", "_", -1) + ".go"
		fileName = filepath.Join(dir, fileName)
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		_, err = writer.WriteString("package " + filepath.Base(dir) + "\nconst (\n")
		if err != nil {
			return err
		}
		idList := make([]int, 0, len)
		for id, _ := range protoIDNameMap {
			idList = append(idList, int(id))
		}
		sort.Ints(idList)
		for _, id := range idList {
			_, err = writer.WriteString(fmt.Sprintf("%s =uint16(%d)\n",
				strings.ToUpper(protoIDNameMap[uint16(id)]), id))
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString(")\n\n\n")
		if err != nil {
			return err
		}

		for _, id := range idList {
			name := protoIDNameMap[uint16(id)]
			protoMsg := protoNameMsgMap[name]
			err = protoMsg.write(writer)
			if err != nil {
				return err
			}
		}
		err = writer.Flush()
		return err
	}
	return nil
}

//检测类型合法性
func check_field_type() error {
	for _, protoMsg := range protoNameMsgMap {
		for _, field := range protoMsg.Fields {
			err := check_field_type_1(&field)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func check_field_type_1(field *ProtoMsgField) error {
	if _, ok := protoFieldType1[field.Type1]; !ok {
		errStr := fmt.Sprintf("type1 error:(%s) %s:only required or optional", field.FilePos, field.Type1)
		return errors.New(errStr)
	}
	if _, ok := protoFieldType2[field.Type2]; !ok {
		if _, ok := protoNameMsgMap[field.Type2]; !ok {
			errStr := fmt.Sprintf("type2 error:(%s) %s:类型错误", field.FilePos, field.Type2)
			return errors.New(errStr)
		}
	}
	return nil
}

//首字母大写
func firstToUpper(str string) string {
	if str[0] >= 'a' && str[0] <= 'z' {
		return string(str[0]-byte(32)) + str[1:]
	}
	return str
}
