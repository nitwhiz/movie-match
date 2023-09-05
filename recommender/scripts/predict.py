from sklearn.neural_network import MLPRegressor
from sklearn.pipeline import make_pipeline
from sklearn.preprocessing import StandardScaler, MinMaxScaler
from sqlalchemy.engine import Connection

from database import get_voted_media, get_all_media
from processing import process_media


def predict_votes(db_connection: Connection, user_id: str):
    media_voted_raw = get_voted_media(db_connection, user_id)

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

    media_all_raw = get_all_media(db_connection)

    media_all = process_media(media_all_raw)

    pred_all = pipe.predict(media_all)

    scaler = MinMaxScaler(feature_range=(-1, 1))
    scaled_pred = scaler.fit_transform(pred_all.reshape(-1, 1)).flatten()

    predicted_media = media_all_raw.copy()

    predicted_media['vote_type'] = scaled_pred
    predicted_media['media_id'] = media_all_raw['media_id']

    # remove media previously voted
    # predicted_media = predicted_media.loc[~predicted_media['media_id'].isin(media_voted_raw['media_id'])]

    # remove _y columns
    return predicted_media.filter(regex='^(?!.*_y$)')
