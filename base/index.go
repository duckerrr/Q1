package base

import "Protection/protectDevices"

func LineIndex(device *protectDevices.Device)	{
	device.LineMap = map[string]int{}
	for i:=0; i<len(device.Lines); i++{
		f_bus, t_bus := device.Lines[i].F_bus, device.Lines[i].T_bus
		device.LineMap[f_bus + "_" + t_bus] = i
		device.LineMap[t_bus + "_" + f_bus] = i
	}
}

func BusIndex(device *protectDevices.Device)	{
	device.BusMap = map[string]int{}
	for i:=0; i<len(device.Buses); i++{
		device.BusMap[device.Buses[i].Name] = i
	}
}

func ShuntIndex(device *protectDevices.Device)	{
	device.ShuntMap = map[string]int{}
	for i:=0; i<len(device.Shunts); i++{
		device.ShuntMap[device.Shunts[i].Name] = i
	}
}

func SynIndex(device *protectDevices.Device)	{
	device.SynMap = map[string]int{}
	for i:=0; i<len(device.Syns); i++{
		device.SynMap[device.Syns[i].Name] = i
	}
}

func TransIndex(device *protectDevices.Device)	{
	device.TransMap = map[string]int{}
	for i:=0; i<len(device.Transes); i++{
		device.TransMap[device.Transes[i].Name] = i
	}
}

func Breaker_NameIndex(device *protectDevices.Device){
	device.BreakerMap_name = map[string]int{}
	for i:=0; i<len(device.Breakers); i++{
		device.BreakerMap_name[device.Breakers[i].Name] = i
	}
}

func Breaker_BusLineIndex(device *protectDevices.Device){
	device.BreakerMap_busLine = map[string]int{}
	for i:=0; i<len(device.Breakers); i++{
		bus2 := device.Breakers[i].Bus_close + "+" + device.Breakers[i].Line
		device.BreakerMap_busLine[bus2] = i
	}
}