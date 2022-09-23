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
