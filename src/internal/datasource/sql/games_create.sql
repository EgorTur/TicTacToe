INSERT INTO games (id, board, player_x, player_o, game_type, status)
VALUES ($1, $2::jsonb, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE SET
    board = EXCLUDED.board,
    player_x = EXCLUDED.player_x,
    player_o = EXCLUDED.player_o,
    game_type = EXCLUDED.game_type,
    status = EXCLUDED.status;