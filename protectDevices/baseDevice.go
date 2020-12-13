package protectDevices

type BaseDevice struct {
	Name string						// transformer的名字
	Status      bool    				// 设备连接状态
	Breaker_name [][]string			// 配置断路器名称
	SettingType []string			// 对应的整定值的类型(电流/电压/百分比……)
	SettingVal [][]float64			// 对应的整定值
	Delay [][]float64				// 对应的延时
	Protection []string			// 保护类型
	IfFault bool					// 设备是否有故障，true为有故障
}
