import { Client } from 'pg';
import * as fs from 'fs';

const enum MediaType {
  TV = 1,
  MOVIE = 2,
}

const run = async () => {
  const client = new Client({
    host: 'localhost',
    user: 'root',
    password: 'root',
    database: 'movie_match',
  });

  await client.connect();

  const result = await client.query(`
    select m.id           as media_id,
           m.type         as media_type,
           m.title        as media_title,
           m.summary      as media_summary,
           m.rating       as media_rating,
           m.release_date as media_release_date,
           g.name         as genre_name
    from media m
             join media_genres mg on m.id = mg.media_id
             join genres g on g.id = mg.genre_id
  `);

  const data: Record<
    string,
    {
      mediaId: string;
      mediaType: MediaType;
      mediaTitle: string;
      mediaSummary: string;
      mediaRating: number;
      mediaReleaseDate: number;
      genre0: string | null;
      genre1: string | null;
      genre2: string | null;
      genre3: string | null;
    }
  > = {};

  for (const row of result.rows) {
    if (!data[row['media_id']]) {
      data[row['media_id']] = {
        mediaId: row['media_id'],
        mediaType:
          row['media_type'] === 'movie' ? MediaType.MOVIE : MediaType.TV,
        mediaTitle: row['media_title'],
        mediaSummary: row['media_summary'],
        mediaRating: Number(row['media_rating']),
        mediaReleaseDate: Date.parse(row['media_release_date']),
        genre0: null,
        genre1: null,
        genre2: null,
        genre3: null,
      };
    }

    // place genre at correct index
    for (let i = 0; i < 4; ++i) {
      const idx = `genre${i}` as 'genre0' | 'genre1' | 'genre2' | 'genre3';

      if (data[row['media_id']][idx] === null) {
        data[row['media_id']][idx] = row['genre_name'];
        break;
      }
    }
  }

  fs.writeFileSync(
    '../data/media_all.json',
    JSON.stringify(Object.values(data)),
    'utf-8'
  );

  await client.end();
};

run().then((r) => console.log('done!'));
