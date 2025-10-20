package encryption

import (
	"errors"

	"strings"

	"github.com/lib/pq"
)

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var pqErr *pq.Error
	if ok := errors.As(err, &pqErr); ok {

		switch string(pqErr.Code) {
		case "40001": // serialization_failure - transaction retry safe
			return true
		case "40P01": // deadlock_detected
			return true
		}
	}

	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "timeout") {
		return true
	}

	return false
}
