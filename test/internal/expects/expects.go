package expects

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExpDeleteOrUpdateDevice(t *testing.T, wantBody string, expectBody string) {
	t.Helper()
	assert.Equal(t, wantBody, expectBody)
}
