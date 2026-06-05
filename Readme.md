# Anime Recommender Training

## Overview
Matrix factorization model for anime recommendations using user–anime ratings.

Prediction:
r̂ = μ + bu + bi + Wu · Hi

## Data
- anime.csv → anime metadata
- rating.csv → user ratings

## Training
- Shuffle data each epoch
- Mini-batch SGD
- Update:
  - user vectors (W)
  - anime vectors (H)
  - biases (bu, bi)

## Params
- k = 50
- lr = 0.01
- reg = 0.02
- batch = 100k

## Loss
MSE (rating prediction error)

## Metrics
- RMSE (error)
- NDCG@10 (ranking quality)

## Save
W.npy, H.npy, bu.npy, bi.npy after each epoch

## Note
Collaborative filtering → learns user behavior, not genre similarity.

## RESULTS
```W: (73517, 50)
H: (50, 34528)
NDCG@10: 0.8552502296775826
```


## Vector Search 
for efficeintly finding top similar vectors instead of brute force evaluation we will be using 
Redis hnsw feature 
for stroing and querying efficently 



