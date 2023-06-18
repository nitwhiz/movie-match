import re
import string

import pandas as pd
from keras import Sequential
from keras.layers import Dense
from keras_preprocessing.text import Tokenizer
from nltk.corpus import stopwords
from sklearn.model_selection import train_test_split


def make_model():
    media = pd.read_json('../data/media-200.json', lines=True)

    # transpose genres into columns

    def extract_names(row):
        if row is None:
            return []
        return sorted([re.sub('\s+', '', obj['name'].lower()) for obj in row])

    max_genres = media['genres'].apply(lambda x: len(extract_names(x))).max()
    column_names = [f'genre_{i + 1}' for i in range(max_genres)]

    media[column_names] = pd.DataFrame(media['genres'].apply(lambda x: extract_names(x)).tolist())

    media[['summary', 'title']] = media[['summary', 'title']].apply(lambda x: x.apply(lambda y: y.lower()))

    media = media.drop(columns=['id', 'genres', 'createdAt', 'updatedAt'])

    # remove punctuation

    media['summary'] = media['summary'].apply(lambda x: re.sub('[^\w\s]', '', x))

    # remove stopwords

    stops = stopwords.words('german')

    media['summary'] = media['summary'].apply(
        lambda x: ' '.join([word for word in x.lower().translate(x.maketrans('', '', string.punctuation)).split() if
                            word not in stops])
    )

    media.fillna('', inplace=True)

    # add vote to movies
    # vote all "historie" movies "positive"

    #  1 = positive
    #  0 = neutral
    # -1 = negative

    def vote(row):
        for i in range(max_genres):
            if row[f'genre_{i + 1}'] == 'historie':
                return 1

        return 0

    media['vote'] = media.apply(lambda x: vote(x), axis=1)

    genre_tok = Tokenizer(oov_token='<UnknownGenre>')

    for i in range(max_genres):
        genre_tok.fit_on_texts(media[f'genre_{i + 1}'])

    model = Sequential()

    model.add(Dense(24, input_dim=max_genres))
    model.add(Dense(1))

    model.summary()

    model.compile(optimizer='adam', loss='mean_squared_error')

    genre_cols = []

    for i in range(max_genres):
        genre_cols.append(f'genre_{i + 1}')

    media[genre_cols] = media[genre_cols].apply(
        lambda x: x.apply(
            lambda y: (lambda z: z[0][0] if z[0] else float(1))(genre_tok.texts_to_sequences([y]))
        )
    )

    x_train, x_test, y_train, y_test = train_test_split(media[genre_cols], media['vote'])

    model.fit(x_train, y_train, epochs=100)

    y_pred = model.predict(x_test)

    print(x_test)
    print(y_pred)


if __name__ == "__main__":
    make_model()
