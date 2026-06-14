UPDATE games SET board = $2::jsonb, player_x = $3, player_o = $4, game_type = $5, status = $6
WHERE id = $1;