package numeric

//import "C"
import (
	"fmt"
	"testing"
	"unsafe"
)

func TestOperation(t *testing.T) {
	rows:=[]int32{0,1,2,0,2}
	cols:=[]int32{0,1,2,2,0}
	val:=[]complex128{1+3i,1+4i,1+5i,1+1.5i,1+1i}
	m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:5,Row_indx:rows,Col_indx:cols,Values:val}
	m.Init()
	var T ZSpmatrix
	//T.CSR=Operation('C',m.CSR)
	T.CSR=Operation('T',Operation('C',m.CSR))
	PrintZMatrix(T.CSR)
	m.Free()
}

func TestSparseDMV(t *testing.T) {
	rows:=[]int32{0,1,2}
	cols:=[]int32{0,1,2}
	val:=[]float64{3.,4.,5.}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows,Col_indx:cols,Values:val}
	var a *DSpmatrix
	a=&m
	a.Init()
	x:=[]float64{1.,2.,3.}
	y:=[]float64{0.,0.,0.}
	SparseDMV('N',1.0,a.CSR,x,0.0,y)
	PrintDMatrix(m.CSR)
	fmt.Println(y)
	a.Free()
}

func TestSparseZMV(t *testing.T) {
	rows:=[]int32{0,1,2}
	cols:=[]int32{0,1,2}
	val:=[]complex128{3+5i,4+6i,5+7i}
	m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows,Col_indx:cols,Values:val}
	var a *ZSpmatrix
	a=&m
	a.Init()
	x:=[]complex128{1.,2.,3.+1i}
	y:=[]complex128{0.,0.,0.}
	SparseZMV('N',1.0,a.CSR,x,0.0,y)
	fmt.Println(y)
	a.Free()
}

func TestSparseDMM(t *testing.T) {
	rows:=[]int32{0,1,2,0,1}
	cols:=[]int32{0,1,2,2,2}
	val:=[]float64{3,4,5,1.5,1}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:5,Row_indx:rows,Col_indx:cols,Values:val}
	var a *DSpmatrix
	a=&m
	a.Init()
	x:=[]float64{1.,2.,3,0,5,1,3,0,6}
	y:=[]float64{0.,0.,0.,0,0,0,0,0,0}
	SparseDMM('N',1.0,a.CSR,x,3,3,3,0.0,y)
	fmt.Println(y)
	a.Free()
}

func TestSparseZMM(t *testing.T) {
	rows:=[]int32{0,1,2,0,1}
	cols:=[]int32{0,1,2,2,2}
	val:=[]complex128{3+5i,4+6i,5+7i,1+2i,1}
	m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:5,Row_indx:rows,Col_indx:cols,Values:val}
	var a *ZSpmatrix
	a=&m
	a.Init()
	x:=[]complex128{1.,2.,3.+1i,0,5,1,3,0,6}
	y:=[]complex128{0.,0.,0.,0,0,0,0,0,0}
	SparseZMM('N',1.0,a.CSR,x,3,3,3,0.0,y)
	fmt.Println(y)
	a.Free()
}

func TestSparseDAdd(t *testing.T) {
	rows_a:=[]int32{0,1,2}
	cols_a:=[]int32{0,1,2}
	val_a:=[]float64{3,4,5}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	var a *DSpmatrix
	a=&m
	a.Init()
	rows_b:=[]int32{0,1,2}
	cols_b:=[]int32{0,0,2}
	val_b:=[]float64{3,4,10}
	b:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	var C DSpmatrix
	C.CSR=SparseDAdd('N',a.CSR,-1.0,b.CSR,C.CSR)
	C.Values,C.Row_indx,C.Col_indx ,_,_= GetCOO_D(C.CSR)
	PrintDMatrix(a.CSR)
	values_out,rows_idx := GetDMatrixvalue(a.CSR)
	fmt.Println(values_out,rows_idx)
	a.Free()
	b.Free()
	//c.Free()
}

func TestSparseZAdd(t *testing.T) {
	rows_a:=[]int32{0,1,2}
	cols_a:=[]int32{0,1,2}
	val_a:=[]complex128{3+5i,4+6i,5+7i}
	a:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	//var a *ZSpmatrix
	//a=&m
	a.Init()
	rows_b:=[]int32{0,1,2}
	cols_b:=[]int32{0,0,2}
	val_b:=[]complex128{3+5i,4+6i,5+7i}
	b:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	var c ZSpmatrix
	c.CSR=SparseZAdd('N',a.CSR,1.0,b.CSR,c.CSR)
	PrintZMatrix(c.CSR)
	value_out,rows_idx,_ := GetZMatrixvalue(c.CSR)
	fmt.Println(value_out)
	fmt.Println(rows_idx)
	a.Free()
	b.Free()
	c.Free()
}

func TestSparseDSPMM(t *testing.T) {
	rows_a:=[]int32{0,1,2}
	cols_a:=[]int32{0,1,2}
	val_a:=[]float64{3,4,5}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	var a *DSpmatrix
	a=&m
	a.Init()
	rows_b:=[]int32{0,1,2,1}
	cols_b:=[]int32{0,0,2,1}
	val_b:=[]float64{3,4,5,3}
	b:=DSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	var c Spmatrix
	c.CSR=SparseDSPMM('N',a.CSR,b.CSR,c.CSR)
	PrintDMatrix(c.CSR)
	m.Free()
	b.Free()
	c.Free()
}

func TestSparseZSPMM(t *testing.T) {
	rows_a := []int32{0, 1, 2}
	cols_a := []int32{0, 1, 2}
	val_a := []complex128{3 + 5i, 4 + 6i, 5 + 7i}
	m := ZSpmatrix{Rows: 3, Cols: 3, Nnz: 3, Row_indx: rows_a, Col_indx: cols_a, Values: val_a}
	var a *ZSpmatrix
	a = &m
	a.Init()
	rows_b := []int32{0, 1, 2, 1}
	cols_b := []int32{0, 0, 2, 1}
	val_b := []complex128{3 + 5i, 4 + 6i, 5 + 7i, 3}
	b := ZSpmatrix{Rows: 3, Cols: 3, Nnz: 4, Row_indx: rows_b, Col_indx: cols_b, Values: val_b}
	b.Init()
	var c Spmatrix
	c.CSR = SparseZSPMM('N',a.CSR, b.CSR, c.CSR)
	print()
	PrintZMatrix(c.CSR)
	m.Free()
	b.Free()
	c.Free()
}


func TestSparseDSPMMD(t *testing.T) {
	rows_a:=[]int32{0,1,2}
	cols_a:=[]int32{0,1,2}
	val_a:=[]float64{3,4,5}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	var a *DSpmatrix
	a=&m
	a.Init()
	rows_b:=[]int32{0,1,2,1}
	cols_b:=[]int32{0,0,2,1}
	val_b:=[]float64{3,4,5,3}
	b:=DSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	c:=[]float64{0.,0.,0.,0,0,0,0,0,0}
	SparseDSPMMD('N',a.CSR,b.CSR,c,3)
	fmt.Println(c)
	a.Free()
	b.Free()
}

func TestSparseZSPMMD(t *testing.T) {
	rows_a:=[]int32{0,1,2}
	cols_a:=[]int32{0,1,2}
	val_a:=[]complex128{3+5i,4+6i,5+7i}
	m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	var a *ZSpmatrix
	a=&m
	a.Init()
	rows_b:=[]int32{0,1,2,1}
	cols_b:=[]int32{0,0,2,1}
	val_b:=[]complex128{3+5i,4+6i,5+7i,3}
	b:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	c:=[]complex128{0.,0.,0.,0,0,0,0,0,0}
	SparseZSPMMD('N',a.CSR,b.CSR,c,3)
	fmt.Println(c)
	a.Free()
	b.Free()
}

func TestSparseDSP2M(t *testing.T) {
	rows_a:=[]int32{0,1,2}
	cols_a:=[]int32{0,1,2}
	val_a:=[]float64{3,4,5}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	var a *DSpmatrix
	a=&m
	a.Init()
	rows_b:=[]int32{0,1,2,1}
	cols_b:=[]int32{0,0,2,1}
	val_b:=[]float64{3,4,5,3}
	b:=DSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	var c Spmatrix
	c.CSR=SparseDSP2M('N','N',a.CSR,b.CSR,c.CSR)
	PrintDMatrix(c.CSR)
	a.Free()
	b.Free()
	c.Free()
}

func TestSparseZSP2M(t *testing.T) {
	rows_a:=[]int32{0,1,2}
	cols_a:=[]int32{0,1,2}
	val_a:=[]complex128{3+5i,4+6i,5+7i}
	m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	var a *ZSpmatrix
	a=&m
	a.Init()
	rows_b:=[]int32{0,1,2,1}
	cols_b:=[]int32{0,0,2,1}
	val_b:=[]complex128{3+5i,4+6i,5+7i,3}
	b:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	var c Spmatrix
	c.CSR=SparseZSP2M('N','N',a.CSR,b.CSR,c.CSR)
	PrintZMatrix(c.CSR)
	a.Free()
	b.Free()
	c.Free()
}

func TestSparseSetValue(t *testing.T) {
	d_rows:=[]int32{0,1,2}
	d_cols:=[]int32{0,1,2}
	d_val:=[]float64{3.,4.,5.}
	d_m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:d_rows,Col_indx:d_cols,Values:d_val}
	var d_a *DSpmatrix
	d_a=&d_m
	d_a.Init()
	SparseDSetValue(d_a.CSR,1,1,5.5)
	PrintDMatrix(d_m.CSR)
	z_rows:=[]int32{0,1,2}
	z_cols:=[]int32{0,1,2}
	z_val:=[]complex128{3+5i,4+6i,5+7i}
	z_m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:z_rows,Col_indx:z_cols,Values:z_val}
	var z_a *ZSpmatrix
	z_a=&z_m
	z_a.Init()
	SparseZSetValue(z_a.CSR,1,1,5.5+1i)
	PrintZMatrix(z_a.CSR)
	d_a.Free()
	z_a.Free()
}

func TestPardiso(t *testing.T) {
	rows:=[]int32{0,1,2,1,0}
	cols:=[]int32{0,1,2,0,2}
	val:=[]float64{3.,4.,5.,6.,7.}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:5,Row_indx:rows,Col_indx:cols,Values:val}
	a:=&m
	a.Init()
	val[1]=2
	b:=[]float64{10.,10.,5.}
	x:=[]float64{0,0,0}
	_,_,rows_start,_,cols_indx,values:=GetIndexValuesDMatrix(m.CSR)
	var iparm[64]int32
	Pardiso(11,13,3,unsafe.Pointer(values),rows_start,cols_indx,1,iparm,unsafe.Pointer(&b[0]),unsafe.Pointer(&x[0]))
	fmt.Println(x)
	z_rows:=[]int32{0,1,2,1}
	z_cols:=[]int32{0,1,2,0}
	z_val:=[]complex128{3+5i,4+6i,5+7i,3+8i}
	z_m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:z_rows,Col_indx:z_cols,Values:z_val}
	z_a:=&z_m
	z_b:=[]complex128{5,10.,5.}
	z_x:=[]complex128{0,0,0}
	z_a.Init()
	_,_,z_rows_start,_,z_cols_indx,z_values:=GetIndexValuesZMatrix(z_a.CSR)
	Pardiso(13,13,3,unsafe.Pointer(z_values),z_rows_start,z_cols_indx,1,iparm,unsafe.Pointer(&z_b[0]),unsafe.Pointer(&z_x[0]))
	fmt.Println(z_x)
	fmt.Println(5/(3+5i),(10-5/(3+5i)*(3+8i))/(4+6i),5/(5+7i))
	a.Free()
	z_a.Free()
}

func TestZ2D(t *testing.T){
	z_rows:=[]int32{0,1,2,1}
	z_cols:=[]int32{0,1,2,0}
	z_val:=[]complex128{3+5i,4+6i,5+7i,3+8i}
	z_m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:z_rows,Col_indx:z_cols,Values:z_val}
	z_a:=&z_m
	z_a.Init()
	var real_matrix,imag_matrix Spmatrix
	real_matrix.CSR,imag_matrix.CSR=Z2DSpmatrix(z_a.CSR)
	PrintDMatrix(real_matrix.CSR)
	PrintDMatrix(imag_matrix.CSR)
	z_a.Free()
	real_matrix.Free()
	imag_matrix.Free()
}

func TestDInv(t *testing.T) {
	rows:=[]int32{0,1,2,0}
	cols:=[]int32{0,1,2,2}
	val:=[]float64{3.,4.,5.,7.}
	m:=DSpmatrix{Rows: 3,Cols: 3,Nnz:4,Row_indx:rows,Col_indx:cols,Values:val}
	a:=&m
	a.Init()
	x:=[]float64{0,0,0,0,0,0,0,0,0}
	_,_,rows_start,_,cols_indx,values:=GetIndexValuesDMatrix(m.CSR)
	DInv(11,3,values,rows_start,cols_indx,x)
	fmt.Println(x)
}

func TestZInv(t *testing.T) {
	z_rows:=[]int32{0,1,2}
	z_cols:=[]int32{0,1,2}
	z_val:=[]complex128{3+5i,4+6i,5+7i}
	z_m:=ZSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:z_rows,Col_indx:z_cols,Values:z_val}
	z_a:=&z_m
	z_a.Init()
	z_x:=[]complex128{0,0,0,0,0,0,0,0,0}
	_,_,z_rows_start,_,z_cols_indx,z_values:=GetIndexValuesZMatrix(z_a.CSR)
	ZInv(13,3,z_values,z_rows_start,z_cols_indx,z_x)
	fmt.Println(z_x)
}

func TestSparseZVMmul_Ele(t *testing.T) {
	var a,b ZSpmatrix
	rows_a := []int32{0, 1, 2}
	cols_a := []int32{0,0,0}
	val_a := []complex128{3 + 5i, 4 + 6i, 5 + 7i}
	a = ZSpmatrix{Rows: 3, Cols: 1, Nnz: 3, Row_indx: rows_a, Col_indx: cols_a, Values: val_a}
	a.Init()
	rows_b := []int32{0, 1, 2}
	cols_b := []int32{0,0,0}
	val_b := []complex128{ 5i,  6i,  7i}
	b = ZSpmatrix{Rows: 3, Cols: 1, Nnz: 3, Row_indx: rows_b, Col_indx: cols_b, Values: val_b}
	b.Init()
	var C ZSpmatrix
	C = SparseZVMmul_Ele(a,b)
	PrintZMatrix(C.CSR)
	a.Free()
	b.Free()
	C.Free()
}

func TestZSparsemul(t *testing.T) {

	var a ZSpmatrix
	rows_a := []int32{0, 1, 2}
	cols_a := []int32{0, 1, 2}
	val_a := []complex128{3 + 5i, 4 + 6i, 5 + 7i}
	a = ZSpmatrix{Rows: 3, Cols: 3, Nnz: 3, Row_indx: rows_a, Col_indx: cols_a, Values: val_a}
	a.Init()
	rows_b := []int32{0, 1, 2, 1}
	cols_b := []int32{0, 0, 2, 1}
	val_b := []complex128{3 + 5i, 4 + 6i, 5 + 7i, 3}
	b := ZSpmatrix{Rows: 3, Cols: 3, Nnz: 4, Row_indx: rows_b, Col_indx: cols_b, Values: val_b}
	b.Init()
	var C ZSpmatrix
	C = ZSparsemul(a,b)
	PrintZMatrix(C.CSR)
	a.Free()
	b.Free()
	C.Free()
}

func TestSparseZMulEle(t *testing.T) {
	rows_a:=[]int32{0,1,2,3,4,5}
	cols_a:=[]int32{0,0,0,0,0,0}
	val_a:=[]complex128{3,4,5,1+1i,3,7}
	m:=ZSpmatrix{Rows: 6,Cols: 1,Nnz:6,Row_indx:rows_a,Col_indx:cols_a,Values:val_a}
	var a *ZSpmatrix
	a=&m
	a.Init()
	rows_b:=[]int32{0,1,2,3,4,5}
	cols_b:=[]int32{0,0,0,0,0,0}
	val_b:=[]complex128{3+5i,4+6i,5+7i,3,1+1i,2}
	b:=ZSpmatrix{Rows: 6,Cols: 1,Nnz:6,Row_indx:rows_b,Col_indx:cols_b,Values:val_b}
	b.Init()
	var c ZSpmatrix
	c.CSR=SparseZMulEle(a.CSR,b.CSR)
	PrintZMatrix(c.CSR)
	a.Free()
	b.Free()
	c.Free()
}

func TestDSpmatrixAugment(t *testing.T) {
	rows:=[]int32{0,1,2}
	cols:=[]int32{0,1,2}
	val:=[]float64{3,4,5}
	a:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows,Col_indx:cols,Values:val}
	a.Init()
	rows1:=[]int32{0,1,2}
	cols1:=[]int32{0,1,2}
	val1:=[]float64{3,4,5}
	b:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows1,Col_indx:cols1,Values:val1}
	b.Init()
	C := DSpmatrixAugment(a,b)
	C.Init()
	PrintDMatrix(C.CSR)
	a.Free()
	b.Free()
	C.Free()
}

func TestDSpmatrixstack(t *testing.T) {
	rows:=[]int32{0,1,2}
	cols:=[]int32{0,1,2}
	val:=[]float64{3,4,5}
	a:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows,Col_indx:cols,Values:val}
	a.Init()
	rows1:=[]int32{0,1,2}
	cols1:=[]int32{0,1,2}
	val1:=[]float64{3,4,5}
	b:=DSpmatrix{Rows: 3,Cols: 3,Nnz:3,Row_indx:rows1,Col_indx:cols1,Values:val1}
	b.Init()
	C := DSpmatrixstack(a,b)
	C.Init()
	PrintDMatrix(C.CSR)
	a.Free()
	b.Free()
	C.Free()
}
