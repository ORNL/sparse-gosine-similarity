package sparsegosinesim

import (
	"context"

	"github.com/bluele/psort"
	"github.com/james-bowman/sparse"
)

type SimilaritySet []SimilarityRow

type SimilarityRow struct {
	Idx    int
	Values []SimilarityValue
}

type SimilarityValue struct {
	Idx int
	S   float64
}

// SparseDotProduct performs efficient sparse matrix multiplication and returns the elements of the resulting array that are above `lowerBound` for each row
// If `onlyResultsAboveDiag` is set to `true` then only the values above the main diagonal of the resulting matrix are calculated and returned, this is more efficient when calculating cosine similarity
//of transposed matrixes because the main diagonal and below will be duplicates and unncessary information.
// Adapted from this Python / C++ code - https://medium.com/wbaa/https-medium-com-ingwbaa-boosting-selection-of-the-most-similar-entities-in-large-scale-datasets-450b3242e618
// All elements above `lowerBound` can be returned by setting `ntop` to 0
func SparseDotProduct(ctx context.Context, A, B *sparse.CSR, onlyResultsAboveDiag bool, lowerBound float64, ntop int) SimilaritySet {

	if ctx == nil {
		ctx = context.Background()
	}

	//Convert to raw to be able to access CSR sparse matrix data structures
	RawA := A.RawMatrix()
	RawB := B.RawMatrix()

	nRow, _ := A.Dims()
	_, nCol := B.Dims()

	//Info regarding CSR format can be found here: https://www.geeksforgeeks.org/sparse-matrix-representations-set-3-csr/
	Ap := RawA.Indptr // stores the cumulative number of non-zero elements upto ( not including) the i-th row
	Aj := RawA.Ind    // stores the column index of each non-zero element in A
	Ax := RawA.Data   // stores the non-zero elements of A

	//Same idea as the above data structs but for B
	Bp := RawB.Indptr
	Bj := RawB.Ind
	Bx := RawB.Data

	//Stores a similarity value which is checked for eligibility against the lowerBound
	type Candidate struct {
		index int
		value float64
	}

	sums := make([]float64, nCol) //sparse vector that records the multiplication result of the current row
	next := make([]int, nCol)     //sparse vector that keeps a linked list of the current row. Every element points to the next column index

	//Initialize the array to all -1
	for i := range next {
		next[i] = -1
	}

	//Loop through each row of A to perform dot product
	var similaritySet SimilaritySet
	for i := 0; i < nRow; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			var candidates []Candidate

			head := -2  //Stores the head of the `next` linked-list slice
			length := 0 //The number of values in the `next` linked-list slice

			//These variables are used as indexes into Aj and Ax which are slices which hold the column numbers and values of non-zero elements
			jjStart := Ap[i] //the cumulative number of non-zero elements in all rows before i
			jjEnd := Ap[i+1] //the cumulative number of non-zero elements in all rows before i and including i

			for jj := jjStart; jj < jjEnd; jj++ {
				j := Aj[jj] //column index of next non-zero element in A[i]
				v := Ax[jj] //value of next non-zero element in A[i]

				//Sum the product of each element of row i of B with element (i,j) of A to get the dot product for this row
				kkStart := Bp[j]
				kkEnd := Bp[j+1]
				for kk := kkStart; kk < kkEnd; kk++ {
					k := Bj[kk]

					//Skip unneeded calculations if only calculating results above the main diagonal of the output matrix
					if onlyResultsAboveDiag && k <= i {
						continue
					}

					sums[k] += v * Bx[kk]

					//Update linked list variables
					if next[k] == -1 {
						next[k] = head
						head = k
						length++
					}
				}
			}

			// Create a new similarity row to hold the data
			similarityRow := &SimilarityRow{
				Idx: i,
			}

			//Traverse all of the resulting values and add all of those above the threshold to the candidate list or directly to the similarityRow if ntop = 0
			for jj := 0; jj < length; jj++ {
				if sums[head] > lowerBound {
					if ntop > 0 {
						c := Candidate{index: head, value: sums[head]}
						candidates = append(candidates, c)
					} else {
						similarityRow.Values = append(similarityRow.Values, SimilarityValue{head, sums[head]})
					}
				}

				temp := head
				head = next[head]

				next[temp] = -1
				sums[temp] = 0
			}

			//Add the top n candidates to the similarityRow
			if ntop > 0 {
				clen := len(candidates)
				psort.Slice(candidates, func(i, j int) bool {
					return candidates[i].value > candidates[j].value
				}, ntop)
				if clen > ntop {
					clen = ntop
				}

				for a := 0; a < clen; a++ {
					similarityRow.Values = append(similarityRow.Values, SimilarityValue{candidates[a].index, candidates[a].value})
				}
			}

			//If the similarityRow is not empty add it to the similaritySet
			if len(similarityRow.Values) > 0 {
				similaritySet = append(similaritySet, *similarityRow)
			}
		}
	}
	return similaritySet
}
