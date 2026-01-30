// Package modelconv provides conversion between proto messages and model types
// using JSON marshaling as a bridge. This approach requires model types to have
// JSON tags matching proto field names (camelCase) and custom JSON marshal/unmarshal
// implementations for types like enums and timestamps.
package modelconv

import (
	"encoding/json"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/luckymaks/bm_backend/backend/internal/model"
)

// FromProto converts a proto message to a model type via JSON marshaling.
// The target type T must have JSON tags matching the proto field names.
func FromProto[T any](msg proto.Message) (T, error) {
	var result T
	jsonBytes, err := protojson.Marshal(msg)
	if err != nil {
		return result, errors.Wrap(err, "failed to marshal proto to JSON")
	}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, errors.Wrap(err, "failed to unmarshal JSON to model")
	}
	return result, nil
}

// ToProto converts a model type to a proto message via JSON marshaling.
// The model type must have JSON tags matching the proto field names.
func ToProto[T proto.Message](mdl any, target T) (T, error) {
	jsonBytes, err := json.Marshal(mdl)
	if err != nil {
		return target, errors.Wrap(err, "failed to marshal model to JSON")
	}
	if err = protojson.Unmarshal(jsonBytes, target); err != nil {
		return target, errors.Wrap(err, "failed to unmarshal JSON to proto")
	}
	return target, nil
}

// ToConnectError converts a model error to a connect error.
func ToConnectError(err error) error {
	if errors.Is(err, model.ErrValidation) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	if errors.Is(err, model.ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}
	if errors.Is(err, model.ErrAlreadyExists) {
		return connect.NewError(connect.CodeAlreadyExists, err)
	}
	return connect.NewError(connect.CodeInternal, err)
}
