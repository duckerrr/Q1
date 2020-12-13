package main

import (
	"Protection/common"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/garyburd/redigo/redis"
)

// 接收一次系统的信息
type Energy_struct struct {
	Voltage []float64	`json:"Voltage"`
	Time float64	`json:"Time"`
}

//发送给一次系统的信息
type Line_action struct {
	Line_name []string	`json:"Line_name"`
	Action []string`json:"Action"`
}

func main() {

	time_now := 0.
	// 读取配置文件
	fp := "grid/CEPRI 36节点系统.xlsx"
	xlsx, err := excelize.OpenFile(fp)
	fp2 := "setting/setting.xlsx"
	xlsx_protection, err := excelize.OpenFile(fp2)

	device := common.Load(xlsx, xlsx_protection)

	breaker_step := map[string]int{}	//记录某个breaker的步数

	Actions , faultDevice_Breakers := []map[string]map[string]float64{}, []map[string]map[string][]string{}		//记录每个产生动作信息的时步的动作信息

	// 用redis接收一次系统的电压信息，及时间
	pubToEnergy,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		panic( err)
		return
	}
	defer pubToEnergy.Close() //程序结束后主动关闭客户端，下同

	subFromEnergy,err := redis.Dial("tcp","127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer subFromEnergy.Close()

	// 监听一次系统
	// chan是用于反馈订阅的信息是否处理好或 是否成功订阅
	chan_SubscribeSuccess1 := make(chan bool)

	// 订阅一次系统的频道
	//接收一次系统数据的结构体
	EnergyStruct := Energy_struct{}
	chan_energyData := make(chan Energy_struct)

	go func(EnergyStruct *Energy_struct) {
		psc := redis.PubSubConn{Conn: subFromEnergy}
		err_sub := psc.Subscribe("energy_pub")
		if err_sub != nil {
			panic(err_sub)
		}else {
			for {//死循环，由chan_SubscribeSuccess控制进行/暂停
				switch v := psc.Receive().(type) {
				case redis.Message:
					err_unmarshal := json.Unmarshal(v.Data, EnergyStruct)
					if err_unmarshal == nil{
						// 执行收到订阅信息后的操作
						fmt.Printf("%s: message: ", v.Channel)
						chan_SubscribeSuccess1 <- true
						chan_energyData <- *EnergyStruct
					}
				case redis.Subscription:
					fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
				case error:
					panic(v)
				}
			}
		}
	}(&EnergyStruct)

	// 用redis接收故障系统的故障信息(待完善)
	// 故障信息传入fault_info
	// 故障清除信息传入fault_disappear_info

	for {
		//将Actions中空的删除，faultDevice_Breakers对应删除
		for i, action := range Actions{
			if len(action) == 0{
				Actions = append(Actions[:i], Actions[i+1:]...)
				faultDevice_Breakers = append(faultDevice_Breakers[:i], faultDevice_Breakers[i+1:]...)
			}
		}

		fault_info := map[string]map[string]float64{}
		//if math.Abs(time_now - 0.05) < 0.001{
		//	fault_info = map[string]map[string]float64{"bus":{"tianhe.220M2":0}, "syn":{"SG1":0},"trans":{"Zzhongxin.B1":0}}
		//}
		//if math.Abs(time_now - 0.15) < 0.001{
		//	fault_info = map[string]map[string]float64{"line":{"CN-49897_CN-84181":0.5}}
		//}
		//fault_disappear_info := map[string]string{"line":"CN-49897_CN-84181"}
		//fault_disappear_info := map[string]string{"bus":"tianhe.220M2"}

		fault_disappear_info := map[string]string{}

		if <-chan_SubscribeSuccess1 == true{
			EnergyStruct = <- chan_energyData
			fmt.Println(EnergyStruct)
		}
		for i:=0; i < len(device.Buses); i++{
			device.Buses[i].Voltage = EnergyStruct.Voltage[i]
		}
		time_next := EnergyStruct.Time
		Action , faultDevice_Breaker := common.Protect(device, fault_info, 0.01,100, time_now)
		// Action非空，才append
		if len(Action) != 0{
			Actions = append(Actions, Action)
			faultDevice_Breakers = append(faultDevice_Breakers, faultDevice_Breaker)
		}

		// 根据故障清除信息，通过faultDevice_Breaker找到对应的断路器，将Action中对应的动作信息删除
		for device_type, device_name := range fault_disappear_info{
			breakers := faultDevice_Breaker[device_type][device_name]
			for _, breaker := range breakers{
				for _, action := range Actions{
					delete(action, breaker)
				}
			}
		}

		for breaker, _ := range Action {
			breaker_step[breaker] = 0
		}
		//breaker名字+ "+" + 动作信息
		br_now := []string{}
		action_now := []string{}

		// 根据Action修改断路器的状态，并记录交互信息
		for _, Action := range Actions{
			for breaker, _ := range Action{
				br_i := device.BreakerMap_name[breaker]

				// 只有一个动作的对应 —— 没有配置重合闸
				if len(Action[breaker]) == 1{

					time := Action[breaker]["break1"]

					// 如果时间动作时间在区间[此刻时间, 此刻时间+步长]，则可动作
					if time >= time_now && time < time_next{
						// 如果该断路器本来就断开或者仍在闭锁，则不执行动作指令
						if device.Breakers[br_i].Status == false && device.Breakers[br_i].IfLock {continue}

						device.Breakers[br_i].Status = false
						br_now = append(br_now, breaker)
						action_now = append(action_now, "break")
						delete(breaker_step, breaker)
						delete(Action, breaker)
					}

				}else if len(Action[breaker]) == 5{
					// step0——“break1” / step1 —— “close” / step2 —— “break2” / step3 —— "lock" / step4 —— "unlock"

					for j:=0; j<5; j++{
						switch breaker_step[breaker] {
						case 0:
							time := Action[breaker]["break1"]
							if time > time_now && time <= time_next{
								// 如果该断路器本来就断开或者仍在闭锁，则不执行动作指令
								if device.Breakers[br_i].Status == false && device.Breakers[br_i].IfLock {continue}

								device.Breakers[br_i].Status = false
								action_now = append(action_now, breaker + "+" + "break")
								breaker_step[breaker] = 1
							}

						case 1:
							time := Action[breaker]["close"]
							if time > time_now && time <= time_next {

								device.Breakers[br_i].Status = true
								action_now = append(action_now, breaker + "+" + "close")
								breaker_step[breaker] = 2
							}

						case 2:
							time := Action[breaker]["break2"]
							if time >= time_now && time < time_next {
								device.Breakers[br_i].Status = false
								action_now = append(action_now, breaker + "+" + "break")
								breaker_step[breaker] = 3
							}

						case 3:
							time := Action[breaker]["lock"]
							if time >= time_now && time < time_next {
								device.Breakers[br_i].IfLock = true
								breaker_step[breaker] = 4
							}

						case 4:
							time := Action[breaker]["unlock"]
							if time >= time_now && time < time_next {
								device.Breakers[br_i].IfLock = false
								delete(breaker_step, breaker)
								delete(Action, breaker)
							}
						}
					}
				}
			}
		}
		time_now = time_next
		// 线路动作信息，发布至一次系统
		line_action := Line_action{}

		for i:=0; i < len(br_now); i++ {
			if action_now[i]=="break1" || action_now[i]=="break2" || action_now[i]=="break"{
				line_action.Line_name = append(line_action.Line_name, device.Breakers[device.BreakerMap_name[br_now[i]]].Line)
				line_action.Action = append(line_action.Action, "break")
			}
			if action_now[i]=="close"{
				line_action.Line_name = append(line_action.Line_name, device.Breakers[device.BreakerMap_name[br_now[i]]].Line)
				line_action.Action = append(line_action.Action, "close")
			}
		}
		value, _ := json.Marshal(line_action)
		_ ,err_pub := pubToEnergy.Do("PUBLISH", "energy_sub", value)
		if err_pub != nil{panic(err_pub)}
	}

	//fmt.Println(time2.Now().Sub(st))

}

