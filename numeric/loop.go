package numeric

// 构造复数向量
func ComplexVector(re []float64, im []float64) []complex128 {
	vec := make([]complex128, len(re))
	for i, _ := range re {
		vec[i] = complex(re[i], im[i])
	}
	return vec
}
