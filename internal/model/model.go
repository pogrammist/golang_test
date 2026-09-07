package model

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Subscription struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=0"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name"`
	Price       *int    `json:"price"`
	EndDate     *string `json:"end_date"`
}

type ListSubscriptionsRequest struct {
	UserID      string `json:"-"`
	ServiceName string `json:"-"`
	Limit       int    `json:"-"`
	Offset      int    `json:"-"`
}

type TotalCostRequest struct {
	UserID      string `json:"-"`
	ServiceName string `json:"-"`
	PeriodStart string `json:"-"`
	PeriodEnd   string `json:"-"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

func NewDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
