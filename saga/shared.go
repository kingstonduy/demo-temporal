package shared

func IsRetryableError(err error) bool {
	return (string)(err.Error()[0]) == RETRYABLE_ERROR
}
