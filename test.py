import numpy as np
import pandas as pd
from collections import defaultdict

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
# -----------------------------
# Load data
# -----------------------------

df1 = pd.read_csv(
    "/home/saidharahas/jupyter_projects/anime_reco/anime_data/anime.csv"
)

df2 = pd.read_csv(
    "/home/saidharahas/jupyter_projects/anime_reco/anime_data/rating.csv"
)
g_mean=df2["rating"].mean()
W = np.load("W.npy")
H = np.load("H.npy")
bu = np.load("bu.npy")
bi = np.load("bi.npy")

print("W:", W.shape)
print("H:", H.shape)


anime_lookup = (
    df1.set_index("anime_id")
       .to_dict("index")
)



def similar_anime(anime_id, top_k=10):

    if anime_id >= H.shape[1]:
        print("Anime ID not present in model.")
        return

    if anime_id not in anime_lookup:
        print("Anime not found in anime.csv")
        return

    target = H[:, anime_id]

    target_norm = np.linalg.norm(target)

    sims = []

    for aid in range(H.shape[1]):

        if aid == anime_id:
            continue

        if aid not in anime_lookup:
            continue

        vec = H[:, aid]

        denom = target_norm * np.linalg.norm(vec)

        if denom == 0:
            continue

        sim = np.dot(target, vec) / denom

        sims.append((sim, aid))

    sims.sort(reverse=True)

  

    target_info = anime_lookup[anime_id]

    print("ID:", anime_id)
    print("Name:", target_info["name"])
    print("Genre:", target_info["genre"])
    print("Rating:", target_info["rating"])

    print("\n" + "=" * 60)
    print(f"TOP {top_k} SIMILAR ANIME")
    print("=" * 60)

    shown = 0

    for sim, aid in sims:

        info = anime_lookup[aid]

        print(
            f"{shown+1:2d}. "
            f"Sim={sim:.4f} | "
            f"ID={aid} | "
            f"{info['name']}"
        )

        shown += 1

        if shown >= top_k:
            break


print("NDCG@10:", ndcg_at_k(df2.iloc[10000:20000], k=10))
print("\n\nGINTAMA°")
similar_anime(28977, top_k=20)

print("\n\nFULLMETAL ALCHEMIST: BROTHERHOOD")
similar_anime(5114, top_k=20)

print("\n\nDEATH NOTE")
similar_anime(1535, top_k=20)

print("\n\nATTACK ON TITAN")
similar_anime(16498, top_k=20)