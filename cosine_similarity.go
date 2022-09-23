package sparsegosinesim

import (
	"context"
	"fmt"

	"github.com/james-bowman/sparse"
)

// CosineSimilarity returns a struct representing the cosine similarity of two matrices with each row containing the top N results. Every row in matrix 1 will be compared to every column in matrix 2. If matrix2 is nil matrix1 will be transposed for the calculation.
// All elements above `lowerBound` can be returned by setting `ntop` to 0
func CosineSimilarity(ctx context.Context, matrix1, matrix2 *sparse.CSR, lowerBound float64, ntop int) (SimilaritySet, error) {
	//Ensure at least 1 matrix was provided
	if matrix1 == nil {
		return nil, fmt.Errorf("first matrix cannot be nil")
	}
	//Normalize the first matrix
	L2Normalize(matrix1)

	//If matrix2 is nil clone matrix2 into matrix1
	equivalent := false
	if matrix2 == nil {
		//Create a clone of matrix1
		r, c := matrix1.Dims()
		matrix2 = sparse.NewCSR(r, c, []int{}, []int{}, []float64{})
		matrix2.Clone(matrix1)

		equivalent = true
	} else {
		//Transpose matrix2 to perform L2 normalization
		Transpose(matrix2)
		L2Normalize(matrix2)
	}

	//Transpose matrix2 so it is compatible for the sparse dot product calculation
	Transpose(matrix2)

	//Return the sparse dot product
	return SparseDotProduct(ctx, matrix1, matrix2, equivalent, lowerBound, ntop), nil
}
