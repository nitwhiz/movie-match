# this does not work at all

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from keras_preprocessing.sequence import pad_sequences
from keras_preprocessing.text import Tokenizer
from tensorflow.python.keras import Input
from tensorflow.python.keras.layers import Embedding, Flatten, Concatenate, Dense
from tensorflow.python.keras.models import Model

from scripts.processing import sanitize_text

print("processing media ...")

media_voted_raw = pd.read_json('../data/media_all_voted.json')

media = sanitize_text(media_voted_raw)

mediaType_input = Input(shape=(1,))
mediaRating_input = Input(shape=(1,))

genre_input = Input(shape=(4,))

summary_input = Input(shape=(100,))
title_input = Input(shape=(5,))

num_genres = 27
max_words = 1000

embedding_dim = 8

genre_emb = Embedding(input_dim=num_genres * 4, output_dim=embedding_dim)(genre_input)

summary_emb = Embedding(input_dim=max_words, output_dim=embedding_dim)(summary_input)
title_emb = Embedding(input_dim=max_words, output_dim=embedding_dim)(title_input)

genre_flat = Flatten()(genre_emb)
summary_flat = Flatten()(summary_emb)
title_flat = Flatten()(title_emb)

mediaType_data = np.array(media['mediaType'])

genre0_data = np.array(media['genre0'])
genre1_data = np.array(media['genre1'])
genre2_data = np.array(media['genre2'])
genre3_data = np.array(media['genre3'])

mediaRating_data = np.array(media['mediaRating'])

summary_data = np.array(media['mediaSummary'])
title_data = np.array(media['mediaTitle'])

target = np.array(media['voteType']).astype(float)  # The target values between -1 and 1

# tokenize

tokenizer_general = Tokenizer(num_words=max_words, oov_token='UNK')
tokenizer_general.fit_on_texts(np.concatenate([summary_data, title_data]))

summary_seq = tokenizer_general.texts_to_sequences(summary_data)
title_seq = tokenizer_general.texts_to_sequences(title_data)

tokenizer_genres = Tokenizer(num_words=num_genres, oov_token='UNK')
tokenizer_genres.fit_on_texts(np.concatenate([genre0_data, genre1_data, genre2_data, genre3_data]))

genre0_seq = tokenizer_genres.texts_to_sequences(genre0_data)
genre1_seq = tokenizer_genres.texts_to_sequences(genre1_data)
genre2_seq = tokenizer_genres.texts_to_sequences(genre2_data)
genre3_seq = tokenizer_genres.texts_to_sequences(genre3_data)

summary_data_padded = pad_sequences(summary_seq, maxlen=100, padding='post')
title_data_padded = pad_sequences(title_seq, maxlen=5, padding='post')

genre0_data_padded = pad_sequences(genre0_seq, maxlen=1, padding='post')
genre1_data_padded = pad_sequences(genre1_seq, maxlen=1, padding='post')
genre2_data_padded = pad_sequences(genre2_seq, maxlen=1, padding='post')
genre3_data_padded = pad_sequences(genre3_seq, maxlen=1, padding='post')

genre_data_combined = np.column_stack([genre0_data_padded, genre1_data_padded, genre2_data_padded, genre3_data_padded])

concatenated = Concatenate()([
    mediaType_input,
    genre_flat,
    mediaRating_input,
    summary_flat,
    title_flat
])

# Define the dense neural network
dense_layer_1 = Dense(128, activation='relu')(concatenated)
dense_layer_2 = Dense(64, activation='relu')(dense_layer_1)
output = Dense(1, activation='sigmoid')(dense_layer_2)  # Sigmoid activation for output between -1 and 1

# Define the model with all inputs and the output
model = Model(
    inputs=[
        mediaType_input,
        genre_input,
        mediaRating_input,
        summary_input,
        title_input
    ],
    outputs=output
)

# Compile the model
model.compile(optimizer='adam', loss='mean_squared_error')

# Print the summary of the model
model.summary()

print("mediaType_data shape:", mediaType_data.shape)
print("genre_data_combined shape:", genre_data_combined.shape)
print("mediaRating_data shape:", mediaRating_data.shape)
print("summary_data_padded shape:", summary_data_padded.shape)
print("title_data_padded shape:", title_data_padded.shape)
print("target shape:", target.shape)

input_data = [
    mediaType_data,
    genre_data_combined,
    mediaRating_data,
    summary_data_padded,
    title_data_padded
]

# Train the model
history = model.fit(
    x=input_data,
    y=target,
    validation_data=(input_data, target),
    epochs=50,
    steps_per_epoch=3000,
    validation_steps=3000
)

plt.plot(history.history['acc'], label='Training Loss')
plt.xlabel('Epoch')
plt.ylabel('Loss')
plt.legend()
plt.show()
