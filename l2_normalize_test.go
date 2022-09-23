package sparsegosinesim

import (
	"math"
	"testing"

	"github.com/james-bowman/sparse"
)

type L2NormalizationTest struct {
	matrix                   *sparse.CSR
	expectedNormalizedMatrix *sparse.CSR
}

var L2NormalizationTests = []L2NormalizationTest{
	{
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
		),
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{0.267, 0.535, 0.802, 0.456, 0.570, 0.684, 0.503, 0.574, 0.646},
		),
	},
}

func TestL2Normalization(t *testing.T) {
	t.Parallel()
	for ti, tt := range L2NormalizationTests {
		L2Normalize(tt.matrix)
		r, c := tt.matrix.Dims()
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				actual := math.Round(tt.matrix.At(i, j)*1000) / 1000
				expected := tt.expectedNormalizedMatrix.At(i, j)
				if actual != expected {
					t.Errorf("TestL2Normalization(test %d): unexpected value at (%d,%d) expected: %f actual: %f", ti, i, j, expected, actual)
				}
			}
		}
	}
}
