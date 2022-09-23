package sparsegosinesim

import (
	"math"

	"github.com/james-bowman/sparse"
)

//L2Normalize - Normalizes a matrix so in each row the sum of the squares is equal to 1
func L2Normalize(matrix *sparse.CSR) {
	rawMatrix := matrix.RawMatrix()
	for i := 0; i < rawMatrix.I; i++ {
		sum := 0.0

		for j := rawMatrix.Indptr[i]; j < rawMatrix.Indptr[i+1]; j++ {
			sum += rawMatrix.Data[j] * rawMatrix.Data[j]
		}
		if sum == 0.0 {
			continue
		}
		sum = math.Sqrt(sum)
		for j := rawMatrix.Indptr[i]; j < rawMatrix.Indptr[i+1]; j++ {
			rawMatrix.Data[j] /= sum
		}
	}
}

//Transpose is a helper function for easily transposing a sparse matrix
func Transpose(matrix *sparse.CSR) {
	matrix.Clone(matrix.T().(*sparse.CSC).ToCSR())
}

//MakeCSR exposes sparse.CSR in this package, making it easier to create a CSR when using Python bindings
func MakeCSR(r int, c int, ia []int, ja []int, data []float64) *sparse.CSR {
	return sparse.NewCSR(r, c, ia, ja, data)
}
