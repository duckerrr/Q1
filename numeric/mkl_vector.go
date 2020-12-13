package numeric

/*串行mkl将-lmkl_intel_thread改为-lmkl_sequential同时去掉-liomp5 -lpthread即可*/

/*
#cgo CFLAGS: -I/opt/intel/mkl/include
#cgo LDFLAGS: -L/opt/intel/mkl/lib/intel64 -L/opt/intel/lib/intel64 -lmkl_intel_lp64 -lmkl_sequential -lmkl_core -lm
#include<mkl.h>
*/
import "C"
import "unsafe"

/*
向量运算
*/

//实数加法
func VdAdd(n int, a, b, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.double)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdAdd(ctype_n, ctype_a, ctype_b, ctype_y)
}

//复数加法
func VzAdd(n int, a, b, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.MKL_Complex16)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzAdd(ctype_n, ctype_a, ctype_b, ctype_y)
}

//实数减法
func VdSub(n int, a, b, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.double)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdSub(ctype_n, ctype_a, ctype_b, ctype_y)
}

//复数减法
func VzSub(n int, a, b, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.MKL_Complex16)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzSub(ctype_n, ctype_a, ctype_b, ctype_y)
}

//实数乘法
func VdMul(n int, a, b, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.double)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdMul(ctype_n, ctype_a, ctype_b, ctype_y)
}

//复数乘法
func VzMul(n int, a, b, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.MKL_Complex16)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzMul(ctype_n, ctype_a, ctype_b, ctype_y)
}

//复数乘法(向量a乘向量b的共轭)
func VzMulByConj(n int, a, b, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.MKL_Complex16)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzMulByConj(ctype_n, ctype_a, ctype_b, ctype_y)
}

//向量a的共轭
func VzConj(n int, a, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzConj(ctype_n, ctype_a, ctype_y)
}

//实数绝对值
func VdAbs(n int, a, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdAbs(ctype_n, ctype_a, ctype_y)
}

//复数绝对值
func VzAbs(n int, a []complex128, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vzAbs(ctype_n, ctype_a, ctype_y)
}

//复数的相角
func VzArg(n int, a []complex128, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vzArg(ctype_n, ctype_a, ctype_y)
}

//实数向量的倒数
func VdInv(n int, a, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdInv(ctype_n, ctype_a, ctype_y)
}

//实数除法
func VdDiv(n int, a, b, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.double)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdDiv(ctype_n, ctype_a, ctype_b, ctype_y)
}

//复数除法
func VzDiv(n int, a, b, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_b := (*C.MKL_Complex16)(unsafe.Pointer(&b[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzDiv(ctype_n, ctype_a, ctype_b, ctype_y)
}

//实数开方
func VdSqrt(n int, a, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdSqrt(ctype_n, ctype_a, ctype_y)
}

//复数开方
func VzSqrt(n int, a, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzSqrt(ctype_n, ctype_a, ctype_y)
}

//实数向量的倒数再开方
func VdInvSqrt(n int, a, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdInvSqrt(ctype_n, ctype_a, ctype_y)
}

//实数指数e^x
func VdExp(n int, a, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdExp(ctype_n, ctype_a, ctype_y)
}

//复数指数e^x
func VzExp(n int, a, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzExp(ctype_n, ctype_a, ctype_y)
}

//实数余弦
func VdCos(n int, a, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdCos(ctype_n, ctype_a, ctype_y)
}

//复数余弦
func VzCos(n int, a, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzCos(ctype_n, ctype_a, ctype_y)
}

//实数正弦
func VdSin(n int, a, y []float64) {
	ctype_n := C.int(n)
	ctype_a := (*C.double)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.double)(unsafe.Pointer(&y[0]))
	C.vdSin(ctype_n, ctype_a, ctype_y)
}

//复数正弦
func VzSin(n int, a, y []complex128) {
	ctype_n := C.int(n)
	ctype_a := (*C.MKL_Complex16)(unsafe.Pointer(&a[0]))
	ctype_y := (*C.MKL_Complex16)(unsafe.Pointer(&y[0]))
	C.vzSin(ctype_n, ctype_a, ctype_y)
}
