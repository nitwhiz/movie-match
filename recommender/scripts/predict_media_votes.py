import pandas as pd
from sklearn.neural_network import MLPRegressor
from sklearn.pipeline import make_pipeline
from sklearn.preprocessing import StandardScaler, MinMaxScaler

from scripts.processing import process_media_flat

print("processing media ...")

media_voted_raw = pd.read_json('../data/media_all_voted.json')
media_voted = process_media_flat(media_voted_raw)

pipe = make_pipeline(
    StandardScaler(),
    MLPRegressor(
        max_iter=800,
        verbose=True,
        random_state=42,
        learning_rate_init=0.0001
    ),
)

x = media_voted.drop(columns=['voteType'])
y = media_voted['voteType'].astype(float)

pipe.fit(x, y)

# vote all media

media_all_raw = pd.read_json('../data/media_all.json')
media_all = process_media_flat(media_all_raw)

pred_all = pipe.predict(media_all)

scaler = MinMaxScaler(feature_range=(-1, 1))
scaled_pred = scaler.fit_transform(pred_all.reshape(-1, 1)).flatten()

predicted_media = media_all_raw.copy()

predicted_media['voteType'] = scaled_pred
predicted_media['mediaId'] = media_all_raw['mediaId']


# predicted_media = predicted_media.loc[~predicted_media['mediaId'].isin(media_voted_raw['mediaId'])]
predicted_media = predicted_media.filter(regex='^(?!.*_y$)')
