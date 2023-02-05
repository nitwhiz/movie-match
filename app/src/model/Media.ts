export interface TMDBMovieData {
  ID: number;
  Genres: { ID: number; Name: string }[];
  Overview: string;
  Title: string;
  vote_average: number;
  poster_path: string;
}

type MediaType = TMDBMovieData;

export interface Media<MT extends MediaType = any> {
  ID: string;
  ForeignID: string;
  Type: string;
  DataSource: string;
  Data: MT;
  CreatedAt: string;
}
