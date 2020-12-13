package numeric

/*
#cgo CFLAGS: -I/home/hjm/intel/mkl/include
#cgo LDFLAGS: -L/home/hjm/intel/mkl/lib/intel64 -L/home/hjm/intel/lib/intel64 -lmkl_intel_lp64 -lmkl_intel_thread -lmkl_core -liomp5 -lpthread -lm
#include<mkl_spblas.h>
#include<mkl.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Spmatrix struct {
	CSR C.sparse_matrix_t
}

func (m *Spmatrix) Free() {
	C.mkl_sparse_destroy(m.CSR)
}

func Operation(operation byte, A C.sparse_matrix_t) C.sparse_matrix_t {
	if operation == 'N' {
		C.mkl_sparse_convert_csr(A, C.SPARSE_OPERATION_NON_TRANSPOSE, &A)
	} else if operation == 'T' {
		C.mkl_sparse_convert_csr(A, C.SPARSE_OPERATION_TRANSPOSE, &A)
	} else if operation == 'C' {
		C.mkl_sparse_convert_csr(A, C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, &A)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return A
}

/*
实数稀疏矩阵
*/

type DSpmatrix struct {
	Rows, Cols, Nnz    int
	Row_indx, Col_indx []int32
	Values             []float64
	COO                C.sparse_matrix_t
	CSR                C.sparse_matrix_t
}

func (m *DSpmatrix) Set(row, col int32, value float64) {
	m.Row_indx = append(m.Row_indx, row)
	m.Col_indx = append(m.Col_indx, col)
	m.Values = append(m.Values, value)
}

//初始化
func (m *DSpmatrix) Init() {
	m.Nnz = len(m.Values)
	rows, cols, nnz := C.int(m.Rows), C.int(m.Cols), C.int(m.Nnz)
	row_indx := (*C.int)(unsafe.Pointer(&m.Row_indx[0]))
	col_indx := (*C.int)(unsafe.Pointer(&m.Col_indx[0]))
	val := (*C.double)(unsafe.Pointer(&m.Values[0]))
	C.mkl_sparse_d_create_coo(&m.COO, C.SPARSE_INDEX_BASE_ZERO, rows, cols, nnz, row_indx, col_indx, val)
	C.mkl_sparse_convert_csr(m.COO, C.SPARSE_OPERATION_NON_TRANSPOSE, &m.CSR)
}

//释放内存
func (m *DSpmatrix) Free() {
	C.mkl_sparse_destroy(m.COO)
	C.mkl_sparse_destroy(m.CSR)
}

//改变某个值（改之前也要是非零的）
func SparseDSetValue(A C.sparse_matrix_t, row_indx int, col_indx int, value float64) {
	row, col, val := C.int(row_indx), C.int(col_indx), C.double(value)
	C.mkl_sparse_d_set_value(A, row, col, val)
}

//稀疏矩阵和向量的乘法运算
func SparseDMV(operation byte, alpha float64, A C.sparse_matrix_t, x []float64, beta float64, y []float64) []float64 {
	ctype_x := (*C.double)(unsafe.Pointer(&x[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	ctype_alpha := C.double(alpha)
	ctype_beta := C.double(beta)
	descr := C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}
	if operation == 'N' {
		C.mkl_sparse_d_mv(C.SPARSE_OPERATION_NON_TRANSPOSE, ctype_alpha, A, descr, ctype_x, ctype_beta, ctype_y)
	} else if operation == 'T' {
		C.mkl_sparse_d_mv(C.SPARSE_OPERATION_TRANSPOSE, ctype_alpha, A, descr, ctype_x, ctype_beta, ctype_y)
	} else if operation == 'C' {
		C.mkl_sparse_d_mv(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, ctype_alpha, A, descr, ctype_x, ctype_beta, ctype_y)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return y
}

//稀疏矩阵和稠密矩阵的乘法运算
func SparseDMM(operation byte, alpha float64, A C.sparse_matrix_t, x []float64, row_x, row_y, col_y int, beta float64, y []float64) []float64 {
	ctype_x := (*C.double)(unsafe.Pointer(&x[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	ctype_row_x := C.int(row_x)
	ctype_row_y := C.int(row_y)
	ctype_col_y := C.int(col_y)
	ctype_alpha := C.double(alpha)
	ctype_beta := C.double(beta)
	descr := C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}
	if operation == 'N' {
		C.mkl_sparse_d_mm(C.SPARSE_OPERATION_NON_TRANSPOSE, ctype_alpha, A, descr, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_x, ctype_col_y, ctype_row_x, ctype_beta, ctype_y, ctype_row_y)
	} else if operation == 'T' {
		C.mkl_sparse_d_mm(C.SPARSE_OPERATION_TRANSPOSE, ctype_alpha, A, descr, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_x, ctype_col_y, ctype_row_x, ctype_beta, ctype_y, ctype_row_y)
	} else if operation == 'C' {
		C.mkl_sparse_d_mm(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, ctype_alpha, A, descr, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_x, ctype_col_y, ctype_row_x, ctype_beta, ctype_y, ctype_row_y)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return y
}

//稀疏矩阵的加法
func SparseDAdd(operation byte, A C.sparse_matrix_t, alpha float64, B C.sparse_matrix_t, c C.sparse_matrix_t) C.sparse_matrix_t {
	ctype_alpha := C.double(alpha)
	if operation == 'N' {
		C.mkl_sparse_d_add(C.SPARSE_OPERATION_NON_TRANSPOSE, A, ctype_alpha, B, &c)
	} else if operation == 'T' {
		C.mkl_sparse_d_add(C.SPARSE_OPERATION_TRANSPOSE, A, ctype_alpha, B, &c)
	} else if operation == 'C' {
		C.mkl_sparse_d_add(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, A, ctype_alpha, B, &c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//稀疏矩阵和稀疏矩阵的乘法（结果存为稀疏矩阵）
func SparseDSPMM(operation byte, A C.sparse_matrix_t, B C.sparse_matrix_t, c C.sparse_matrix_t) C.sparse_matrix_t {
	if operation == 'N' {
		C.mkl_sparse_spmm(C.SPARSE_OPERATION_NON_TRANSPOSE, A, B, &c)
	} else if operation == 'T' {
		C.mkl_sparse_spmm(C.SPARSE_OPERATION_TRANSPOSE, A, B, &c)
	} else if operation == 'C' {
		C.mkl_sparse_spmm(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, A, B, &c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//稀疏矩阵和稀疏矩阵的乘法（结果存为稠密矩阵）
func SparseDSPMMD(operation byte, A C.sparse_matrix_t, B C.sparse_matrix_t, c []float64, row_c int) []float64 {
	ctype_c := (*C.double)(unsafe.Pointer(&c[0]))
	ctype_row_c := C.int(row_c)
	if operation == 'N' {
		C.mkl_sparse_d_spmmd(C.SPARSE_OPERATION_NON_TRANSPOSE, A, B, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_c, ctype_row_c)
	} else if operation == 'T' {
		C.mkl_sparse_d_spmmd(C.SPARSE_OPERATION_TRANSPOSE, A, B, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_c, ctype_row_c)
	} else if operation == 'C' {
		C.mkl_sparse_d_spmmd(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, A, B, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_c, ctype_row_c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//稀疏矩阵和稀疏矩阵的乘法（两个稀疏矩阵均可做转置、共轭转置变换，结果存为稀疏矩阵）
func SparseDSP2M(operation_A, operation_B byte, A C.sparse_matrix_t, B C.sparse_matrix_t, c C.sparse_matrix_t) C.sparse_matrix_t {
	descr_A, descr_B := C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}, C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}
	if operation_A == 'N' && operation_B == 'N' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_NON_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_NON_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'N' && operation_B == 'T' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_NON_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'N' && operation_B == 'C' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_NON_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'T' && operation_B == 'N' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_NON_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'T' && operation_B == 'T' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'T' && operation_B == 'C' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'C' && operation_B == 'N' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_NON_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'C' && operation_B == 'T' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'C' && operation_B == 'C' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//输出CSR格式稀疏矩阵的阶数、行起始索引、行终止索引、列索引、对应非零元素
func PrintDMatrix(m C.sparse_matrix_t) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var nnz int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.double
	C.mkl_sparse_d_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	fmt.Println("实数矩阵CSR存储格式的数据：")
	fmt.Printf("矩阵阶数:%vX%v\n", rows, cols)
	gotype_values := (unsafe.Pointer(values))
	gotype_cols_indx := (unsafe.Pointer(cols_indx))
	gotype_rows_start := (unsafe.Pointer(rows_start))
	gotype_rows_end := (unsafe.Pointer(rows_end))
	nnz = int(*(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(int(rows)-1))))
	rows_indx := make([]C.int, nnz)
	for i := 0; i < int(rows); i++ {
		row_start := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start) + 4*uintptr(i)))
		row_end := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(i)))
		for j := int(*row_start); j < int(*row_end); j++ {
			rows_indx[j] = C.int(i)
		}
	}
	fmt.Printf("非零元素个数:%v\n", nnz)
	for i := 0; i < nnz; i++ {
		values_double := (*C.double)(unsafe.Pointer(uintptr(gotype_values) + 8*uintptr(i)))
		cols_index := (*C.int)(unsafe.Pointer(uintptr(gotype_cols_indx) + 4*uintptr(i)))
		fmt.Printf("值: %v  ", *values_double)
		fmt.Printf("行: %v  ", rows_indx[i])
		fmt.Printf("列: %v\n", *cols_index)
	}
}

func GetDMatrixvalue(m C.sparse_matrix_t) ([]float64, []int32) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var nnz int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.double
	C.mkl_sparse_d_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	//fmt.Println("实数矩阵CSR存储格式的数据：")
	//fmt.Printf("矩阵阶数:%vX%v\n",rows,cols)
	gotype_values := (unsafe.Pointer(values))
	//gotype_cols_indx:=(unsafe.Pointer(cols_indx))
	gotype_rows_start := (unsafe.Pointer(rows_start))
	gotype_rows_end := (unsafe.Pointer(rows_end))
	nnz = int(*(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(int(rows)-1))))
	rows_indx := make([]int32, nnz)
	values_out := make([]float64, nnz)
	for i := 0; i < int(rows); i++ {
		row_start := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start) + 4*uintptr(i)))
		row_end := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(i)))
		for j := int(*row_start); j < int(*row_end); j++ {
			rows_indx[j] = int32(i)
		}
	}
	//fmt.Printf("非零元素个数:%v\n",nnz)
	for i := 0; i < nnz; i++ {
		values_double := (*C.double)(unsafe.Pointer(uintptr(gotype_values) + 8*uintptr(i)))
		values_out[i] = float64(*values_double)
		//cols_index:=(*C.int)(unsafe.Pointer(uintptr(gotype_cols_indx) + 4*uintptr(i)))
		//fmt.Printf("值: %v  ",*values_double)
		//fmt.Printf("行: %v  ",rows_indx[i])
		//fmt.Printf("列: %v\n",*cols_index)
	}

	return values_out, rows_indx
}
func GetCOO_D(m C.sparse_matrix_t) ([]float64, []int32, []int32, int, int) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var nnz int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.MKL_Complex16
	// 此处有误
	C.mkl_sparse_z_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	gotype_values := (unsafe.Pointer(values))
	gotype_cols_indx := (unsafe.Pointer(cols_indx))
	gotype_rows_start := (unsafe.Pointer(rows_start))
	gotype_rows_end := (unsafe.Pointer(rows_end))
	nnz = int(*(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(int(rows)-1))))
	rows_indx := make([]int32, nnz)
	cols_indx1 := make([]int32, nnz)
	values_out := make([]float64, nnz)
	rows1 := int(rows)
	cols1 := int(cols)
	for i := 0; i < int(rows); i++ {
		row_start := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start) + 4*uintptr(i)))
		row_end := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(i)))
		for j := int(*row_start); j < int(*row_end); j++ {
			rows_indx[j] = int32(i)
		}
	}

	for i := 0; i < nnz; i++ {
		values_double := (*C.double)(unsafe.Pointer(uintptr(gotype_values) + 8*uintptr(i)))
		cols_index := (*C.int)(unsafe.Pointer(uintptr(gotype_cols_indx) + 4*uintptr(i)))
		values_out[i] = float64(*values_double)
		cols_indx1[i] = int32(*cols_index)

	}
	return values_out, rows_indx, cols_indx1, rows1, cols1
}

func GetZMatrixvalue(m C.sparse_matrix_t) ([]complex128, []int32, []int32) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var nnz int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.MKL_Complex16
	C.mkl_sparse_z_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	//fmt.Println("复数矩阵CSR存储格式的数据：")
	//fmt.Printf("矩阵阶数:%vX%v\n",rows,cols)
	gotype_values := (unsafe.Pointer(values))
	gotype_cols_indx := (unsafe.Pointer(cols_indx))
	gotype_rows_start := (unsafe.Pointer(rows_start))
	gotype_rows_end := (unsafe.Pointer(rows_end))
	nnz = int(*(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(int(rows)-1))))
	rows_indx := make([]int32, nnz)
	cols_indx1 := make([]int32, nnz)
	values_out := make([]complex128, nnz)
	for i := 0; i < int(rows); i++ {
		row_start := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start) + 4*uintptr(i)))
		row_end := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(i)))
		for j := int(*row_start); j < int(*row_end); j++ {
			rows_indx[j] = int32(i)
		}
	}
	//fmt.Printf("非零元素个数:%v\n",nnz)
	for i := 0; i < nnz; i++ {
		values_complex128 := (*complex128)(unsafe.Pointer(uintptr(gotype_values) + 16*uintptr(i)))
		values_out[i] = complex128(*values_complex128)
		cols_index := (*C.int)(unsafe.Pointer(uintptr(gotype_cols_indx) + 4*uintptr(i)))
		cols_indx1[i] = int32(*cols_index)
		//fmt.Printf("值: %v  ",*values_complex128)
		//fmt.Printf("行: %v  ",rows_indx[i])
		//fmt.Printf("列: %v\n",*cols_index)
	}
	return values_out, rows_indx, cols_indx1
}

//返回CSR格式稀疏矩阵的阶数，行起始索引、列索引、对应非零元素的指针
func GetIndexValuesDMatrix(m C.sparse_matrix_t) (C.int, C.int, *C.int, *C.int, *C.int, *C.double) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.double
	C.mkl_sparse_d_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	return rows, cols, rows_start, rows_end, cols_indx, values
}

/*
复数稀疏矩阵
*/

type ZSpmatrix struct {
	Rows, Cols, Nnz    int
	Row_indx, Col_indx []int32
	Values             []complex128
	COO                C.sparse_matrix_t
	CSR                C.sparse_matrix_t
}

func (m *ZSpmatrix) Set(row, col int32, value complex128) {
	m.Row_indx = append(m.Row_indx, row)
	m.Col_indx = append(m.Col_indx, col)
	m.Values = append(m.Values, value)
}

//初始化
func (m *ZSpmatrix) Init() {
	rows, cols, nnz := C.int(m.Rows), C.int(m.Cols), C.int(m.Nnz)
	row_indx := (*C.int)(unsafe.Pointer(&m.Row_indx[0]))
	col_indx := (*C.int)(unsafe.Pointer(&m.Col_indx[0]))
	val := (*C.MKL_Complex16)(unsafe.Pointer(&m.Values[0]))
	C.mkl_sparse_z_create_coo(&m.COO, C.SPARSE_INDEX_BASE_ZERO, rows, cols, nnz, row_indx, col_indx, val)
	C.mkl_sparse_convert_csr(m.COO, C.SPARSE_OPERATION_NON_TRANSPOSE, &m.CSR)
}

//将CSR格式复数矩阵内部的值输至ZSpmatrix结构体内
func (m *ZSpmatrix)GetCOO_Z() {
	var indexing C.sparse_index_base_t
	var rows,cols C.int
	var nnz int
	var cols_indx,rows_start,rows_end *C.int
	var values *C.MKL_Complex16
	C.mkl_sparse_z_export_csr(m.CSR,&indexing,&rows,&cols,&rows_start,&rows_end,&cols_indx,&values)
	gotype_values:=(unsafe.Pointer(values))
	gotype_cols_indx:=(unsafe.Pointer(cols_indx))
	gotype_rows_start:=(unsafe.Pointer(rows_start))
	gotype_rows_end:=(unsafe.Pointer(rows_end))
	nnz=int(*(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(int(rows)-1))))
	m.Row_indx=make([]int32,nnz)
	m.Col_indx=make([]int32,nnz)
	m.Values= make([]complex128,nnz)
	m.Rows= int(rows)
	m.Cols= int(cols)
	m.Nnz = nnz
	for i:=0;i<int(rows);i++{
		row_start:=(*C.int)(unsafe.Pointer(uintptr(gotype_rows_start) + 4*uintptr(i)))
		row_end:=(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(i)))
		for j:=int(*row_start);j<int(*row_end);j++{
			m.Row_indx[j]=int32(i)
		}
	}

	for i:=0;i<nnz;i++{
		values_complex:=(*complex128)(unsafe.Pointer(uintptr(gotype_values) + 16*uintptr(i)))
		cols_index:=(*int32)(unsafe.Pointer(uintptr(gotype_cols_indx) + 4*uintptr(i)))
		m.Values[i] = *values_complex
		m.Col_indx[i] = *cols_index

	}
}

//释放内存
func (m *ZSpmatrix) Free() {
	C.mkl_sparse_destroy(m.COO)
	C.mkl_sparse_destroy(m.CSR)
}

//直接创建复数CSR格式矩阵
func SparseZCreateCSR(CSR *C.sparse_matrix_t,rows,cols int,rowStart,rowEnd,colIndx *int32,Value *complex128) {
	C.mkl_sparse_z_create_csr(CSR,C.SPARSE_INDEX_BASE_ZERO,C.int(rows),C.int(cols),(*C.int)(unsafe.Pointer(rowStart)),(*C.int)(unsafe.Pointer(rowEnd)),(*C.int)(unsafe.Pointer(colIndx)),(*C.MKL_Complex16)(unsafe.Pointer(Value)))
}

//改变某个值（改之前也要是非零的）
func SparseZSetValue(A C.sparse_matrix_t, row_indx int, col_indx int, value complex128) {
	row, col := C.int(row_indx), C.int(col_indx)
	val := C.MKL_Complex16{C.double(real(value)), C.double(imag(value))}
	C.mkl_sparse_z_set_value(A, row, col, val)
}

//稀疏矩阵和向量的乘法运算
func SparseZMV(operation byte, alpha complex128, A C.sparse_matrix_t, x []complex128, beta complex128, y []complex128) []complex128 {
	ctype_x := (*C.MKL_Complex16)(unsafe.Pointer(&x[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	ctype_alpha := C.MKL_Complex16{C.double(real(alpha)), C.double(imag(alpha))}
	ctype_beta := C.MKL_Complex16{C.double(real(beta)), C.double(imag(beta))}
	descr := C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}
	if operation == 'N' {
		C.mkl_sparse_z_mv(C.SPARSE_OPERATION_NON_TRANSPOSE, ctype_alpha, A, descr, ctype_x, ctype_beta, ctype_y)
	} else if operation == 'T' {
		C.mkl_sparse_z_mv(C.SPARSE_OPERATION_TRANSPOSE, ctype_alpha, A, descr, ctype_x, ctype_beta, ctype_y)
	} else if operation == 'C' {
		C.mkl_sparse_z_mv(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, ctype_alpha, A, descr, ctype_x, ctype_beta, ctype_y)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}

	return y
}

//稀疏矩阵和稠密矩阵的乘法运算
func SparseZMM(operation byte, alpha complex128, A C.sparse_matrix_t, x []complex128, row_x, row_y, col_y int, beta complex128, y []complex128) []complex128 {
	ctype_x := (*C.MKL_Complex16)(unsafe.Pointer(&x[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	ctype_alpha := C.MKL_Complex16{C.double(real(alpha)), C.double(imag(alpha))}
	ctype_beta := C.MKL_Complex16{C.double(real(beta)), C.double(imag(beta))}
	ctype_row_x := C.int(row_x)
	ctype_row_y := C.int(row_y)
	ctype_col_y := C.int(col_y)
	descr := C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}
	if operation == 'N' {
		C.mkl_sparse_z_mm(C.SPARSE_OPERATION_NON_TRANSPOSE, ctype_alpha, A, descr, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_x, ctype_col_y, ctype_row_x, ctype_beta, ctype_y, ctype_row_y)
	} else if operation == 'T' {
		C.mkl_sparse_z_mm(C.SPARSE_OPERATION_TRANSPOSE, ctype_alpha, A, descr, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_x, ctype_col_y, ctype_row_x, ctype_beta, ctype_y, ctype_row_y)
	} else if operation == 'C' {
		C.mkl_sparse_z_mm(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, ctype_alpha, A, descr, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_x, ctype_col_y, ctype_row_x, ctype_beta, ctype_y, ctype_row_y)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}

	return y
}

//稀疏矩阵的加法
func SparseZAdd(operation byte, A C.sparse_matrix_t, alpha complex128, B C.sparse_matrix_t, c C.sparse_matrix_t) C.sparse_matrix_t {
	ctype_alpha := C.MKL_Complex16{C.double(real(alpha)), C.double(imag(alpha))}
	if operation == 'N' {
		C.mkl_sparse_z_add(C.SPARSE_OPERATION_NON_TRANSPOSE, A, ctype_alpha, B, &c)
	} else if operation == 'T' {
		C.mkl_sparse_z_add(C.SPARSE_OPERATION_TRANSPOSE, A, ctype_alpha, B, &c)
	} else if operation == 'C' {
		C.mkl_sparse_z_add(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, A, ctype_alpha, B, &c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//稀疏矩阵和稀疏矩阵的乘法（结果存为稀疏矩阵）
func SparseZSPMM(operation byte, A C.sparse_matrix_t, B C.sparse_matrix_t, c C.sparse_matrix_t) C.sparse_matrix_t {
	if operation == 'N' {
		C.mkl_sparse_spmm(C.SPARSE_OPERATION_NON_TRANSPOSE, A, B, &c)
	} else if operation == 'T' {
		C.mkl_sparse_spmm(C.SPARSE_OPERATION_TRANSPOSE, A, B, &c)
	} else if operation == 'C' {
		C.mkl_sparse_spmm(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, A, B, &c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//稀疏矩阵和稀疏矩阵的乘法（结果存为稠密矩阵）
func SparseZSPMMD(operation byte, A C.sparse_matrix_t, B C.sparse_matrix_t, c []complex128, row_c int) []complex128 {
	ctype_c := (*C.MKL_Complex16)(unsafe.Pointer(&c[0]))
	ctype_row_c := C.int(row_c)
	if operation == 'N' {
		C.mkl_sparse_z_spmmd(C.SPARSE_OPERATION_NON_TRANSPOSE, A, B, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_c, ctype_row_c)
	} else if operation == 'T' {
		C.mkl_sparse_z_spmmd(C.SPARSE_OPERATION_TRANSPOSE, A, B, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_c, ctype_row_c)
	} else if operation == 'C' {
		C.mkl_sparse_z_spmmd(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, A, B, C.SPARSE_LAYOUT_COLUMN_MAJOR, ctype_c, ctype_row_c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//稀疏矩阵和稀疏矩阵的乘法（两个稀疏矩阵均可做转置、共轭转置变换，结果存为稀疏矩阵）
func SparseZSP2M(operation_A, operation_B byte, A C.sparse_matrix_t, B C.sparse_matrix_t, c C.sparse_matrix_t) C.sparse_matrix_t {
	descr_A, descr_B := C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}, C.struct_matrix_descr{C.SPARSE_MATRIX_TYPE_GENERAL, 0, 0}
	if operation_A == 'N' && operation_B == 'N' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_NON_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_NON_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'N' && operation_B == 'T' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_NON_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'N' && operation_B == 'C' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_NON_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'T' && operation_B == 'N' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_NON_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'T' && operation_B == 'T' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'T' && operation_B == 'C' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'C' && operation_B == 'N' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_NON_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'C' && operation_B == 'T' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else if operation_A == 'C' && operation_B == 'C' {
		C.mkl_sparse_sp2m(C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_A, A, C.SPARSE_OPERATION_CONJUGATE_TRANSPOSE, descr_B, B, C.SPARSE_STAGE_FULL_MULT, &c)
	} else {
		panic("请输入正确的字符，N表示不操作，T表示转置，C表示共轭转置")
	}
	return c
}

//输出COO格式稀疏矩阵的阶数、行起始索引、行终止索引、列索引、对应非零元素
func PrintZMatrix(m C.sparse_matrix_t) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var nnz int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.MKL_Complex16
	C.mkl_sparse_z_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	fmt.Println("复数矩阵CSR存储格式的数据：")
	fmt.Printf("矩阵阶数:%vX%v\n", rows, cols)
	gotype_values := (unsafe.Pointer(values))
	gotype_cols_indx := (unsafe.Pointer(cols_indx))
	gotype_rows_start := (unsafe.Pointer(rows_start))
	gotype_rows_end := (unsafe.Pointer(rows_end))
	nnz = int(*(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(int(rows)-1))))
	rows_indx := make([]C.int, nnz)
	for i := 0; i < int(rows); i++ {
		row_start := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start) + 4*uintptr(i)))
		row_end := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_end) + 4*uintptr(i)))
		for j := int(*row_start); j < int(*row_end); j++ {
			rows_indx[j] = C.int(i)
		}
	}
	fmt.Printf("非零元素个数:%v\n", nnz)
	for i := 0; i < nnz; i++ {
		values_complex128 := (*complex128)(unsafe.Pointer(uintptr(gotype_values) + 16*uintptr(i)))
		cols_index := (*C.int)(unsafe.Pointer(uintptr(gotype_cols_indx) + 4*uintptr(i)))
		fmt.Printf("值: %v  ", *values_complex128)
		fmt.Printf("行: %v  ", rows_indx[i])
		fmt.Printf("列: %v\n", *cols_index)
	}
}

//返回CSR格式稀疏矩阵的阶数，行起始索引、列索引、对应非零元素的指针
func GetIndexValuesZMatrix(m C.sparse_matrix_t) (C.int, C.int, *C.int, *C.int, *C.int, *C.MKL_Complex16) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.MKL_Complex16
	C.mkl_sparse_z_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	return rows, cols, rows_start, rows_end, cols_indx, values
}

//复数稀疏矩阵取实部/虚部
func Z2DSpmatrix(m C.sparse_matrix_t) (C.sparse_matrix_t, C.sparse_matrix_t) {
	var indexing C.sparse_index_base_t
	var rows, cols C.int
	var nnz C.int
	var cols_indx, rows_start, rows_end *C.int
	var values *C.MKL_Complex16
	C.mkl_sparse_z_export_csr(m, &indexing, &rows, &cols, &rows_start, &rows_end, &cols_indx, &values)
	gotype_values := (unsafe.Pointer(values))
	gotype_rows_start := (unsafe.Pointer(rows_start))
	nnz_indx := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start) + 4*uintptr(rows)))
	nnz = *(nnz_indx)
	values_real, values_imag := make([]float64, nnz), make([]float64, nnz)
	for i := 0; i < int(nnz); i++ {
		values_complex128 := (*complex128)(unsafe.Pointer(uintptr(gotype_values) + 16*uintptr(i)))
		values_real[i] = real(*values_complex128)
		values_imag[i] = imag(*values_complex128)
	}
	ctype_real_values := (*C.double)(unsafe.Pointer(&values_real[0]))
	ctype_imag_values := (*C.double)(unsafe.Pointer(&values_imag[0]))
	var real_CSR, imag_CSR C.sparse_matrix_t
	C.mkl_sparse_d_create_csr(&real_CSR, C.SPARSE_INDEX_BASE_ZERO, rows, cols, rows_start, rows_end, cols_indx, ctype_real_values)
	C.mkl_sparse_d_create_csr(&imag_CSR, C.SPARSE_INDEX_BASE_ZERO, rows, cols, rows_start, rows_end, cols_indx, ctype_imag_values)
	return real_CSR, imag_CSR
}

/*
稀疏求解器
*/

//整合了实数/复数
func Pardiso(mtype C.int, phase C.int, n int, values unsafe.Pointer, rows_start *C.int, cols_indx *C.int, nrhs int, iparm [64]int32, b unsafe.Pointer, x unsafe.Pointer) unsafe.Pointer {
	var maxfct, mnum C.int = 1, 1
	var msglvl C.int = 1
	var perm, Error C.int
	var pt [64]unsafe.Pointer
	//C.mkl_set_dynamic(0)
	//C.mkl_set_num_threads(8)
	ctype_n := C.int(n)
	ctype_pt := unsafe.Pointer(&pt[0])
	ctype_iparm := (*C.int)(unsafe.Pointer(&iparm[0]))
	ctype_nrhs := C.int(nrhs)
	C.pardisoinit(C._MKL_DSS_HANDLE_t(ctype_pt), &mtype, ctype_iparm)
	iparm[34] = 1
	C.pardiso(C._MKL_DSS_HANDLE_t(ctype_pt), &maxfct, &mnum, &mtype, &phase, &ctype_n, values, rows_start, cols_indx, &perm, &ctype_nrhs, ctype_iparm, &msglvl, b, x, &Error)
	phase = -1
	var d unsafe.Pointer
	C.pardiso(C._MKL_DSS_HANDLE_t(ctype_pt), &maxfct, &mnum, &mtype, &phase, &ctype_n, values, rows_start, cols_indx, &perm, &ctype_nrhs, ctype_iparm, &msglvl, d, d, &Error)
	return x
}

//求实数N*N矩阵的逆
func DInv(mtype C.int, n int, values *C.double, rows_start *C.int, cols_indx *C.int, x []float64) []float64 {
	var maxfct, mnum C.int = 1, 1
	var msglvl C.int = 1
	var perm, Error C.int
	var pt [64]unsafe.Pointer
	var iparm [64]int32
	b := make([]float64, n*n)
	for i := 0; i < n; i++ {
		b[(n+1)*i] = 1
	}
	phase := C.int(13)
	ctype_n := C.int(n)
	ctype_pt := unsafe.Pointer(&pt[0])
	ctype_iparm := (*C.int)(unsafe.Pointer(&iparm[0]))
	C.pardisoinit(C._MKL_DSS_HANDLE_t(ctype_pt), &mtype, ctype_iparm)
	iparm[34] = 1
	C.pardiso(C._MKL_DSS_HANDLE_t(ctype_pt), &maxfct, &mnum, &mtype, &phase, &ctype_n, unsafe.Pointer(values), rows_start, cols_indx, &perm, &ctype_n, ctype_iparm, &msglvl, unsafe.Pointer(&b[0]), unsafe.Pointer(&x[0]), &Error)
	phase = -1
	var d unsafe.Pointer
	C.pardiso(C._MKL_DSS_HANDLE_t(ctype_pt), &maxfct, &mnum, &mtype, &phase, &ctype_n, unsafe.Pointer(values), rows_start, cols_indx, &perm, &ctype_n, ctype_iparm, &msglvl, d, d, &Error)
	return x
}

//求复数N*N矩阵的逆
func ZInv(mtype C.int, n int, values *C.MKL_Complex16, rows_start *C.int, cols_indx *C.int, x []complex128) []complex128 {
	var maxfct, mnum C.int = 1, 1
	var msglvl C.int = 1
	var perm, Error C.int
	var pt [64]unsafe.Pointer
	var iparm [64]int32
	b := make([]complex128, n*n)
	for i := 0; i < n; i++ {
		b[(n+1)*i] = 1
	}
	phase := C.int(13)
	ctype_n := C.int(n)
	ctype_pt := unsafe.Pointer(&pt[0])
	ctype_iparm := (*C.int)(unsafe.Pointer(&iparm[0]))
	C.pardisoinit(C._MKL_DSS_HANDLE_t(ctype_pt), &mtype, ctype_iparm)
	iparm[34] = 1
	C.pardiso(C._MKL_DSS_HANDLE_t(ctype_pt), &maxfct, &mnum, &mtype, &phase, &ctype_n, unsafe.Pointer(values), rows_start, cols_indx, &perm, &ctype_n, ctype_iparm, &msglvl, unsafe.Pointer(&b[0]), unsafe.Pointer(&x[0]), &Error)
	phase = -1
	var d unsafe.Pointer
	C.pardiso(C._MKL_DSS_HANDLE_t(ctype_pt), &maxfct, &mnum, &mtype, &phase, &ctype_n, unsafe.Pointer(values), rows_start, cols_indx, &perm, &ctype_n, ctype_iparm, &msglvl, d, d, &Error)
	return x
}

/*
矩阵合并
*/
func Concatenate(A, B, c, D C.sparse_matrix_t) C.sparse_matrix_t {
	/*
		var Extend_A,Extend_B,Extend_C,Extend_D,Matrix C.sparse_matrix_t
		rows_A,cols_A,rows_start_A,rows_end_A,cols_indx_A,values_A:=Get_index_values_d_matrix(A)
		gotype_rows_start:=(unsafe.Pointer(rows_start_A))
		var i C.int
		for i=rows_A;i<rows_A+rows_B;i++{

		}

		C. mkl_sparse_d_create_csr(
	*/
	var matrix C.sparse_matrix_t
	return matrix
}

/*
矩阵元素相乘(CSR格式）
*/
func SparseZMulEle(A, B C.sparse_matrix_t) C.sparse_matrix_t {
	var indexing_A, indexing_B C.sparse_index_base_t
	var rows_A, rows_B, cols_A, cols_B C.int
	var nnz_A int
	var cols_indx_A, cols_indx_B, rows_start_A, rows_start_B, rows_end_A, rows_end_B *C.int
	var values_A, values_B *C.MKL_Complex16
	C.mkl_sparse_z_export_csr(A, &indexing_A, &rows_A, &cols_A, &rows_start_A, &rows_end_A, &cols_indx_A, &values_A)
	C.mkl_sparse_z_export_csr(B, &indexing_B, &rows_B, &cols_B, &rows_start_B, &rows_end_B, &cols_indx_B, &values_B)
	gotype_values_A := (unsafe.Pointer(values_A))
	gotype_cols_indx_A := (unsafe.Pointer(cols_indx_A))
	gotype_rows_start_A := (unsafe.Pointer(rows_start_A))
	gotype_rows_end_A := (unsafe.Pointer(rows_end_A))
	nnz_A = int(*(*C.int)(unsafe.Pointer(uintptr(gotype_rows_end_A) + 4*uintptr(int(rows_A)-1))))
	gotype_values_B := (unsafe.Pointer(values_B))
	gotype_cols_indx_B := (unsafe.Pointer(cols_indx_B))
	gotype_rows_start_B := (unsafe.Pointer(rows_start_B))
	gotype_rows_end_B := (unsafe.Pointer(rows_end_B))
	row_indx_C, col_index_C, values_C := make([]int32, nnz_A), make([]int32, nnz_A), make([]complex128, nnz_A)
	var num int
	for i := 0; i < int(rows_A); i++ {
		row_start_A := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start_A) + 4*uintptr(i)))
		row_end_A := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_end_A) + 4*uintptr(i)))
		row_start_B := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_start_B) + 4*uintptr(i)))
		row_end_B := (*C.int)(unsafe.Pointer(uintptr(gotype_rows_end_B) + 4*uintptr(i)))
		for j := *row_start_A; j < *row_end_A; j++ {
			cols_index_A := (*C.int)(unsafe.Pointer(uintptr(gotype_cols_indx_A) + 4*uintptr(j)))
			for k := *row_start_B; k < *row_end_B; k++ {
				cols_index_B := (*C.int)(unsafe.Pointer(uintptr(gotype_cols_indx_B) + 4*uintptr(k)))
				if *cols_index_A == *cols_index_B {
					values_complex128_A := (*complex128)(unsafe.Pointer(uintptr(gotype_values_A) + 16*uintptr(j)))
					values_complex128_B := (*complex128)(unsafe.Pointer(uintptr(gotype_values_B) + 16*uintptr(k)))
					values_C[num] = (*values_complex128_A) * (*values_complex128_B)
					row_indx_C[num] = int32(i)
					col_index_C[num] = int32(*cols_index_A)
					num++
					break
				}
				if *cols_index_A < *cols_index_B {
					break
				}
			}

		}
	}
	if num == 0 {
		panic("乘积为全零矩阵")
	}
	c := ZSpmatrix{Rows: int(rows_A), Cols: int(cols_A), Nnz: num, Row_indx: row_indx_C, Col_indx: col_index_C, Values: values_C}
	c.Init()
	return c.CSR
}

//实现元素乘
func SparseZVMmul_Ele(A ZSpmatrix, B ZSpmatrix) ZSpmatrix {
	v := A.Values
	var rows, cols []int32
	for i := 0; i < len(v); i++ {
		rows = append(rows, int32(i))
		cols = append(cols, int32(i))
	}
	diagV := ZSpmatrix{Rows: len(v), Cols: len(v), Nnz: len(v), Row_indx: rows, Col_indx: cols, Values: v}
	diagV.Init()
	var c ZSpmatrix
	c.CSR = SparseZSPMM('N', diagV.CSR, B.CSR, c.CSR)
	return c
}

func ZSparsemul(A ZSpmatrix, B ZSpmatrix) ZSpmatrix {
	var c ZSpmatrix
	c.CSR = SparseZSPMM('N', A.CSR, B.CSR, c.CSR)
	return c
}

func DSpmatrixAugment(A DSpmatrix, B DSpmatrix) DSpmatrix {
	//var C DSpmatrix
	var rows, cols []int32
	var values []float64
	A.Values, A.Row_indx, A.Col_indx, A.Rows, A.Cols = GetCOO_D(A.CSR)
	B.Values, B.Row_indx, B.Col_indx, B.Rows, B.Cols = GetCOO_D(B.CSR)
	values = append(A.Values, B.Values...)
	for i := 0; i < len(B.Col_indx); i++ {
		B.Col_indx[i] += int32(A.Cols)
	}
	rows = append(A.Row_indx, B.Row_indx...)
	cols = append(A.Col_indx, B.Col_indx...)
	c := DSpmatrix{Rows: A.Rows, Cols: A.Cols + B.Cols, Nnz: len(values), Row_indx: rows, Col_indx: cols, Values: values}
	c.Init()
	return c
}

func DSpmatrixstack(A DSpmatrix, B DSpmatrix) DSpmatrix {
	//var C DSpmatrix
	var rows, cols []int32
	var values []float64
	A.Values, A.Row_indx, A.Col_indx, A.Rows, A.Cols = GetCOO_D(A.CSR)
	B.Values, B.Row_indx, B.Col_indx, B.Rows, B.Cols = GetCOO_D(B.CSR)
	values = append(A.Values, B.Values...)
	for i := 0; i < len(B.Row_indx); i++ {
		B.Row_indx[i] += int32(A.Rows)
	}
	rows = append(A.Row_indx, B.Row_indx...)
	cols = append(A.Col_indx, B.Col_indx...)
	c := DSpmatrix{Rows: A.Rows + B.Rows, Cols: A.Cols, Nnz: (len(values)), Row_indx: rows, Col_indx: cols, Values: values}
	c.Init()
	return c
}
