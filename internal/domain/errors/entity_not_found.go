package errors

import (
	"fmt"

	"github.com/google/uuid"
)

type EntityNotFoundError struct {
	Id uuid.UUID
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("EntityNotFound error for id %s", e.Id)
}
