import re
import string

from keras_preprocessing.sequence import pad_sequences
from keras_preprocessing.text import Tokenizer
from nltk.corpus import stopwords


def sanitize_text(media):
    # fill empty values

    media = media.fillna('')

    # remove punctuation

    media['mediaSummary'] = media['mediaSummary'].apply(lambda x: re.sub('[^\w\s]', '', x))
    media['mediaTitle'] = media['mediaTitle'].apply(lambda x: re.sub('[^\w\s]', '', x))

    media['genre0'] = media['genre0'].apply(lambda x: re.sub('\s', '', x))
    media['genre1'] = media['genre1'].apply(lambda x: re.sub('\s', '', x))
    media['genre2'] = media['genre2'].apply(lambda x: re.sub('\s', '', x))
    media['genre3'] = media['genre3'].apply(lambda x: re.sub('\s', '', x))

    # remove stop words

    stops = stopwords.words('german')

    media['mediaSummary'] = media['mediaSummary'].apply(
        lambda x: ' '.join([word for word in x.lower().translate(x.maketrans('', '', string.punctuation)).split() if
                            word not in stops])
    )

    media['mediaTitle'] = media['mediaTitle'].apply(
        lambda x: ' '.join([word for word in x.lower().translate(x.maketrans('', '', string.punctuation)).split() if
                            word not in stops])
    )

    return media


def process_media_flat(media):
    # remove unneeded columns
    # not sure how to process mediaReleaseDate for now - maybe use unix seconds since 1970

    media = media.drop(columns=['mediaId', 'mediaReleaseDate'])

    # remove stopwords

    media = sanitize_text(media)

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

    summary_word_count = 60

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
