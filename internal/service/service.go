package service

import (
	"context"
	"errors"
	"time"

	"github.com/pogrammist/golang_test/internal/model"
	"github.com/pogrammist/golang_test/internal/repository"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req *model.CreateSubscriptionRequest) (*model.Subscription, error)
	GetSubscription(ctx context.Context, id int) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context, filter *model.ListSubscriptionsRequest) ([]*model.Subscription, int, error)
	UpdateSubscription(ctx context.Context, id int, req *model.UpdateSubscriptionRequest) (*model.Subscription, error)
	DeleteSubscription(ctx context.Context, id int) error
	CalculateTotalCost(ctx context.Context, req *model.TotalCostRequest) (*model.TotalCostResponse, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format, expected MM-YYYY")
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id int) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, filter *model.ListSubscriptionsRequest) ([]*model.Subscription, int, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	filter.Limit = limit

	subscriptions, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, count, nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id int, req *model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	if req.EndDate != nil && *req.EndDate != "" {
		_, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end_date format, expected MM-YYYY")
		}
	}
	return s.repo.Update(ctx, id, req)
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, req *model.TotalCostRequest) (*model.TotalCostResponse, error) {
	if req.PeriodStart == "" || req.PeriodEnd == "" {
		return nil, errors.New("period_start and period_end are required")
	}

	totalCost, err := s.repo.CalculateTotalCost(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.TotalCostResponse{TotalCost: totalCost}, nil
}
