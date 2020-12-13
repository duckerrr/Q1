package logic

import "Protection/protectDevices"

func ShuntProtect(device *protectDevices.Device,beta, theta float64) ([]string, []float64){
	breaker_start := []string{}
	delay := []float64{}

	// 1. 遍历全体shunt对应母线的电压，与整定值进行比较，找出将被启动的断路器
	for _, shunt := range device.Shunts{
		busIndex := device.BusMap[shunt.Bus]
		for j, settingVals := range shunt.SettingVal{
			for k, settingVal := range settingVals{
				//若母线电压低于整定值，且保护类型为shuntVL，且此时断路器状态为"闭合"，则该断路器将会被启动
				breaker := shunt.Breaker_name[j][k]
				if device.Buses[busIndex].Voltage < settingVal && shunt.Protection[j] == "shuntVL" && device.Breakers[device.BreakerMap_name[breaker]].Status{
					breaker_start = append(breaker_start, breaker)
					delay = append(delay, shunt.Delay[j][k])
				}
			}
		}
	}

	// 2. 根据延时从小到大的顺序，对延时及断路器切片进行了排序
	delay, breaker_start = Sort2fs(delay,breaker_start,"up")

	// 3. 分别按时间先后运行断路器
	breaker_Run := []string{}
	time_Run := []float64{}

	for i, breaker := range breaker_start{
		ifRun := Weibull(delay[i], beta, theta)
		if ifRun{
			breaker_Run = append(breaker_Run, breaker)
			time_Run = append(time_Run, delay[i])
		}
	}
	return breaker_Run, time_Run
}
