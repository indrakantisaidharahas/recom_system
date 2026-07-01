import pandas as pd
import numpy as np



df1 = pd.read_csv(
    "/home/saidharahas/jupyter_projects/anime_reco/anime_data/anime.csv"
)

df2 = pd.read_csv(
    "/home/saidharahas/jupyter_projects/anime_reco/anime_data/rating.csv"
)



t_users = df2["user_id"].max()
t_anime = df1["anime_id"].max()

print("Max User ID:", t_users)
print("Max Anime ID:", t_anime)



k = 50
lr = 0.01
reg = 0.02

epochs = 10
batch_size = 100000



W = np.random.normal(0, 0.01, (t_users + 1, k))
H = np.random.normal(0, 0.01, (k, t_anime + 1))

bu = np.zeros(t_users + 1)
bi = np.zeros(t_anime + 1)

g_mean = df2["rating"].mean()



def update(st: int, bs: int):

    ldf = df2.iloc[st:st + bs]

    total_err = 0.0

    for row in ldf.itertuples():

        i = row.user_id
        j = row.anime_id
        r = row.rating

        if i > t_users:
            continue

        if j > t_anime:
            continue

        pred = (
            g_mean
            + bu[i]
            + bi[j]
            + W[i] @ H[:, j]
        )

        err = r - pred

        total_err += err * err

        old_w = W[i].copy()

        W[i] += lr * (
            err * H[:, j]
            - reg * W[i]
        )

        H[:, j] += lr * (
            err * old_w
            - reg * H[:, j]
        )

        bu[i] += lr * (
            err
            - reg * bu[i]
        )

        bi[j] += lr * (
            err
            - reg * bi[j]
        )

    return total_err


def rmse(df):

    se = 0.0
    cnt = 0

    for row in df.itertuples():

        i = row.user_id
        j = row.anime_id
        r = row.rating

        if i > t_users:
            continue

        if j > t_anime:
            continue

        pred = (
            g_mean
            + bu[i]
            + bi[j]
            + W[i] @ H[:, j]
        )

        err = r - pred

        se += err * err
        cnt += 1

    return np.sqrt(se / cnt)



df2 = df2.sample(frac=1, random_state=42).reset_index(drop=True)

split = int(0.9 * len(df2))

train_df = df2.iloc[:split]
test_df = df2.iloc[split:]

df2 = train_df



for epoch in range(epochs):

    print(f"\n===== Epoch {epoch+1}/{epochs} =====")

    df2 = df2.sample(frac=1).reset_index(drop=True)

    epoch_loss = 0.0

    for st in range(0, len(df2), batch_size):

        loss = update(st, batch_size)

        epoch_loss += loss

    print("Epoch Loss:", epoch_loss)

    test_rmse = rmse(test_df)

    print("Test RMSE:", test_rmse)

    # np.save("W2.npy", W)
    # np.save("H2.npy", H)
    # np.save("bu.npy", bu)
    # np.save("bi.npy", bi)

    print("Checkpoint Saved")

print("\nTraining Complete")