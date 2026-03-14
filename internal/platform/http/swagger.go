package platformhttp

import (
	nethttp "net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewSwaggerHandler() nethttp.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)
}
