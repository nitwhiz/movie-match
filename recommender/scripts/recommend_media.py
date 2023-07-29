import re
import string

import pandas as pd
from keras_preprocessing.sequence import pad_sequences
from keras_preprocessing.text import Tokenizer
from nltk.corpus import stopwords
from sklearn.naive_bayes import GaussianNB


def process_media(media):
    # remove unneeded columns
    # not sure how to process mediaReleaseDate for now - maybe use unix seconds since 1970

    media = media.drop(columns=['mediaId', 'mediaReleaseDate'])

    # remove punctuation

    media['mediaSummary'] = media['mediaSummary'].apply(lambda x: re.sub('[^\w\s]', '', x))

    # remove stopwords

    stops = stopwords.words('german')

    media['mediaSummary'] = media['mediaSummary'].apply(
        lambda x: ' '.join([word for word in x.lower().translate(x.maketrans('', '', string.punctuation)).split() if
                            word not in stops])
    )

    media['mediaTitle'] = media['mediaTitle'].apply(
        lambda x: ' '.join([word for word in x.lower().translate(x.maketrans('', '', string.punctuation)).split() if
                            word not in stops])
    )

    # fill null with empty strings

    media.fillna('', inplace=True)

    # fit tokenizer for genres

    max_genres = 4

    genre_tok = Tokenizer(oov_token='<UnknownGenre>')

    for i in range(4):
        genre_tok.fit_on_texts(media[f'genre{i}'])

    # map genres to tokens

    genre_cols = []

    for i in range(max_genres):
        genre_cols.append(f'genre{i}')

    media[genre_cols] = media[genre_cols].apply(
        lambda x: x.apply(
            lambda y: (lambda z: z[0][0] if z[0] else float(0))(genre_tok.texts_to_sequences([y]))
        )
    )

    # fit tokenizer on summaries and titles

    prosa_tokenizer = Tokenizer(oov_token='<OOV>')

    prosa_tokenizer.fit_on_texts(media['mediaSummary'])
    prosa_tokenizer.fit_on_texts(media['mediaTitle'])

    # map prosa to tokens

    summary_word_count = 20

    summary_seqs = pad_sequences(prosa_tokenizer.texts_to_sequences(media['mediaSummary']), maxlen=summary_word_count,
                                 padding='post')

    for i in range(summary_word_count):
        media[f'summary{i}'] = summary_seqs[:, i]

    title_word_count = 5

    title_seqs = pad_sequences(prosa_tokenizer.texts_to_sequences(media['mediaTitle']), maxlen=title_word_count,
                               padding='post')

    for i in range(title_word_count):
        media[f'title{i}'] = title_seqs[:, i]

    return media.drop(columns=['mediaSummary', 'mediaTitle'])


def main():
    print("hello world!")

    media_voted_raw = pd.read_json('../data/media_all_voted.json')
    media_voted = process_media(media_voted_raw)

    # x_train, x_test, y_train, y_test = train_test_split(media, media['voteType'], test_size=.25, random_state=0)
    # cm = confusion_matrix(y_test, y_pred)

    classifier = GaussianNB()

    x = media_voted.drop(columns=['voteType'])
    y = media_voted['voteType']

    classifier.fit(x, y)

    # vote all media

    media_all_raw = pd.read_json('../data/media_all.json')
    media_all = process_media(media_all_raw)

    pred_all = classifier.predict(media_all)

    media_all_raw['voteType'] = pred_all

    recommended_media = media_all_raw[media_all_raw['voteType'] == 1]

    merged_df = recommended_media.merge(media_voted_raw, on='mediaId', how='left', indicator=True)

    recommended_media = merged_df[merged_df['_merge'] == 'left_only']
    recommended_media.drop(columns=['_merge'], inplace=True)

    media_voted.head()


if __name__ == "__main__":
    main()
