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

The matrix factorization model produces a 50-dimensional embedding for every anime.

To efficiently retrieve similar anime, these embeddings are stored in Redis Stack using its HNSW (Hierarchical Navigable Small World) vector index.

Why HNSW?

A brute-force search compares a query against every anime:

O(N)

HNSW constructs a multi-layer graph where each anime vector is connected to nearby neighbors. During search, the algorithm traverses this graph to quickly locate the closest vectors without scanning the entire dataset.

Advantages:

Fast approximate nearest-neighbor search
Scales well to large embedding collections
Low query latency
Integrates directly with Redis


