package base

import (
	"Protection/numeric"
	"Protection/protectDevices"
)

//构造邻接矩阵
func Build_Adjacency(device *protectDevices.Device){
	num_bus := len(device.Buses)
	rows, cols := []int32{}, []int32{}
	values := []float64{}

	// 两个母线有连接的，对称地置为1
	for _, line := range device.Lines{
		idx_from := device.BusMap[line.F_bus]
		idx_to := device.BusMap[line.T_bus]

		rows = append(rows, int32(idx_from))
		cols = append(cols, int32(idx_to))
		values = append(values, 1)

		rows = append(rows, int32(idx_to))
		cols = append(cols, int32(idx_from))
		values = append(values, 1)
	}

	// 对角元素置为1
	for i:=0; i<num_bus; i++{
		rows = append(rows, int32(i))
		cols = append(cols, int32(i))
		values = append(values, 1)
	}

	device.Adjacency_matric = &numeric.DSpmatrix{Rows: num_bus, Cols:num_bus, Nnz:len(values), Row_indx:rows, Col_indx:cols, Values:values}


}

//depth为需要寻找的最大深度，由全部断路器的属性为百分比的整定值最大值决定。
//如最大值为2.3，则向上取整，再减2；depth=1
//构造可达性矩阵
func Build_Reachable(device *protectDevices.Device){
	Adjacency := device.Adjacency_matric
	Adjacency.Init()

	//集中所有breakers的整定值
	all_settingVal := []float64{}
	for i:=0; i<len(device.Breakers); i++{
		all_settingVal = append(all_settingVal, device.Breakers[i].SettingVal...)
	}

	//定义depth
	_, depth_f := floats_Max(all_settingVal)
	depth := int(depth_f-1)

	mul_matrix := Adjacency
	mul_matrix.Init()
	device.Reachable_matrixs = append(device.Reachable_matrixs,Adjacency)

	for i:=1; i < depth+1; i++{
		mul_matrix.CSR = numeric.SparseDSPMM('N',mul_matrix.CSR,Adjacency.CSR,mul_matrix.CSR)
		temp := numeric.DSpmatrix{Rows: Adjacency.Rows, Cols:Adjacency.Cols, Nnz:1, Row_indx:[]int32{0}, Col_indx:[]int32{0}, Values:[]float64{0}}
		temp.Init()

		temp.CSR = numeric.SparseDAdd('N', mul_matrix.CSR,1, Adjacency.CSR, temp.CSR)
		temp.Values, temp.Row_indx, temp.Col_indx, temp.Rows, temp.Cols = numeric.GetCOO_D(temp.CSR)
		temp.Nnz = len(temp.Values)

		device.Reachable_matrixs = append(device.Reachable_matrixs, &temp)
		//device.Reachable_matrixs[i].Values, device.Reachable_matrixs[i].Row_indx, device.Reachable_matrixs[i].Col_indx,
		//	device.Reachable_matrixs[i].Rows, device.Reachable_matrixs[i].Cols = numeric.GetCOO_D(device.Reachable_matrixs[i].CSR)
	}

	Adjacency.Free()
	//temp_matrix.Free()
}

// 找稀疏矩阵中的某个元素
func Find_element(matrix *numeric.DSpmatrix, row, col int) float64{
	for m:=0; m<matrix.Nnz; m++{
		i := matrix.Row_indx[m]
		j := matrix.Col_indx[m]
		if i == int32(row) && j == int32(col){
			return matrix.Values[m]
		}
	}
	return 0
}

// 找稀疏矩阵中的某行/列元素
func Find_vec(matrix *numeric.DSpmatrix, rowORcol string, rowcol int) []float64{
	vec_value := []float64{}
	vec := []float64{}
	pos := []int{}
	if rowORcol == "row"{
		for m:=0; m<matrix.Nnz; m++{
			i := int(matrix.Row_indx[m])
			if i == rowcol{
				vec_value = append(vec_value, matrix.Values[m])
				pos = append(pos, int(matrix.Col_indx[m]))
			}
		}
		for i:=0; i<matrix.Cols;i++{
			vec = append(vec,0)
		}

		for i:=0; i<len(pos);i++{
			vec[pos[i]] = vec_value[i]
		}

	}else if rowORcol == "col"{
		for m:=0; m<matrix.Nnz; m++{
			i := int(matrix.Col_indx[m])
			if i == rowcol{
				vec_value = append(vec_value, matrix.Values[m])
				pos = append(pos, int(matrix.Row_indx[m]))
			}
		}

		for i:=0; i<matrix.Rows;i++{
			vec = append(vec,0)
		}

		for i:=0; i<len(pos);i++{
			vec[pos[i]] = vec_value[i]
		}

	}else{
		panic("请输入正确的字符，row表示取行向量，col表示取列向量")
	}


	return vec
}


//找float64切片的最大值
func floats_Max(fs []float64) (int,float64){
	max_idx := 0
	max_val := fs[0]
	for i, val := range fs {
		if val > max_val{
			max_idx = i
			max_val = val
		}
	}
	return max_idx, max_val
}