export interface Genre {
  id: number;
  name: string;
}

export interface Media {
  id: string;
  foreignId: string;
  type: string;
  provider: string;
  title: string;
  summary: string;
  genres: Genre[];
  releaseDate: string;
  createdAt: string;
}
