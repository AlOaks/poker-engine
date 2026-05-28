-- Add migration script here
-- Add migration script here
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    street INT,
    done BOOLEAN,
    winners JSON,
    deck JSON,
    board JSON
);

CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    hole_cards JSON,
    folded BOOLEAN,
    game_id INT NOT NULL REFERENCES games(id),
    player_name VARCHAR(100)
);

