package platformhttp

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"effective_mobile_test/internal/platform/auth"
	"effective_mobile_test/internal/subscriptions/controller"
)

func NewRouter(
	logger *slog.Logger,
	subscriptionController *controller.SubscriptionController,
	swaggerHandler http.Handler,
	apiKey string,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(LoggingMiddleware(logger))

	r.Route("/subscriptions", func(sr chi.Router) {
		sr.Use(auth.Middleware(apiKey))
		subscriptionController.Register(sr)
	})

	if swaggerHandler != nil {
		r.Get("/swagger/*", swaggerHandler.ServeHTTP)
	}

	return r
}
