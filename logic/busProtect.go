package logic

import (
	"Protection/protectDevices"
)

func BusProtect(device *protectDevices.Device, fault_bus string, beta, theta float64) ([]string,[]float64){
	busIndex := device.BusMap[fault_bus]

	// 1. 找出所有需要启动的断路器及其延时,并根据其出线(bus_far)分组
	breaker_start_sep := [][]string{}
	delay_sep := [][]float64{}

	line_sep := []string{}	//辅助记录breaker_sep每一组对应的线路是什么

	for i, breakers := range device.Buses[busIndex].Breaker_name{
		for j, breaker := range breakers{
			m := device.BreakerMap_name[breaker]
			line := device.Breakers[m].Line
			contain, index := StringssContain(line_sep, line)
			if !contain {
				//如果该断路器对应的母线未进行分组或远端无母线，则新建一组;
				line_sep = append(line_sep, line)
				breaker_start_sep = append(breaker_start_sep, []string{breaker})
				delay_sep = append(delay_sep, []float64{device.Buses[busIndex].Delay[i][j]})
			}else {
				//如果该断路器对应的母线已进行分组，则将断路器加入至原来一组中
				breaker_start_sep[index[0]] = append(breaker_start_sep[index[0]], breaker)
				delay_sep[index[0]] = append(delay_sep[index[0]], device.Buses[busIndex].Delay[i][j])
			}
		}
	}

	// 2. 根据延时从小到大的顺序，对延时及断路器切片进行了排序
	for i, _ := range breaker_start_sep{
		delay_sep[i], breaker_start_sep[i] = Sort2fs(delay_sep[i],breaker_start_sep[i],"up")
	}

	// 3. 各组分别按时间先后运行断路器
	breaker_Run := []string{}
	time_Run := []float64{}

	for i, breakers := range breaker_start_sep{
		for j, breaker := range breakers{
			ifRun := Weibull(delay_sep[i][j], beta, theta)
			if ifRun{
				breaker_Run = append(breaker_Run, breaker)
				time_Run = append(time_Run, delay_sep[i][j])
				break
			}
		}
	}

	return breaker_Run, time_Run
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