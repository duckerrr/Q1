package numeric

/*
#cgo linux CFLAGS: -O3 -I/home/hjm/文档/Go_test/frame/include
#cgo linux LDFLAGS: -L/home/hjm/文档/Go_test/frame/lib
#cgo linux LDFLAGS: -lklu -lamd -lcholmod -lcolamd -lsuitesparseconfig -lbtf
#include<klu.h>
 */
import "C"
import "unsafe"

type KLUZSolver struct {
	N int
	Common C.klu_common
	Symbolic *C.klu_symbolic
	Numeric *C.klu_numeric
}

func (Solver *KLUZSolver)GetLU(n int,Ap,Ai *int32,Ax *complex128){
	Solver.N=n
	C.klu_defaults(&Solver.Common)
	Solver.Common.btf=0
	Solver.Common.scale=0
	Solver.Symbolic=C.klu_analyze(C.int(n),(*C.int)(unsafe.Pointer(Ap)),(*C.int)(unsafe.Pointer(Ai)),&Solver.Common)
	Solver.Numeric=C.klu_z_factor((*C.int)(unsafe.Pointer(Ap)),(*C.int)(unsafe.Pointer(Ai)),(*C.double)(unsafe.Pointer(Ax)),Solver.Symbolic,&Solver.Common)
}

func (Solver *KLUZSolver)Solve(d,nrhs int,B []complex128) []complex128{
	X:=make([]complex128,len(B))
	copy(X,B)
	C.klu_z_solve(Solver.Symbolic,Solver.Numeric,C.int(d),C.int(nrhs),(*C.double)(unsafe.Pointer(&X[0])),&Solver.Common)
	return X
}

func (Solver *KLUZSolver)GetLciUic(Yic,Yci *ZSpmatrix) ([]complex128,[]complex128,[]ZSpmatrix){
	if Yic.Rows<=0||Yic.Cols<=0{
		Yic.GetCOO_Z()
		panic(Yic.Rows>0&&Yic.Cols>0)
	}
	if Yci.Rows<=0||Yci.Cols<=0{
		Yci.GetCOO_Z()
		panic(Yci.Rows>0&&Yci.Cols>0)
	}
	YicOnes:=make([]complex128,Yic.Cols*Yic.Cols)
	YciOnes:=make([]complex128,Yci.Cols*Yci.Cols)
	for i:=0;i<Yic.Cols;i++{
		YicOnes[(Yic.Cols+1)*i]=1
	}
	for i:=0;i<Yci.Cols;i++{
		YciOnes[(Yci.Cols+1)*i]=1
	}
	YicVector:=make([]complex128,Yic.Rows*Yic.Cols)
	YciVector:=make([]complex128,Yci.Rows*Yci.Cols)
	SparseZMM('N',1.0,Yic.CSR,YicOnes,Yic.Cols,Yic.Rows,Yic.Cols,1.0,YicVector)
	SparseZMM('N',1.0,Yci.CSR,YciOnes,Yci.Cols,Yci.Rows,Yci.Cols,1.0,YciVector)
	UicVector:=make([]complex128,Yic.Rows*Yic.Cols)
	LciVector:=make([]complex128,Yci.Rows*Yci.Cols)
	C.klu_z_get_Uic(Solver.Symbolic,Solver.Numeric,(*C.double)(unsafe.Pointer(&YicVector[0])),C.int(Yic.Rows),C.int(Yic.Cols),&Solver.Common,(*C.double)(unsafe.Pointer(&UicVector[0])))
	C.klu_z_get_Lci(Solver.Symbolic,Solver.Numeric,(*C.double)(unsafe.Pointer(&YciVector[0])),C.int(Yci.Cols),C.int(Yci.Rows),&Solver.Common,(*C.double)(unsafe.Pointer(&LciVector[0])))
	var UicRowIndx,UicColIndx,LciRowIndx,LciColIndx []int32
	var UicValue,LciValue []complex128
	var UicNnz,LciNnz int
	for i:=0;i<Yic.Rows;i++{
		for j:=0;j<Yic.Cols; j++{
			if UicVector[i+Yic.Rows*j]!=0{
				UicRowIndx=append(UicRowIndx,int32(i))
				UicColIndx=append(UicColIndx,int32(j))
				UicValue=append(UicValue,UicVector[i+Yic.Rows*j])
				UicNnz++
			}
			if LciVector[j+Yic.Cols*i]!=0{
				LciRowIndx=append(LciRowIndx,int32(j))
				LciColIndx=append(LciColIndx,int32(i))
				LciValue=append(LciValue,LciVector[j+Yic.Cols*i])
				LciNnz++
			}
		}
	}
	UicLci:=make([]ZSpmatrix,2)
	UicLci[0]=ZSpmatrix{Rows:Yic.Rows,Cols:Yic.Cols,Nnz:UicNnz,Row_indx:UicRowIndx,Col_indx:UicColIndx,Values:UicValue}
	UicLci[0].Init()
	UicLci[1]=ZSpmatrix{Rows:Yci.Rows,Cols:Yci.Cols,Nnz:LciNnz,Row_indx:LciRowIndx,Col_indx:LciColIndx,Values:LciValue}
	UicLci[1].Init()
	return UicVector,LciVector,UicLci
}

func (Solver *KLUZSolver)GetWi(I []complex128,nrhs int) []complex128{
	Wi:=make([]complex128,Solver.N)
	C.klu_z_get_Uic(Solver.Symbolic,Solver.Numeric,(*C.double)(unsafe.Pointer(&I[0])),C.int(Solver.N),C.int(nrhs),&Solver.Common,(*C.double)(unsafe.Pointer(&Wi[0])))
	return Wi
}

func (Solver *KLUZSolver)SolveVi(Wi []complex128,nrhs int)[]complex128{
	Vi:=make([]complex128,Solver.N)
	C.klu_z_solve_Vi(Solver.Symbolic,Solver.Numeric,(*C.double)(unsafe.Pointer(&Wi[0])),C.int(Solver.N),C.int(nrhs),&Solver.Common,(*C.double)(unsafe.Pointer(&Vi[0])))
	return Vi
}

/*  这个函数可能不需要
func (Solver *KLUZSolver)PermutationLciUic(Lci,Uic []complex128) ([]complex128,[]complex128){
	PermutationLci:=make([]complex128,len(Lci))
	PermutationUic:=make([]complex128,len(Uic))
	var oldrow,oldcol int
	for i:=0;i<Solver.N;i++{
		oldrow=int(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(Solver.Numeric.Pnum))+4*uintptr(i))))
		oldcol=int(*(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(Solver.Symbolic.Q))+4*uintptr(i))))
		for j:=0;j<len(Lci)/Solver.N;j++{
			PermutationLci[i+Solver.N*j]=Lci[oldrow+Solver.N*j]
			PermutationUic[j+len(Lci)/Solver.N*i]=Uic[j+len(Lci)/Solver.N*oldcol]
		}
	}
	return PermutationUic,PermutationLci
}
 */

func (Solver *KLUZSolver)Free() {
	C.klu_free_symbolic(&Solver.Symbolic,&Solver.Common)
	C.klu_z_free_numeric(&Solver.Numeric,&Solver.Common)
}

func (admit *ZSpmatrix)GetZLU() (*ZSpmatrix,*ZSpmatrix,[]int32,[]int32){
	var Y,L,U,diagL,diagU ZSpmatrix
	var Li,Lp,Ui,Up *C.int
	var Lx,Ux,Udiag *C.double
	Y.CSR=Operation('T',admit.CSR)
	n,_,Ap,_,Ai,Ax:=GetIndexValuesZMatrix(Y.CSR)
	P,Q:=make([]int32,int(n)),make([]int32,int(n))
	var Common C.klu_common
	C.klu_defaults(&Common)
	Common.btf=0
	Common.scale=0
	Symbolic:=C.klu_analyze(n,Ap,Ai,&Common)
	C.klu_z_get_LU(Ap,Ai,(*C.double)(unsafe.Pointer(Ax)),Symbolic,&Common,(*C.int)(unsafe.Pointer(&P[0])),(*C.int)(unsafe.Pointer(&Q[0])),&Li,&Lp,&Ui,&Up,&Lx,&Ux,&Udiag)
	C.klu_free_symbolic(&Symbolic,&Common)
	SparseZCreateCSR(&L.CSR,int(n),int(n),(*int32)(unsafe.Pointer(Lp)),(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(Lp))+4)),(*int32)(unsafe.Pointer(Li)),(*complex128)(unsafe.Pointer(Lx)))
	SparseZCreateCSR(&U.CSR,int(n),int(n),(*int32)(unsafe.Pointer(Up)),(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(Up))+4)),(*int32)(unsafe.Pointer(Ui)),(*complex128)(unsafe.Pointer(Ux)))
	L.CSR=Operation('T',L.CSR)
	U.CSR=Operation('T',U.CSR)
	diagRowStart,diagColIndx,diagValue:=make([]int32,int(n)+1),make([]int32,int(n)),make([]complex128,int(n))
	diagRowStart[0]=0
	for i:=0;i<int(n);i++{
		diagRowStart[i+1],diagColIndx[i],diagValue[i]=int32(i+1),int32(i),1
	}
	SparseZCreateCSR(&diagL.CSR,int(n),int(n),&diagRowStart[0],&diagRowStart[1],&diagColIndx[0],&diagValue[0])
	SparseZCreateCSR(&diagU.CSR,int(n),int(n),&diagRowStart[0],&diagRowStart[1],&diagColIndx[0],(*complex128)(unsafe.Pointer(Udiag)))
	L.CSR=SparseZAdd('N',L.CSR,1.0,diagL.CSR,L.CSR)
	U.CSR=SparseZAdd('N',U.CSR,1.0,diagU.CSR,U.CSR)
	return &L,&U,P,Q
}




