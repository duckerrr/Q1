package protectDevices

type Bus struct {
	BaseDevice
	Voltage float64
	Num_breaker int 		//断路器数量
}
