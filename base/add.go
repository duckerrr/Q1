package base

import (
	"Protection/protectDevices"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

func BusAdd(xlsx *excelize.File) []protectDevices.Bus {

	// 1.获取表格
	busRows := xlsx.GetRows("Bus")
	busRows = busRows[1:]

	// 2.解析数据, 添加bus
	buses := make([]protectDevices.Bus, 0)
	for i, _ := range busRows {
		row := busRows[i]
		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		name := row[0]
		status := true
		if row[2] == "0" {
			continue
		}

		bus := protectDevices.Bus{}
		bus.Name = name
		bus.Status = status
		bus.Voltage = 1			//母线电压初始化为1
		bus.IfFault = false		//初始化为无故障
		buses = append(buses, bus)
	}
	return buses
}

func LineAdd(xlsx *excelize.File) []protectDevices.Line {

	// 1.获取表格
	lineRows := xlsx.GetRows("Line")
	lineRows = lineRows[1:]

	// 2.解析数据, 添加line
	lines := make([]protectDevices.Line, 0)
	for i, _ := range lineRows {
		row := lineRows[i]
		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}

		line := protectDevices.Line{}
		line.Name, line.F_bus, line.T_bus = row[0], row[1], row[2]
		line.Status = true
		if row[3] == "0" {
			line.Status = false
			continue
		}
		line.IfFault = false	//初始化为无故障
		lines = append(lines, line)
	}
	return lines
}

func TransAdd(xlsx *excelize.File) ([]protectDevices.Trans, []protectDevices.Bus) {

	// 1.获取表格
	transRows := xlsx.GetRows("Trans")
	transRows = transRows[1:]

	// 2.解析数据, 添加trans
	transes := make([]protectDevices.Trans, 0)
	for i, _ := range transRows {
		row := transRows[i]
		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		trans := protectDevices.Trans{}
		name := row[0]
		busStr := strings.Split(row[1], "/")
		fromBus, toBus := busStr[0], busStr[1]

		status := true
		if row[2] == "0" {
			status = false
		}

		trans.Name = name
		trans.Type_trans = 2
		trans.IfFault = false	//初始化为无故障
		trans.Bus = append(trans.Bus, fromBus)
		trans.Bus = append(trans.Bus, toBus)
		trans.Status = status
		trans.IfExternalPro = false
		trans.IfInternalPro = false
		trans.Gas_density = 0		//瓦斯密度初始化为0

		transes = append(transes, trans)
	}

	// 1.获取表格
	transRows = xlsx.GetRows("Trans3")
	transRows = transRows[1:]
	buses := []protectDevices.Bus{}

	// 2.解析数据, 添加trans
	for i, _ := range transRows {
		row := transRows[i]
		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		trans := protectDevices.Trans{}
		name := row[0]
		trans_buses := strings.Split(row[1], "/")

		status := true
		if row[2] == "0" {
			status = false
		}

		trans.Name = name
		trans.Type_trans = 3
		trans.IfFault = false	//初始化为无故障
		trans.Bus = append(trans.Bus, trans_buses[0])
		trans.Bus = append(trans.Bus, trans_buses[1])
		trans.Bus = append(trans.Bus, trans_buses[2])
		trans.Status = status
		trans.IfExternalPro = false
		trans.IfInternalPro = false
		trans.Gas_density = 0		//瓦斯密度初始化为0
		transes = append(transes, trans)

		bus := protectDevices.Bus{}
		bus.Name = name
		bus.Status = status
		bus.Voltage = 1			//母线电压初始化为1
		bus.IfFault = false		//初始化为无故障
		buses = append(buses, bus)

	}
	return transes, buses
}

func ShuntAdd(xlsx *excelize.File) []protectDevices.Shunt {
	// 1.获取表格
	shuntRows := xlsx.GetRows("Shunt")
	shuntRows = shuntRows[1:]

	// 2.解析数据, 添加shunt
	shunts := make([]protectDevices.Shunt, 0)
	for i, _ := range shuntRows {
		row := shuntRows[i]
		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		shunt := protectDevices.Shunt{}
		name := row[0]
		busName := row[1]

		shunt.Name = name
		shunt.Bus = busName
		shunt.Status = true
		shunt.IfFault = false	//初始化为无故障
		shunts = append(shunts, shunt)
	}
	return shunts
}

func SynAdd(xlsx *excelize.File) []protectDevices.Syn {

	// 1.获取表格
	syn6Rows := xlsx.GetRows("Syn6")
	syn6Rows = syn6Rows[1:]

	// 2.解析数据, 添加syn
	syns := make([]protectDevices.Syn, 0)
	for i, _ := range syn6Rows {
		row := syn6Rows[i]

		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		syn := protectDevices.Syn{}

		syn.Name, syn.Bus = row[0], row[1]
		syn.Status = true
		syn.IfFault = false	//初始化为无故障
		syns = append(syns, syn)
	}
	return syns
}

//初始化断路器的名称、所在线路、母线
func BreakerAdd(xlsx *excelize.File) ([]protectDevices.Breaker){
	breakerRows := xlsx.GetRows("Breaker")
	breakerRows = breakerRows[1:]

	breakers := []protectDevices.Breaker{}
	for i, _ := range breakerRows {
		row := breakerRows[i]

		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		breaker := protectDevices.Breaker{}

		breaker.Name, breaker.Bus_close = row[0], row[1]
		breaker.Line = row[2]
		breaker.IfLock = false	//初始化为未闭锁
		if row[3] == "1"{
			breaker.Status = true
		}else{
			breaker.Status = false
		}
		breaker.IfReclosing = false		//先初始化为false
		breakers = append(breakers, breaker)
	}
	return breakers

}


//setting中给断路器赋值
func BreakerAddSetting(device *protectDevices.Device, buses, breaker []string, Protection,settingType string, delay,settingVal []float64, protectionType string){

	for i, name := range breaker {
		//如果该breaker已录入，则直接找索引赋值
		br := &device.Breakers[device.BreakerMap_name[name]]

		br.Type_protection = append(br.Type_protection, protectionType)
		br.Protection = append(br.Protection, Protection)
		br.Delay = append(br.Delay, delay[i])
		br.SettingType = append(br.SettingType, settingType)
		br.SettingVal = append(br.SettingVal, settingVal[i])

		// 如果buses是有两个bus，则录入breaker.Line
		if protectionType == "line"{
			if i == 0 {
				br.Bus_close = buses[0]
				br.Bus_far = buses[1]
			} else {
				br.Bus_close = buses[1]
				br.Bus_far = buses[0]
			}
			br.Line = device.Lines[device.LineMap[buses[0] + "_" + buses[1]]].Name
		} else if protectionType == "trans" {
			br.Bus_close = buses[i]
		} else {
			// "bus"/"syn"/"shunt"
			br.Bus_close = buses[0]
		}

	}
}

func SettingAdd(xlsx *excelize.File, device *protectDevices.Device){

	settingRows := xlsx.GetRows("Setting")
	settingRows = settingRows[1:]


	for i, _ := range settingRows {
		row := settingRows[i]

		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		// 还没有进入device.breakers的断路器集合
		breakers_new := []protectDevices.Breaker{}

		//分类讨论
		// 1.line
		if row[0]=="line" {
			Protection := row[3]
			breaker, delay := splitInfo(row[4])
			_, settingVal := splitInfo(row[5])
			settingType := row[6]

			//从已有的device.Lines中找，插入要的数值
			i := device.LineMap[row[1]+"_"+row[2]]
			device.Lines[i].Protection = append(device.Lines[i].Protection, Protection)
			device.Lines[i].Breaker_name = append(device.Lines[i].Breaker_name, breaker)
			device.Lines[i].Delay = append(device.Lines[i].Delay, delay)
			device.Lines[i].SettingVal = append(device.Lines[i].SettingVal, settingVal)
			device.Lines[i].SettingType = append(device.Lines[i].SettingType, settingType)

			// device.Breaker同步赋值
			BreakerAddSetting(device,[]string{row[1],row[2]},breaker,Protection,settingType,delay,settingVal,"line")
		}

		// 2. bus
		if row[0]=="bus" {
			Protection := row[2]
			breaker, delay := splitInfo(row[3])
			_, settingVal := splitInfo(row[4])
			settingType := row[5]

			//从已有的device.Buses中找，插入要的数值
			i := device.BusMap[row[1]]
			device.Buses[i].Protection = append(device.Buses[i].Protection, Protection)
			device.Buses[i].Breaker_name = append(device.Buses[i].Breaker_name, breaker)
			device.Buses[i].Delay = append(device.Buses[i].Delay, delay)
			device.Buses[i].SettingVal = append(device.Buses[i].SettingVal, settingVal)
			device.Buses[i].SettingType = append(device.Buses[i].SettingType, settingType)

			// device.Breaker同步赋值
			BreakerAddSetting(device,[]string{row[1]},breaker,Protection,settingType,delay,settingVal,"bus")
		}

		// 3. shunt
		if row[0]=="shunt"{
			Protection := row[2]
			breaker, delay := splitInfo(row[3])
			_, settingVal := splitInfo(row[4])
			settingType := row[5]

			//从已有的device.Shunts中找，插入要的数值
			i := device.ShuntMap[row[1]]
			device.Shunts[i].Protection = append(device.Shunts[i].Protection,Protection)
			device.Shunts[i].Breaker_name = append(device.Shunts[i].Breaker_name, breaker)
			device.Shunts[i].Delay = append(device.Shunts[i].Delay, delay)
			device.Shunts[i].SettingVal = append(device.Shunts[i].SettingVal, settingVal)
			device.Shunts[i].SettingType = append(device.Shunts[i].SettingType, settingType)

			// device.Breaker同步赋值
			BreakerAddSetting(device,[]string{device.Shunts[i].Bus},breaker,Protection,settingType,delay,settingVal,"shunt")
		}

		// 4. syn
		if row[0]=="syn"{
			Protection := row[2]
			breaker, delay := splitInfo(row[3])
			_, settingVal := splitInfo(row[4])
			settingType := row[5]

			//从已有的device.Syns中找，插入要的数值
			i := device.SynMap[row[1]]
			device.Syns[i].Protection = append(device.Syns[i].Protection,Protection)
			device.Syns[i].Breaker_name = append(device.Syns[i].Breaker_name, breaker)
			device.Syns[i].Delay = append(device.Syns[i].Delay, delay)
			device.Syns[i].SettingVal = append(device.Syns[i].SettingVal, settingVal)
			device.Syns[i].SettingType = append(device.Syns[i].SettingType, settingType)

			// device.Breaker同步赋值
			BreakerAddSetting(device,[]string{device.Syns[i].Bus},breaker,Protection,settingType,delay,settingVal,"syn")
		}

		// 5. trans
		if row[0]=="trans"{
			Protection := row[2]
			breaker, delay := splitInfo(row[3])
			_, settingVal := splitInfo(row[4])
			settingType := row[5]

			//从已有的device.trans中找，插入要的数值
			i := device.TransMap[row[1]]
			device.Transes[i].Protection = append(device.Transes[i].Protection,Protection)
			device.Transes[i].Breaker_name = append(device.Transes[i].Breaker_name, breaker)
			device.Transes[i].Delay = append(device.Transes[i].Delay, delay)
			device.Transes[i].SettingVal = append(device.Transes[i].SettingVal, settingVal)
			device.Transes[i].SettingType = append(device.Transes[i].SettingType, settingType)

			isInner, _ := StringssContain(Trans_innerProtect, Protection)
			isOuter, _ := StringssContain(Trans_outerProtect, Protection)

			if isInner{
				device.Transes[i].IfExternalPro = true
			} else if isOuter{
				device.Transes[i].IfInternalPro = true
			}else{
				panic("setting.txt中的变压器保护类型不符合要求！（gas/transdifferential）")
			}

			// device.Breaker同步赋值
			BreakerAddSetting(device,device.Transes[i].Bus,breaker,Protection,settingType,delay,settingVal,"trans")
		}
		device.Breakers = append(device.Breakers,breakers_new...)
	}

	// 为母线的断路器数量属性赋值
	for _, breaker := range device.Breakers{
		device.Buses[device.BusMap[breaker.Bus_close]].Num_breaker +=1
	}

}

func ReclosingAdd(xlsx *excelize.File, device *protectDevices.Device) {

	reclosingRows := xlsx.GetRows("Reclosing")
	reclosingRows = reclosingRows[1:]

	for i, _ := range reclosingRows {
		row := reclosingRows[i]

		// 忽略"#"开头的数据
		if row[0][0] == '#' {
			continue
		}
		breaker_name := row[0]
		i := device.BreakerMap_name[breaker_name]
		device.Breakers[i].IfReclosing = true
		device.Breakers[i].Delay_recl, _ = strconv.ParseFloat(row[1],64)
		device.Breakers[i].Lock_recl, _ = strconv.ParseFloat(row[2],64)
	}
}


func splitInfo(s string) ([]string, []float64){
	breaker := []string{}
	value := []string{}

	st := strings.Trim(s,"[")
	st = strings.Trim(st,"]")

	sts := strings.Split(st,"/")
	for _,val := range sts{
		temp_sp := strings.Split(val,":")
		breaker = append(breaker, temp_sp[0])
		value = append(value, temp_sp[1])
	}

	return breaker,ss2fs(value)
}

//将[]string转为[]float64
func ss2fs(ss []string) []float64{
	fs := []float64{}

	for _,val := range ss{
		if n, err := strconv.ParseFloat(val, 64); err == nil {
			fs = append(fs, n)
		}
	}
	return fs
}


func StringssContain(ss []string, s string) (bool, []int){
	contain_flag := false
	index := []int{}
	for i, val := range ss {
		if val == s {
			contain_flag = true
			index = append(index, i)
		}
	}
	return contain_flag,index
}

func str2float(str ...string) []float64 {
	flo := make([]float64, 0)
	for _, eachStr := range str {
		fl, _ := strconv.ParseFloat(eachStr, 64)
		flo = append(flo, fl)
	}
	return flo
}

func str2int32(str string) int32 {
	in, _ := strconv.ParseInt(str, 0, 32)
	return int32(in)
}

func strSplit(str ...string) [][]string {
	ss := make([][]string, 0)
	for _, eachStr := range str {
		s := strings.Split(eachStr, "/")
		ss = append(ss, s)
	}
	return ss
}


//func LineAdd(fileName string) ([]protectDevices.Line){
//
//	file, err := os.Open(fileName)
//	if err != nil {
//		fmt.Println("Open file error!", err)
//		return []protectDevices.Line{}
//	}
//	defer file.Close()
//
//	buf := bufio.NewReader(file)
//
//	lines := []protectDevices.Line{}
//	for {
//		line_read, _, err := buf.ReadLine()
//
//		line_sl := strings.Split(string(line_read), "	")	//tab分隔
//
//		if err != nil {
//			if err == io.EOF {
//				fmt.Println("File", fileName, "read ok!")
//				break
//			} else {
//				fmt.Println("Read file error!", err)
//				return []protectDevices.Line{}
//			}
//		}
//		line := protectDevices.Line{}
//		line.F_bus, line.T_bus = line_sl[0], line_sl[1]
//		line.Name = line_sl[0] + "_" + line_sl[1]
//		line.IfFault = false	//初始化为无故障
//		if line_sl[2] == "1"{
//			line.Status = true
//		}else{
//			line.Status = false
//		}
//		lines = append(lines, line)
//	}
//
//	return lines
//}

//func BusAdd(fileName string) ([]protectDevices.Bus){
//
//	file, err := os.Open(fileName)
//	if err != nil {
//		fmt.Println("Open file error!", err)
//		return []protectDevices.Bus{}
//	}
//	defer file.Close()
//
//	buf := bufio.NewReader(file)
//
//	buses := []protectDevices.Bus{}
//	for {
//		line_read, _, err := buf.ReadLine()
//
//		line_sl := strings.Split(string(line_read), "	")	//tab分隔
//
//		if err != nil {
//			if err == io.EOF {
//				fmt.Println("File", fileName, "read ok!")
//				break
//			} else {
//				fmt.Println("Read file error!", err)
//				return []protectDevices.Bus{}
//			}
//		}
//		bus := protectDevices.Bus{}
//		bus.Name = line_sl[0]
//		bus.Voltage = 1			//母线电压初始化为1
//		bus.IfFault = false	//初始化为无故障
//		if line_sl[1] == "1"{
//			bus.Status = true
//		}else{
//			bus.Status = false
//		}
//		buses = append(buses, bus)
//	}
//
//	return buses
//}


//func ShuntAdd(fileName string) ([]protectDevices.Shunt){
//
//	file, err := os.Open(fileName)
//	if err != nil {
//		fmt.Println("Open file error!", err)
//		return []protectDevices.Shunt{}
//	}
//	defer file.Close()
//
//	buf := bufio.NewReader(file)
//
//	shunts := []protectDevices.Shunt{}
//	for {
//		line_read, _, err := buf.ReadLine()
//
//		line_sl := strings.Split(string(line_read), "	")	//tab分隔
//
//		if err != nil {
//			if err == io.EOF {
//				fmt.Println("File", fileName, "read ok!")
//				break
//			} else {
//				fmt.Println("Read file error!", err)
//				return []protectDevices.Shunt{}
//			}
//		}
//		shunt := protectDevices.Shunt{}
//		shunt.Name = line_sl[0]
//		shunt.Bus = line_sl[1]
//		shunt.IfFault = false	//初始化为无故障
//
//		if line_sl[2] == "1"{
//			shunt.Status = true
//		}else{
//			shunt.Status = false
//		}
//		shunts = append(shunts, shunt)
//	}
//
//	return shunts
//}

//func SynAdd(fileName string) ([]protectDevices.Syn){
//
//	file, err := os.Open(fileName)
//	if err != nil {
//		fmt.Println("Open file error!", err)
//		return []protectDevices.Syn{}
//	}
//	defer file.Close()
//
//	buf := bufio.NewReader(file)
//
//	syns := []protectDevices.Syn{}
//	for {
//		line_read, _, err := buf.ReadLine()
//
//		line_sl := strings.Split(string(line_read), "	")	//tab分隔
//
//		if err != nil {
//			if err == io.EOF {
//				fmt.Println("File", fileName, "read ok!")
//				break
//			} else {
//				fmt.Println("Read file error!", err)
//				return []protectDevices.Syn{}
//			}
//		}
//		syn := protectDevices.Syn{}
//		syn.Name = line_sl[0]
//		syn.Bus = line_sl[1]
//		syn.IfFault = false	//初始化为无故障
//
//		if line_sl[2] == "1"{
//			syn.Status = true
//		}else{
//			syn.Status = false
//		}
//		syns = append(syns, syn)
//	}
//
//	return syns
//}

//func TransAdd(fileName string) ([]protectDevices.Trans){
//
//	file, err := os.Open(fileName)
//	if err != nil {
//		fmt.Println("Open file error!", err)
//		return []protectDevices.Trans{}
//	}
//	defer file.Close()
//
//	buf := bufio.NewReader(file)
//
//	transes := []protectDevices.Trans{}
//	for {
//		line_read, _, err := buf.ReadLine()
//
//		line_sl := strings.Split(string(line_read), "	")	//tab分隔
//
//		if err != nil {
//			if err == io.EOF {
//				fmt.Println("File", fileName, "read ok!")
//				break
//			} else {
//				fmt.Println("Read file error!", err)
//				return []protectDevices.Trans{}
//			}
//		}
//		trans := protectDevices.Trans{}
//		trans.Name = line_sl[0]
//		trans.Type_trans,_ = strconv.Atoi(line_sl[1])	//string转int
//		trans.IfFault = false	//初始化为无故障
//
//		buses := strings.Split(line_sl[2],"/")
//		for _, bus := range buses{
//			trans.Bus = append(trans.Bus, bus)
//		}
//		if line_sl[3] == "1"{
//			trans.Status = true
//		}else{
//			trans.Status = false
//		}
//		trans.IfExternalPro = false
//		trans.IfInternalPro = false
//		trans.Gas_density = 0		//瓦斯密度初始化为0
//		transes = append(transes, trans)
//	}
//
//	return transes
//}

