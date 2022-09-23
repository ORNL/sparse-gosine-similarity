package sparsegosinesim

import (
	"context"
	"math"
	"testing"

	"github.com/james-bowman/sparse"
)

type CosineSimilarityTest struct {
	matrix1    *sparse.CSR
	matrix2    *sparse.CSR
	ntop       int
	lowerBound float64

	expectedSimsMap map[int]map[int]float64
}

var CosineSimilarityTests = []CosineSimilarityTest{
	{
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 1, 2, 3, 1, 2, 4},
		),
		nil,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				1: 1.000,
				2: 0.991,
			},
			1: {
				2: 0.991,
			},
		},
	},
	{
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
		),
		nil,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				1: 0.975,
				2: 0.959,
			},
			1: {
				2: 0.998,
			},
		},
	},
	{
		sparse.NewCSR(
			2, 3,
			[]int{0, 3, 6},
			[]int{0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6},
		),
		sparse.NewCSR(
			3, 2,
			[]int{0, 2, 4, 6},
			[]int{0, 1, 0, 1, 0, 1},
			[]float64{4, 7, 5, 8, 6, 9},
		),
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				0: 0.975,
				1: 0.959,
			},
			1: {
				0: 1.000,
				1: 0.998,
			},
		},
	},
	{
		sparse.NewCSR(
			2, 3,
			[]int{0, 3, 6},
			[]int{0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6},
		),
		sparse.NewCSR(
			3, 2,
			[]int{0, 2, 4, 6},
			[]int{0, 1, 0, 1, 0, 1},
			[]float64{4, 7, 5, 8, 6, 9},
		),
		1,
		0.0,
		map[int]map[int]float64{
			0: {
				0: 0.975,
			},
			1: {
				0: 1.000,
			},
		},
	},
	{
		sparse.NewCSR(
			2, 3,
			[]int{0, 3, 6},
			[]int{0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6},
		),
		sparse.NewCSR(
			3, 2,
			[]int{0, 2, 4, 6},
			[]int{0, 1, 0, 1, 0, 1},
			[]float64{4, 7, 5, 8, 6, 9},
		),
		0,
		0.98,
		map[int]map[int]float64{
			1: {
				0: 1.000,
				1: 0.998,
			},
		},
	},
	{
		sparse.NewCSR(
			2, 3,
			[]int{0, 3, 6},
			[]int{0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6},
		),
		sparse.NewCSR(
			3, 2,
			[]int{0, 2, 4, 6},
			[]int{0, 1, 0, 1, 0, 1},
			[]float64{4, 7, 5, 8, 6, 9},
		),
		1,
		0.98,
		map[int]map[int]float64{
			1: {
				0: 1.000,
			},
		},
	},
}

func TestCosineSimilarity(t *testing.T) {
	t.Parallel()
	for ti, tt := range CosineSimilarityTests {
		actualSimSet, err := CosineSimilarity(context.Background(), tt.matrix1, tt.matrix2, tt.lowerBound, tt.ntop)
		if err != nil {
			t.Errorf("TestCosineSimilarity(Test%d): error while computing cosine similarity: %e", ti, err)
		}
		if len(tt.expectedSimsMap) != len(actualSimSet) {
			t.Errorf("TestCosineSimilarity(Test%d): sim sets are not the same length. expected length: (%d) actual length: (%d)", ti, len(tt.expectedSimsMap), len(actualSimSet))
		}
		for _, simRow := range actualSimSet {
			for _, simVal := range simRow.Values {
				if expectedSimRowValMap, exists := tt.expectedSimsMap[simRow.Idx]; exists {
					if len(simRow.Values) != len(expectedSimRowValMap) {
						t.Errorf("TestCosineSimilarity(Test%d): sim rows are not the same length (row: %d). expected: %d actual %d", ti, simRow.Idx, len(expectedSimRowValMap), len(simRow.Values))
					}
					if expectedScore, exists := expectedSimRowValMap[simVal.Idx]; exists {
						actualScore := math.Round(simVal.S*1000) / 1000
						if expectedScore != actualScore {
							t.Errorf("TestCosineSimilarity(Test%d): unexpected score at (%d, %d). expected: %f actual %f", ti, simRow.Idx, simVal.Idx, expectedScore, actualScore)
						}
					} else {
						t.Errorf("TestCosineSimilarity(Test%d): unexpected column in row. row: %d column: %d", ti, simRow.Idx, simVal.Idx)
					}
				} else {
					t.Errorf("TestCosineSimilarity(Test%d): unexpected row index (%d)", ti, simRow.Idx)
				}
			}
		}
	}
}
