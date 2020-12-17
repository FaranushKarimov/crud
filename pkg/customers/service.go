package customers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

//ErrNotFound var
var (
	ErrNotFound = errors.New("item not found")
	ErrInternal = errors.New("internal error")
)

//Service type
type Service struct {
	pool *pgxpool.Pool
}

//NewService type
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{
		pool: pool,
	}
}

//Customer struct
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

//ByID func
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		SELECT id, name, phone, active, created
		FROM customers
		WHERE id = $1;
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

//All func
func (s *Service) All(ctx context.Context) ([]*Customer, error) {
	items := []*Customer{}

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created FROM customers;
	`)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

//AllActive func
func (s *Service) AllActive(ctx context.Context) ([]*Customer, error) {
	items := []*Customer{}

	rows, err := s.pool.Query(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE active = $1;
	`, true)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

//Save func
func (s *Service) Save(ctx context.Context, customer *Customer) (*Customer, error) {
	item := &Customer{}

	if customer.ID == 0 {
		err := s.pool.QueryRow(ctx, `
			INSERT INTO customers(name, phone)
			VALUES($1, $2)
			RETURNING id, name, phone, active, created; 
		`, customer.Name, customer.Phone).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)

		if err != nil {
			log.Print(err)
			return nil, ErrInternal
		}
		return item, nil
	}

	err := s.pool.QueryRow(ctx, `
			UPDATE customers
			SET name=$1, phone=$2 where id=$3
			RETURNING id, name, phone, active, created; 
		`, customer.Name, customer.Phone, customer.ID).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created,
	)

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil
}

//RemoveByID func
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		DELETE FROM customers
		WHERE id = $1
		RETURNING id, name, phone, active, created;
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

//BlockByID func
func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		UPDATE customers
		SET active = false
		WHERE id = $1
		RETURNING id, name, phone, active, created;
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

//UnblockByID func
func (s *Service) UnblockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		UPDATE customers
		SET active = true
		WHERE id = $1
		RETURNING id, name, phone, active, created;
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}
