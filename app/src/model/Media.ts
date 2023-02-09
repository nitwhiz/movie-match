export interface Genre {
  id: number;
  name: string;
}

export const enum MediaType {
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
  releaseDate: string;
  createdAt: string;
}
