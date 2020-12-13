package common

import (
	"Protection/logic"
	"Protection/protectDevices"
)

//保护功能仿真
// fault_info的外层key(string)为故障设备类型，内层key(string)为设备名称，value(float64)为故障位置
// 由于服务每个时步都会调用protect，因此需要传入此刻时间time_now
func Protect(device *protectDevices.Device, fault_info map[string]map[string]float64, beta, theta, time_now float64) (map[string]map[string]float64,map[string]map[string][]string){

	// 记录故障的设备与对应动作的断路器(用于故障清除时，重合闸中断)
	// 第一层key是设备类型； 第二层key是设备名称，value是断路器集合
	faultDevice_Breaker := map[string]map[string][]string{}

	//动作的断路器及动作时间
	breaker_Run := []string{}
	time_Run := []float64{}

	// 一. 根据fault_info判断起动什么保护，得到动作的断路器及动作时间
	// 对fault_info第一层key_value拆解
	for fault_type, sub_map := range fault_info{
		if fault_type == "line"{
			// 对fault_info第二层key_value拆解
			for fault_line, fault_pos := range sub_map{
				device.Lines[device.LineMap[fault_line]].IfFault = true
				breaker_run, time_run := logic.LineProtect(device,fault_line, fault_pos, beta, theta)
				breaker_Run = append(breaker_Run,breaker_run...)
				time_Run = append(time_Run,time_run...)
				faultDevice_Breaker["line"] = map[string][]string{fault_line:breaker_run}
			}
		}

		if fault_type == "bus"{
			for fault_bus, _ := range sub_map{
				device.Lines[device.LineMap[fault_bus]].IfFault = true
				breaker_run, time_run := logic.BusProtect(device, fault_bus, beta, theta)
				breaker_Run = append(breaker_Run,breaker_run...)
				time_Run = append(time_Run,time_run...)
				faultDevice_Breaker["bus"] = map[string][]string{fault_bus:breaker_run}
			}
		}


		if fault_type == "syn"{
			for fault_syn, _ := range sub_map{
				device.Lines[device.LineMap[fault_syn]].IfFault = true
				breaker_run, time_run := logic.SynProtect(device, fault_syn, beta, theta)
				breaker_Run = append(breaker_Run,breaker_run...)
				time_Run = append(time_Run,time_run...)
				faultDevice_Breaker["syn"] = map[string][]string{fault_syn:breaker_run}

			}
		}

		if fault_type == "trans"{
			for fault_trans, _ := range sub_map{
				device.Lines[device.LineMap[fault_trans]].IfFault = true
				breaker_run, time_run := logic.TransProtect(device, fault_trans, beta, theta)
				breaker_Run = append(breaker_Run,breaker_run...)
				time_Run = append(time_Run,time_run...)
				faultDevice_Breaker["trans"] = map[string][]string{fault_trans:breaker_run}

			}
		}
	}

	// 电容电抗低压保护是通过device.Buses[i].voltage进行判断
	breaker_run, time_run := logic.ShuntProtect(device, beta, theta)
	breaker_Run = append(breaker_Run,breaker_run...)
	time_Run = append(time_Run,time_run...)

	Action := logic.Reclosing(device, breaker_Run, time_Run, time_now)

	return Action, faultDevice_Breaker
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