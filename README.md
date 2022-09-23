# sparse-gosine-similarity

**sparse-gosine-similarity** provides a fast way to perform a sparse matrix multiplication followed by top-n multiplication result selection as well as functionality for using the matrix multiplication to calculate cosine similarity. This package is a pure Go port of a [Python package](https://github.com/ing-bank/sparse_dot_topn) developed by ING Bank which uses Cython to execute the matrix multiplication in C++. Additional details about the implementation and efficiency of the algorithm as well as blog posts describing it can be found at the [ing-bank/sparse_dot_topn](https://github.com/ing-bank/sparse_dot_topn) repository.

## Functions

- SparseDotProduct - Performs efficient matrix multiplication by multipling each row of matrix A against each column of matrix B
- L2Normalization - Normalizes the provided matrix so that in each row the sum of squares will always add up to 1. This can be utilized to calculate the cosine similarity by only performing a dot product calculation on the normalized matrix.
- CosineSimilarity - Given two matrices this function performs L2 normalization on the rows of matrix A and the columns of matrix B and then computes the dot product which will then be equivalent to the cosine similarity. Given one matrix the rows of the given matrix are normalized and the multiplication is perform on matrix A and its transposition.

## Python 3 Compatability

This package has been designed so it can be run on Python using [gopy](https://github.com/go-python/gopy). Once the bindings have been generated they can be imported and used in Python:

```py
from gobindings import sparse-cossim
sparse-cossim.X
```

### Automated Binding Generation

The `generate-bindings.sh` script will generate the bindings needed to run the package in Python by creating a `build` folder where Go will be downloaded, a Python virtual environment will be created, and all needed dependencies will be installed. This script enables generation of the bindings without installing Go at the system level. Python 3 must be installed and accessible via `python3` for the script to work.

### Manual Binding Generation

Assuming you already have Go installed you can also install `gopy` manually by doing:

```sh
python3 -m pip install pybindgen
go get golang.org/x/tools/cmd/goimports
go install golang.org/x/tools/cmd/goimports
go get github.com/go-python/gopy
go install github.com/go-python/gopy
python3 -m pip install --upgrade setuptools wheel
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:.
```

Once gopy is installed you can manually generate the bindings for this package by doing:

```sh
go get github.com/ORNL/sparse-gosine-similarity
gopy build -output=gobindings -vm=python3 github.com/ORNL/sparse-gosine-similarity
```

## Go example

```go
package main

import (
	"context"

	sgs "github.com/ORNL/sparse-gosine-similarity"
	"github.com/james-bowman/sparse"
)

func main() {
	n := 10

	//Calculate dot product only
	a := sparse.Random(sparse.CSRFormat, 100, 1000000, 0.01).(*sparse.CSR)
	b := sparse.Random(sparse.CSRFormat, 1000000, 200, 0.01).(*sparse.CSR)

	_ = sgs.SparseDotProduct(context.Background(), a, b, false, 0.01, n)

	//Calculate cosine similarity by manually normalizing
	a = sparse.Random(sparse.CSRFormat, 100, 1000000, 0.005).(*sparse.CSR)
	b = sparse.Random(sparse.CSRFormat, 1000000, 200, 0.005).(*sparse.CSR)

	sgs.L2Normalize(a)

	sgs.Transpose(b)
	sgs.L2Normalize(b)
	sgs.Transpose(b)

	_ = sgs.SparseDotProduct(context.Background(), a, b, false, 0.001, n)

	//Calculate cosine similarity directly
	a = sparse.Random(sparse.CSRFormat, 100, 1000000, 0.005).(*sparse.CSR)
	b = sparse.Random(sparse.CSRFormat, 1000000, 200, 0.005).(*sparse.CSR)

	_, _ = sgs.CosineSimilarity(context.Background(), a, b, 0.001, n)
}
```

## Python example using Go bindings
```py
from gobindings import sparsegosinesim, go

#Define the dimensions of a CSR matrix
ia=[0,3,6]
ja=[0,1,2,0,1,2]
data=[1,2,3,4,5,6]

#Create Go versions of the CSR matrices
mat=sparsegosinesim.MakeCSR(2,3,go.Slice_int(ia), go.Slice_int(ja),go.Slice_float64(data))
mat2=sparsegosinesim.MakeCSR(2,3,go.Slice_int(ia), go.Slice_int(ja),go.Slice_float64(data))

#Transpose one and calculate all cosine similarities
sparsegosinesim.Transpose(mat2)
ctx=go.context_Context()
simSet=sparsegosinesim.CosineSimilarity(ctx, mat, mat2, 0.0, 0)

#Print the results
for i in range(len(simSet)):
	for j in range(len(simSet[i].Values)):
		print(simSet[i].Idx, simSet[i].Values[j].Idx, simSet[i].Values[j].S)
	
# Alternatively, compare mat to itself and only print results above the resulting matrix's main diagonal
nil=go.nil
simSet=sparsegosinesim.CosineSimilarity(ctx, mat, nil, 0.0, 0)

for i in range(len(simSet)):
	for j in range(len(simSet[i].Values)):
		print(simSet[i].Idx, simSet[i].Values[j].Idx, simSet[i].Values[j].S)
```

## Testing

There are tests for each go file provided. These tests can be run by doing

```sh
go test -v ./...
```