SELECT
    u.id AS user_id,
    u.login,
    CASE
        WHEN COUNT(g.id) = 0 THEN 0
        ELSE COUNT(g.id) FILTER (
            WHERE (g.player_x = u.id AND g.status = 'win_x')
            OR (g.player_o = u.id AND g.status = 'win_o')
        )::float / COUNT(g.id)
    END AS win_ratio
FROM users u
LEFT JOIN games g ON (g.player_x = u.id OR g.player_o = u.id)
    AND g.status IN ('win_x', 'win_o', 'draw')
GROUP BY u.id, u.login
ORDER BY win_ratio DESC
LIMIT $1;