package common

import (
	"Protection/base"
	"Protection/protectDevices"
	"github.com/360EntSecGroup-Skylar/excelize"
)


func Load(xlsx_grid, xlsx_protection *excelize.File) (*protectDevices.Device){
	device :=  protectDevices.Device{}
	device.Lines = base.LineAdd(xlsx_grid)
	device.Buses = base.BusAdd(xlsx_grid)
	device.Shunts = base.ShuntAdd(xlsx_grid)
	device.Syns = base.SynAdd(xlsx_grid)
	buses := []protectDevices.Bus{}
	device.Transes, buses = base.TransAdd(xlsx_grid)
	device.Buses = append(device.Buses, buses...)
	device.Breakers = base.BreakerAdd(xlsx_protection)

	//要对line/shunt/syn/trans/breaker建立索引后，才能继续settingAdd
	base.LineIndex(&device)
	base.BusIndex(&device)
	base.ShuntIndex(&device)
	base.SynIndex(&device)
	base.TransIndex(&device)
	base.Breaker_NameIndex(&device)
	base.SettingAdd(xlsx_protection,&device)

	//已经录入所有breakers，先对breakers建立索引
	base.Breaker_BusLineIndex(&device)
	base.ReclosingAdd(xlsx_protection, &device)

	//构造可达性矩阵
	base.Build_Adjacency(&device)
	base.Build_Reachable(&device)

	return &device
}

