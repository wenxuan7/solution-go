package utils

import (
	"fmt"
	"github.com/google/uuid"
)

func UUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("utils: fail to uuid.NewUUID in UUID; %w", err)
	}
	return id.String(), nil
}
