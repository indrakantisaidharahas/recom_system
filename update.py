import numpy as np
import pandas as pd 



df1=pd.read_csv("/home/saidharahas/jupyter_projects/anime_reco/anime_data/anime.csv")
df2=pd.read_csv("/home/saidharahas/jupyter_projects/anime_reco/anime_data/rating.csv")


##batch incrementation of ratings training 
t_users=df2["user_id"].nunique()
print(t_users)
t_anime=df1["anime_id"].nunique()
print(t_anime)

t_anime= df1["anime_id"].max()

k=10
W = np.random.normal(0, 0.01, (t_users+1, k))
H = np.random.normal(0, 0.01, (k, t_anime+1))
lr = 0.05
reg = 0.02

def update(st: int, bs: int):

    ldf = df2.iloc[st:st + bs]

    for epochs in range(50):

        total_err = 0

        for row in ldf.itertuples():

            j = row.anime_id
            i = row.user_id
            r = row.rating

            if j > 73515:
                continue

            if i > 12294:
                continue

            pred = W[i] @ H[:, j]
            err = r - pred

            total_err += err * err

            old_w = W[i].copy()

            W[i] += lr * (err * H[:, j] - reg * W[i])
            H[:, j] += lr * (err * old_w - reg * H[:, j])

        print(f"Epoch {epochs + 1}")
        print("Loss:", total_err)
        print("-" * 30)
            
         
                                 
update(0,1000)

