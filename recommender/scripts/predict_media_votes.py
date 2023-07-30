import pandas as pd
import sqlalchemy
from sklearn.neural_network import MLPRegressor
from sklearn.pipeline import make_pipeline
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from sqlalchemy import create_engine

from scripts.processing import process_media

alchemyEngine = create_engine('postgresql+psycopg2://root:root@localhost/movie_match', pool_recycle=3600)
dbConnection = alchemyEngine.connect()

with open('media_voted.sql', 'r') as f:
    media_voted_raw = pd.read_sql_query(
        sqlalchemy.text(f.read()),
        dbConnection,
        params={
            "user_id": "f5319987-66cb-40a1-b12c-0a3360db7819"
        }
    )

media_voted = process_media(media_voted_raw)

pipe = make_pipeline(
    StandardScaler(),
    MLPRegressor(
        max_iter=5000,
        verbose=True,
        random_state=42,
        learning_rate_init=0.0001
    ),
)

x = media_voted.drop(columns=['vote_type'])
y = media_voted['vote_type'].astype(float)

pipe.fit(x, y)

# vote all media
with open('media_all.sql', 'r') as f:
    media_all_raw = pd.read_sql_query(
        f.read(),
        dbConnection
    )

media_all = process_media(media_all_raw)

pred_all = pipe.predict(media_all)

scaler = MinMaxScaler(feature_range=(-1, 1))
scaled_pred = scaler.fit_transform(pred_all.reshape(-1, 1)).flatten()

predicted_media = media_all_raw.copy()

predicted_media['vote_type'] = scaled_pred
predicted_media['media_id'] = media_all_raw['media_id']

# predicted_media = predicted_media.loc[~predicted_media['mediaId'].isin(media_voted_raw['mediaId'])]
predicted_media = predicted_media.filter(regex='^(?!.*_y$)')
