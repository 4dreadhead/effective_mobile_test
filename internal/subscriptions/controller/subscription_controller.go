package controller

import (
	apperrors "effective_mobile_test/internal/platform/errors"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"effective_mobile_test/internal/platform/httputil"
	"effective_mobile_test/internal/subscriptions/presenter"
	"effective_mobile_test/internal/subscriptions/usecase"

	"github.com/go-chi/chi/v5"
)

type SubscriptionController struct {
	usecase *usecase.SubscriptionUsecase
	logger  *slog.Logger
}

func NewSubscriptionController(uc *usecase.SubscriptionUsecase, logger *slog.Logger) *SubscriptionController {
	return &SubscriptionController{usecase: uc, logger: logger}
}

func (c *SubscriptionController) Register(r chi.Router) {
	r.Get("/total", c.total)
	r.Post("/", c.create)
	r.Get("/", c.list)
	r.Get("/{id}", c.get)
	r.Put("/{id}", c.update)
	r.Delete("/{id}", c.delete)
}

type createRequest struct {
	ServiceName string `json:"service_name"`
	MonthlyCost int    `json:"monthly_cost"`
	UserID      string `json:"user_id"`
	From        string `json:"from"`
	To          string `json:"to"`
}

type updateRequest struct {
	ServiceName *string `json:"service_name"`
	MonthlyCost *int    `json:"monthly_cost"`
	From        *string `json:"from"`
	To          *string `json:"to"`
}

// @Security ApiKeyAuth
// @Summary Create subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Param request body createRequest true "Create subscription"
// @Success 201 {object} view.SubscriptionResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Router /subscriptions [post]
func (c *SubscriptionController) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	sub, err := c.usecase.Create(r.Context(), usecase.CreateSubscriptionInput{
		ServiceName: req.ServiceName,
		MonthlyCost: req.MonthlyCost,
		UserID:      req.UserID,
		From:        req.From,
		To:          req.To,
	})
	if err != nil {
		c.writeUsecaseError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, presenter.ToSubscriptionResponse(sub))
}

// @Security ApiKeyAuth
// @Summary Get subscription
// @Description Retrieve a subscription by ID
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 200 {object} view.SubscriptionResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 404 {object} httputil.ErrorResponse
// @Router /subscriptions/{id} [get]
func (c *SubscriptionController) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_id", "id must be numeric")
		return
	}
	sub, err := c.usecase.Get(r.Context(), id)
	if err != nil {
		c.writeUsecaseError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, presenter.ToSubscriptionResponse(sub))
}

// @Security ApiKeyAuth
// @Summary Update subscription
// @Description Update an existing subscription by ID
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Param request body updateRequest true "Update subscription"
// @Success 200 {object} view.SubscriptionResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 404 {object} httputil.ErrorResponse
// @Router /subscriptions/{id} [put]
func (c *SubscriptionController) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_id", "id must be numeric")
		return
	}
	var req updateRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	sub, err := c.usecase.Update(r.Context(), id, usecase.UpdateSubscriptionInput{
		ServiceName: req.ServiceName,
		MonthlyCost: req.MonthlyCost,
		From:        req.From,
		To:          req.To,
	})
	if err != nil {
		c.writeUsecaseError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, presenter.ToSubscriptionResponse(sub))
}

// @Security ApiKeyAuth
// @Summary Delete subscription
// @Description Delete an existing subscription by ID
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 404 {object} httputil.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (c *SubscriptionController) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid_id", "id must be numeric")
		return
	}
	if err = c.usecase.Delete(r.Context(), id); err != nil {
		c.writeUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Security ApiKeyAuth
// @Summary List subscriptions
// @Description List subscriptions with optional filters
// @Tags subscriptions
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Success 200 {array} view.SubscriptionResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Router /subscriptions [get]
func (c *SubscriptionController) list(w http.ResponseWriter, r *http.Request) {
	var filter usecase.ListFilter
	filter.UserID = r.URL.Query().Get("user_id")
	filter.ServiceName = r.URL.Query().Get("service_name")

	subs, err := c.usecase.List(r.Context(), filter)
	if err != nil {
		c.writeUsecaseError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, presenter.ToSubscriptionList(subs))
}

// @Security ApiKeyAuth
// @Summary Total subscription cost
// @Description Calculate total subscription cost for a period
// @Tags subscriptions
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Param from query string false "From period (MM.YYYY)"
// @Param to query string false "To period (MM.YYYY)"
// @Success 200 {object} view.TotalResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Router /subscriptions/total [get]
func (c *SubscriptionController) total(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")
	from := query.Get("from")
	to := query.Get("to")

	total, hasData, err := c.usecase.TotalCost(r.Context(), usecase.TotalFilter{
		UserID:      userID,
		ServiceName: serviceName,
		From:        from,
		To:          to,
	})
	if err != nil {
		c.writeUsecaseError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, presenter.ToTotalResponse(total, hasData))
}

func (c *SubscriptionController) writeUsecaseError(w http.ResponseWriter, err error) {
	c.logger.Error(err.Error())

	switch {
	case errors.Is(err, apperrors.ErrInvalidFields):
		httputil.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
	case errors.Is(err, apperrors.ErrRecordNotFound):
		httputil.WriteError(w, http.StatusNotFound, "not_found", "resource not found")
	default:
		httputil.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func parseID(value string) (uint64, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
