package node

import (
	"fmt"
	"reflect"
)

type DuplicateServiceError struct {
	Kind reflect.Type
}

func (e *DuplicateServiceError) Error() string {
	return fmt.Sprintf("duplicate service: %v", e.Kind)
}
