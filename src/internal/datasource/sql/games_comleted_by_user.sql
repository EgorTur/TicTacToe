SELECT id, board, player_x, player_o, game_type, status, created_at
FROM games
WHERE (player_x = $1 OR player_o = $1)
    AND status IN ('win_x', 'win_o', 'draw')
ORDER BY created_at DESC;