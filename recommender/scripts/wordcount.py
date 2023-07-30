import matplotlib.pyplot as plt
import pandas as pd

from scripts.processing import sanitize_text

media_raw = pd.read_json('../data/media_all.json')

media = sanitize_text(media_raw)

media['title_word_count'] = media['mediaTitle'].apply(
    lambda x: len(x.split(sep=' '))
)

media['summary_word_count'] = media['mediaSummary'].apply(
    lambda x: len(x.split(sep=' '))
)

plt.figure()
media.boxplot(column='title_word_count', whis=[5, 95])
plt.show()

plt.figure()
media.boxplot(column='summary_word_count', whis=[5, 95])
plt.show()

