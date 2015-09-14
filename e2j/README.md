特定格式excel文件导出为特定(json)配置的工具

	默认导出excel中第一个sheet表

	第一行为名称   没有名称默认不导出  名称不能为空
	第二行为server导出信息
	第三行为client导出信息
	其余行可以为空行

导出信息结构说明
	|:不导出
	n_开头:导出为整数
	nl_开头:导出为整数列表
	s_开头:导出为字符串
	sl_开头:导出为字符串列表
	t_开头:导出为时间点(秒) 格式 hh:mm
	tl_开头:导出为时间点列表
	不为空:导出为字符串
	
	
	
命令格式:e2j -inputDir InputDir -outDir OutDir -et ET -sc SC
		InputDir:要导出的excel文件目录
		OutDir:导出*.config的目标目录
		ET:导出类型  先只支持json
		SC:导出server或client
