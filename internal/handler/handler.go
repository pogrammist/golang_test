package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pogrammist/golang_test/internal/model"
	"github.com/pogrammist/golang_test/internal/service"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *zap.Logger
}

func NewSubscriptionHandler(service service.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, logger: logger}
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body model.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.CreateSubscription(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("subscription created", zap.Int("id", sub.ID), zap.String("user_id", sub.UserID))
	c.JSON(http.StatusCreated, sub)
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Description Get subscription details by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get subscription", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Description List all subscriptions with optional filtering
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Param limit query int false "Page size" default(10)
// @Param offset query int false "Page offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	filter := &model.ListSubscriptionsRequest{
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("service_name"),
		Limit:       10,
		Offset:      0,
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}
	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	subscriptions, count, err := h.service.ListSubscriptions(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscriptions": subscriptions,
		"total":         count,
		"limit":         filter.Limit,
		"offset":        filter.Offset,
	})
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body model.UpdateSubscriptionRequest true "Update data"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.UpdateSubscription(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("failed to update subscription", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteSubscription(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete subscription", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("subscription deleted", zap.Int("id", id))
	c.Status(http.StatusNoContent)
}

// CalculateTotalCost godoc
// @Summary Calculate total subscription cost
// @Description Calculate total cost for a selected period with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Param period_start query string true "Period start (MM-YYYY)"
// @Param period_end query string true "Period end (MM-YYYY)"
// @Success 200 {object} model.TotalCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	req := &model.TotalCostRequest{
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("service_name"),
		PeriodStart: c.Query("period_start"),
		PeriodEnd:   c.Query("period_end"),
	}

	resp, err := h.service.CalculateTotalCost(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to calculate total cost", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HealthCheck godoc
// @Summary Health check
// @Description Service health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *SubscriptionHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
