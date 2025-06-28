DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'microcinema') THEN
        EXECUTE 'CREATE DATABASE microcinema';
    END IF;
END$$;

\c microcinema;


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    rating DECIMAL(3, 1) CHECK (rating >= 0 AND rating <= 10),
    release_date DATE,
    duration INTEGER, 
    poster_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS movie_genres (
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    genre VARCHAR(50) NOT NULL,
    PRIMARY KEY (movie_id, genre)
);

INSERT INTO movies (title, description, rating) VALUES
    ('The Shawshank Redemption', 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.', 9.3),
    ('The Godfather', 'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.', 9.2),
    ('The Dark Knight', 'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.', 9.0),
    ('Pulp Fiction', 'The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.', 8.9),
    ('Forrest Gump', 'The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.', 8.8)
ON CONFLICT (id) DO NOTHING;

INSERT INTO movie_genres (movie_id, genre) VALUES
    (1, 'Drama'),
    (2, 'Crime'),
    (2, 'Drama'),
    (3, 'Action'),
    (3, 'Crime'),
    (3, 'Drama'),
    (4, 'Crime'),
    (4, 'Drama'),
    (5, 'Drama'),
    (5, 'Romance')
ON CONFLICT (movie_id, genre) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title);
