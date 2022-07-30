package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"storage/internal/entity"

	httptransport "github.com/go-kit/kit/transport/http"
)

// DecodeRequestWithoutBody ...
func DecodeRequestWithoutBody() httptransport.DecodeRequestFunc {
	return func(_ context.Context, _ *http.Request) (any, error) {
		var request entity.EmptyRequest

		return request, nil
	}
}

// DecodeRequest ...
func DecodeRequest[req entity.IDRequest |
	entity.UsernamePasswordRequest |
	entity.UsernameRequest |
	entity.UsernamePasswordEmailRequest](request req,
) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (any, error) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, fmt.Errorf("failed to decode request: %w", err)
		}

		return request, nil
	}
}

// EncodeResponse ...
func EncodeResponse(_ context.Context, w http.ResponseWriter, response any) error {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
