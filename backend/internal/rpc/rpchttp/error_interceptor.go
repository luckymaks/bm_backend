package rpchttp

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"

	"github.com/luckymaks/bm_backend/backend/internal/model"
)

func NewErrorInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err != nil {
				err = mapModelError(err)
			}

			return resp, err
		}
	}
}

func mapModelError(err error) error {
	var connectErr *connect.Error
	if errors.As(err, &connectErr) {
		return err
	}

	switch {
	case errors.Is(err, model.ErrNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, model.ErrAlreadyExists):
		return connect.NewError(connect.CodeAlreadyExists, err)
	case errors.Is(err, model.ErrValidation):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return err
	}
}
