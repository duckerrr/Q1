package protectDevices

import "Protection/numeric"

type Device struct {
	//线路故障看Line类型中断路器的配置，母线故障看Bus类型中断路器的配置……
	Lines []Line
	Buses []Bus
	Transes []Trans
	Shunts []Shunt
	Syns []Syn
	Breakers []Breaker

	//索引
	LineMap map[string]int 		//线路索引的string是f_bus+"_"+t_bus 和 t_bus+"_"+f_bus ,两个key都能找到同一个索引
	BusMap map[string]int 		//母线索引的string是bus的name
	TransMap map[string]int		//变压器索引的string是 trans的name
	ShuntMap map[string]int		//电容电抗索引的string是 shunt的name
	SynMap map[string]int		//发电机索引的string是 syn的name
	BreakerMap_name map[string]int	//断路器索引的string是 breaker的name,
	BreakerMap_busLine map[string]int	//断路器索引的string是 bus + "+" + line，两个bus唯一确定一个breaker

	//线路和母线构成的邻接矩阵
	Adjacency_matric *numeric.DSpmatrix			//母线i和母线j相连，则邻接矩阵Adjacency_matrix[i][j]=true；
	Reachable_matrixs []*numeric.DSpmatrix 		//可达性矩阵Reachable是邻接矩阵的1次方到n次方的累加，用于判断某个母线在跨越n个母线后，与哪些母线有连接关系
												//[0]是Adjacency，[1]是邻接矩阵的1次方到2次方的累加.……
}