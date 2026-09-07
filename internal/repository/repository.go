package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pogrammist/golang_test/internal/model"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id int) (*model.Subscription, error)
	List(ctx context.Context, filter *model.ListSubscriptionsRequest) ([]*model.Subscription, error)
	Update(ctx context.Context, id int, req *model.UpdateSubscriptionRequest) (*model.Subscription, error)
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context, filter *model.ListSubscriptionsRequest) (int, error)
	CalculateTotalCost(ctx context.Context, req *model.TotalCostRequest) (int, error)
	RunMigrations(migrationsPath, databaseURL string) error
}

type repository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &repository{db: db}
}

func (r *repository) RunMigrations(migrationsPath, databaseURL string) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	m, err := migrate.New(
		"file://"+absPath,
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (r *repository) Create(ctx context.Context, sub *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)
}

func (r *repository) GetByID(ctx context.Context, id int) (*model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`
	sub := &model.Subscription{}
	var endDate sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if endDate.Valid {
		t := endDate.Time
		sub.EndDate = &t
	}
	return sub, nil
}

func (r *repository) List(ctx context.Context, filter *model.ListSubscriptionsRequest) ([]*model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions WHERE 1=1
	`
	args := []interface{}{}

	if filter.UserID != "" {
		query += " AND user_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, filter.UserID)
	}
	if filter.ServiceName != "" {
		query += " AND service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, filter.ServiceName)
	}

	query += " ORDER BY id DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		var endDate sql.NullTime
		if err := rows.Scan(
			&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
			&sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if endDate.Valid {
			t := endDate.Time
			sub.EndDate = &t
		}
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, rows.Err()
}

func (r *repository) Update(ctx context.Context, id int, req *model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	query := "UPDATE subscriptions SET updated_at = NOW()"
	args := []interface{}{}

	if req.ServiceName != nil {
		query += ", service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, *req.ServiceName)
	}
	if req.Price != nil {
		query += ", price = $" + strconv.Itoa(len(args)+1)
		args = append(args, *req.Price)
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			query += ", end_date = NULL"
		} else {
			query += ", end_date = $" + strconv.Itoa(len(args)+1)
			args = append(args, *req.EndDate)
		}
	}

	query += " WHERE id = $" + strconv.Itoa(len(args)+1) + " RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at"
	args = append(args, id)

	sub := &model.Subscription{}
	var endDate sql.NullTime
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID,
		&sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if endDate.Valid {
		t := endDate.Time
		sub.EndDate = &t
	}
	return sub, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	return err
}

func (r *repository) Count(ctx context.Context, filter *model.ListSubscriptionsRequest) (int, error) {
	query := "SELECT COUNT(*) FROM subscriptions WHERE 1=1"
	args := []interface{}{}

	if filter.UserID != "" {
		query += " AND user_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, filter.UserID)
	}
	if filter.ServiceName != "" {
		query += " AND service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, filter.ServiceName)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *repository) CalculateTotalCost(ctx context.Context, req *model.TotalCostRequest) (int, error) {
	periodStart, err := time.Parse("01-2006", req.PeriodStart)
	if err != nil {
		return 0, err
	}
	periodEnd, err := time.Parse("01-2006", req.PeriodEnd)
	if err != nil {
		return 0, err
	}

	query := `
		SELECT price, start_date, end_date
		FROM subscriptions WHERE 1=1
	`
	args := []interface{}{}

	if req.UserID != "" {
		query += " AND user_id = $" + strconv.Itoa(len(args)+1)
		args = append(args, req.UserID)
	}
	if req.ServiceName != "" {
		query += " AND service_name = $" + strconv.Itoa(len(args)+1)
		args = append(args, req.ServiceName)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	totalCost := 0
	for rows.Next() {
		var price int
		var startDate time.Time
		var endDateNull sql.NullTime
		if err := rows.Scan(&price, &startDate, &endDateNull); err != nil {
			return 0, err
		}

		endDate := periodEnd
		if endDateNull.Valid {
			endDate = endDateNull.Time
		}

		overlapStart := maxDate(startDate, periodStart)
		overlapEnd := minDate(endDate, periodEnd)
		if overlapStart.Before(overlapEnd) || overlapStart.Equal(overlapEnd) {
			months := countMonths(overlapStart, overlapEnd)
			totalCost += months * price
		}
	}
	return totalCost, rows.Err()
}

func maxDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minDate(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func countMonths(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	return years*12 + months + 1
}
