SELECT id, board, player_x, player_o, game_type, status
FROM games
WHERE status = 'waiting';