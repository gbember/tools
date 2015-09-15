自定义protobuf文件导出工具

protobuf文件格式:
	message MESSAGE_NAME[id=MESSAGE_ID]{
		TYPE1	TYPE2	FIELD_NAME	=1;
	}
	
	MESSAGE_NAME:协议名字
	MESSAGE_ID:协议ID号   当前版本只支持协议号2字节 从二进制解析成协议的时候为了速度根据MESSAGE_ID/100来switch处理
	TYPE1: required原生 optional数组
	TYPE2: bool,int8,uint8,int16,uint16,int32,string,uint32,int64,uint64或者MESSAGE_NAME
	FIELD_NAME:字符串  符合struct里的字段名定义
	
命令格式
	cproto -inputDir InputDir -outDir OutDir
		InputDir:*.proto文件所在目录 默认当前目录
		OutDir:*_proto.go输出文件所在目录 默认当前目录