package customers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ErrNotFound ...
var ErrNotFound = errors.New("item not found")

// ErrInternal ...
var ErrInternal = errors.New("internal error")

// Service ...
type Service struct {
	pool *pgxpool.Pool
}

// NewService ...
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// Customer ...
type Customer struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

// ByID function
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
	SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	// log.Print("ByID ", item)
	return item, nil
}

// All ...
func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)

	rows, err := s.pool.Query(ctx, `
	SELECT id, name, phone, active, created FROM customers;
	`)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, err

	}

	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return items, nil
}

// AllActive ...
func (s *Service) AllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)

	rows, err := s.pool.Query(ctx, `
	SELECT id, name, phone, active, created FROM customers WHERE active;
	`)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, err

	}

	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return items, nil
}

// RemoveByID ...
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
	DELETE FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

// BlockByID ...
func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	active := "FALSE"
	err := s.pool.QueryRow(ctx, `
	UPDATE customers SET active = $2 WHERE id = $1
	`, id, active).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

// UnblockByID ...
func (s *Service) UnblockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	active := "TRUE"
	err := s.pool.QueryRow(ctx, `
	UPDATE customers SET active = $2 WHERE id = $1
	`, id, active).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

// Save ...
func (s *Service) Save(ctx context.Context, item *Customer) (*Customer, error) {

	if item.ID == 0 {

		err := s.pool.QueryRow(ctx, `
		INSERT INTO customers (name, phone, password) VALUES ($1, $2, $3)
		RETURNING id, name, phone, password, active, created;
		-- ON CONFLICT (phone)  DO UPDATE SET name = excluded.name;  `,
			item.Name, item.Phone, item.Password).Scan(&item.ID, &item.Name, &item.Phone, &item.Password, &item.Active, &item.Created)

		if err != nil {
			log.Print(err)
			return nil, err
		}

		if errors.Is(err, pgx.ErrNoRows) {
			log.Print("No rows")
			return nil, ErrNotFound
		}

		if err != nil {
			log.Print(err)
			return nil, err

		}

		// log.Print("Save ", item)
		return item, nil
	}

	rowsDB, err := s.All(ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for _, value := range rowsDB {
		if value.ID == item.ID {
			_, err = s.pool.Exec(ctx, `
				UPDATE customers SET name = $1 WHERE phone = $2, active = $3, created = $4;
				`, item.Name, item.Phone, item.Active, item.Created)
			log.Print("item.Created 205: ", item.Created)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			return item, nil
		}
	}

	return nil, ErrNotFound
}
