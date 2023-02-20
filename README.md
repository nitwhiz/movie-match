# movie-match

Helping you find matching movies for you and your SO since 2023!

## Getting Started

### Configuration

Example `config.yaml`, suitable for usage with the `docker-compose.yml` from the next section:

```yaml
database:
  host: "db"
  port: "5432"
  user: "root"
  password: "root"

media_providers:
  tmdb:
    language: "de"
    region: "de-DE"
    api_key: "<YOUR TMDB API KEY>"
    poster_base_url: "https://image.tmdb.org/t/p/w780"

poster:
  fs_base_path: "/opt/movie-match/posters"

login:
  jwt_key: "<something secret>"
  users:
    - username: "user1"
      display_name: "User 1"
      password: "<PASSWORD HASH, SEE 'Passwords'>"
    - username: "user2"
      display_name: "User 2"
      password: "<PASSWORD HASH, SEE 'Passwords'>"

```

### Passwords

To generate passwords for your user config, run the `hash` command:

```shell
 docker run --rm -it ghcr.io/nitwhiz/movie-match-server:latest hash 
```

You should generate passwords with the same version of the server that's going to consume the password.

### Running with `docker-compose`

Example `docker-compose.yml`:

```yaml
services:

  db:
    image: postgres:15.1-alpine3.17
    environment:
      POSTGRES_USER: "root"
      POSTGRES_PASSWORD: "root"
      POSTGRES_DB: "movie_match"
    ports:
      - "5432:5432"
    volumes:
      - "db_data:/var/lib/postgresql/data"
    networks:
      - db

  server:
    image: ghcr.io/nitwhiz/movie-match-server:latest
    volumes:
      - "./config.yaml:/opt/movie-match/config.yaml:ro" # mount your config
      - "server_data_posters:/opt/movie-match/posters" # mount a directory to store media posters
    ports:
      - "6445:6445"
    depends_on:
      - db
    networks:
      - db

  app:
    image: ghcr.io/nitwhiz/movie-match-app:latest
    environment:
      MOVIEMATCH_API_SERVER_BASE_URL: "http://localhost:6445/"
    ports:
      - "8080:80"

networks:
  db:

volumes:
  db_data:
  server_data_posters:
```
