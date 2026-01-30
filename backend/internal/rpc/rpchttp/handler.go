package rpchttp

import (
	"net/http"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"

	"github.com/luckymaks/bm_backend/backend/proto/bm/v1/bmv1connect"
)

func NewHandler(rpcHandler bmv1connect.ApiServiceHandler) http.Handler {
	mux := http.NewServeMux()

	interceptors := connect.WithInterceptors(
		NewErrorInterceptor(),
	)

	apiPath, apiHandler := bmv1connect.NewApiServiceHandler(rpcHandler, interceptors)
	mux.Handle("/v1/rpc"+apiPath, http.StripPrefix("/v1/rpc", apiHandler))

	return WithCORS(mux)
}

func WithCORS(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   append(connectcors.AllowedHeaders(), "X-API-Key"),
		ExposedHeaders:   connectcors.ExposedHeaders(),
		MaxAge:           7200,
		AllowCredentials: false,
	})
	return c.Handler(h)
}
