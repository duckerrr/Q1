package logic

import (
	"Protection/protectDevices"
)

// i, j为序号
func CurrentI(device *protectDevices.Device, i, j int, faultpos, beta, theta float64) bool {
	line := device.Lines[i]
	breaker := line.Breaker_name[j]
	settingval  := line.SettingVal[j]
	delay := line.Delay[j]
	index_br_from := device.BreakerMap_name[breaker[0]]
	index_br_to := device.BreakerMap_name[breaker[1]]

	// 拒动模拟
	ifRun_from := Weibull(beta,theta, delay[0])
	ifRun_to := Weibull(beta,theta, delay[1])

	//f_bus,且避开bus故障
	if faultpos > settingval[0] || faultpos == 0{
		ifRun_from = false
	}

	//to_bus,且避开bus故障
	if 1 - faultpos > settingval[1] || faultpos == 1{
		ifRun_to = false
	}

	//哪个breaker的ifRun为true，则状态(status)变为false
	if ifRun_from{
		device.Breakers[index_br_from].Status = false
		device.Breakers[index_br_from].Log_step = map[float64]bool{delay[0]:false}
	}

	if ifRun_to{
		device.Breakers[index_br_to].Status = false
		device.Breakers[index_br_to].Log_step = map[float64]bool{delay[1]:false}
	}

	//两个breaker都动作了，故障清除了
	if ifRun_from && ifRun_to{
		return false
	}
	return true
}

//func CurrentII(Breaker breaker, breaker string, breakerG1 []string, delay float64, settingVal float64, cb chan bool, cs chan string, chan_ifOpen chan bool, faultStartTime time.Time, faultpos float64, faultLine string,beta float64, theta float64){
//
//	time.Sleep(time.Duration(delay*1000) * time.Millisecond)
//	_, ifOpen := <- chan_ifOpen
//	ifRun := Weibull(beta,theta, faultStartTime)
//	if settingVal == -1{
//		//如果没有设定值，默认1.5
//		settingVal = 1.2
//	}
//
//	idx := Breaker.index_map[breaker + "_" + "currentII"]
//	line := Breaker.line[idx]
//	bus := Breaker.bus[idx]
//	f_bus_fault := strings.Split(faultLine, "_")[0]
//
//	if ifOpen{
//		//如果channel没关，则运行
//		if line == faultLine{
//			//如果该breaker是在故障线路上
//			if bus == f_bus_fault{
//				//如果该breaker靠近f_bus
//				if settingVal > faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//			}else{
//				//如果该breaker靠近to_bus
//				if settingVal > 1 - faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//
//			}
//		}else{
//			//如果该breaker不在故障线路上
//			ifInG1, _ := SSIfInclude(breaker, breakerG1)
//			if ifInG1{
//				//如果该breaker是G1
//				if settingVal > 1 + faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//			}else{
//				//如果该breaker是G2
//				if settingVal > 2 - faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//			}
//		}
//	}
//
//}
//
//func CurrentIII(Breaker breaker, breaker string, breakerG1 []string, delay float64, settingVal float64, cb chan bool, cs chan string, chan_ifOpen chan bool, faultStartTime time.Time, faultpos float64, faultLine string, beta float64, theta float64){
//
//	time.Sleep(time.Duration(delay*1000) * time.Millisecond)
//	_, ifOpen := <- chan_ifOpen
//	ifRun := Weibull(beta,theta, faultStartTime)
//
//	if settingVal == -1{
//		//如果没有设定值，默认2
//		settingVal = 2
//	}
//
//	idx := Breaker.index_map[breaker + "_" + "III"]
//	line := Breaker.line[idx]
//	bus := Breaker.bus[idx]
//	f_bus_fault := strings.Split(faultLine, "_")[0]
//	if ifOpen{
//		//如果channel没关，则运行
//		if line == faultLine{
//			//如果该breaker是在故障线路上
//			if bus == f_bus_fault{
//				//如果该breaker靠近f_bus
//				if settingVal > faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//			}else{
//				//如果该breaker靠近to_bus
//				if settingVal > 1 - faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//
//			}
//		}else{
//			//如果该breaker不在故障线路上
//			ifInG1, _ := SSIfInclude(breaker, breakerG1)
//			if ifInG1{
//				//如果该breaker是G1
//				if settingVal > 1 + faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//			}else{
//				//如果该breaker是G2
//				if settingVal > 2 - faultpos{
//					cb <- ifRun
//					cs <- breaker
//				}else{
//					cb <- false
//					cs <- breaker
//				}
//			}
//		}
//	}
//
//}
