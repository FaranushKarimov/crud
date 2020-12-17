package security

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4/pgxpool"
)

//Service type struct
type Service struct {
	pool *pgxpool.Pool
}

//NewService type struct
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		pool: pool,
	}
}

//Auth func
func (s *Service) Auth(ctx context.Context, login string, password string) bool {
	var id sql.NullInt64
	err := s.pool.QueryRow(ctx, `
		SELECT id
		FROM managers
		WHERE login = $1 AND password = $2;
	`, login, password).Scan(&id)

	if err != nil {
		return false
	}
	return true
}
