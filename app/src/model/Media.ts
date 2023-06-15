import { VoteType } from './Vote';

export interface Genre {
  id: number;
  name: string;
}

export const enum MediaType {
  ALL = 'all',
  TV = 'tv',
  MOVIE = 'movie',
}

export interface Media {
  id: string;
  foreignId: string;
  type: MediaType;
  provider: string;
  title: string;
  summary: string;
  genres: Genre[];
  runtime: number;
  rating: number;
  releaseDate: string;
  createdAt: string;
  updatedAt: string;
}

export interface RecommendedMedia extends Media {
  score: string;
  seen: boolean;
  voteType: VoteType;
}
