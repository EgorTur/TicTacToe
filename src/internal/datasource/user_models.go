package datasource

import "github.com/google/uuid"

type UserModel struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Login        string    `db:"login" json:"login"`
	PasswordHash string    `db:"password_hash" json:"-"`
}
