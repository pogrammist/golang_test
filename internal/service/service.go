package service

import (
	"context"

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
	return nil, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id int) (*model.Subscription, error) {
	return nil, nil
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, filter *model.ListSubscriptionsRequest) ([]*model.Subscription, int, error) {
	return nil, 0, nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id int, req *model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	return nil, nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return nil
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, req *model.TotalCostRequest) (*model.TotalCostResponse, error) {
	return nil, nil
}
