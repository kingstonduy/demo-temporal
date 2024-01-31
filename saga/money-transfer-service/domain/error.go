package domain

const RETRYABLE_ERROR = "1"
const NONRETRYABLE_ERROR = "0"

func IsRetryableError(err error) bool {
	return (string)(err.Error()[0]) == RETRYABLE_ERROR
}
