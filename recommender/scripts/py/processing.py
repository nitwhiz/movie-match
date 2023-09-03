import re
import string

from keras_preprocessing.sequence import pad_sequences
from keras_preprocessing.text import Tokenizer
from nltk.corpus import stopwords
from pandas import DataFrame


def sanitize_text(media: DataFrame):
    # fill empty values

    media = media.fillna('')

    # remove punctuation

    media['media_summary'] = media['media_summary'].apply(lambda x: re.sub('[^\w\s]', '', x))
    media['media_title'] = media['media_title'].apply(lambda x: re.sub('[^\w\s]', '', x))

    media['media_genre_0'] = media['media_genre_0'].apply(lambda x: re.sub('\s', '', x))
    media['media_genre_1'] = media['media_genre_1'].apply(lambda x: re.sub('\s', '', x))
    media['media_genre_2'] = media['media_genre_2'].apply(lambda x: re.sub('\s', '', x))
    media['media_genre_3'] = media['media_genre_3'].apply(lambda x: re.sub('\s', '', x))

    # remove stop words

    stops = stopwords.words('german')

    media['media_summary'] = media['media_summary'].apply(
        lambda x: ' '.join([word for word in x.lower().translate(x.maketrans('', '', string.punctuation)).split() if
                            word not in stops])
    )

    media['media_title'] = media['media_title'].apply(
        lambda x: ' '.join([word for word in x.lower().translate(x.maketrans('', '', string.punctuation)).split() if
                            word not in stops])
    )

    return media


def process_media(media: DataFrame):
    # remove unneeded columns
    # not sure how to process mediaReleaseDate for now - maybe use unix seconds since 1970

    media = media.drop(columns=['media_id', 'media_release_date'])

    # remove stopwords

    media = sanitize_text(media)

    # process vote_type

    if 'vote_type' in media:
        vote_type_dict = {
            'positive': 1,
            'neutral': 0,
            'negative': -1
        }

        media['vote_type'] = media['vote_type'].apply(lambda x: vote_type_dict.get(x))

    # process media_type

    media_type_dict = {
        'movie': 1,
        'tv': 2,
    }

    media['media_type'] = media['media_type'].apply(lambda x: media_type_dict.get(x))

    # fit tokenizer for genres

    max_genres = 4

    genre_tok = Tokenizer(oov_token='<UnknownGenre>')

    for i in range(4):
        genre_tok.fit_on_texts(media[f'media_genre_{i}'])

    # map genres to tokens

    genre_cols = []

    for i in range(max_genres):
        genre_cols.append(f'media_genre_{i}')

    media[genre_cols] = media[genre_cols].apply(
        lambda x: x.apply(
            lambda y: (lambda z: z[0][0] if z[0] else float(0))(genre_tok.texts_to_sequences([y]))
        )
    )

    # fit tokenizer on summaries and titles

    prosa_tokenizer = Tokenizer(oov_token='<OOV>')

    prosa_tokenizer.fit_on_texts(media['media_summary'])
    prosa_tokenizer.fit_on_texts(media['media_title'])

    # map prosa to tokens

    summary_word_count = 60

    summary_seqs = pad_sequences(prosa_tokenizer.texts_to_sequences(media['media_summary']), maxlen=summary_word_count,
                                 padding='post')

    for i in range(summary_word_count):
        media[f'media_summary_{i}'] = summary_seqs[:, i]

    title_word_count = 5

    title_seqs = pad_sequences(prosa_tokenizer.texts_to_sequences(media['media_title']), maxlen=title_word_count,
                               padding='post')

    for i in range(title_word_count):
        media[f'media_title_{i}'] = title_seqs[:, i]

    return media.drop(columns=['media_summary', 'media_title'])
