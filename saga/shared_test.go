package shared

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRetryableError(t *testing.T) {
	myerr := fmt.Errorf("%s%s", RETRYABLE_ERROR, "something bad happened")

	assert.Equal(t, IsRetryableError(myerr), true)

	myerr = fmt.Errorf("%s%s", NONRETRYABLE_ERROR, "something bad happened")
	assert.Equal(t, IsRetryableError(myerr), false)
}
