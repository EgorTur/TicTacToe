-- DROP TABLE IF EXISTS games CASCADE;

CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY,
    board JSONB NOT NULL DEFAULT '[]',
    player_x UUID REFERENCES users(id),
    player_o UUID REFERENCES users(id),
    game_type TEXT NOT NULL DEFAULT 'player',
    status TEXT NOT NULL DEFAULT 'waiting',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);