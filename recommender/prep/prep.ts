import { Client } from 'pg';

const enum VoteType {
  NEGATIVE = -1,
  NEUTRAL = 0,
  POSITIVE = 1,
}

const voteTypeMap: Record<string, VoteType> = {
  negative: VoteType.NEGATIVE,
  neutral: VoteType.NEUTRAL,
  positive: VoteType.POSITIVE,
};

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
    select v.type         as vote_type,
           m.id           as media_id,
           m.type         as media_type,
           m.title        as media_title,
           m.summary      as media_summary,
           m.rating       as media_rating,
           m.release_date as media_release_date,
           g.name         as genre_name
    from votes v
             join media m on v.media_id = m.id
             join media_genres mg on m.id = mg.media_id
             join genres g on g.id = mg.genre_id
    where v.user_id = 'f5319987-66cb-40a1-b12c-0a3360db7819'
  `);

  const data: Record<
    string,
    {
      mediaId: string;
      voteType: VoteType;
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
        voteType: voteTypeMap[row['vote_type']],
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

    for (let i = 0; i < 4; ++i) {
      const idx = `genre${i}` as 'genre0' | 'genre1' | 'genre2' | 'genre3';

      if (data[row['media_id']][idx] === null) {
        data[row['media_id']][idx] = row['genre_name'];
        break;
      }
    }
  }

  console.log(data);

  await client.end();
};

run();
