package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// ErrNoSuchUser variable
var ErrNoSuchUser = errors.New("no such user")

// ErrInvalidPassword variable
var ErrInvalidPassword = errors.New("invalid password")

// ErrInternal variable
var ErrInternal = errors.New("internal error")

// ErrExpireToken ...
var ErrExpireToken = errors.New("token expired")

// Service ...
type Service struct {
	pool *pgxpool.Pool
}

// NewService ...
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// Auth ...
func (s *Service) Auth(login string, password string) bool {

	err := s.pool.QueryRow(context.Background(), `
	SELECT login, password FROM managers WHERE login = $1 AND password = $2
	`, login, password).Scan(&login, &password)

	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

// TokenForCustomer ...
func (s *Service) TokenForCustomer(
	ctx context.Context,
	phone string,
	password string,
) (string, error) {
	var hash string
	var id int64

	err := s.pool.QueryRow(ctx, `SELECT id, password FROM customers WHERE phone = $1`, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token := hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO customers_tokens(token, customer_id) VALUES($1, $2)`, token, id)
	if err != nil {
		return "", ErrInternal
	}

	return token, nil
}

// AuthenticateCustomer ...
func (s *Service) AuthenticateCustomer(
	ctx context.Context,
	token string,
) (int64, error) {
	var id int64
	var expire time.Time

	err := s.pool.QueryRow(ctx, `SELECT customer_id, expire FROM customers_tokens WHERE token = $1`, token).Scan(&id, &expire)

	if err == pgx.ErrNoRows {
		return 0, ErrNoSuchUser
	}

	if err != nil {
		return 0, ErrInternal
	}

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	timeEnd := expire.Format("2006-01-02 15:04:05")

	if timeNow > timeEnd {
		return 0, ErrExpireToken
	}

	return id, nil
}
