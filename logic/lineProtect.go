package logic

import (
	"Protection/base"
	"Protection/protectDevices"
)

//var lineProtectType = []string{"currentI","currentII","currentIII","distanceI","distanceII","distanceIII","linePilot"}

func LineProtect(device *protectDevices.Device, fault_line string, fault_pos, beta, theta float64) ([]string, []float64){
	//for fault_line, _ := range val_map{
	idx := device.LineMap[fault_line]
	//f_bus和t_bus的序号
	f_bus := device.Lines[idx].F_bus
	t_bus := device.Lines[idx].T_bus
	fi := device.BusMap[f_bus]
	ti := device.BusMap[t_bus]

	// 定义“深度"：故障线路上的断路器深度定义为0，相邻线路上的断路器深度定义为1，以此规律逐渐往外扩展，并加1.
	// 从可达性矩阵中，找出所有可能会动作的断路器，并对断路器进行分层(不同深度)，再对整定值做判断，最后对满足要求的整定值的断路器按延时由小到大排序。
	// 由于可由其靠近的母线+所在线路上的另一母线，两个条件唯一确定一个断路器，因此先寻找可能动作的断路器对应的母线集合。

	//1. 对所有可能动作的断路器对应的母线进行分层（bus_level_sep）。bus_indexVecs最后一个切片是全体可能会动作的断路器对应的母线集合
	bus_level_sep := [][]int{}
	bus_indexVecs := [][]int{}
	////第一级是faultLine的两个母线
	bus_indexVecs = append(bus_indexVecs, []int{fi,ti})
	for i:=0; i<len(device.Reachable_matrixs); i++{
		row_vec := base.Find_vec(device.Reachable_matrixs[i],"row",fi)
		col_vec := base.Find_vec(device.Reachable_matrixs[i],"col",ti)

		// 找出故障线路对应的行和列中不为0元素的序号并集(即两个向量中，都为0的位置才为0)
		onesVec := []int{}
		for j:=0; j<len(row_vec); j++{
			if row_vec[j] != 0 || col_vec[j] != 0{
				onesVec = append(onesVec, j)
			}
		}
		bus_indexVecs = append(bus_indexVecs, onesVec)
	}

	for i:=0; i<len(bus_indexVecs); i++{
		if i == 0{
			//第一级是faultLine的两个母线
			bus_level_sep = append(bus_level_sep, bus_indexVecs[0])
		}else {
			//其它级要与上一级比较，差集即为该级的母线
			bus_level_sep = append(bus_level_sep, differenceSet(bus_indexVecs[i-1],bus_indexVecs[i]))
		}
	}

	//2.根据母线的分级情况，推得断路器的分级情况，得到所有可能动作的断路器+对应的深度。
	already_busIndex := []int{fi, ti}
	// 所有可能动作的断路器的序号, breaker_depth为对应位置上母线深度
	breaker_allIndex := []int{device.BreakerMap_busLine[f_bus+"+"+fault_line], device.BreakerMap_busLine[t_bus+"+"+fault_line]}
	breaker_depth := []int{0,0}
	for i:=1; i<len(bus_level_sep); i++{
		for _, busIndexi := range bus_level_sep[i]{
			//判断该母线是否与already_bus集合中的母线有连接关系(根据邻接矩阵判断)
			ifConnect := false
			for _, busIndexj := range already_busIndex{
				if base.Find_element(device.Adjacency_matric, busIndexi, busIndexj) != 0{
					//如果有连接关系
					ifConnect = true

					busLine_key := device.Buses[busIndexi].Name + "+" + device.Lines[device.LineMap[device.Buses[busIndexi].Name+
											"_"+device.Buses[busIndexj].Name]].Name
					breaker_allIndex = append(breaker_allIndex, device.BreakerMap_busLine[busLine_key])
					breaker_depth = append(breaker_depth, i)
				}
			}
			if ifConnect{
				already_busIndex = append(already_busIndex, busIndexi)
			}
		}
	}

	//3.得到所有可能动作的断路器集合后，根据它们的整定值判断是否在动作范围内，并且对满足要求整定值的断路器按延时由小到大排序。

	//故障位置(fault_pos)是基于f_bus而言的
	//根据它们的整定值判断是否在动作范围内,注意：需要将断路器分成两个集合(靠近f_bus和靠近t_bus的)
	//定义：CF:close_from,CT: close_to
	//3.1 将断路器与对应深度的两个切片分别分成两个集合(靠近f_bus和靠近t_bus的)；利用断路器对应层级的可达性矩阵进行判断——
	//		通过判断断路器对应的母线所在的行，f_bus和t_bus位置上的数值，若为0则不是靠近该侧母线
	//		深度1对应reachable[0],深度2对应reachable[1]……
	brIndex_CF, brIndex_CT := []int{breaker_allIndex[0]}, []int{breaker_allIndex[1]}	//记录对应断路器的序号
	depth_CF, depth_CT := []int{0}, []int{0}	//记录对应断路器的深度
	for i:=0; i<len(breaker_allIndex); i++{
		if breaker_depth[i] ==0{
			//深度为0不需要讨论(就是故障线路两端的母线)
			continue
		}
		depth := breaker_depth[i]

		// 由于bus_indexVecs[-1]与breaker_allIndex是一一对应
		vec := base.Find_vec(device.Reachable_matrixs[depth-1],"row",bus_indexVecs[len(bus_indexVecs)-1][i])
		//向量中f_bus和t_bus分别对应位置的值
		f_val, t_val := vec[fi], vec[ti]
		if f_val == 0 {
			brIndex_CT = append(brIndex_CT, breaker_allIndex[i])
			depth_CT = append(depth_CT,breaker_depth[i])
		}
		if t_val == 0 {
			brIndex_CF = append(brIndex_CF, breaker_allIndex[i])
			depth_CF = append(depth_CF,breaker_depth[i])
		}
	}

	// 3.2 在CF和CT两个集合中，都分别找出故障位置在整定值范围内的断路器
	breaker_start_CF, breaker_start_CT := []int{}, []int{}		//启动的断路器的序号
	delay_CF, delay_CT := []float64{}, []float64{}
	// 靠近f_bus一侧的断路器，是用depth+fault_pos和settingVal进行比较
	for i:=0; i<len(brIndex_CF); i++{
		br_i := brIndex_CF[i]
		for j, settingVal := range device.Breakers[br_i].SettingVal{
			if float64(depth_CF[i]) + fault_pos > settingVal{
				continue
			}

			// 判断断路器的保护是属于线路保护
			if device.Breakers[br_i].Type_protection[j] == "line"{
				breaker_start_CF = append(breaker_start_CF, brIndex_CF[i])
				delay_CF = append(delay_CF, device.Breakers[br_i].Delay[j])
			}
		}
	}

	// 靠近t_bus一侧的断路器，是用depth+1-fault_pos和settingVal进行比较
	for i:=0; i<len(brIndex_CT); i++{
		br_i := brIndex_CT[i]
		for j, settingVal := range device.Breakers[br_i].SettingVal{
			if float64(depth_CT[i]) + fault_pos > settingVal{
				continue
			}

			// 判断断路器的保护是属于线路保护
			if device.Breakers[br_i].Type_protection[j] == "line"{
				breaker_start_CT = append(breaker_start_CT, brIndex_CT[i])
				delay_CT = append(delay_CT, device.Breakers[br_i].Delay[j])
			}
		}
	}

	// 3.3 根据延时从小到大的顺序，对延时及断路器切片进行了排序
	delay_CF, breaker_start_CF = Sort2fi(delay_CF,breaker_start_CF,"up")
	delay_CT, breaker_start_CT = Sort2fi(delay_CT,breaker_start_CT,"up")

	// 3.4 CF,CT两边都依次按顺序执行指令，若拒动，则到下一个断路器；两边都有断路器动作时，故障切除
	//ifRun_CF, ifRun_CT := false, false
	breaker_Run := []string{}
	time_Run := []float64{}

	for i:=0; i<len(breaker_start_CF); i++{
		ifRun := Weibull(delay_CF[i],beta, theta)
		if ifRun{
			//ifRun_CF = true
			breaker_Run = append(breaker_Run, device.Breakers[breaker_start_CF[i]].Name)
			time_Run = append(time_Run, delay_CF[i])
			break
		}
	}

	for i:=0; i<len(breaker_start_CT); i++{
		ifRun := Weibull(delay_CT[i],beta, theta)
		if ifRun{
			//ifRun_CT = true
			breaker_Run = append(breaker_Run, device.Breakers[breaker_start_CT[i]].Name)
			time_Run = append(time_Run, delay_CT[i])
			break
		}
	}
	return breaker_Run, time_Run
}



//寻找切片中等于某个值的所有元素位置
func floats_find(fs []float64, f float64) (bool, []int){
	exist := false
	indexs := []int{}
	for i, val := range fs{
		if val == f{
			exist = true
			indexs = append(indexs, i)
		}
	}
	return exist, indexs
}

func intssContain(fs []int, f int) (bool, []int){
	contain_flag := false
	index := []int{}
	for i, val := range fs {
		if val == f {
			contain_flag = true
			index = append(index, i)
		}
	}
	return contain_flag,index
}

func differenceSet(s_short, s_long []int) []int{
	difference := []int{}
	for _, val := range s_long{
		contain, _ := intssContain(s_short,val)
		if !contain{
			difference = append(difference, val)
		}
	}
	return difference
}

//两个slice，相同位置的数据是对应的，对一个进行排序，另一个slice也随之变动
// []int随着[]float64变，s string 为"up"时，为对float64升序；"down"为降序
func Sort2fi(fs []float64, is []int, s string) ([]float64, []int){
	switch s {
	case "up":
		for i:=0; i<len(fs)-1; i++{
			for j:=i+1; j<len(fs); j++{
				if fs[j] < fs[i]{
					fs[i], fs[j] = fs[j], fs[i]
					is[i], is[j] = is[j], is[i]
				}
			}
		}
	case "down":
		for i:=0; i<len(fs)-1; i++{
			for j:=i+1; j<len(fs); j++{
				if fs[j] > fs[i]{
					fs[i], fs[j] = fs[j], fs[i]
					is[i], is[j] = is[j], is[i]
				}
			}
		}
	default:
		panic("输入up为升序，输入down为降序")
		return []float64{}, []int{}
	}

	return fs, is
}

// []string随着[]float64变
func Sort2fs(fs []float64, ss []string, s string) ([]float64, []string){
	switch s {
	case "up":
		for i:=0; i<len(fs)-1; i++{
			for j:=i+1; j<len(fs); j++{
				if fs[j] < fs[i]{
					fs[i], fs[j] = fs[j], fs[i]
					ss[i], ss[j] = ss[j], ss[i]
				}
			}
		}
	case "down":
		for i:=0; i<len(fs)-1; i++{
			for j:=i+1; j<len(fs); j++{
				if fs[j] > fs[i]{
					fs[i], fs[j] = fs[j], fs[i]
					ss[i], ss[j] = ss[j], ss[i]
				}
			}
		}
	default:
		panic("输入up为升序，输入down为降序")
		return []float64{}, []string{}
	}

	return fs, ss
}


