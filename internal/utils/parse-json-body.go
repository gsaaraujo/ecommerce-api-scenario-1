package utils

import (
	"encoding/json"
	"io"
)

func ParseJSONBody[T any](body io.ReadCloser) (T, error) {
	var v T

	defer func() {
		_ = body.Close()
	}()

	data, err := io.ReadAll(body)
	if err != nil {
		return v, err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return v, err
	}

	return v, nil
}
