package sparsegosinesim

import (
	"context"
	"math"
	"testing"

	"github.com/james-bowman/sparse"
)

type DotProductTest struct {
	matrix1              *sparse.CSR
	matrix2              *sparse.CSR
	onlyResultsAboveDiag bool
	ntop                 int
	lowerBound           float64

	expectedSimsMap map[int]map[int]float64
}

var DotProductTests = []DotProductTest{
	{
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 1, 2, 3, 1, 2, 4},
		),
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 1, 1, 2, 2, 2, 3, 3, 4},
		),
		true,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				1: 14,
				2: 17,
			},
			1: {
				2: 17,
			},
		},
	},
	{
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 1, 2, 3, 1, 2, 4},
		),
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 1, 1, 2, 2, 2, 3, 3, 4},
		),
		false,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				0: 14,
				1: 14,
				2: 17,
			},
			1: {
				0: 14,
				1: 14,
				2: 17,
			},
			2: {
				0: 17,
				1: 17,
				2: 21,
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
		sparse.NewCSR(
			3, 3,
			[]int{0, 3, 6, 9},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 4, 7, 2, 5, 8, 3, 6, 9},
		),
		true,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				1: 32,
				2: 50,
			},
			1: {
				2: 122,
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
		false,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				0: 32,
				1: 50,
			},
			1: {
				0: 77,
				1: 122,
			},
		},
	},
	{
		sparse.NewCSR(
			5, 3,
			[]int{0, 3, 6, 9, 12, 15},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		sparse.NewCSR(
			3, 5,
			[]int{0, 5, 10, 15},
			[]int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2, 3, 4},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		false,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				0: 46,
				1: 52,
				2: 58,
				3: 64,
				4: 70,
			},
			1: {
				0: 100,
				1: 115,
				2: 130,
				3: 145,
				4: 160,
			},
			2: {
				0: 154,
				1: 178,
				2: 202,
				3: 226,
				4: 250,
			},
			3: {
				0: 208,
				1: 241,
				2: 274,
				3: 307,
				4: 340,
			},
			4: {
				0: 262,
				1: 304,
				2: 346,
				3: 388,
				4: 430,
			},
		},
	},
	{
		sparse.NewCSR(
			5, 3,
			[]int{0, 3, 6, 9, 12, 15},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		sparse.NewCSR(
			3, 5,
			[]int{0, 5, 10, 15},
			[]int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2, 3, 4},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		true,
		0,
		0.0,
		map[int]map[int]float64{
			0: {
				1: 52,
				2: 58,
				3: 64,
				4: 70,
			},
			1: {
				2: 130,
				3: 145,
				4: 160,
			},
			2: {
				3: 226,
				4: 250,
			},
			3: {
				4: 340,
			},
		},
	},
	{
		sparse.NewCSR(
			5, 3,
			[]int{0, 3, 6, 9, 12, 15},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		sparse.NewCSR(
			3, 5,
			[]int{0, 5, 10, 15},
			[]int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2, 3, 4},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		false,
		3,
		0.0,
		map[int]map[int]float64{
			0: {
				2: 58,
				3: 64,
				4: 70,
			},
			1: {
				2: 130,
				3: 145,
				4: 160,
			},
			2: {
				2: 202,
				3: 226,
				4: 250,
			},
			3: {
				2: 274,
				3: 307,
				4: 340,
			},
			4: {
				2: 346,
				3: 388,
				4: 430,
			},
		},
	},
	{
		sparse.NewCSR(
			5, 3,
			[]int{0, 3, 6, 9, 12, 15},
			[]int{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 2},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		sparse.NewCSR(
			3, 5,
			[]int{0, 5, 10, 15},
			[]int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2, 3, 4},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		),
		false,
		3,
		150,
		map[int]map[int]float64{
			1: {
				4: 160,
			},
			2: {
				2: 202,
				3: 226,
				4: 250,
			},
			3: {
				2: 274,
				3: 307,
				4: 340,
			},
			4: {
				2: 346,
				3: 388,
				4: 430,
			},
		},
	},
}

func TestDotProduct(t *testing.T) {
	t.Parallel()
	for ti, tt := range DotProductTests {
		actualSimSet := SparseDotProduct(context.Background(), tt.matrix1, tt.matrix2, tt.onlyResultsAboveDiag, tt.lowerBound, tt.ntop)
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
