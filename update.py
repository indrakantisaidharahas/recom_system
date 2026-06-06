import pandas as pd 
import numpy as np
from collections import defaultdict


df1=pd.read_csv("/home/saidharahas/jupyter_projects/anime_reco/anime_data/anime.csv")
df2=pd.read_csv("/home/saidharahas/jupyter_projects/anime_reco/anime_data/rating.csv")


##batch incrementation of ratings training 
t_users=df2["user_id"].nunique()
print(t_users)
t_anime=df1["anime_id"].nunique()


t_anime= df1["anime_id"].max()
print(t_anime)

k=10
W = np.random.normal(0, 0.01, (t_users+1, k))
H = np.random.normal(0, 0.01, (k, t_anime+1))
bi=np.random.normal(0, 0.01,  t_anime+1)
bu=np.random.normal(0, 0.01, t_users+1)
lr = 0.05

reg = 0.02
g_mean=df2["rating"].mean()

def update(st: int, bs: int):

    ldf = df2.iloc[st:st + bs]

    for epochs in range(100):

        total_err = 0

        for row in ldf.itertuples():

            j = row.anime_id
            i = row.user_id
            r = row.rating

            if j > 73515:
                continue

            if i > 12294:
                continue

            pred = g_mean + bu[i] + bi[j] + W[i] @ H[:, j]
            err = r - pred

            total_err += err * err

            old_w = W[i].copy()

            W[i] += lr * (err * H[:, j] - reg * W[i])
            H[:, j] += lr * (err * old_w - reg * H[:, j])
            bu[i] += lr * (err - reg * bu[i])
            bi[j] += lr * (err - reg * bi[j])

        print(f"Epoch {epochs + 1}")
        print("Loss:", total_err)
        print("-" * 30)
            
def rmse(st :int):
                 ldf=df2.iloc[st:st+1000]
                 rmse=0
                 for row in ldf.itertuples():
                        
                                j = row.anime_id
                                i = row.user_id
                                r = row.rating

                                if j > 73515:
                                    continue

                                if i > 12294:
                                    continue

                                pred = g_mean + bu[i] + bi[j] + W[i] @ H[:, j]
                                err = r - pred

                                rmse+= err * err
                 rmse=rmse/1000
                 rmse=rmse**0.5 
                 print(rmse)


def ndcg_at_k(df, k=10):
    user_items = defaultdict(list)

    # Collect predictions and true ratings
    for row in df.itertuples():
        i = row.user_id
        j = row.anime_id
        r = row.rating

        if j > 73515:
            continue

        if i > 12294:
            continue

        pred = g_mean + bu[i] + bi[j] + W[i] @ H[:, j]

        user_items[i].append((pred, r))

    ndcg_scores = []

    for user, items in user_items.items():

        if len(items) < 2:
            continue

        # Sort by predicted score
        ranked = sorted(items, key=lambda x: x[0], reverse=True)

        dcg = 0.0
        for idx, (_, rel) in enumerate(ranked[:k]):
            dcg += (2**rel - 1) / np.log2(idx + 2)

        # Ideal ranking
        ideal = sorted(items, key=lambda x: x[1], reverse=True)

        idcg = 0.0
        for idx, (_, rel) in enumerate(ideal[:k]):
            idcg += (2**rel - 1) / np.log2(idx + 2)

        if idcg > 0:
            ndcg_scores.append(dcg / idcg)

    return np.mean(ndcg_scores)
df2 = df2.sample(frac=1, random_state=42).reset_index(drop=True)                                 
update(0,100000)
rmse(5001)
print("NDCG@10:", ndcg_at_k(df2.iloc[10000:20000], k=10))


# choose reference anime (example: index 0)
target = H[:, 28977]

best_sim = -1
best_ind = -1

for i in range(0,10000):  # 3400 anime
    if i==28977:
         continue
    vec = H[:, i]

    sim = np.dot(vec, target) / (
        np.linalg.norm(vec) * np.linalg.norm(target)
    )

    if sim > best_sim:
        best_sim = sim
        best_ind = i

print("Most similar anime ID:", best_ind)

print(df1[df1["anime_id"] == best_ind])
print(df1[df1["anime_id"]==28977])
          
