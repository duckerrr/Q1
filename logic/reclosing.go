package logic

import "Protection/protectDevices"

func Reclosing(device *protectDevices.Device, breaker_Run []string, time_Run []float64, time_now float64) map[string]map[string]float64{
	// breaker_Run中逐个断路器判断是否配置有重合闸，若有启动重合闸

	//action的第一层key是：breaker的名字；第二层key是：动作类型——第一次断开(break1)/合上(close)/
	//		第二次断开(break2)/闭锁(lock)/解除闭锁(unlock),value是：执行该动作的时间
	// 对于重合闸的断路器，先把故障未消除的全流程
	action := map[string]map[string]float64{}

	setReclosings_brIndex := []int{}	// 记录设置了重合闸的断路器的序号
	for i, breaker := range breaker_Run{
		br_i := device.BreakerMap_name[breaker]	//序号
		if device.Breakers[br_i].IfReclosing{
			setReclosings_brIndex = append(setReclosings_brIndex, br_i)
			//配置了重合闸的动作后，先立即合上,若故障未消除，延时后断开，同时闭锁，一段时间后解除闭锁
			temp_map := map[string]float64{"break1":time_now + time_Run[i]}
			temp_map["close"] = time_now + time_Run[i]
			temp_map["break2"] = time_now + time_Run[i] + device.Breakers[br_i].Delay_recl
			temp_map["lock"] = time_now + time_Run[i] + device.Breakers[br_i].Delay_recl
			temp_map["unlock"] = time_now + time_Run[i] + device.Breakers[br_i].Delay_recl +
										device.Breakers[br_i].Lock_recl
			action[breaker] = temp_map
		}else {
			// 没有配置重合闸的直接动作
			action[breaker] = map[string]float64{"break1":time_now + time_Run[i]}
		}
	}

	return action
}
