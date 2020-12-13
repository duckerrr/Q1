package protectDevices

//全网络断路器的信息，不基于baseDevice
type Breaker struct{
	Name string						// breaker的名字
	Status bool						//开合状态，1为闭合，0为断开
	Type_protection []string  		// 保护类型：线路保护/母线保护/变压器保护/……
	Protection []string				// 具体保护
	Delay []float64					// 保护动作延时
	SettingType	[]string			// 电流值(current)/电压值(voltage)/保护范围百分比(percentage)
	SettingVal []float64			// 整定值

	Bus_close string 				// breaker所靠近的bus(包括trans和syn)
	Bus_far string 					// breaker所在线路上的另一端的母线
	Line string						// breaker所在的线路

	IfReclosing bool 				//是否有配置重合闸
	Delay_recl float64				//重合闸的延时，无则-1
	Lock_recl float64				//重合闸的闭锁时间，无则-1
	IfLock bool						//true为正在闭锁，false为未闭锁

	Log_step map[float64]bool		//记录每个时步的动作情况，每个时步开始前先清空；key为时间，value为动作情况(对应status,false为闭合->断开)
}